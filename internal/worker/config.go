package worker

import (
	"ebus/internal/broker"
	"ebus/pkg/redis"
)

type Config struct {
	Redis  redis.Config  `json:"redis"`
	Broker broker.Config `json:"broker"`
}
