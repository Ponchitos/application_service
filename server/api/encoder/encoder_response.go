package encoder

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseType struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

func Response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		Error(ctx, err, w)

		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := &responseType{
		Status:   "success",
		Response: response,
	}

	return json.NewEncoder(w).Encode(data)
}
