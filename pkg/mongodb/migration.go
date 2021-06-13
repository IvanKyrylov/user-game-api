package mongo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserJSON struct {
	Email     string `json:"email" binding:"required" validate:"email"`
	LastName  string `json:"last_name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	City      string `json:"city" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
	BirthDate string `json:"birth_date" binding:"required"`
}

type UserGameJSON struct {
	PointsGained int    `json:"points_gained,string"`
	WinStatus    int8   `json:"win_status,string"`
	GameType     int8   `json:"game_type,string"`
	Created      string `json:"created"`
}

type usersJSONRes struct {
	Objects []UserJSON `json:"objects"`
}

type userGameJSONRes struct {
	Objects []UserGameJSON `json:"objects"`
}

func parseUsersJSON(log *log.Logger) []UserJSON {
	jsonFile, err := os.Open("resources/users_go.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var fileRes usersJSONRes
	err = json.Unmarshal(byteValue, &fileRes)
	if err != nil {
		log.Fatal(err)
	}

	return fileRes.Objects
}

func parseUserGamesJSON(log *log.Logger) []UserGameJSON {
	jsonFile, err := os.Open("resources/games.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var fileRes userGameJSONRes
	err = json.Unmarshal(byteValue, &fileRes)
	if err != nil {
		log.Fatal(err)
	}

	return fileRes.Objects
}

func insertUsers(client *mongo.Database, log *log.Logger) []interface{} {
	users := parseUsersJSON(log)
	userCollection := client.Collection("users")
	log.Println("Started inserting users...")
	usersMongo := make([]interface{}, 0, len(users))
	for _, user := range users {
		birthDate, err := time.Parse("Monday, January 2, 2006 3:04 PM", user.BirthDate)
		if err != nil {
			log.Println(err)
		}
		validYear := rand.Intn(2015-1980) + 1980
		validBirthDate := time.Date(validYear, birthDate.Month(), birthDate.Day(), birthDate.Hour(), birthDate.Minute(), birthDate.Second(), birthDate.Nanosecond(), birthDate.Location())
		newUser := bson.D{
			{"email", user.Email},
			{"last_name", user.LastName},
			{"country", user.Country},
			{"city", user.City},
			{"gender", user.Gender},
			{"birth_date", primitive.NewDateTimeFromTime(validBirthDate)},
			{"rating", int64(0)},
		}
		usersMongo = append(usersMongo, newUser)
	}
	insertRes, _ := userCollection.InsertMany(context.Background(), usersMongo)
	log.Printf("insertedUsers: %v", len(insertRes.InsertedIDs))
	log.Println("Finished inserting users")
	return insertRes.InsertedIDs
}

func insertUserGames(client *mongo.Database, log *log.Logger, userIds []interface{}) {
	userGames := parseUserGamesJSON(log)
	userGamesCollection := client.Collection("user_games")

	log.Println("Started inserting user games...")
	var foundUserIds = make([]primitive.ObjectID, 0, len(userIds))
	for _, userID := range userIds {
		idToObjectID, ok := userID.(primitive.ObjectID)
		if !ok {
			log.Println("Cannot cast userId to ObjectID")
			continue
		}
		foundUserIds = append(foundUserIds, idToObjectID)
	}

	for _, foundUserID := range foundUserIds {
		var randGames = make([]interface{}, 0)

		randCof := rand.Intn(1000)

		for i := 0; i < randCof+5000; i++ {
			randGame := userGames[rand.Intn(len(userGames))]
			created, err := time.Parse("1/2/2006 3:04 PM", randGame.Created)
			if err != nil {
				log.Fatal(err)
				continue
			}
			var newUserGame = bson.D{
				{"_id", primitive.NewObjectID()},
				{"points_gained", randGame.PointsGained},
				{"win_status", randGame.WinStatus},
				{"game_type", randGame.GameType},
				{"user_id", foundUserID},
				{"created", created},
			}
			randGames = append(randGames, newUserGame)
		}
		_, err := userGamesCollection.InsertMany(context.Background(), randGames)
		if err != nil {
			log.Fatal(err)
		}

		userCollection := client.Collection("users")
		replacement := bson.D{{"$set", bson.D{{"rating", int64(len(randGames))}}}}
		_, err = userCollection.UpdateOne(context.Background(), bson.D{{"_id", foundUserID}}, replacement)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Inserted user games")
}

func createIndexes(client *mongo.Database) (err error) {

	_, err = client.Collection("user_games").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"user_id": 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	return nil
}

func Migrate(client *mongo.Database, log *log.Logger) {
	rand.Seed(time.Now().Unix())
	log.Println("Started DB initialization...")

	err := createIndexes(client)
	if err != nil {
		log.Fatal(err)
		return
	}
	insertedIds := insertUsers(client, log)
	insertUserGames(client, log, insertedIds)
}
