package databases

import (
	"context"
	"gin-auth-mongo/utils/datetime"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db := os.Getenv("REDIS_DB")

	// convert db string to integer
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalf("Invalid Redis DB number: %v", err)
	}

	// configure Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port, // Redis address
		Password: password,          // no password if not set
		DB:       dbInt,             // default DB
	})

	// Optionally, you can check the connection
	if err := RedisClient.Ping(GetRedisContext()).Err(); err != nil {
		panic(err) // Handle error appropriately
	}
}

// get context
func GetRedisContext() context.Context {
	return context.Background()
}

func RedisGet(key string) (string, error) {
	result, err := RedisClient.Get(GetRedisContext(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", nil
		}
		return "", err
	}
	// log.Println("result: ", result)
	return result, nil
}

func RedisSetWithoutExpiry(key string, value string) error {
	return RedisClient.Set(GetRedisContext(), key, value, 0).Err()
}

/* func RedisSet(key string, value string, expiry time.Duration) error {
	return RedisClient.Set(GetRedisContext(), key, value, expiry).Err()
} */

func RedisSet(key string, value string, expiry int, unit datetime.TIME_UNIT) error {
	return RedisClient.Set(GetRedisContext(), key, value, time.Duration(expiry)*time.Duration(unit)).Err()
}

func RedisDel(key string) error {
	return RedisClient.Del(GetRedisContext(), key).Err()
}

func RedisExists(key string) (bool, error) {
	_, err := RedisClient.Get(GetRedisContext(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// unit: seconds
func RedisTTL(key string) (time.Duration, error) {
	return RedisClient.TTL(GetRedisContext(), key).Result()
}

// unit: milliseconds
func RedisPTTL(key string) (time.Duration, error) {
	return RedisClient.PTTL(GetRedisContext(), key).Result()
}

func RedisExpire(key string, expiry int, unit datetime.TIME_UNIT) error {
	return RedisClient.Expire(GetRedisContext(), key, time.Duration(expiry)*time.Duration(unit)).Err()
}

func RedisIncr(key string) (int64, error) {
	count, err := RedisClient.Incr(GetRedisContext(), key).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}
