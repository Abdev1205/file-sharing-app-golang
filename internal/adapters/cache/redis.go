package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // e.g. "localhost:6379"
		Password: "",                      // No password set
		DB:       0,                       // Default DB
	})

	_, err := client.Ping(client.Context()).Result()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to redis")

	return client

}
