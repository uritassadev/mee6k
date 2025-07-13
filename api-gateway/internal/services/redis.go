package services

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisService() (*RedisService, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://:redis_password_123@localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Cache operations
func (s *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, key, data, expiration).Err()
}

func (s *RedisService) Get(key string, dest interface{}) error {
	data, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (s *RedisService) Delete(key string) error {
	return s.client.Del(s.ctx, key).Err()
}

func (s *RedisService) Exists(key string) bool {
	count, _ := s.client.Exists(s.ctx, key).Result()
	return count > 0
}

// Real-time data operations
func (s *RedisService) PublishAlert(channel string, alert interface{}) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	return s.client.Publish(s.ctx, channel, data).Err()
}

func (s *RedisService) Subscribe(channel string) *redis.PubSub {
	return s.client.Subscribe(s.ctx, channel)
}

// Rate limiting
func (s *RedisService) IsRateLimited(key string, limit int, window time.Duration) (bool, error) {
	pipe := s.client.TxPipeline()
	
	// Increment counter
	incr := pipe.Incr(s.ctx, key)
	pipe.Expire(s.ctx, key, window)
	
	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return false, err
	}
	
	count := incr.Val()
	return count > int64(limit), nil
}

// Session management
func (s *RedisService) SetSession(sessionID string, data interface{}, expiration time.Duration) error {
	return s.Set("session:"+sessionID, data, expiration)
}

func (s *RedisService) GetSession(sessionID string, dest interface{}) error {
	return s.Get("session:"+sessionID, dest)
}

func (s *RedisService) DeleteSession(sessionID string) error {
	return s.Delete("session:" + sessionID)
}

// Metrics and counters
func (s *RedisService) IncrementCounter(key string) error {
	return s.client.Incr(s.ctx, key).Err()
}

func (s *RedisService) GetCounter(key string) (int64, error) {
	return s.client.Get(s.ctx, key).Int64()
}

func (s *RedisService) SetCounter(key string, value int64) error {
	return s.client.Set(s.ctx, key, value, 0).Err()
}