package mongo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/IvanKyrylov/user-game-api/internal/game"
	"github.com/IvanKyrylov/user-game-api/internal/user"
	"github.com/IvanKyrylov/user-game-api/pkg/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	minSize int = 5000
	// maxSize int = minSize
	maxChan int = 16
)

type UserJson struct {
	Objects []user.User `json:"objects"`
}

type GameJson struct {
	Objects []game.Game `json:"objects"`
}

func Migrate(client *mongo.Database, logger *log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	userCollection := client.Collection("users")

	var userJson UserJson

	usersFile, err := os.Open("resources/users_go.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer usersFile.Close()

	userBytes, err := ioutil.ReadAll(usersFile)
	if err != nil {
		logger.Fatal(err)
	}

	if err := json.Unmarshal(userBytes, &userJson); err != nil {
		logger.Fatal(err)
	}

	usersArray := make([]interface{}, len(userJson.Objects))
	for i, user := range userJson.Objects {
		usersArray[i] = user
	}
	userInserts, err := userCollection.InsertMany(ctx, usersArray)
	if err != nil {
		logger.Fatal(err)
	}

	gameCollection := client.Collection("user_games")

	var gameJson GameJson

	gamesFile, err := os.Open("resources/games.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer usersFile.Close()

	gamesBytes, err := ioutil.ReadAll(gamesFile)
	if err != nil {
		logger.Fatal(err)
	}

	if err := json.Unmarshal(gamesBytes, &gameJson); err != nil {
		logger.Fatal(err)
	}

	rand.Seed(time.Now().Unix())

	wg := sync.WaitGroup{}

	wg.Add(len(userInserts.InsertedIDs))
	ch := make(chan struct{}, maxChan)

	for _, userUUID := range userInserts.InsertedIDs {
		oid := userUUID.(primitive.ObjectID)
		ch <- struct{}{}
		go func() {
			defer wg.Done()

			quantity := rand.Intn(minSize) + minSize

			gamesArray := make([]interface{}, quantity)

			for i := 0; i < quantity; i++ {
				game := gameJson.Objects[rand.Intn(len(gameJson.Objects))]

				game.UserUUID = oid.Hex()
				gamesArray[i] = game
			}
			_, err := gameCollection.InsertMany(ctx, gamesArray)
			if err != nil {
				log.Fatal(err)
			}
			<-ch
		}()
	}

	if err != nil {
		logger.Fatal(err)
	}

	wg.Wait()
	logging.CommonLog.Println("Files have been loaded successfully")

	// gameCollection := client.Collection("user_games")

	// var gameJson GameJson

	// gamesFile, err := os.Open("resources/games.json")
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// defer gamesFile.Close()

	// gamesBytes, err := ioutil.ReadAll(gamesFile)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// if err := json.Unmarshal(gamesBytes, &gameJson); err != nil {
	// 	logger.Fatal(err)
	// }

	// gamesArray := make([]interface{}, len(gameJson.Objects))
	// for i, game := range gameJson.Objects {
	// 	gamesArray[i] = game
	// }
	// gameInserts, err := gameCollection.InsertMany(ctx, gamesArray)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// gamesFile.Close()

	// userCollection := client.Collection("users")

	// var userJson UserJson

	// usersFile, err := os.Open("resources/users_go.json")
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// defer usersFile.Close()

	// userBytes, err := ioutil.ReadAll(usersFile)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// if err := json.Unmarshal(userBytes, &userJson); err != nil {
	// 	logger.Fatal(err)
	// }

	// usersArray := make([]interface{}, len(userJson.Objects))

	// rand.Seed(time.Now().Unix())

	// wg := sync.WaitGroup{}

	// wg.Add(len(userJson.Objects))
	// // ch := make(chan struct{}, maxChan)

	// for index, userObj := range userJson.Objects {
	// 	go func(index int, user user.User) {
	// 		defer wg.Done()

	// 		quantity := rand.Intn(minSize) + minSize
	// 		user.GamesId = make([]string, quantity)
	// 		for i := 0; i < quantity; i++ {
	// 			id := gameInserts.InsertedIDs[rand.Intn(len(gameJson.Objects))].(primitive.ObjectID)
	// 			user.GamesId[i] = id.Hex()
	// 		}
	// 		usersArray[index] = user
	// 	}(index, userObj)

	// }
	// wg.Wait()
	// _, err = userCollection.InsertMany(ctx, usersArray)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// usersFile.Close()

	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// logging.CommonLog.Println("Files have been loaded successfully")

}
