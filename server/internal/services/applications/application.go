package applications

import (
	"context"
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/internal/broker"
	"github.com/Ponchitos/application_service/server/internal/handlers"
	"github.com/Ponchitos/application_service/server/internal/models"
	"github.com/Ponchitos/application_service/server/internal/repositories"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"mime/multipart"
)

type Service interface {
	UploadApkFile(ctx context.Context, apk multipart.File) (uuid string, err error)
	CancelApkUpload(ctx context.Context, apkUUID string) error
	CompleteUploadFile(ctx context.Context, apkUUID, enterpriseID string) (applicationUUID, applicationVersionUUID string, err error)
	GetApplicationInfo(ctx context.Context, applicationUUID, applicationVersionUUID, enterpriseID string) (*models.Application, error)
	GetApplications(ctx context.Context, enterpriseID string, offset, limit int) (count int, applications []*models.BasicApplication, err error)
	GetApplicationVersions(ctx context.Context, applicationUUID, enterpriseID string) ([]*models.BasicApplication, error)
	DeleteVersionApplication(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error
	ChangeApplicationStatus(ctx context.Context, versionUUID, enterpriseID, status string) error
	UninstallApplicationComplete(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) error
	DownloadApkFile(ctx context.Context, applicationUUID, versionUUID, enterpriseID string) ([]byte, string, error)
	DownloadApkFileByChunk(ctx context.Context, applicationUUID, versionUUID, enterpriseID string, offset, whence int64) ([]byte, int64, string, error)
}

type application struct {
	lgr        logger.Logger
	store      repositories.ApplicationRepository
	apkHandler handlers.ApkFileHandler
	config     *config.Config
	broker     broker.ApplicationBroker
}

func NewApplicationService(lgr logger.Logger, store repositories.ApplicationRepository, apkHandler handlers.ApkFileHandler, config *config.Config, broker broker.ApplicationBroker) Service {
	return &application{lgr, store, apkHandler, config, broker}
}
