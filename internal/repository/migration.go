package repository

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(databaseConnection string) (bool, error) {
	m, err := migrate.New("file://db/migrations", databaseConnection)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}
	}

	return true, nil
}
