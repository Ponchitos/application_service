package application

import (
	"context"
	"encoding/json"
)

type deleteApplicationVersionRecordMessage struct {
	VersionUUID string `json:"versionUId"`
}

func (b *broker) DeleteApplicationVersionRecord(ctx context.Context, topic, versionUUID string) error {
	messageOfBytes, err := json.Marshal(&deleteApplicationVersionRecordMessage{versionUUID})
	if err != nil {
		b.lgr.Errorw("Serialize error: ",
			"operation", "DeleteApplicationVersionRecord",
			"topic", topic,
			"error", err,
			"versionUUID", versionUUID,
		)

		return err
	}

	err = b.publisher.Publish(topic, messageOfBytes)
	if err != nil {
		b.lgr.Errorw("Message didn't send: ",
			"operation", "DeleteApplicationVersionRecord",
			"topic", topic,
			"error", err,
			"versionUUID", versionUUID,
		)

		return err
	}

	return nil
}
