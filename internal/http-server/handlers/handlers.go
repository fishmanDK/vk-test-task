package handlers

import (
	"github.com/fishmanDK/vk-test-task/internal/service"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"

	_ "github.com/fishmanDK/vk-test-task/docs"
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

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	authRouter := r.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/sign-in", h.signIn).Methods("POST")
		authRouter.HandleFunc("/sign-up", h.signUp).Methods("POST")
	}

	filmsRouter := r.PathPrefix("/films").Subrouter()
	filmsRouter.Use(h.authMiddleware(logger))
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
	actorsRouter.Use(h.authMiddleware(logger))
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
