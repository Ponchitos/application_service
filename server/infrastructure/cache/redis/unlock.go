package redis

import (
	"context"
	"fmt"
	"strings"
)

func (c *client) Unlock(ctx context.Context, key, lock string) {
	value, err := c.conn.Get(ctx, fmt.Sprintf("%s#lock", key)).Result()
	if err == nil {
		compare := strings.Compare(value, lock)
		if compare == 0 {
			c.conn.Del(ctx, fmt.Sprintf("%s#lock", key)).Result()
		}
	}
}
