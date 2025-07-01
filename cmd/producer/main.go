package main

import (
	"bytes"
	"ebus/internal/domain"
	"ebus/pkg/logging"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, err := logging.InitLogger()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Recovered from panic", zap.Error(err.(error)))
		}
		logger.Info("Stopping producer")
	}()

	logger.Info("Sending request")
	publishMessage(logger, *http.DefaultClient, nil)
}

type Payload struct {
	Data  string `json:"data"`
	Value string `json:"value"`
}

func publishMessage(logger *zap.Logger, httpClient http.Client, _ *domain.RawEvent) {
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

	logger.Info("Message published", zap.String("body", string(body)))

	if response.StatusCode != http.StatusOK {
		panic("Invalid status code" + response.Status)
	}
}
