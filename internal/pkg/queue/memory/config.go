package memory

import "time"

type Config struct {
	BatchSize     int           `config:"batch_size"`
	BufferSize    int           `config:"buffer_size"`
	FlushInterval time.Duration `config:"flush_interval"`
}
