package handlers

import (
	"encoding/json"
	"fmt"
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"github.com/fishmanDK/vk-test-task/internal/service"
	"github.com/fishmanDK/vk-test-task/internal/service/response"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

// @Summary GetAllActors
// @Security ApiKeyAuth
// @Tags actors
// @Description Getting a list of actors with all its films.
// @ID get-all-actors
// @Accept  json
// @Produce  json
// @Success 200 {array} storage.Actor "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /actors [get]
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

// @Summary GetActorById
// @Security ApiKeyAuth
// @Tags actors
// @Description Get actor by id.
// @ID get-actor-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Actor ID"
// @Success 200 {array} storage.Film "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /actors/{id} [get]
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

// @Summary CreateActor
// @Security ApiKeyAuth
// @Tags actors
// @Description Create new actor. All fields are required except Films and Surname.
// @ID create-actor
// @Accept  json
// @Produce  json
// @Param input body vk_test_task.NewActor true "new actor"
// @Success 200 {array} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /actors/ [post]
func (h *Handlers) createActor(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.createActor"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	dataUser := r.Context().Value(parseDataUser).(*service.ParseDataUser)
	if dataUser.Role != admin {
		logger.Error("op", slog.String(accessDenied, fmt.Sprintf("access denied to user with ID: %d, email: %s", dataUser.ID, dataUser.Email)))
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

// @Summary ChangeActor
// @Security ApiKeyAuth
// @Tags actors
// @Description Change actor by id.<br /> Date format <yyyy-mm-dd>;<br />The films field must contain film IDs. If you want to remove information about a specific film, then put "-" in front of his ID <"-id">
// @ID change-actor
// @Accept  json
// @Produce  json
// @Param id path string true "Actor ID"
// @Param input body vk_test_task.ChangeDataActor true "new data actor"
// @Success 200 {array} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /actors/{id} [patch]
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

// @Summary DeleteActor
// @Security ApiKeyAuth
// @Tags actors
// @Description Delete actor by id.
// @ID delete-actor
// @Accept  json
// @Produce  json
// @Param id path string true "Actor ID"
// @Success 200 {array} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /actors/{id} [delete]
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
