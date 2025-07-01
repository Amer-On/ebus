package main

import (
	"context"
	"crypto/tls"
	"ebus/internal/broker"
	"ebus/internal/worker"
	"ebus/pkg/config"
	"ebus/pkg/logging"
	"ebus/pkg/redis"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     1000,
		IdleConnTimeout:     30 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Second,
	}

	var wg sync.WaitGroup

	workerCount := 4
	for i := range workerCount {
		wg.Add(1)
		go runWorker(ctx, &wg, i, httpClient)
	}

	<-sigs
	cancel()

	wg.Wait()
}

func runWorker(ctx context.Context, wg *sync.WaitGroup, id int, httpClient *http.Client) {
	logger, err := logging.InitLogger()
	if err != nil {
		panic("Error initializing logger")
	}

	logger.Info("Worker started", zap.Int("id", id))

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic", zap.Any("error", r))
		}
		wg.Done()
		logger.Info("Stopping worker", zap.Int("id", id))
	}()

	// Read config
	config, err := config.ReadConfig[worker.Config]("internal/worker/config.yaml")
	if err != nil {
		panic("Error reading config")
	}

	// Create Redis Client
	redisClient := redis.NewClient(config.Redis)

	// Create Broker
	broker := broker.NewRedisStreamsBroker(logger, redisClient)

	worker := worker.NewWorker(logger, broker, httpClient)

	worker.Start(ctx)
}
