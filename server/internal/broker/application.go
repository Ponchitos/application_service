package broker

import "context"

type ApplicationBroker interface {
	CompleteUploadFileRecord(ctx context.Context, topic, packageName string) error
	DeleteApplicationVersionRecord(ctx context.Context, topic, versionUUID string) error
}
