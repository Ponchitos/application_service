package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (str *storage) checkBucketExist(ctx context.Context) (bool, error) {
	if str.sess == nil {
		return false, nil
	}

	listBuckets, err := s3.New(str.sess).ListBucketsWithContext(ctx, nil)
	if err != nil {
		return false, err
	}

	for _, bucket := range listBuckets.Buckets {
		if bucket != nil && aws.StringValue(bucket.Name) == str.config.S3Bucket {
			return true, nil
		}
	}

	return false, nil
}
