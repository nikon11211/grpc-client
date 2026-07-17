package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestWithTracerProvider(t *testing.T) {
	opts := &clientOptions{}
	tp := noop.NewTracerProvider()

	opt := WithTracerProvider(tp)
	opt(opts)

	assert.NotNil(t, opts.tracerProvider)
	assert.Equal(t, tp, opts.tracerProvider)
}

func TestWithDialOptions(t *testing.T) {
	opts := &clientOptions{}

	dialOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	opt := WithDialOptions(dialOpt)
	opt(opts)

	assert.Len(t, opts.dialOptions, 1)
}

func TestWithInsecureSkipVerify(t *testing.T) {
	opts := &clientOptions{}

	opt := WithInsecureSkipVerify()
	opt(opts)

	assert.Len(t, opts.dialOptions, 1)
}
