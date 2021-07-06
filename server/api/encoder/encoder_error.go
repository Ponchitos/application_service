package encoder

import (
	"context"
	"encoding/json"
	errorConstants "github.com/Ponchitos/application_service/server/errors"
	"net/http"
)

func Error(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch err {
	case errorConstants.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := &responseType{
		Status:   "error",
		Response: err.Error(),
	}

	json.NewEncoder(w).Encode(response)
}
