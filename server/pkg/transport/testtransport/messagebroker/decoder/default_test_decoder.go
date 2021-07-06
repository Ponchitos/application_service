package decoder

import (
	"encoding/json"
)

func DefaultTestDecoderPublisherFunc(topic string, message interface{}) (interface{}, error) {
	messageOfBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return messageOfBytes, nil
}
