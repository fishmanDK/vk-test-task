package storage

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
	vk_test_task "vk-test-task"
)

type Actors interface {
	GetAllActors() (*[]Actor, error)
	GetActorByID(id string) (*Actor, error)
	CreateActor(newActor vk_test_task.NewActor) error
	ChangeActor(idActor string, changeDataActor vk_test_task.ChangeDataActor) error
	DeleteActor(idActor string) error
}

type ActorsStorage struct {
	db *sqlx.DB
}

func MustActorsStorage(db *sqlx.DB) *ActorsStorage {
	return &ActorsStorage{
		db: db,
	}
}

type Actor struct {
	IdActor   string `db:"id,omitempty"`
	FirstName string `db:"first_name,omitempty"`
	LastName  string `db:"last_name,omitempty"`
	Surname   string `db:"surname,omitempty"`
	Sex       string `db:"sex,omitempty"`
	Birthday  string `db:"birthday,omitempty"`
	Films     []Film `db:"films,omitempty"`
}

type customActor struct {
	IdActor   sql.NullString `db:"id,omitempty"`
	FirstName sql.NullString `db:"first_name,omitempty"`
	LastName  sql.NullString `db:"last_name,omitempty"`
	Surname   sql.NullString `db:"surname,omitempty"`
	Sex       sql.NullString `db:"sex,omitempty"`
	Birthday  sql.NullString `db:"birthday,omitempty"`
	Films     []Film         `db:"films,omitempty"`
}

