package main

import (
	"bytes"
	"ebus/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Subscription matches the domain.Subscription structure
type Subscription struct {
	Topic           string `json:"topic"`
	CallbackAddress string `json:"callback_address"`
}

func main() {
	// Subscribe to the gateway
	err := subscribe("http://localhost:10000/subscribe", domain.Subscription{
		Subscriber: domain.Subscriber{
			CallbackAddress: "http://localhost:10001/callback",
			Name:            "first_consumer",
		},
		Topic: "ebus",
		Event: "*",
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Start HTTP server to receive callbacks
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Read the incoming message body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Log the message
		log.Printf("Received message: %s", string(body))

		// Respond
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Message received")
	})

	fmt.Println("Consumer started on :10001")
	if err := http.ListenAndServe(":10001", r); err != nil {
		log.Fatalf("Consumer server error: %v", err)
	}
}

// subscribe sends a POST request to the /subscribe endpoint
func subscribe(url string, sub domain.Subscription) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to send subscription request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("subscription failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Subscribed successfully: topic=%s, callback=%s", sub.Topic, sub.CallbackAddress)
	return nil
}
