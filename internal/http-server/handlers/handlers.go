package handlers

import (
	"github.com/gorilla/mux"
	"log/slog"
	"vk-test-task/internal/service"
)

type Handlers struct {
	service *service.Service
}

func MustHandlers(srvc *service.Service) (*Handlers, error) {
	return &Handlers{
		service: srvc,
	}, nil
}

func (h *Handlers) InitRouts(logger *slog.Logger) *mux.Router {
	r := mux.NewRouter()
	r.Use(h.loggerMiddleware(logger))

	authRouter := r.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/signIn", h.signIn).Methods("POST")
		authRouter.HandleFunc("/signUp", h.signUp).Methods("POST")
	}

	filmsRouter := r.PathPrefix("/films").Subrouter()
	filmsRouter.Use(h.authMiddleware)
	{
		filmsRouter.HandleFunc("", h.getAllFilms).Methods("GET")

		onlyFilmRouter := filmsRouter.PathPrefix("/").Subrouter()
		{

			onlyFilmRouter.HandleFunc("/", h.searchFilm).Methods("GET")
			onlyFilmRouter.HandleFunc("/", h.createFilm).Methods("POST")
			onlyFilmRouter.HandleFunc("/{id}", h.getFilmByID).Methods("GET")
			onlyFilmRouter.HandleFunc("/{id}", h.changeFilm).Methods("PATCH")
			onlyFilmRouter.HandleFunc("/{id}", h.deleteFilm).Methods("DELETE")
		}
	}

	actorsRouter := r.PathPrefix("/actors").Subrouter()
	actorsRouter.Use(h.authMiddleware)
	{
		actorsRouter.HandleFunc("", h.getAllActors).Methods("Get")

		onlyActorRouter := actorsRouter.PathPrefix("/").Subrouter()
		{
			onlyActorRouter.HandleFunc("/", h.createActor).Methods("POST")
			onlyActorRouter.HandleFunc("/{id}", h.getActorById).Methods("GET")
			onlyActorRouter.HandleFunc("/{id}", h.changeActor).Methods("PATCH")
			onlyActorRouter.HandleFunc("/{id}", h.deleteActor).Methods("DELETE")
		}

	}

	return r
}
