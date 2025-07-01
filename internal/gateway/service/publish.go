package service

import (
	"context"
	"ebus/internal/broker"
	"ebus/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PublicationService struct {
	logger              *zap.Logger
	subscriptionService *SubscribtionService
	broker              broker.Broker
}

func NewPublishService(
	logger *zap.Logger,
	subscriptionService *SubscribtionService,
	broker broker.Broker,
) *PublicationService {
	return &PublicationService{
		logger:              logger,
		subscriptionService: subscriptionService,
		broker:              broker,
	}
}

func (p *PublicationService) Publish(ctx context.Context, topic string, event domain.RawEvent) error {
	p.logger.Info("Publishing event", zap.Any("event", event))
	subscribers, err := p.subscriptionService.GetSubscribers(event)
	if err != nil {
		p.logger.Error("Failed to get subscribers")
		return err
	}
	p.logger.Info("Got subscribers", zap.Any("subscribers", subscribers))

	// Send messages via fanout mechanism
	for _, subscriber := range subscribers {
		idempotencyKey := uuid.New().String()

		message := domain.NewMessage(topic, event.Name, idempotencyKey, subscriber, event.Payload)
		err := p.broker.Publish(ctx, topic, message)
		if err != nil {
			p.logger.Error("Failed to publish message", zap.Any("message", message)) // actually needs retry logic - for now skip
			continue
		}
	}
	p.logger.Info("All messages sent")

	return nil
}
