package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	kafka "github.com/segmentio/kafka-go"
)

var (
	uniqueRequestTracker = NewUniqueRequestTracker() // Tracker for unique requests
	kafkaProducer        = InitKafkaProducer()       // Kafka producer initialization
	kafkaTopic           = "unique-requests-count"   // Kafka topic to send the count
)

// Request structure
type Request struct {
	ID       int    `form:"id" binding:"required"`
	Endpoint string `form:"endpoint"`
}

// UniqueRequestTracker struct to track unique requests with a lock
type UniqueRequestTracker struct {
	uniqueRequests sync.Map   // In-memory store for deduplication
	mu             sync.Mutex // Mutex to ensure safe access
}

// NewUniqueRequestTracker initializes and returns a UniqueRequestTracker
func NewUniqueRequestTracker() *UniqueRequestTracker {
	return &UniqueRequestTracker{}
}

// AddIfUnique adds a request ID if it's unique, returns true if it's unique
func (tracker *UniqueRequestTracker) AddIfUnique(id int) bool {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	// Try to store the ID, check if it's already loaded
	if _, loaded := tracker.uniqueRequests.LoadOrStore(id, true); !loaded {
		return true
	}
	return false
}

// Reset resets the in-memory store of unique requests
func (tracker *UniqueRequestTracker) Reset() {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	tracker.uniqueRequests = sync.Map{}
}

// Count returns the number of unique requests
func (tracker *UniqueRequestTracker) Count() int {
	count := 0
	tracker.uniqueRequests.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
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

	// Check if the ID is unique in the current minute using UniqueRequestTracker
	if uniqueRequestTracker.AddIfUnique(req.ID) {
		if req.Endpoint != "" {
			sendHTTPRequest(req.Endpoint) // Send HTTP request to provided endpoint with count
		}
	}
	c.String(http.StatusOK, "ok")
}

// Send the unique request count to Kafka every minute
func sendUniqueRequestCountToKafkaEveryMinute() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		count := uniqueRequestTracker.Count()

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

		// Reset unique requests for the next minute
		uniqueRequestTracker.Reset()
	}
}

// Send HTTP GET request to an endpoint
func sendHTTPRequest(endpoint string) {
	count := uniqueRequestTracker.Count()
	resp, err := http.Get(fmt.Sprintf("%s?count=%d", endpoint, count))
	if err != nil {
		log.Printf("Failed to send HTTP request: %v", err)
		return
	}
	log.Printf("Sent request to %s, Status Code: %d", endpoint, resp.StatusCode)
}

// Send HTTP POST request (Extension 1)
func sendHTTPPostRequest(endpoint string) {
	count := uniqueRequestTracker.Count()
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
