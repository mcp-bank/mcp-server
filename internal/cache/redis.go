package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func New() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         "cache:6379",
		Username:     "",
		Password:     "",
		DB:           0,
		MaxRetries:   3,
		ReadTimeout:  3,
		WriteTimeout: 3,
		MinIdleConns: 5,
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()

	ping := client.Ping(timeout)
	if err := ping.Err(); err != nil {
		slog.Error("New:",
			"err", err)
		return nil
	}
	return client
}
