package s3

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (str *storage) GetByChunk(ctx context.Context, key string, offset, whence int64) ([]byte, int64, error) {
	if str.sess == nil {
		return nil, 0, nil
	}

	bucketExist, err := str.checkBucketExist(ctx)
	if err != nil {
		return nil, 0, err
	}

	if !bucketExist {
		return nil, 0, errors.NewError("Bucket don't exist", "Bucket не существует")
	}

	info, err := s3.New(str.sess).HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(str.config.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, 0, errors.NewErrorf("Cannot get apk info: %v", "Не удалось получить информацию о apk файле: %v", err)
	}

	response := &aws.WriteAtBuffer{}

	_, err = s3manager.NewDownloader(str.sess).DownloadWithContext(ctx, response, &s3.GetObjectInput{
		Bucket: aws.String(str.config.S3Bucket),
		Key:    aws.String(key),
		Range:  aws.String(fmt.Sprintf("bytes=%v-%v", whence, offset)),
	})
	if err != nil {
		return nil, 0, errors.NewErrorf("Cannot download apk file: %v", "Не удалось скачать apk файл: %v", err)
	}

	return response.Bytes(), aws.Int64Value(info.ContentLength), nil
}
