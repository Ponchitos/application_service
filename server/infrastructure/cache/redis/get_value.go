package redis

import (
	"context"
	"fmt"
)

func (c *client) GetValue(ctx context.Context, key string) (string, bool, error) {
	value, err := c.conn.Get(ctx, fmt.Sprintf("%s#1", key)).Result()
	if err != nil {
		return "", false, err
	}

	if len(value) == 0 {
		value, err = c.conn.Get(ctx, fmt.Sprintf("%s#2", key)).Result()
		if err != nil {
			return "", false, err
		}

		return value, true, nil
	}

	return value, false, nil
}
