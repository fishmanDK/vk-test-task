package storage

import (
	"fmt"
	"github.com/fishmanDK/vk-test-task/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const timeout = 4 * time.Second

type StorageServ struct {
	Actors
	Films
	Auth
	MemoryStorage
}

type Storage struct {
	Storage *sqlx.DB
}

func ConnectStorage(db_cnf *config.ConfigPostgresDB) (*sqlx.DB, error) {
	const op = "storage.MustStorage"

	db, err := sqlx.Open("postgres", db_cnf.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func MustStorage(db_cnf *config.ConfigPostgresDB) (*StorageServ, error) {
	const op = "storage.MustStorageServ"

	db, err := ConnectStorage(db_cnf)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &StorageServ{
		Actors:        MustActorsStorage(db),
		Films:         MustFilmsStorage(db),
		Auth:          NewAuthDb(db),
		MemoryStorage: MustMyStorage(),
	}, nil
}
