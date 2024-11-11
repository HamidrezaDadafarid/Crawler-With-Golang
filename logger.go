package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func logRequestStatus(userID string, allowed bool, requestCount int, maxRequests int) {
	if LogEnabled {
		status := "allowed"
		if !allowed {
			status = "blocked"
		}
		log.Printf("User %s request %s (count: %d/%d) at %v\n", userID, status, requestCount, maxRequests, time.Now())
		log.Printf("Telegram notification sent for %s", userID)
	}
}

func notifyUserRateLimitExceeded(userID string, requestCount int, maxRequests int) {
	message := fmt.Sprintf("User %s has exceeded the rate limit. Requests made: %d out of allowed %d.", userID, requestCount, maxRequests)
	log.Println(message)
	switch NotifyService {
	case "Telegram":
		log.Printf("Telegram notification sent for %s", userID)
	case "Email":
		log.Printf("Email notification sent for %s", userID)
	default:
		log.Printf("Notification sent for %s via %s", userID, NotifyService)
	}
}

func cleanupExpiredKeys(client *redis.Client, userID string) {
	key := fmt.Sprintf("rate_limit:%s", userID)
	if _, err := client.Del(ctx, key).Result(); err != nil {
		log.Printf("Error cleaning up expired keys for user %s: %v", userID, err)
	} else {
		log.Printf("Expired keys cleaned for user %s", userID)
	}
}
