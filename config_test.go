package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, 4*1024*1024, cfg.MaxCallRecvMsgSize)
	assert.Equal(t, 30*time.Second, cfg.RequestTimeout)
	assert.Empty(t, cfg.Service)
	assert.Empty(t, cfg.Address)
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: Config{
				Address:            "localhost:50051",
				MaxCallRecvMsgSize: 1024,
				RequestTimeout:     10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing address",
			cfg: Config{
				MaxCallRecvMsgSize: 1024,
			},
			wantErr: true,
			errMsg:  "address is required",
		},
		{
			name: "negative max msg size",
			cfg: Config{
				Address:            "localhost:50051",
				MaxCallRecvMsgSize: -1,
			},
			wantErr: true,
			errMsg:  "max call recv msg size must be non-negative",
		},
		{
			name: "negative timeout",
			cfg: Config{
				Address:        "localhost:50051",
				RequestTimeout: -1 * time.Second,
			},
			wantErr: true,
			errMsg:  "request timeout must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
