package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{redisClient: redisClient}
}

func (redisRepo *RedisRepository) InsertUserToken(key string, value string, expires time.Duration) error {
	return redisRepo.redisClient.Set(context.Background(), key, value, expires).Err()
}

func (redisRepo *RedisRepository) ReplaceToken(currentToken, newToken string, expires time.Duration) error {
	email, err := redisRepo.deleteToken(currentToken)
	if err != nil {
		return err
	}
	return redisRepo.redisClient.Set(context.Background(), newToken, email, expires).Err()
}

func (redisRepo *RedisRepository) deleteToken(token string) (string, error) {
	return redisRepo.redisClient.GetDel(context.Background(), token).Result()
}
