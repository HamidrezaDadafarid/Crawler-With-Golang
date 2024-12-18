package rate_limiter

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

func InitRedis(ctx context.Context) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    })

    if _, err := rdb.Ping(ctx).Result(); err != nil {
        panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
    }
    return rdb
}

func CheckRateLimit(ctx context.Context, rdb *redis.Client, userID string, role string) (bool, error) {
    key := fmt.Sprintf("rate_limit:%s", userID)
    limit := getLimitByRole(role)

    count, err := rdb.Incr(ctx, key).Result()
    if err != nil {
        return false, fmt.Errorf("failed to increment rate limit counter: %w", err)
    }

    if err := rdb.Expire(ctx, key, time.Minute).Err(); err != nil {
        return false, fmt.Errorf("failed to set expiration for rate limit key: %w", err)
    }

    return count <= int64(limit), nil
}

func getLimitByRole(role string) int {
    switch role {
    case "admin":
        return 100
    case "super_user":
        return 50
    case "user":
        return 20
    default:
        return 10
    }
}
