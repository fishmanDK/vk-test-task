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

type Actor struct {
	Name  string    `json:"name"`
	Male  string    `json:"male"`
	Birth time.Time `json:"birth"`
	Films []Film    `json:"films"`
}

func (h *Handlers) getAllActors(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getAllActors"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	actors, err := h.service.Actors.GetAllActors()

	w.Header().Set("Content-Type", "encoding/json")

	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)
}

func (h *Handlers) getActorById(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getActorById"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "encoding/json")

	actor, err := h.service.Actors.GetActorByID(id)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actor)
}

func (h *Handlers) createActor(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.createActor"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	if pDataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser); pDataUser.Role != admin {
		logger.Error(op, slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", pDataUser.ID, pDataUser.Email)))
		newErrorResponse(w, http.StatusForbidden, accessDenied)
		return
	}

	var input vk_test_task.NewActor
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Actors.CreateActor(input)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.OK(createSUCCESS))
}

func (h *Handlers) changeActor(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.changeActor"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	if pDataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser); pDataUser.Role != admin {
		logger.Error(op, slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", pDataUser.ID, pDataUser.Email)))
		newErrorResponse(w, http.StatusForbidden, accessDenied)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var input vk_test_task.ChangeDataActor
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Actors.ChangeActor(id, input)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.OK(changeSUCCESS))
}

func (h *Handlers) deleteActor(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.deleteActor"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")
	if pDataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser); pDataUser.Role != admin {
		logger.Error(op, slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", pDataUser.ID, pDataUser.Email)))
		newErrorResponse(w, http.StatusForbidden, accessDenied)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteActor(id)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.OK(deleteSUCCESS))
}
