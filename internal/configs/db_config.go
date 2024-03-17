package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Pgconfig struct {
	user, database, host, password, port, sslMode string
}

func PgConfigFromEnv() (Pgconfig, error) {
	const op = "configs.PgConfigFromEnv"

	var missing []string

	get := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			missing = append(missing, key)
		}
		return val
	}

	cfg := Pgconfig{}
	if err := godotenv.Load(); err != nil {
		return cfg, fmt.Errorf("%s: %w", op, err)
	}

	cfg = Pgconfig{
		user:     get("PG_USER"),
		database: get("PG_DATABASE"),
		host:     get("PG_HOST"),
		password: get("PG_PASSWORD"),
		port:     get("PG_PORT"),
		sslMode:  os.Getenv("PG_SSLMODE"),
	}

	fmt.Println(cfg)
	switch cfg.sslMode {
	case "", "disable", "allow", "require", "verify-ca", "verify-full":
		// valid sslmode
	default:
		return cfg, fmt.Errorf(`invalid sslmode "%s": expected one of "", "disable", "allow", "require", "verify-ca", or "verify-full"`, cfg.sslMode)
	}

	if len(missing) > 0 {
		return cfg, fmt.Errorf("missing required environment variables: %v", missing)
	}
	return cfg, nil
}

func (pg Pgconfig) String() string {
	s := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", pg.host, pg.port, pg.user, pg.database, pg.password, pg.sslMode)

	return s
}
