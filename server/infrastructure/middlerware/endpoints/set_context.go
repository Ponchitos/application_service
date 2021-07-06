package endpoints

import (
	"context"
	"github.com/Ponchitos/application_service/server/pkg"
	"time"
)

func ContextInit(next pkg.Endpoint, timeout time.Duration) pkg.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return next(ctx, request)
	}
}
