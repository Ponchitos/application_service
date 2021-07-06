package application

import (
	"context"
	"encoding/json"
)

type completeUploadFileRecordMessage struct {
	PackageName string `json:"packageName"`
}

func (b *broker) CompleteUploadFileRecord(ctx context.Context, topic, packageName string) error {
	messageOfBytes, err := json.Marshal(&completeUploadFileRecordMessage{packageName})
	if err != nil {
		b.lgr.Errorw("Serialize error: ",
			"operation", "CompleteUploadFileRecord",
			"topic", topic,
			"error", err,
			"packageName", packageName,
		)

		return err
	}

	err = b.publisher.Publish(topic, messageOfBytes)
	if err != nil {
		b.lgr.Errorw("Message didn't send: ",
			"operation", "CompleteUploadFileRecord",
			"topic", topic,
			"error", err,
			"packageName", packageName,
		)

		return err
	}

	return nil
}
