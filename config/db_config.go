package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ConfigPostgresDB struct {
	PostgresDB `yaml:"postgres_db"`
}

type PostgresDB struct {
	PgUser     string `yaml:"pg_user"`
	PgDatabase string `yaml:"pg_database"`
	PgHost     string `yaml:"pg_host"`
	PgPort     string `yaml:"pg_port"`
	PgSslmode  string `yaml:"pg_sslmode"`
	PgPassword string `yaml:"pg_password"`
}

func MustLoadConfigPostgresDB() *ConfigPostgresDB {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("%w", err)
	}

	configPath := os.Getenv("CONFIG_PATH_DB")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg ConfigPostgresDB

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func (pg ConfigPostgresDB) String() string {
	s := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", pg.PgHost, pg.PgPort, pg.PgUser, pg.PgDatabase, pg.PgPassword, pg.PgSslmode)
	return s
}
