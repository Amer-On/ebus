package broker

import (
	"context"
	"ebus/internal/domain"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	idempotencyKey string = "idempotencyKey"
	topic          string = "ebus"
	consumerGroup  string = "ebus_workers"
	dataKey        string = "data"
)

type RedisStreamsBroker struct {
	logger        *zap.Logger
	redisClient   *redis.Client
	consumerGroup string
}

func NewRedisStreamsBroker(logger *zap.Logger, redisClient *redis.Client) *RedisStreamsBroker {

	return &RedisStreamsBroker{
		logger:        logger,
		redisClient:   redisClient,
		consumerGroup: consumerGroup,
	}
}

func (b *RedisStreamsBroker) Publish(ctx context.Context, topic string, message *domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = b.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		Values: map[string]any{dataKey: data},
	}).Result()
	if err != nil {
		return fmt.Errorf("error publishing message to stream: %v", err)
	}

	return nil
}

func (b *RedisStreamsBroker) Subscribe(ctx context.Context, _ string, callbackFunction CallbackFunction, consumerName string) error {
	_, err := b.redisClient.XGroupCreateMkStream(ctx, topic, b.consumerGroup, "0").Result()
	// Make this check more elegant. Now it feels pretty weird.
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("could not create consumer group: %v", err)
	}
	return b.processMessages(ctx, callbackFunction, consumerName)
}

func (b *RedisStreamsBroker) processMessages(ctx context.Context, callbackFunction CallbackFunction, consumerName string) error {
	for {
		select {
		case <-ctx.Done():
			b.logger.Info("Subscription Stopped")
			return nil
		default:
			streams, err := b.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    b.consumerGroup,
				Consumer: consumerName,
				Streams:  []string{topic, ">"},
				Count:    10,
				Block:    5 * time.Millisecond,
			}).Result()

			switch err {
			case redis.Nil:
				b.logger.Debug("No messages available")
			default:
				b.logger.Error("Error reading from stream", zap.Error(err))
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					data, ok := message.Values[dataKey]
					if !ok {
						b.logger.Error(
							"Error getting message content",
							zap.Error(err),
						)
						b.redisClient.XAck(ctx, topic, b.consumerGroup, message.ID)
						continue
					}

					var bytes []byte
					switch v := data.(type) {
					case []byte:
						bytes = v
					case string:
						bytes = []byte(v)
					default:
						b.logger.Error(
							"Unsupported type for message content",
							zap.Any("type", reflect.TypeOf(data)),
						)
						b.redisClient.XAck(ctx, topic, b.consumerGroup, message.ID)
						continue
					}

					var domainMessage domain.Message
					if err := json.Unmarshal(bytes, &domainMessage); err != nil {
						b.logger.Error(
							"Error deserializing message",
							zap.Error(err),
							zap.String("message", string(bytes)),
						)
						b.redisClient.XAck(ctx, topic, b.consumerGroup, message.ID)
						continue
					}
					domainMessage.ID = message.ID

					// Idempotency Key check
					// exists, _ := b.redisClient.SIsMember(ctx, idempotencyKey, message.IdempotencyKey).Result()
					// if exists {
					//  b.redisClient.XAck(ctx, topic, b.consumerGroup, message.ID)
					// 	b.logger.Error("Message has already been processed")
					// }

					b.logger.Info(
						"Received message",
						zap.String(idempotencyKey, domainMessage.IdempotencyKey),
						zap.String(idempotencyKey, message.ID),
						zap.String("topic", domainMessage.Topic),
						zap.String("event", domainMessage.Event),
					)

					// Process the message here
					err := callbackFunction(ctx, &domainMessage)
					if err != nil {
						b.logger.Error(
							"Failed to process message",
							zap.String("ID", message.ID),
							zap.String(idempotencyKey, domainMessage.IdempotencyKey),
						)
					}

					// Add idempotency key to Redis to prevent duplicates
					// _, err = b.redisClient.SAdd(ctx, idempotencyKey, message.IdempotencyKey).Result()

					// Acknoledge message
					if _, err := b.redisClient.XAck(ctx, topic, b.consumerGroup, message.ID).Result(); err != nil {
						b.logger.Error(
							"Could not acknoledge message",
							zap.String("message_id", message.ID),
							zap.Error(err),
						)
					}
				}
			}
		}
	}
}

// Function to be called on worker stop for graceful shutdown
func (b *RedisStreamsBroker) Unsubscribe(ctx context.Context, consumerName string) error {
	err := b.redisClient.Close()
	if err != nil {
		panic(fmt.Errorf("error closing Redis client: %v", err))
	}

	return nil
}
