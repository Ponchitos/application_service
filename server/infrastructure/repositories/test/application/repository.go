package application

import (
	"github.com/Ponchitos/application_service/server/internal/repositories"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

type testRepository struct {
	lgr logger.Logger
}

func NewTestApplicationRepository(lgr logger.Logger) repositories.ApplicationRepository {
	return &testRepository{lgr}
}
