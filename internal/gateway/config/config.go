package config

import (
	"ebus/internal/broker"
	"ebus/internal/gateway/server"
	"ebus/pkg/redis"
)

type Config struct {
	Redis  redis.Config  `yaml:"redis"`
	Broker broker.Config `yaml:"broker"`
	Server server.Config `yaml:"server"`
}
