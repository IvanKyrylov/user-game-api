package game

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/IvanKyrylov/user-game-api/internal/apperror"
)

const (
	gamesURL = "/api/games"
	gameURL  = "/api/game/"
	// gameByUserURL = "/api/games/user/"
)

type Handler struct {
	Logger      *log.Logger
	GameService Service
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(gameURL, apperror.Middleware(h.GetGame))
	// router.HandleFunc(gameByUserURL, apperror.Middleware(h.GetGamesByPlayer))
	router.HandleFunc(gamesURL, apperror.Middleware(h.GetAllGames))
}

func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}
	h.Logger.Println("GET GAME")

	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get id from path")
	id := r.URL.Path[len(gameURL):]
	if id == "" {
		return apperror.BadRequestError("id query parameter is required and must be a comma separated integers")
	}

	game, err := h.GameService.GetById(r.Context(), id)
	if err != nil {
		return err
	}

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(gameBytes)

	return nil
}

func (h *Handler) GetGamesByPlayer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}

	h.Logger.Println("GET GAMES BY PLAYER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get player id from URL")
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	h.Logger.Println("get limit from URL")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		return apperror.BadRequestError("limit query parameter is required positive integers")
	}

	h.Logger.Println("get limit from URL")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		return apperror.BadRequestError("page query parameter is required positive integers")
	}

	games, err := h.GameService.GetByPlayer(r.Context(), uuid, int64(limit), int64(page))
	if err != nil {
		return err
	}

	gamesBytes, err := json.Marshal(games)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(gamesBytes)

	return nil
}

func (h *Handler) GetAllGames(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}

	h.Logger.Println("GET ALL GAMES")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get limit from URL")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		return apperror.BadRequestError("limit query parameter is required positive integers")
	}

	h.Logger.Println("get limit from URL")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		return apperror.BadRequestError("page query parameter is required positive integers")
	}

	games, err := h.GameService.GetAll(r.Context(), int64(limit), int64(page))
	if err != nil {
		return err
	}

	gamesBytes, err := json.Marshal(games)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(gamesBytes)

	return nil
}
