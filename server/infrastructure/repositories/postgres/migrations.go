package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Ponchitos/application_service/server/config"
	customError "github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4"
)

func runMigrations(config *config.Config, lgr logger.Logger) error {
	lgr.Info("start migrations...")

	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", config.DBuser, config.DBpass, config.DBhost, config.DBport, config.DBname))
	if err != nil {
		lgr.Info("failed migration")

		return customError.NewErrorf("Connection error: %v", "Connection error: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		lgr.Info("failed migration")

		return customError.NewErrorf("Error while configuring the driver: %v", "Error while configuring the driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", config.MigrationPath), "postgres", driver)
	if err != nil {
		lgr.Info("failed migration")

		return customError.NewErrorf("Error while configuring the migrate instance: %v", "Error while configuring the migrate instance: %v", err)
	}
	defer migration.Close()

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		lgr.Error("Migration UP error: ", err)

		if err := migration.Down(); err != nil {
			lgr.Info("failed migration")

			return customError.NewErrorf("Migration error: %v", "Migration error: %v", err)
		}

		lgr.Info("failed migration")

		return err
	}

	lgr.Info("successful migrations")

	return nil
}
