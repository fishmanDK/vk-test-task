package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	vk_test_task "vk-test-task"
)

const (
	sortASC    = "asc"
	sortDESC   = "desc"
	sortTITLE  = "title"
	sortRATING = "rating"
	sortDATE   = "date"
)

type Films interface {
	GetAllFilms(orderBy, q string) (*[]Film, error)
	GetFilm(id string) (*Film, error)
	SearchFilm(byTitle, byActor string) (*[]Film, error)
	CreateFilm(newFilm vk_test_task.CreateFilm) error
	DeleteFilmById(id string) error
	ChangeFilm(id string, changedDataFilm vk_test_task.ChangeDataFilm) error
}

type FilmsStorage struct {
	db *sqlx.DB
}

func MustFilmsStorage(db *sqlx.DB) *FilmsStorage {
	return &FilmsStorage{
		db: db,
	}
}

type Film struct {
	ID          string  `db:"id,omitempty"`
	Title       string  `db:"title,omitempty"`
	Description string  `db:"description,omitempty"`
	Date        string  `db:"release_date,omitempty"`
	Actors      []Actor `db:"actors,omitempty"`
	Rating      string  `db:"rating,omitempty"`
}

type customFilm struct {
	ID          sql.NullString `db:"id,omitempty"`
	Title       sql.NullString `db:"title,omitempty"`
	Description sql.NullString `db:"description,omitempty"`
	Date        sql.NullString `db:"release_date,omitempty"`
	Actors      []Actor        `db:"actors,omitempty"`
	Rating      sql.NullString `db:"rating,omitempty"`
}

