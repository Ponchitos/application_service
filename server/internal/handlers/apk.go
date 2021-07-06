package handlers

import (
	"context"
	"github.com/Ponchitos/application_service/server/internal/models"
	"mime/multipart"
)

const (
	LocalMode    = "LOCAL"
	ExternalMode = "EXTERNAL"
)

type ApkFileHandler interface {
	ApkFileSave(ctx context.Context, apkFile multipart.File, name, mode string) (path string, err error)

	GetMetadataOfAPK(filePath, name string) (*models.ApplicationMetadata, error)

	DeleteApkFile(ctx context.Context, filePath, mode string) error

	DownloadApkFile(ctx context.Context, name, mode string) ([]byte, error)

	DownloadApkFileByChunk(ctx context.Context, key string, offset, whence int64, mode string) ([]byte, int64, error)
}
