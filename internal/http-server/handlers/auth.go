package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	vk_test_task "vk-test-task"
	"vk-test-task/internal/service/response"
)

func (h *Handlers) signIn(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.signIn"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	var input vk_test_task.User
	w.Header().Set("Content-Type", "encoding/json")
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	r.Body.Close()

	tokens, err := h.service.Auth.Authentication(input)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func (h *Handlers) signUp(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.signUp"

	logger := r.Context().Value(loggerKey).(*slog.Logger)

	w.Header().Set("Content-Type", "encoding/json")

	var input vk_test_task.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.Auth.CreateUser(input)
	if err != nil {
		logger.Error(op, slog.String("err", err.Error()))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.OK(createSUCCESS))
}