func (f *FilmsStorage) GetAllFilms(orderBy, q string) (*[]Film, error) {
	const op = "storage.GetAllFilms"

	var query string
	query = "SELECT films.id, films.title, films.description, films.release_date, films.rating, actors.id, actors.first_name, actors.last_name, actors.surname, actors.birthday, actors.sex FROM films_actors JOIN films ON films_actors.film_id = films.id JOIN actors ON films_actors.actor_id = actors.id "

	switch orderBy {
	case sortRATING:
		query += "ORDER BY rating"
	case sortDATE:
		query += "ORDER BY date"
	case sortTITLE:
		query += "ORDER BY title"
	default:
		query += "ORDER BY rating"
	}
	switch q {
	case sortASC:
		query += " ASC"
	case sortDESC:
		query += " DESC"
	default:
		query += " DESC"
	}

	tx, err := f.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	films := []Film{}
	last_film := Film{}
	for rows.Next() {
		film := Film{}
		actor := Actor{}

		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.Date, &film.Rating, &actor.IdActor, &actor.FirstName, &actor.LastName, &actor.Surname, &actor.Birthday, &actor.Sex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		//err = film.formatDate()
		//if err != nil {
		//	return nil, fmt.Errorf("%s: %w", op, err)
		//}

		//err = actor.formatBirth()
		//if err != nil {
		//	return nil, fmt.Errorf("%s: %w", op, err)
		//}

		if last_film.equals(film) {
			films[len(films)-1].Actors = append(films[len(films)-1].Actors, actor)

		} else {
			film.Actors = append(film.Actors, actor)
			films = append(films, film)
			last_film = film
		}

	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &films, nil
}

func (f *FilmsStorage) SearchFilm(byTitle, byActor string) (*[]Film, error) {
	const op = "storage.SearchFilm"

	query := "SELECT films.id, films.title, films.description, films.release_date, films.rating, actors.first_name, actors.last_name, actors.surname, actors.birthday, actors.sex FROM films_actors JOIN films ON films_actors.film_id = films.id JOIN actors ON films_actors.actor_id = actors.id WHERE films.title LIKE '%' || $1 || '%' AND actors.first_name LIKE '%' || $2 || '%';"

	tx, err := f.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.Query(query, byTitle, byActor)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	films := []Film{}
	last_film := Film{}
	for rows.Next() {
		film := Film{}
		actor := Actor{}
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.Date, &film.Rating, &actor.FirstName, &actor.LastName, &actor.Surname, &actor.Birthday, &actor.Sex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		//err = film.formatDate()
		//if err != nil {
		//	return nil, fmt.Errorf("%s: %w", op, err)
		//}

		//err = actor.formatBirth()
		//if err != nil {
		//	return nil, fmt.Errorf("%s: %w", op, err)
		//}

		if last_film.equals(film) {
			films[len(films)-1].Actors = append(films[len(films)-1].Actors, actor)

		} else {
			film.Actors = append(film.Actors, actor)
			films = append(films, film)
			last_film = film
		}

	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &films, nil

}

func (f *FilmsStorage) CreateFilm(newFilm vk_test_task.CreateFilm) error {
	const op = "storage.CreateFilm"

	query := "INSERT INTO films (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id"

	tx, err := f.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	var newIdFilm int
	err = tx.QueryRow(query, newFilm.Title, newFilm.Description, newFilm.Date, newFilm.Rating).Scan(&newIdFilm)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, idActor := range newFilm.Actors {
		nn, _ := strconv.Atoi(idActor)
		query := "INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)"
		_, err := tx.Exec(query, newIdFilm, nn)
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

func (f *FilmsStorage) GetFilm(id string) (*Film, error) {
	const op = "storage.GetFilm"

	tx, err := f.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query := "SELECT films.id, films.title, films.description, films.release_date, films.rating, actors.id, actors.first_name, actors.last_name, actors.surname, actors.birthday, actors.sex FROM films_actors JOIN films ON films_actors.film_id = films.id JOIN actors ON films_actors.actor_id = actors.id WHERE films.id = $1"

	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	film := Film{}
	for rows.Next() {
		actor := Actor{}
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.Date, &film.Rating, &actor.IdActor, &actor.FirstName, &actor.LastName, &actor.Surname, &actor.Birthday, &actor.Sex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		//err = film.formatDate()
		//if err != nil {
		//	return nil, fmt.Errorf("%s: %w", op, err)
		//}

		//err = actor.formatBirth()
		//if err != nil {
		//	return nil, fmt.Errorf("%s: %w", op, err)
		//}

		film.Actors = append(film.Actors, actor)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.New("user ID with data does not exist"))
	}

	if film.equals(Film{}) {
		return nil, fmt.Errorf("%s: %w", op, errors.New("user ID with data does not exist"))
	}

	return &film, nil
}

func (f *FilmsStorage) DeleteFilmById(id string) error {
	const op = "storage.DeleteFilmById"

	tx, err := f.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "delete from films_actors where film_id = $1"
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = "delete from films where id = $1"
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (f *FilmsStorage) ChangeFilm(id string, changedDataFilm vk_test_task.ChangeDataFilm) error {
	const op = "storage.ChangeFilm"

	tx, err := f.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if changedDataFilm.Title != "" {
		query := "UPDATE films SET title = $1 WHERE id = $2"
		_, err = tx.Exec(query, changedDataFilm.Title, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changedDataFilm.Date != "" {
		query := "UPDATE films SET release_date = $1 WHERE id = $2"
		_, err = tx.Exec(query, changedDataFilm.Date, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changedDataFilm.Description != "" {
		query := "UPDATE films SET description = $1 WHERE id = $2"
		_, err = tx.Exec(query, changedDataFilm.Description, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if changedDataFilm.Rating != "" {
		query := "UPDATE films SET rating = $1 WHERE id = $2"
		_, err = tx.Exec(query, changedDataFilm.Rating, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	for _, idActor := range changedDataFilm.Actors {
		if idActor[0] == '-' {
			query := "DELETE FROM films_actors WHERE film_id = $1 AND actor_id = $2"
			_, err = tx.Exec(query, id, string(idActor[1]))
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		} else {
			query := "INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)"
			_, err := tx.Exec(query, id, idActor)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	tx.Commit()

	return nil
}

//func (f *Film) formatDate() error {
//	const op = "storage.formatDate"
//
//	date, err := time.Parse("2006-01-02", f.Date)
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//	f.Date = date.Format("2006-01-02")
//	return nil
//}

func (f Film) equals(other Film) bool {
	return f.ID == other.ID && f.Title == other.Title && f.Description == other.Description && f.Date == other.Date && f.Rating == other.Rating
}

func (f customFilm) eq(other customFilm) bool {
	return f.ID == other.ID && f.Title == other.Title && f.Description == other.Description && f.Date == other.Date && f.Rating == other.Rating
}
