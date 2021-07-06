package filestorages

import (
	"context"
	"mime/multipart"
)

type FileStorage interface {
	Save(ctx context.Context, apkFile multipart.File, name string) (path string, err error)
	Delete(ctx context.Context, path string) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetByChunk(ctx context.Context, key string, offset, whence int64) ([]byte, int64, error)
}
