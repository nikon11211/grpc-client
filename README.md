<h1 align="center">gRPC Client - Enterprise-Grade gRPC Client Library for Go</h1>

<p align="center">
  <a href="https://pkg.go.dev/github.com/nikon11211/grpc-client">
    <img src="https://pkg.go.dev/badge/github.com/nikon11211/grpc-client.svg" alt="Go Reference"/>
  </a>
  <a href="https://goreportcard.com/report/github.com/nikon11211/grpc-client">
    <img src="https://goreportcard.com/badge/github.com/nikon11211/grpc-client" alt="Go Report Card"/>
  </a>
  <a href="https://github.com/nikon11211/grpc-client/actions/workflows/test.yaml">
    <img src="https://github.com/nikon11211/grpc-client/actions/workflows/test.yaml/badge.svg" alt="Tests"/>
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"/>
  </a>
  <a href="https://golang.org/">
    <img src="https://img.shields.io/badge/Go-%3E%3D%201.21-blue" alt="Go Version"/>
  </a>
</p>

<p align="center">
  <b>A high-performance, production-ready gRPC client library for Go microservices</b><br/>
  <i>OpenTelemetry tracing • Custom logger interface • Functional options • Connection management</i>
</p>

---

## ✨ Why gRPC Client?

This library provides a robust, type-safe wrapper around Google's gRPC client with built-in support for distributed tracing, structured logging, and flexible configuration. Designed for microservices architectures where observability and reliability are paramount.

```go
// Simple, yet powerful
client, err := grpcclient.New(cfg, logger, constructor,
grpcclient.WithTracerProvider(tp),
grpcclient.WithDialOptions(grpc.WithInsecure()),
)

// Type-safe client access
exampleClient := client.Client.(ExampleServiceClient)
```

---

## 🎯 Features

<table>
<tr>
<td width="50%">

### 🚀 Core Features
- Type-safe client wrapper with generic client constructor pattern
- Functional options for flexible configuration
- Built-in OpenTelemetry tracing support
- Custom logger interface - bring your own logger
- Connection lifecycle management with graceful shutdown
- Message size configuration** for large payloads
- Request timeout support

### 📊 Observability
- Automatic trace context propagation via OpenTelemetry
- Structured logging at all lifecycle events
- Connection state logging for debugging
- Error tracing with detailed context

</td>
<td width="50%">

### 🔒 Enterprise Ready
- **TLS support via gRPC dial options
- **Custom dial options for advanced configuration
- **Nil-safe design with noop logger fallback
- **Validation of all configuration parameters
- **100% test coverage

</td>
</tr>
</table>

## 📦 Installation

```bash
go get github.com/nikon11211/grpc-client
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Your Application                      │
├─────────────────────────────────────────────────────────┤
│                    gRPC Client Library                  │
│                                                         │
│  ┌──────────┐ ┌──────────┐ ┌────────────────────────┐   │
│  │  Config  │ │ Options  │ │  OpenTelemetry Tracer  │   │
│  └────┬─────┘ └────┬─────┘ └───────────┬────────────┘   │
│       └────────────┴───────────────────┘                │
│                        │                                │
│                   Client Wrapper                        │
│                        │                                │
│              google.golang.org/grpc                     │
└─────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Basic Example

```go
package main

import (
	grpcclient "github.com/nikon11211/grpc-client"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	*grpc.ClientConn
}

func main() {
	cfg := &grpcclient.Config{
		Address:            "user-service:50051",
		MaxCallRecvMsgSize: 4 * 1024 * 1024,
	}

	logger := grpcclient.NoopLogger{}

	constructor := func(conn *grpc.ClientConn) interface{} {
		return &UserServiceClient{ClientConn: conn}
	}

	client, err := grpcclient.New(cfg, logger, constructor)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Use the typed client
	userClient := client.Client.(*UserServiceClient)
	// Make gRPC calls...
}
```

### With OpenTelemetry Tracing

```go
import (
"go.opentelemetry.io/otel/trace/noop"
)

tp := noop.NewTracerProvider()

client, err := grpcclient.New(cfg, logger, constructor,
grpcclient.WithTracerProvider(tp),
)
```

### Distributed Tracing with OpenTelemetry

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func handleOrder(ctx context.Context, orderID string) {
    tracer := otel.Tracer("order-service")
    ctx, span := tracer.Start(ctx, "handleOrder")
    defer span.End()
    
    // TraceID and SpanID automatically injected into logs
    log.InfoCtx(ctx, "Processing order")
    
    // All downstream logs will include trace context
    processPayment(ctx, orderID)
}
```

### With Custom Dial Options

```go
client, err := grpcclient.New(cfg, logger, constructor,
grpcclient.WithDialOptions(
grpc.WithInsecure(),
grpc.WithBlock(),
grpc.WithTimeout(5 * time.Second),
),
grpcclient.WithInsecureSkipVerify(),
)
```

## 📚 Advanced Usage

### Functional Options Pattern

```go
client, err := grpcclient.New(
cfg,
logger,
constructor,
grpcclient.WithTracerProvider(tp),
grpcclient.WithDialOptions(
grpc.WithKeepaliveParams(keepalive.ClientParameters{
Time:    10 * time.Second,
Timeout: 3 * time.Second,
}),
),
)
```

### Error Handling

```go
client, err := grpcclient.New(cfg, logger, constructor)
if err != nil {
switch {
case strings.Contains(err.Error(), "config is required"):
panic("Configuration must be provided")
case strings.Contains(err.Error(), "invalid config"):
panic("Configuration validation failed")
default:
panic(fmt.Sprintf("Unexpected error: %v", err))
}
}
```

## 🔧 Configuration Reference

```go
type Config struct {
// Service is the name of the gRPC service
Service string

// Address is the target address of the gRPC server (required)
Address string

// MaxCallRecvMsgSize is the maximum message size the client can receive (default: 4MB)
MaxCallRecvMsgSize int

// RequestTimeout is the default timeout for gRPC requests (default: 30s)
RequestTimeout time.Duration
}
```

## 🧪 Testing

```go
// Run all tests
go test ./...

// Run with race detection
go test -race ./...
 
// Run with coverage
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt

// Run benchmarks
go test -bench=. -benchmem
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🌟 Show Your Support

Give a ⭐️ if this project helped you! Share it with your team to improve logging across your microservices.

## 🙏 Acknowledgments


- [gRPC-Go](https://google.golang.org/grpc/) - The official Go gRPC implementation
- [OpenTelemetry](https://opentelemetry.io/) - Distributed tracing standard
- [Logger] - Our companion logging library
---

<p align="center">
  <b>Made with ❤️ for the Go community</b><br/>
  <sub>Built for performance, designed for reliability</sub>
</p>