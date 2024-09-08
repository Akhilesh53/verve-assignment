package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8" // Redis client
	kafka "github.com/segmentio/kafka-go"
)

var (
	redisClient    = InitRedis()        // Redis client for storing unique IDs
	kafkaProducer  = InitKafkaProducer() // Kafka producer initialization
	kafkaTopic     = "unique-requests-count" // Kafka topic to send the count
	redisKey       = "unique-requests"      // Redis key for tracking unique IDs
	redisCtx       = context.Background()   // Context for Redis operations
)

// Request structure
type Request struct {
	ID       int    `form:"id" binding:"required"`
	Endpoint string `form:"endpoint"`
}

func main() {
	r := gin.Default()
	r.GET("/api/verve/accept", handleRequest)

	go sendUniqueRequestCountToKafkaEveryMinute() // Send unique request count to Kafka every minute
	r.Run(":8080")                                // Start server on port 8080
}

// Handle incoming requests
func handleRequest(c *gin.Context) {
	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.String(http.StatusBadRequest, "failed")
		return
	}

	// Check if the ID is unique using Redis
	if AddIfUniqueInRedis(req.ID) {
		if req.Endpoint != "" {
			sendHTTPRequest(req.Endpoint) // Send HTTP request to provided endpoint with count
		}
	}
	c.String(http.StatusOK, "ok")
}

// AddIfUniqueInRedis adds the request ID to Redis and checks if it's unique
func AddIfUniqueInRedis(id int) bool {
	// SADD adds the ID to the set if it's not already present
	result, err := redisClient.SAdd(redisCtx, redisKey, id).Result()
	if err != nil {
		log.Printf("Failed to add ID to Redis: %v", err)
		return false
	}
	// If result == 1, it means the ID was added for the first time
	return result == 1
}

// GetUniqueRequestCountFromRedis fetches the unique count from Redis
func GetUniqueRequestCountFromRedis() int {
	count, err := redisClient.SCard(redisCtx, redisKey).Result() // SCARD returns the number of elements in the set
	if err != nil {
		log.Printf("Failed to get unique request count from Redis: %v", err)
		return 0
	}
	return int(count)
}

// ResetUniqueRequestsInRedis resets the set in Redis by deleting the key
func ResetUniqueRequestsInRedis() {
	err := redisClient.Del(redisCtx, redisKey).Err()
	if err != nil {
		log.Printf("Failed to reset unique requests in Redis: %v", err)
	}
}

// Send the unique request count to Kafka every minute
func sendUniqueRequestCountToKafkaEveryMinute() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		count := GetUniqueRequestCountFromRedis()

		// Create a Kafka message with the unique request count
		message := kafka.Message{
			Key:   []byte("unique-request-count"),
			Value: []byte(fmt.Sprintf("%d", count)),
		}

		// Send the message to the Kafka topic
		err := kafkaProducer.WriteMessages(context.Background(), message)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %v", err)
		} else {
			log.Printf("Sent unique request count to Kafka: %d", count)
		}

		// Reset unique requests in Redis for the next minute
		ResetUniqueRequestsInRedis()
	}
}

// Send HTTP GET request to an endpoint
func sendHTTPRequest(endpoint string) {
	count := GetUniqueRequestCountFromRedis()
	resp, err := http.Get(fmt.Sprintf("%s?count=%d", endpoint, count))
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		return
	}
	log.Printf("Sent request to %s, Status Code: %d", endpoint, resp.StatusCode)
}

// Send HTTP POST request (Extension 1)
func sendHTTPPostRequest(endpoint string) {
	count := GetUniqueRequestCountFromRedis()
	payload := map[string]int{"count": count}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send HTTP POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Sent POST request to %s, Status Code: %d", endpoint, resp.StatusCode)
}

// InitKafkaProducer initializes the Kafka producer
func InitKafkaProducer() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"), // Replace with your Kafka broker address
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

// InitRedis initializes and returns a Redis client
func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
}
