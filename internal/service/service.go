package service

import (
	"github.com/fishmanDK/vk-test-task/internal/storage"
)

type Service struct {
	Actors
	Films
	Auth
}

func MustService(db *storage.StorageServ) (*Service, error) {
	return &Service{
		Actors: MustActorsService(db),
		Films:  MustServiceFilms(db),
		Auth:   NewAuthService(db),
	}, nil
}
