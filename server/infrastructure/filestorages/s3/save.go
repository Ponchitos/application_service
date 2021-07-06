package s3

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"mime/multipart"
)

func (str *storage) Save(ctx context.Context, apkFile multipart.File, name string) (string, error) {
	if str.sess == nil {
		return "test local link", nil
	}

	bucketExist, err := str.checkBucketExist(ctx)
	if err != nil {
		return "", err
	}

	if !bucketExist {
		err := str.createBucket(ctx)
		if err != nil {
			return "", errors.NewErrorf("Cannot create bucket: %v", "Не удалось создать bucket: %v", err)
		}
	}

	result, err := s3manager.NewUploader(str.sess).UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(str.config.S3Bucket),
		Key:    aws.String(name),
		Body:   apkFile,
	})

	if err != nil {
		return "", errors.NewErrorf("Cannot upload file to s3 store: %v", "Не удалось загрузить файл в s3 хранилище: %v", err)
	}

	if _, err := apkFile.Seek(0, io.SeekStart); err != nil {
		return "", errors.NewErrorf("Cannot set seek: %v", "Не удалось выставить начало файла: %v", err)
	}

	return result.Location, nil
}
