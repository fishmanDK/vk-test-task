package service

import (
	"fmt"
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"github.com/fishmanDK/vk-test-task/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Actors
type Actors interface {
	GetAllActors() (*[]storage.Actor, error)
	GetActorByID(id string) (*storage.Actor, error)
	CreateActor(newActor vk_test_task.NewActor) error
	ChangeActor(idActor string, changeDataActor vk_test_task.ChangeDataActor) error
	DeleteActor(idActor string) error
}

type ActorsService struct {
	db *storage.StorageServ
}

func MustActorsService(db *storage.StorageServ) *ActorsService {
	return &ActorsService{
		db: db,
	}
}

func (a *ActorsService) GetAllActors() (*[]storage.Actor, error) {
	const op = "actors.GetAllActors"

	actors, err := a.db.Actors.GetAllActors()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}

func (a *ActorsService) GetActorByID(id string) (*storage.Actor, error) {
	const op = "actors.GetActorByID"

	actor, err := a.db.Actors.GetActorByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actor, nil
}

func (a *ActorsService) CreateActor(newActor vk_test_task.NewActor) error {
	const op = "service.CreateActor"

	err := a.db.Actors.CreateActor(newActor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *ActorsService) ChangeActor(idActor string, changeDataActor vk_test_task.ChangeDataActor) error {
	const op = "service.ChangeActor"

	err := a.db.Actors.ChangeActor(idActor, changeDataActor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *ActorsService) DeleteActor(idActor string) error {
	const op = "service.DeleteActor"

	err := a.db.Actors.DeleteActor(idActor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
