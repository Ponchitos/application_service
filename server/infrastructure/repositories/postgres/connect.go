package postgres

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/postgres/application"
	serviceRepository "github.com/Ponchitos/application_service/server/internal/repositories"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type dataBase struct {
	client *pgxpool.Pool
	lgr    logger.Logger
}

func OpenConnection(conf *config.Config, lgr logger.Logger) (repositories.DataBase, error) {
	ctx := context.Background()

	dbConnect := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", conf.DBuser, conf.DBpass, conf.DBhost, conf.DBport, conf.DBname)

	connectConfig, err := pgx.ParseConfig(dbConnect)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(connectConfig.ConnString())
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConnLifetime = time.Second * 10

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)

	if err != nil {
		return nil, err
	}

	db := &dataBase{
		client: pool,
		lgr:    lgr,
	}

	//err = runMigrations(conf, lgr)
	//if err != nil {
	//	db.Close()
	//
	//	return nil, err
	//}

	return db, nil
}

func (db *dataBase) Close() {
	db.client.Close()

	db.lgr.Info("Postgres stopped")
}

func (db *dataBase) GetApplicationRepository() serviceRepository.ApplicationRepository {
	return application.NewApplicationRepository(db.client, db.lgr)
}
