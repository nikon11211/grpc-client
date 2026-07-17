package grpc

import (
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConstructor func(*grpc.ClientConn) interface{}

type ClientOption func(*clientOptions)

type clientOptions struct {
	tracerProvider trace.TracerProvider
	dialOptions    []grpc.DialOption
}

func WithTracerProvider(tp trace.TracerProvider) ClientOption {
	return func(o *clientOptions) {
		o.tracerProvider = tp
	}
}

func WithDialOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.dialOptions = append(o.dialOptions, opts...)
	}
}

func WithInsecureSkipVerify() ClientOption {
	return func(o *clientOptions) {
		o.dialOptions = append(o.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}
