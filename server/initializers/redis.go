package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisMustConnect() *redis.Client {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	context := context.Background()
	if err := client.Ping(context).Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the Redis")
	return client
}
