package keymanager

import "context"

type Manager interface {
	GetKey(ctx context.Context, enterpriseID string) (token string, err error)
}
