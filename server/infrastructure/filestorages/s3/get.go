package s3

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (str *storage) Get(ctx context.Context, key string) ([]byte, error) {
	if str.sess == nil {
		return nil, nil
	}

	bucketExist, err := str.checkBucketExist(ctx)
	if err != nil {
		return nil, err
	}

	if !bucketExist {
		return nil, errors.NewError("Bucket don't exist", "Bucket не существует")
	}

	response := &aws.WriteAtBuffer{}

	_, err = s3manager.NewDownloader(str.sess).DownloadWithContext(ctx, response, &s3.GetObjectInput{
		Bucket: aws.String(str.config.S3Bucket),
		Key:    aws.String(key),
	})

	return response.Bytes(), err
}
