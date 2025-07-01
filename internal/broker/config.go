package broker

type Config struct {
	ReadBatchSize int `yaml:"read_batch_size"`
	SleepTimeMS   int `yaml:"sleep_time_ms"`
}
