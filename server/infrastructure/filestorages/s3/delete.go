package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (str *storage) Delete(ctx context.Context, key string) error {
	if str.sess == nil {
		return nil
	}

	_, err := s3.New(str.sess).DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{Bucket: aws.String(str.config.S3Bucket), Key: aws.String(key)})

	return err
}
