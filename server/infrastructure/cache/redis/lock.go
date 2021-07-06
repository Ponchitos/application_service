package redis

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"time"
)

func (c *client) Lock(ctx context.Context, key, lock string) error {
	set, err := c.conn.SetNX(ctx, fmt.Sprintf("%s#lock", key), lock, time.Millisecond*50).Result()
	if err != nil {
		return err
	}

	if !set {
		return errors.NewError("Can't set value", "Не удалось добавить значение")
	}

	return nil
}
