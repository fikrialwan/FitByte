package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fikrialwan/FitByte/config"
	"github.com/fikrialwan/FitByte/internal/dto"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	SetUserProfile(userID string, profile dto.UserResponse, ttl time.Duration) error
	GetUserProfile(userID string) (dto.UserResponse, error)
	DeleteUserProfile(userID string) error
	SetJWTBlacklist(token string, ttl time.Duration) error
	IsJWTBlacklisted(token string) bool
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type cacheService struct {
	client *redis.Client
}

func NewCacheService(config *config.Config) CacheService {
	redisAddr := config.RedisAddr
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     config.RedisPassword,
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
		MaxRetries:   3,
	})

	return &cacheService{
		client: rdb,
	}
}

func (c *cacheService) SetUserProfile(userID string, profile dto.UserResponse, ttl time.Duration) error {
	key := fmt.Sprintf("user:profile:%s", userID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), key, data, ttl).Err()
}

func (c *cacheService) GetUserProfile(userID string) (dto.UserResponse, error) {
	key := fmt.Sprintf("user:profile:%s", userID)
	data, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return dto.UserResponse{}, err
	}

	var profile dto.UserResponse
	err = json.Unmarshal([]byte(data), &profile)
	return profile, err
}

func (c *cacheService) DeleteUserProfile(userID string) error {
	key := fmt.Sprintf("user:profile:%s", userID)
	return c.client.Del(context.Background(), key).Err()
}

func (c *cacheService) SetJWTBlacklist(token string, ttl time.Duration) error {
	key := fmt.Sprintf("jwt:blacklist:%s", token)
	return c.client.Set(context.Background(), key, "1", ttl).Err()
}

func (c *cacheService) IsJWTBlacklisted(token string) bool {
	key := fmt.Sprintf("jwt:blacklist:%s", token)
	_, err := c.client.Get(context.Background(), key).Result()
	return err == nil
}

func (c *cacheService) Set(key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(context.Background(), key, value, ttl).Err()
}

func (c *cacheService) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *cacheService) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}
