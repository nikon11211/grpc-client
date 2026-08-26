package grpc

import (
	"errors"
	"time"
)

type Config struct {
	Service            string        `mapstructure:"service" yaml:"service"`
	Address            string        `mapstructure:"address" yaml:"address" validate:"required"`
	MaxCallRecvMsgSize int           `mapstructure:"max_call_recv_msg_size" yaml:"max_call_recv_msg_size"`
	RequestTimeout     time.Duration `mapstructure:"request_timeout" yaml:"request_timeout"`
}

func DefaultConfig() Config {
	return Config{
		MaxCallRecvMsgSize: 4 * 1024 * 1024,
		RequestTimeout:     30 * time.Second,
	}
}

func (c Config) Validate() error {
	if c.Address == "" {
		return errors.New("address is required")
	}
	if c.MaxCallRecvMsgSize < 0 {
		return errors.New("max call recv msg size must be non-negative")
	}
	if c.RequestTimeout < 0 {
		return errors.New("request timeout must be non-negative")
	}
	return nil
}
