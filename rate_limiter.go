package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func isAllowedSlidingWindow(client *redis.Client, userID string, maxRequests int, windowSize time.Duration) (bool, int, error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	key := fmt.Sprintf("rate_limit:%s", userID)

	if _, err := client.ZAdd(ctx, key, &redis.Z{Score: float64(timestamp), Member: timestamp}).Result(); err != nil {
		log.Println("Error adding timestamp to Redis:", err)
		return false, 0, err
	}
	client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", timestamp-int64(windowSize/time.Millisecond)))

	requestCount, err := client.ZCard(ctx, key).Result()
	if err != nil {
		log.Println("Error retrieving request count:", err)
		return false, 0, err
	}

	if requestCount > int64(maxRequests) {
		return false, int(requestCount), nil
	}
	client.Expire(ctx, key, windowSize)
	return true, int(requestCount), nil
}

func getUserRateLimit(role string) (int, time.Duration) {
	switch role {
	case "VIP":
		return VIPMaxRequests, SlidingWindowSize
	case "Admin":
		return AdminMaxRequests, SlidingWindowSize
	default:
		return DefaultMaxRequests, SlidingWindowSize
	}
}
