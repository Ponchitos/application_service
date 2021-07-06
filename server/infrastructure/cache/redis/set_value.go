package redis

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"time"
)

func (c *client) SetValue(ctx context.Context, key string, value interface{}) error {
	set, err := c.conn.SetNX(ctx, fmt.Sprintf("%s#1", key), value, time.Minute*50).Result()
	if err != nil {
		return err
	}

	if !set {
		return errors.NewError("Can't set value", "Не удалось добавить значение")
	}

	set, err = c.conn.SetNX(ctx, fmt.Sprintf("%s#2", key), value, time.Minute*60).Result()
	if err != nil {
		return err
	}

	if !set {
		return errors.NewError("Can't set value", "Не удалось добавить значение")
	}

	return nil
}
