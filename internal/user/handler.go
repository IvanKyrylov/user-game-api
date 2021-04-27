package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/IvanKyrylov/user-game-api/internal/apperror"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/user/"
)

type Handler struct {
	Logger      *log.Logger
	UserService Service
}

func (h *Handler) Register(router *http.ServeMux) {
	router.HandleFunc(userURL, apperror.Middleware(h.GetUser))
	// router.HandleFunc(userURL, apperror.Middleware(h.GetUserByName))
	router.HandleFunc(usersURL, apperror.Middleware(h.GetAllUsers))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}
	h.Logger.Println("GET USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get uuid from path")
	uuid := r.URL.Path[len(userURL):]
	if uuid == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	user, err := h.UserService.GetById(r.Context(), uuid)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

func (h *Handler) GetUserByName(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}
	h.Logger.Println("GET USER BY NAME")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get name from URL")
	uuid := r.URL.Query().Get("name")
	if uuid == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	user, err := h.UserService.GetByName(r.Context(), uuid)
	if err != nil {
		return err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apperror.BadRequestError("metod GET")
	}

	h.Logger.Println("GET USERS")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Println("get limit from URL")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		return apperror.BadRequestError("limit query parameter is required positive integers")
	}

	h.Logger.Println("get limit from URL")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 0 {
		return apperror.BadRequestError("page query parameter is required positive integers")
	}

	users, err := h.UserService.GetAll(r.Context(), int64(limit), int64(page))
	if err != nil {
		return err
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(usersBytes)
	return nil
}
