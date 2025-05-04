package cfg

import (
	"fmt"
	"time"

	"github.com/alexflint/go-arg"
)

type (
	Config struct {
		Server    Server
		DragonFly DragonFly
	}

	Server struct {
		Port string `arg:"env:SERVER_PORT,required"`
	}

	DragonFly struct {
		Address         string        `arg:"env:DRAGONFLY_ADDRESS,required"`
		Password        string        `arg:"env:DRAGONFLY_PASSWORD,required"`
		DB              int           `arg:"env:DRAGONFLY_DB,required"`
		ReadTimeout     time.Duration `arg:"env:DRAGONFLY_READ_TIMEOUT" default:"3s"`
		WriteTimeout    time.Duration `arg:"env:DRAGONFLY_WRITE_TIMEOUT" default:"3s"`
		MaxRetries      int           `arg:"env:DRAGONFLY_MAX_RETRIES" default:"3"`
		MinRetryBackoff time.Duration `arg:"env:DRAGONFLY_MIN_RETRY_BACKOFF" default:"500ms"`
		MaxRetryBackoff time.Duration `arg:"env:DRAGONFLY_MAX_RETRY_BACKOFF" default:"5s"`
	}
)

func NewConfig() (*Config, error) {
	var cfg Config

	if err := arg.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("arg.Parse: %w", err)
	}

	return &cfg, nil
}
