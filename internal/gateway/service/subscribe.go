package service

import (
	"ebus/internal/domain"

	"go.uber.org/zap"
)

type SubscribtionService struct {
	logger        *zap.Logger
	subscriptions []*domain.Subscription
}

// STUB REALIZATION
func NewSubscriptionService(logger *zap.Logger) *SubscribtionService {
	return &SubscribtionService{
		logger:        logger,
		subscriptions: make([]*domain.Subscription, 0),
	}
}

// ATM stored in the memory, later to be stored in redis/postgres
func (s *SubscribtionService) Subscribe(subscription domain.Subscription) error {
	// TODO: Validate the callback address by sending an http request to it (e.g. ping pong)
	s.subscriptions = append(s.subscriptions, &subscription)
	s.logger.Info("Subscription added", zap.Any("subscription", subscription))

	return nil
}

func (s *SubscribtionService) Unsubscribe(topic string, event string, callbackAddress string) error {
	for idx, subscription := range s.subscriptions {
		if subscription.Topic != topic {
			continue
		}
		if event == "*" || subscription.Event != event {
			continue
		}

		s.subscriptions[idx] = s.subscriptions[len(s.subscriptions)-1]
		s.subscriptions = s.subscriptions[:len(s.subscriptions)-1]
	}

	return nil
}

// ATM stored in the memory, later to be stored in redis/postgres
func (s *SubscribtionService) GetSubscribers(event domain.RawEvent) ([]domain.Subscriber, error) {
	subscribersSet := make(map[string]string)
	subscribers := make([]domain.Subscriber, 0, 12) // 12 от балды

	for _, subscription := range s.subscriptions {
		if subscription.Topic != event.Topic {
			continue
		} else if subscription.Event == "*" || subscription.Event == event.Name {
			if _, ok := subscribersSet[subscription.CallbackAddress]; !ok {
				subscribers = append(subscribers, subscription.Subscriber)
			}
			subscribersSet[subscription.CallbackAddress] = ""
		}
	}

	return subscribers, nil
}
