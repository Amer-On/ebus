package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"ebus/internal/domain"
	"ebus/pkg/logging" // замени на актуальный путь
)

const (
	totalMessages = 50000
	sendDuration  = 10 * time.Second
	workerCount   = 100
)

func main() {
	logger, err := logging.InitLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting producer")

	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     1000,
	}
	client := http.Client{
		Transport: transport,
	}

	messageCh := make(chan int, 1000) // буферизированный канал

	// Стартуем воркеры
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for range messageCh {
				publishMessage(logger, client)
			}
		}(i)
	}

	// Генерация сообщений с равномерным интервалом
	interval := sendDuration / totalMessages
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < totalMessages; i++ {
		<-ticker.C
		messageCh <- i
	}

	// Завершаем
	close(messageCh)
	logger.Info("Finished sending all messages")
}

type Payload struct {
	Data  string `json:"data"`
	Value string `json:"value"`
}

func publishMessage(logger *zap.Logger, httpClient http.Client) {
	payload := Payload{Data: "Hello", Value: "World"}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	event := domain.NewRawEvent("test_event", "ebus", rawPayload)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPost, "http://localhost:10000/publish", bytes.NewBuffer(eventBytes))
	if err != nil {
		panic(err)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic("Error reading body")
	}

	logger.Debug("Message published", zap.String("body", string(body)))

	if response.StatusCode != http.StatusOK {
		panic("Invalid status code" + response.Status)
	}
}
