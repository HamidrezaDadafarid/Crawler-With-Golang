
package main

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestUserRateLimit(t *testing.T) {
    client := NewRedisClient()
    userID := "testUser1"
    role := "basic"
    maxRequests := 10
    duration := time.Minute

    limiter := NewRateLimiter(client, maxRequests, duration)

    for i := 0; i < maxRequests; i++ {
        allowed, err := limiter.AllowRequest(userID, role)
        assert.NoError(t, err)
        assert.True(t, allowed)
    }

    allowed, err := limiter.AllowRequest(userID, role)
    assert.NoError(t, err)
    assert.False(t, allowed)

    client.Del(userID)
}
