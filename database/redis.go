package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var (
	RedisAddr             = getEnv("REDIS_ADDR", "localhost:6379")
	SlidingWindowSize     = getEnvDuration("SLIDING_WINDOW_SIZE", time.Minute)
	DefaultMaxRequests    = getEnvInt("DEFAULT_MAX_REQUESTS", 10)
	VIPMaxRequests        = getEnvInt("VIP_MAX_REQUESTS", 20)
	AdminMaxRequests      = getEnvInt("ADMIN_MAX_REQUESTS", 50)
	RetryInterval         = getEnvDuration("RETRY_INTERVAL", 5*time.Second)
	LogEnabled            = getEnvBool("LOG_ENABLED", true)
	NotifyService         = getEnv("NOTIFY_SERVICE", "Telegram")
)

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1"
	}
	return fallback
}

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

func getUserRateLimit(userID string) (int, time.Duration) {
	switch userID {
	case "VIP_user":
		return VIPMaxRequests, SlidingWindowSize
	case "Admin_user":
		return AdminMaxRequests, SlidingWindowSize
	default:
		return DefaultMaxRequests, SlidingWindowSize
	}
}

func logRequestStatus(userID string, allowed bool, requestCount int, maxRequests int) {
	if LogEnabled {
		status := "allowed"
		if !allowed {
			status = "blocked"
		}
		log.Printf("User %s request %s (count: %d/%d) at %v
", userID, status, requestCount, maxRequests, time.Now())
		fmt.Printf("Logged request for %s with status: %s (count: %d/%d)
", userID, status, requestCount, maxRequests)
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

func testRateLimiting(client *redis.Client, users []string) {
	for _, userID := range users {
		maxRequests, windowSize := getUserRateLimit(userID)

		allowed, requestCount, err := isAllowedSlidingWindow(client, userID, maxRequests, windowSize)
		if err != nil {
			log.Printf("Error in rate limiting for user %s: %v", userID, err)
			continue
		}
		logRequestStatus(userID, allowed, requestCount, maxRequests)

		if allowed {
			fmt.Printf("Request allowed for user %s
", userID)
		} else {
			notifyUserRateLimitExceeded(userID, requestCount, maxRequests)
			cleanupExpiredKeys(client, userID)
		}
	}
}

func main() {
	client := createRedisClient()
	defer client.Close()

	users := []string{"user1", "VIP_user", "Admin_user", "user2"}
	testRateLimiting(client, users)
}
