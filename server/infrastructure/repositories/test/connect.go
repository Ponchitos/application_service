package test

import (
	"github.com/Ponchitos/application_service/server/infrastructure/repositories"
	"github.com/Ponchitos/application_service/server/infrastructure/repositories/test/application"
	serviceRepository "github.com/Ponchitos/application_service/server/internal/repositories"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type testDataBase struct {
	lgr logger.Logger
}

func OpenTestConnection(lgr logger.Logger) (repositories.DataBase, error) {
	return &testDataBase{lgr}, nil
}

func (db *testDataBase) Close() {
	db.lgr.Info("Test store stopped")
}

func (db *testDataBase) GetApplicationRepository() serviceRepository.ApplicationRepository {
	return application.NewTestApplicationRepository(db.lgr)
}
