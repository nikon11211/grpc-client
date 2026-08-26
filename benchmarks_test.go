package grpc

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	benchSinkClient *Client
	benchSinkErr    error
	benchSinkInt    int
)

func benchConfig() *Config {
	return &Config{
		Service:            "bench-service",
		Address:            "localhost:50051",
		MaxCallRecvMsgSize: 4 * 1024 * 1024,
		RequestTimeout:     30 * time.Second,
	}
}

func benchConstructor(conn *grpc.ClientConn) any {
	return conn
}

func BenchmarkConfigValidate(b *testing.B) {
	cfg := benchConfig()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchSinkErr = cfg.Validate()
	}
}

func BenchmarkWithTracerProvider(b *testing.B) {
	var tp trace.TracerProvider = noop.NewTracerProvider()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		opt := WithTracerProvider(tp)
		o := &clientOptions{}
		opt(o)
		benchSinkInt += len(o.dialOptions)
	}
}

func BenchmarkOptionsApply(b *testing.B) {
	opts := []ClientOption{
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithInsecureSkipVerify(),
	}
	tp := noop.NewTracerProvider()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		o := &clientOptions{}
		for _, opt := range opts {
			opt(o)
		}
		WithTracerProvider(tp)(o)
		benchSinkInt += len(o.dialOptions)
	}
}

func BenchmarkNewClient(b *testing.B) {
	cfg := benchConfig()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		client, err := New(cfg, NoopLogger{}, benchConstructor)
		if err != nil {
			benchSinkErr = err
			continue
		}
		benchSinkClient = client
	}
}
