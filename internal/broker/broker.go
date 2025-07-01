package broker

import (
	"context"

	"ebus/internal/domain"
)

type CallbackFunction func(context.Context, *domain.Message) error

type Broker interface {
	Publish(ctx context.Context, topic string, message *domain.Message) error
	Subscribe(ctx context.Context, topic string, callbackFunction CallbackFunction, consumerName string) error
	Unsubscribe(ctx context.Context, consumerName string) error
}
