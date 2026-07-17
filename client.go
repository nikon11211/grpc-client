package grpc

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	logger Logger
	*grpc.ClientConn
	Client interface{}
}

func New(cfg *Config, logger Logger, constructor ClientConstructor, opts ...ClientOption) (*Client, error) {
	const trace = "grpc.NewClient"

	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if logger == nil {
		logger = NoopLogger{}
	}

	if constructor == nil {
		return nil, fmt.Errorf("client constructor is required")
	}

	options := &clientOptions{}
	for _, opt := range opts {
		opt(options)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if cfg.MaxCallRecvMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.MaxCallRecvMsgSize),
		))
	}

	if options.tracerProvider != nil {
		statsHandler := otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(options.tracerProvider),
			otelgrpc.WithPropagators(propagation.TraceContext{}),
		)
		dialOpts = append(dialOpts, grpc.WithStatsHandler(statsHandler))
		logger.Info(fmt.Sprintf("(%s) OpenTelemetry tracing enabled", trace))
	}

	dialOpts = append(dialOpts, options.dialOptions...)

	logger.DebugF("(%s) connecting to %s", trace, cfg.Address)

	conn, err := grpc.NewClient(cfg.Address, dialOpts...)
	if err != nil {
		logger.Error(fmt.Sprintf("(%s) failed to create gRPC client: %s", trace, err.Error()))
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	client := constructor(conn)

	logger.Info(fmt.Sprintf("(%s) client created successfully", trace))

	return &Client{
		logger:     logger,
		ClientConn: conn,
		Client:     client,
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.ClientConn == nil {
		return nil
	}
	return c.ClientConn.Close()
}

func (c *Client) GetLogger() Logger {
	if c == nil {
		return NoopLogger{}
	}
	return c.logger
}
