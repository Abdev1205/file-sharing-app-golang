package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to redis")

	return client

}