func (a *ActorsStorage) GetAllActors() (*[]Actor, error) {
	const op = "storage.GetAllActors"

	tx, err := a.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var query string
	query = "select films.id, films.title, films.description, films.release_date, films.rating, actors.id, actors.first_name, actors.last_name, actors.surname, actors.birthday, actors.sex from films_actors join films on films_actors.film_id = films.id right join actors on films_actors.actor_id = actors.id;"

	rows, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	actors := []Actor{}
	last_actor := customActor{}
	for rows.Next() {
		custom_film := customFilm{}
		custom_actor := customActor{}

		err := rows.Scan(&custom_film.ID, &custom_film.Title, &custom_film.Description, &custom_film.Date, &custom_film.Rating, &custom_actor.IdActor, &custom_actor.FirstName, &custom_actor.LastName, &custom_actor.Surname, &custom_actor.Birthday, &custom_actor.Sex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		film := Film{
			ID:          custom_film.ID.String,
			Title:       custom_film.Title.String,
			Description: custom_film.Description.String,
			Date:        custom_film.Date.String,
			Actors:      custom_film.Actors,
			Rating:      custom_film.Rating.String,
		}
		actor := Actor{
			IdActor:   custom_actor.IdActor.String,
			FirstName: custom_actor.FirstName.String,
			LastName:  custom_actor.LastName.String,
			Surname:   custom_actor.Surname.String,
			Sex:       custom_actor.Sex.String,
			Birthday:  custom_actor.Birthday.String,
		}

		if last_actor.eq(custom_actor) {
			if !custom_film.eq(customFilm{}) {
				actors[len(actors)-1].Films = append(actors[len(actors)-1].Films, film)
			}
		} else {
			if !custom_film.eq(customFilm{}) {
				actor.Films = append(actor.Films, film)
			}

			actors = append(actors, actor)
			last_actor = custom_actor
		}

	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &actors, nil
}

func (a *ActorsStorage) GetActorByID(id string) (*Actor, error) {
	const op = "storage.GetActorByID"

	tx, err := a.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var query string
	query = "select films.id, films.title, films.description, films.release_date, films.rating, actors.id, actors.first_name, actors.last_name, actors.surname, actors.birthday, actors.sex from films_actors join films on films_actors.film_id = films.id right join actors on films_actors.actor_id = actors.id WHERE actor_id = $1"

	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resActor := Actor{}
	last_actor := customActor{}
	for rows.Next() {
		custom_film := customFilm{}
		custom_actor := customActor{}

		err := rows.Scan(&custom_film.ID, &custom_film.Title, &custom_film.Description, &custom_film.Date, &custom_film.Rating, &custom_actor.IdActor, &custom_actor.FirstName, &custom_actor.LastName, &custom_actor.Surname, &custom_actor.Birthday, &custom_actor.Sex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		film := Film{
			ID:          custom_film.ID.String,
			Title:       custom_film.Title.String,
			Description: custom_film.Description.String,
			Date:        custom_film.Date.String,
			Actors:      custom_film.Actors,
			Rating:      custom_film.Rating.String,
		}
		actor := Actor{
			IdActor:   custom_actor.IdActor.String,
			FirstName: custom_actor.FirstName.String,
			LastName:  custom_actor.LastName.String,
			Surname:   custom_actor.Surname.String,
			Sex:       custom_actor.Sex.String,
			Birthday:  custom_actor.Birthday.String,
		}

		if last_actor.eq(custom_actor) {
			if !custom_film.eq(customFilm{}) {
				actor.Films = append(actor.Films, film)
			}
		} else {
			if !custom_film.eq(customFilm{}) {
				actor.Films = append(actor.Films, film)
			}

			resActor = actor
			last_actor = custom_actor
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &resActor, nil
}

func (a *ActorsStorage) CreateActor(newActor vk_test_task.NewActor) error {
	const op = "storage.CreateActor"

	query := "INSERT INTO actors (first_name, last_name, surname, sex, birthday) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	var newIdActor int
	err = tx.QueryRow(query, newActor.FirstName, newActor.LatName, newActor.Surname, newActor.Sex, newActor.Birthday).Scan(&newIdActor)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, idFilm := range newActor.Films {
		nn, _ := strconv.Atoi(idFilm)
		query := "INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)"
		_, err := tx.Exec(query, newIdActor, nn)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *ActorsStorage) ChangeActor(idActor string, changeDataActor vk_test_task.ChangeDataActor) error {
	const op = "storage.ChangeActor"

	tx, err := a.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if changeDataActor.FirstName != "" {
		query := "UPDATE actors SET first_name = $1 WHERE id = $2"
		_, err = tx.Exec(query, changeDataActor.FirstName, idActor)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changeDataActor.LastName != "" {
		query := "UPDATE actors SET last_name = $1 WHERE id = $2"
		_, err = tx.Exec(query, changeDataActor.LastName, idActor)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changeDataActor.Surname != "" {
		query := "UPDATE actors SET surname = $1 WHERE id = $2"
		_, err = tx.Exec(query, changeDataActor.Surname, idActor)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changeDataActor.Birthday != "" {
		query := "UPDATE actors SET birthday = $1 WHERE id = $2"
		_, err = tx.Exec(query, changeDataActor.Birthday, idActor)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changeDataActor.Sex != "" {
		query := "UPDATE actors SET sex = $1 WHERE id = $2"
		_, err = tx.Exec(query, changeDataActor.Sex, idActor)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	for _, idFilm := range changeDataActor.Films {
		if idFilm[0] == '-' {
			query := "DELETE FROM films_actors WHERE film_id = $1 AND actor_id = $2"
			_, err = tx.Exec(query, string(idFilm[1]), idActor)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		} else {
			query := "INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)"
			_, err := tx.Exec(query, idFilm, idActor)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	tx.Commit()
	return nil
}

func (a *ActorsStorage) DeleteActor(idActor string) error {
	const op = "storage.DeleteActor"

	tx, err := a.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "delete from films_actors where actor_id = $1"
	_, err = tx.Exec(query, idActor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = "delete from actors where id = $1"
	_, err = tx.Exec(query, idActor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Actor) formatBirth() error {
	const op = "storage.formatBirth"

	birth, err := time.Parse("2006-01-02", a.Birthday)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.Birthday = birth.Format("2006-01-02")
	return nil
}

func (a Actor) equals(other Actor) bool {
	return a.IdActor == other.IdActor && a.FirstName == other.FirstName && a.LastName == other.LastName && a.Birthday == other.Birthday && a.Surname == other.Surname && a.Sex == other.Sex
}

func (a customActor) eq(other customActor) bool {
	return a.IdActor.String == other.IdActor.String &&
		a.FirstName.String == other.FirstName.String &&
		a.LastName.String == other.LastName.String &&
		a.Birthday.String == other.Birthday.String &&
		a.Surname.String == other.Surname.String &&
		a.Sex.String == other.Sex.String
}
