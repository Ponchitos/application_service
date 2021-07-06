package repositories

import "github.com/Ponchitos/application_service/server/internal/repositories"

type DataBase interface {
	Close()

	GetApplicationRepository() repositories.ApplicationRepository
}
