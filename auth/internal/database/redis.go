package database

import (
	"context"
	"fmt"
	"inkstack-auth/internal/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection established successfully")
	return nil
}

// GetRedis returns the Redis client instance
func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if redisClient != nil {
		log.Println("Redis connection closed")
		return redisClient.Close()
	}
	return nil
}

// BlacklistToken adds a token to the blacklist
func BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return redisClient.Set(ctx, key, "revoked", expiry).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted
func IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	result, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// IncrementLoginAttempts increments failed login attempts counter
func IncrementLoginAttempts(ctx context.Context, identifier string) (int64, error) {
	key := fmt.Sprintf("login_attempts:%s", identifier)
	count, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiry on first attempt
	if count == 1 {
		redisClient.Expire(ctx, key, 15*time.Minute)
	}

	return count, nil
}

// ResetLoginAttempts resets the login attempts counter
func ResetLoginAttempts(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("login_attempts:%s", identifier)
	return redisClient.Del(ctx, key).Err()
}

// GetLoginAttempts gets the current login attempts count
func GetLoginAttempts(ctx context.Context, identifier string) (int, error) {
	key := fmt.Sprintf("login_attempts:%s", identifier)
	count, err := redisClient.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}
