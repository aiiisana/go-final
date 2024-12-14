package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

func SetCacheRedis(cacheKey string, cacheValue string, expiration time.Duration) error {
	err := RedisClient.Set(ctx, cacheKey, cacheValue, expiration).Err()
	if err != nil {
		log.Println("Error setting cache in Redis:", err)
		return err
	}
	log.Println("Cache set successfully in Redis")
	return nil
}

func GetCacheRedis(cacheKey string) (string, error) {
	cacheValue, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		log.Println("Cache key not found in Redis")
		return "", nil
	} else if err != nil {
		log.Println("Error fetching cache from Redis:", err)
		return "", err
	}
	return cacheValue, nil
}

func DeleteCacheRedis(cacheKey string) error {
	err := RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Println("Error deleting cache in Redis:", err)
		return err
	}
	log.Println("Cache deleted successfully from Redis")
	return nil
}
