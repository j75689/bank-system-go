package wireset

import (
	"bank-system-go/internal/config"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient(config config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.Redis.DialTimeout)
	defer cancel()

	return client, client.Ping(ctx).Err()
}
