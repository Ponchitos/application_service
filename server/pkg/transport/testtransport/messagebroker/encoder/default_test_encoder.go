package encoder

import "encoding/json"

func DefaultTestEncoderSubscriberFunc(msg []byte) (interface{}, error) {
	var result interface{}

	err := json.Unmarshal(msg, &result)

	return result, err
}
