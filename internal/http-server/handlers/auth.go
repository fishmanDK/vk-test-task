package handlers

import (
	"encoding/json"
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"github.com/fishmanDK/vk-test-task/internal/service/response"
	"log/slog"
	"net/http"
)

// @Summary SignIn
// @Tags auth
// @Description Login with email and password
// @ID login
// @Accept  json
// @Produce  json
// @Param input body vk_test_task.User true "User credentials"
// @Success 200 {object} vk_test_task.Tokens "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /auth/sign-in [post]
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

// @Summary SignUp
// @Tags auth
// @Description Сreate account. The field "role" must have one of two values (ordinary, admin)
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body vk_test_task.CreateUser true "account info"
// @Success 200 {object} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /auth/sign-up [post]
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
