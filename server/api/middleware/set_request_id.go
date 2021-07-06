package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type ctx int8

const (
	CtxKeyRequestID ctx = iota
	CtxKeyUserName
	CtxKeyUserRole
)

func SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyRequestID, id)))
	})
}
