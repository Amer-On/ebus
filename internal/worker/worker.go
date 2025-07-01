package worker

import (
	"bytes"
	"context"
	"ebus/internal/broker"
	"ebus/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Message for final subscriber
type Message struct {
	ID             string `json:"id"`
	IdempotencyKey string `json:"idempotency_key"`
	Topic          string `json:"topic"`
	Event          string `json:"event"`
	Payload        []byte `json:"payload"`
}

type Worker struct {
	ID             string
	httpClient     *http.Client
	logger         *zap.Logger
	broker         broker.Broker
	failedMessages int
}

func NewWorker(logger *zap.Logger, broker broker.Broker, httpClient *http.Client) *Worker {
	return &Worker{
		ID:         uuid.NewString(),
		httpClient: httpClient,
		logger:     logger,
		broker:     broker,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		if err := recover(); err != nil {
			w.logger.Error("Worker recovered from panic", zap.Error(err.(error)))
		}

		w.httpClient.CloseIdleConnections()
		cancel()
	}()

	w.logger.Info("Worker started", zap.String("ID", w.ID))
	w.broker.Subscribe(ctx, "eventbus", w.HandleMessage, fmt.Sprintf("worker_%s", w.ID))
	w.logger.Info("Worker stopped", zap.String("ID", w.ID))
}

func (w *Worker) HandleMessage(ctx context.Context, message *domain.Message) error {
	w.logger.Info("Received message", zap.Any("message", message))
	payload, err := json.Marshal(message)
	if err != nil {
		w.logger.Error("Failed to marshal message: %v", zap.Error(err))
		return err
	}

	err = w.sendHttpEvent(ctx, message.Subscriber.CallbackAddress, payload, message.IdempotencyKey)
	if err != nil {
		w.logger.Error("Failed to send message via http: %v", zap.Error(err))
		w.failedMessages++
		return err
	}
	return nil
}

func (w *Worker) sendHttpEvent(ctx context.Context, callbackAddress string, payload []byte, idempotencyKey string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		callbackAddress,
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Idempotency-Key", idempotencyKey)

	response, err := w.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send http event: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("HTTP error: %s", response.Status)
}
