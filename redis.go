package main

import (
    "github.com/go-redis/redis/v8"
    "context"
    "log"
)

var ctx = context.Background()

// Initialize Redis client
func InitRedis() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Assuming Redis is running locally
        Password: "",                // No password set
        DB:       0,                 // Use default DB
    })
    _, err := client.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }
    return client
}
