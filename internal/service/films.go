package service

import (
	"fmt"
	vk_test_task "vk-test-task"
	"vk-test-task/internal/storage"
)

type FilmsService struct {
	repo *storage.StorageServ
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Films
type Films interface {
	GetAllFilms(orderBy, q string) (*[]storage.Film, error)
	SearchFilm(byTitle, byActor string) (*[]storage.Film, error)
	CreateFilm(newFilm vk_test_task.CreateFilm) error
	GetFilm(id string) (*storage.Film, error)
	DeleteFilmById(id string) error
	ChangeFilm(id string, changedDataFilm vk_test_task.ChangeDataFilm) error
}

func MustServiceFilms(repo *storage.StorageServ) *FilmsService {
	return &FilmsService{
		repo: repo,
	}
}

func (f *FilmsService) GetAllFilms(orderBy, q string) (*[]storage.Film, error) {
	const op = "service.GetAllFilms"

	films, err := f.repo.Films.GetAllFilms(orderBy, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return films, nil
}

func (f *FilmsService) SearchFilm(byTitle, byActor string) (*[]storage.Film, error) {
	const op = "service.SearchFilm"

	film_s, err := f.repo.Films.SearchFilm(byTitle, byActor)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return film_s, nil
}

func (f *FilmsService) CreateFilm(newFilm vk_test_task.CreateFilm) error {
	const op = "service.CreateFilm"

	err := f.repo.Films.CreateFilm(newFilm)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (f *FilmsService) GetFilm(id string) (*storage.Film, error) {
	const op = "service.GetFilm"

	film, err := f.repo.Films.GetFilm(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return film, nil
}

func (f *FilmsService) DeleteFilmById(id string) error {
	const op = "service.DeleteFilmById"

	err := f.repo.Films.DeleteFilmById(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (f *FilmsService) ChangeFilm(id string, changedDataFilm vk_test_task.ChangeDataFilm) error {
	const op = "service.ChangeFilm"

	err := f.repo.Films.ChangeFilm(id, changedDataFilm)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
