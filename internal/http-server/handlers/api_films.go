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

// @Summary GetAllFilms
// @Security ApiKeyAuth
// @Tags films
// @Description Getting a list of movies with all its actors.<br />The list of films can be sorted using the following fields:<br />----order - can take values "rating", "name", "date";<br /> ----q - can take values "asc", "desc";<br /> By default the list is sorted by rating in descending order.
// @ID get-all-films
// @Accept  json
// @Produce  json
// @Param        order    query     string  false  "name search by byTitle"
// @Param        q    query     string  false  "name search by byActor"
// @Success 200 {array} storage.Film "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /films [get]
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

// @Summary GetFilmByID
// @Security ApiKeyAuth
// @Tags films
// @Description Get film by id
// @ID get-film-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Film ID"
// @Success 200 {object} storage.Film "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /films/{id} [get]
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

// @Summary SearchFilm
// @Security ApiKeyAuth
// @Tags films
// @Description Search Film(s) by title or actor(first_name)
// @ID search-film
// @Accept  json
// @Produce  json
// @Param        byTitle    query     string  false  "name search by byTitle"  Format(email)
// @Param        byActor    query     string  false  "name search by byActor"
// @Success 200 {array} storage.Film "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /films/ [get]
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

// @Summary CreateFilm
// @Security ApiKeyAuth
// @Tags films
// @Description Create new film.<br />1 <= len(title) >= 150; 0 <= len(description) >= 1000; 0 <= rating >= 10; date format <yyyy-mm-dd>;<br />The actors field must contain actor IDs.
// @ID search-film
// @Accept  json
// @Produce  json
// @Param input body vk_test_task.CreateFilm true "new film"
// @Success 200 {object} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /films/ [post]
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

// @Summary ChangeFilm
// @Security ApiKeyAuth
// @Tags films
// @Description Change film by id.<br />1 <= len(title) >= 150; 0 <= len(description) >= 1000; 0 <= rating >= 10; date format <yyyy-mm-dd>;<br />The actors field must contain actor IDs. If you want to remove information about a specific actor, then put - in front of his ID <"-id">
// @ID change-film
// @Accept  json
// @Produce  json
// @Param input body vk_test_task.ChangeDataFilm true "new data film"
// @Param id path string true "Film ID"
// @Success 200 {object} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /films/{id} [PATCH]
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

// @Summary DeleteFilm
// @Security ApiKeyAuth
// @Tags films
// @Description Delete film by id.
// @ID delete-film
// @Accept  json
// @Produce  json
// @Param id path string true "Film ID"
// @Success 200 {object} response.Response "Success response"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /films/{id} [delete]
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
