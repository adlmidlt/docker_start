package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

var Client *redis.Client
var ctx = context.Background()

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0, // БД по умолчанию.
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis")
}

func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

func Set(key string, value string, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

func Delete(key string) error {
	return Client.Del(ctx, key).Err()
}
