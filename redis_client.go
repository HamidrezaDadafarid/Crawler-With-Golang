package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func createRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})
	for i := 0; i < 3; i++ {
		_, err := client.Ping(ctx).Result()
		if err == nil {
			log.Println("Connected to Redis successfully")
			return client
		}
		log.Printf("Attempt %d: Could not connect to Redis. Retrying in %v...", i+1, RetryInterval)
		time.Sleep(RetryInterval)
	}
	log.Fatal("Could not connect to Redis after multiple attempts")
	return client
}
