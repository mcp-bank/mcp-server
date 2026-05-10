package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         "cache:6379",
		Username:     "",
		Password:     "",
		DB:           0,
		MaxRetries:   3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		MinIdleConns: 5,
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	ping := client.Ping(timeout)
	if err := ping.Err(); err != nil {
		return nil, err
	}
	return client, nil
}
