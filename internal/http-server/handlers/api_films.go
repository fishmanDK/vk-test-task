package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"time"
	vk_test_task "vk-test-task"
	"vk-test-task/internal/service"
	"vk-test-task/internal/service/response"
)

type Film struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Actors      []Actor   `json:"actors"`
}

func (h *Handlers) getAllFilms(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getAllFilms"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	queryParams := r.URL.Query()
	order := queryParams.Get("order")
	q := queryParams.Get("q")

	w.Header().Set("Content-Type", "encoding/json")

	films, err := h.service.Films.GetAllFilms(order, q)

	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(films)
}

func (h *Handlers) getFilmByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getFilmByID"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	vars := mux.Vars(r)
	id := vars["id"]
	film, err := h.service.Films.GetFilm(id)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(film)
}

func (h *Handlers) searchFilm(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.searchFilm"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	queryParams := r.URL.Query()
	byTitle := queryParams.Get("byTitle")
	byActor := queryParams.Get("byActor")

	w.Header().Set("Content-Type", "encoding/json")

	film_s, err := h.service.Films.SearchFilm(byTitle, byActor)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(film_s)
}

func (h *Handlers) createFilm(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.createFilm"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	if pDataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser); pDataUser.Role != admin {
		logger.Error(op, slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", pDataUser.ID, pDataUser.Email)))
		newErrorResponse(w, http.StatusForbidden, accessDenied)
		return
	}

	var input vk_test_task.CreateFilm
	w.Header().Set("Content-Type", "encoding/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Films.CreateFilm(input)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.OK(createSUCCESS))
}

func (h *Handlers) changeFilm(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.changeFilm"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	if pDataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser); pDataUser.Role != admin {
		logger.Error(op, slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", pDataUser.ID, pDataUser.Email)))
		newErrorResponse(w, http.StatusForbidden, accessDenied)
		return
	}

	var input vk_test_task.ChangeDataFilm

	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "encoding/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Films.ChangeFilm(id, input)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response.OK(changeSUCCESS))
}

func (h *Handlers) deleteFilm(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.deleteFilm"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	if pDataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser); pDataUser.Role != admin {
		logger.Error(op, slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", pDataUser.ID, pDataUser.Email)))
		newErrorResponse(w, http.StatusForbidden, accessDenied)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteFilmById(id)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.OK(deleteSUCCESS))
}
