package s3

import (
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/filestorages"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"sync"
)

type storage struct {
	sync.Mutex

	config *config.Config
	lgr    logger.Logger
	sess   *session.Session
}

func NewS3Storage(conf *config.Config, lgr logger.Logger) (filestorages.FileStorage, error) {
	var err error

	storage := &storage{config: conf, lgr: lgr}

	if conf.Env == config.Local {
		storage.sess = nil
	} else {
		storage.sess, err = session.NewSession(
			&aws.Config{
				Region:           aws.String(conf.S3Region),
				Credentials:      credentials.NewStaticCredentials(conf.S3AccountKey, conf.S3Secret, ""),
				Endpoint:         aws.String(conf.S3Endpoint),
				S3ForcePathStyle: aws.Bool(true),
			},
		)

		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}
