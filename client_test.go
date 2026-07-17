package grpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
)

type mockLogger struct {
	debugFCount int
	debugCount  int
	infoCount   int
	warnCount   int
	errorCount  int
	messages    []string
}

func (m *mockLogger) DebugF(format string, args ...any) {
	m.debugFCount++
	m.messages = append(m.messages, format)
}

func (m *mockLogger) Debug(msg string) {
	m.debugCount++
	m.messages = append(m.messages, msg)
}

func (m *mockLogger) Info(msg string) {
	m.infoCount++
	m.messages = append(m.messages, msg)
}

func (m *mockLogger) Warn(msg string) {
	m.warnCount++
	m.messages = append(m.messages, msg)
}

func (m *mockLogger) Error(msg string) {
	m.errorCount++
	m.messages = append(m.messages, msg)
}

func TestNoopLogger(t *testing.T) {
	logger := NoopLogger{}

	assert.NotPanics(t, func() {
		logger.DebugF("test %s", "arg")
		logger.Debug("test")
		logger.Info("test")
		logger.Warn("test")
		logger.Error("test")
	})
}

func TestNewClientNilConfig(t *testing.T) {
	client, err := New(nil, NoopLogger{}, nil)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "config is required")
}

func TestNewClientInvalidConfig(t *testing.T) {
	cfg := &Config{
		Address: "",
	}

	client, err := New(cfg, NoopLogger{}, nil)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "invalid config")
}

func TestNewClientNilConstructor(t *testing.T) {
	cfg := &Config{
		Address: "localhost:50051",
	}

	client, err := New(cfg, NoopLogger{}, nil)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "client constructor is required")
}

func TestNewClientNilLogger(t *testing.T) {
	cfg := &Config{
		Address: "localhost:50051",
	}

	constructor := func(conn *grpc.ClientConn) any {
		return conn
	}

	client, err := New(cfg, nil, constructor)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing client: %v\n", err)
		}
	}(client)

	assert.IsType(t, NoopLogger{}, client.logger)
}

func TestNewClientWithTracerProvider(t *testing.T) {
	cfg := &Config{
		Address: "localhost:50051",
	}

	logger := &mockLogger{}
	tp := noop.NewTracerProvider()

	constructor := func(conn *grpc.ClientConn) any {
		return conn
	}

	client, err := New(cfg, logger, constructor, WithTracerProvider(tp))
	require.NoError(t, err)
	require.NotNil(t, client)
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing client: %v\n", err)
		}
	}(client)

	assert.Greater(t, logger.infoCount, 0)
}

func TestNewClientWithMaxMsgSize(t *testing.T) {
	cfg := &Config{
		Address:            "localhost:50051",
		MaxCallRecvMsgSize: 1024 * 1024,
	}

	constructor := func(conn *grpc.ClientConn) any {
		return conn
	}

	client, err := New(cfg, NoopLogger{}, constructor)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing client: %v\n", err)
		}
	}(client)
}

func TestNewClientWithDialOptions(t *testing.T) {
	cfg := &Config{
		Address: "localhost:50051",
	}

	constructor := func(conn *grpc.ClientConn) any {
		return conn
	}

	client, err := New(cfg, NoopLogger{}, constructor,
		WithInsecureSkipVerify(),
	)

	require.NoError(t, err)
	require.NotNil(t, client)
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing client: %v\n", err)
		}
	}(client)
}

func TestClientClose(t *testing.T) {
	t.Run("close valid client", func(t *testing.T) {
		cfg := &Config{
			Address: "localhost:50051",
		}

		constructor := func(conn *grpc.ClientConn) any {
			return conn
		}

		client, err := New(cfg, NoopLogger{}, constructor)
		require.NoError(t, err)
		require.NotNil(t, client)

		err = client.Close()
		assert.NoError(t, err)
	})

	t.Run("close nil client", func(t *testing.T) {
		var client *Client
		err := client.Close()
		assert.NoError(t, err)
	})

	t.Run("close client with nil connection", func(t *testing.T) {
		client := &Client{}
		err := client.Close()
		assert.NoError(t, err)
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("valid client", func(t *testing.T) {
		cfg := &Config{
			Address: "localhost:50051",
		}

		logger := &mockLogger{}
		constructor := func(conn *grpc.ClientConn) any {
			return conn
		}

		client, err := New(cfg, logger, constructor)
		require.NoError(t, err)
		defer func(client *Client) {
			err := client.Close()
			if err != nil {
				fmt.Printf("Error closing client: %v\n", err)
			}
		}(client)

		assert.Equal(t, logger, client.GetLogger())
	})

	t.Run("nil client", func(t *testing.T) {
		var client *Client
		logger := client.GetLogger()
		assert.IsType(t, NoopLogger{}, logger)
	})
}

func TestClientConnection(t *testing.T) {
	cfg := &Config{
		Address: "localhost:50051",
	}

	constructor := func(conn *grpc.ClientConn) any {
		return conn
	}

	client, err := New(cfg, NoopLogger{}, constructor)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer func(client *Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing client: %v\n", err)
		}
	}(client)

	assert.NotNil(t, client.ClientConn)
}

func TestNewClientConnectionError(t *testing.T) {
	cfg := &Config{
		Address: "invalid-address-with-special-chars://",
	}

	logger := &mockLogger{}
	constructor := func(conn *grpc.ClientConn) any {
		return conn
	}

	client, err := New(cfg, logger, constructor)
	if err == nil {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing client: %v\n", err)
			return
		}
	}
}

func TestMockLogger(t *testing.T) {
	logger := &mockLogger{}

	logger.DebugF("debugf %s", "arg")
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	assert.Equal(t, 1, logger.debugFCount)
	assert.Equal(t, 1, logger.debugCount)
	assert.Equal(t, 1, logger.infoCount)
	assert.Equal(t, 1, logger.warnCount)
	assert.Equal(t, 1, logger.errorCount)
	assert.Len(t, logger.messages, 5)
}
