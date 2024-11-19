package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckRateLimit_NormalUser(t *testing.T) {
	rdb := rate_limiter.InitRedis()
	defer rdb.Close()

	userID := "test_normal_user"
	role := "normal_user"

	for i := 0; i < 20; i++ {
		allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestCheckRateLimit_PremiumUser(t *testing.T) {
	rdb := rate_limiter.InitRedis()
	defer rdb.Close()

	userID := "test_premium_user"
	role := "premium_user"

	for i := 0; i < 50; i++ {
		allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestCheckRateLimit_AdminUser(t *testing.T) {
	rdb := rate_limiter.InitRedis()
	defer rdb.Close()

	userID := "test_admin_user"
	role := "admin"

	for i := 0; i < 100; i++ {
		allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestRateLimitReset(t *testing.T) {
	rdb := rate_limiter.InitRedis()
	defer rdb.Close()

	userID := "test_reset_user"
	role := "normal_user"

	for i := 0; i < 20; i++ {
		allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rate_limiter.CheckRateLimit(rdb, userID, role)
	assert.NoError(t, err)
	assert.False(t, allowed)

	time.Sleep(time.Minute)

	allowed, err = rate_limiter.CheckRateLimit(rdb, userID, role)
	assert.NoError(t, err)
	assert.True(t, allowed)
}
