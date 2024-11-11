package main

import (
<<<<<<< HEAD
	database "main/database"
	"main/models"
)

func main() {
	// For testing db review for final codes
	dbManager := database.GetInstnace()
	dbManager.Db.AutoMigrate(&models.Ads{}, &models.Filters{}, &models.Users{}, &models.Users_Ads{}, &models.WatchList{})
=======
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func testRateLimiting(client *redis.Client, users map[string]string) {
	for userID, role := range users {
		maxRequests, windowSize := getUserRateLimit(role)

		allowed, requestCount, err := isAllowedSlidingWindow(client, userID, maxRequests, windowSize)
		if err != nil {
			log.Printf("Error in rate limiting for user %s: %v", userID, err)
			continue
		}
		logRequestStatus(userID, allowed, requestCount, maxRequests)

		if allowed {
			fmt.Printf("Request allowed for user %s with role %s\n", userID, role)
		} else {
			notifyUserRateLimitExceeded(userID, requestCount, maxRequests)
			cleanupExpiredKeys(client, userID)
		}
	}
}

func main() {
	client := createRedisClient()
	defer client.Close()

	users := map[string]string{
		"user1":      "Standard",
		"VIP_user":   "VIP",
		"Admin_user": "Admin",
		"user2":      "Standard",
	}
	testRateLimiting(client, users)
>>>>>>> ccc663f (Replace environment variables with hardcoded values in redis.go)
}
