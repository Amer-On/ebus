package redis

import (
	"fmt"

	redisLib "github.com/redis/go-redis/v9"
)

func NewClient(config Config) *redisLib.Client {
	return redisLib.NewClient(&redisLib.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
}
