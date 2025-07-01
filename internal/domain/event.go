package domain

import "time"

type RawMeta map[string]string

type Meta struct {
	Timestamp time.Time `json:"timestamp"`
}

func NewMeta(rawMeta RawMeta) *Meta {
	// convert from raw meta

	return &Meta{
		Timestamp: time.Now(),
	}
}

type RawEvent struct {
	Payload []byte `json:"payload"`
	Name    string `json:"name"`
	Topic   string `json:"topic"`
	Meta    RawMeta
}

func NewRawEvent(name, topic string, payload []byte) *RawEvent {
	return &RawEvent{
		Payload: payload,
		Name:    name,
		Topic:   topic,
		// Meta: ,
	}
}

type Event[T any] struct {
	Topic   string
	Name    string
	Payload *T
	Meta    Meta
}

func NewEvent[T any](payload *T) *Event[T] {
	return &Event[T]{
		Payload: payload,
	}
}
