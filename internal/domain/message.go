package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	ID             string     `json:"id"`
	IdempotencyKey string     `json:"idempotency_key"`
	Topic          string     `json:"topic"`
	Subscriber     Subscriber `json:"subscriber"`
	Event          string     `json:"event"`
	Payload        []byte     `json:"payload"`
}

func NewMessage(topic, event, idempotencyKey string, subscriber Subscriber, payload []byte) *Message {
	return &Message{
		ID:             uuid.New().String(),
		IdempotencyKey: idempotencyKey,
		Topic:          topic,
		Subscriber:     subscriber,
		Event:          event,
		Payload:        payload,
	}
}

func UnmarshallMessage(data []byte) (*Message, error) {
	var message *Message

	err := json.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
