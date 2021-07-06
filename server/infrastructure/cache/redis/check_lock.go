package redis

import "context"

func (c *client) CheckLock(ctx context.Context, key string) bool {
	exist, err := c.conn.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	if exist == 1 {
		return true
	}

	return false
}
