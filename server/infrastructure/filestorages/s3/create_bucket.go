package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (str *storage) createBucket(ctx context.Context) error {
	if str.sess == nil {
		return nil
	}

	_, err := s3.New(str.sess).CreateBucketWithContext(ctx, &s3.CreateBucketInput{Bucket: aws.String(str.config.S3Bucket)})

	return err
}
