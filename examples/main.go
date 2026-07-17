package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	grpcclient "github.com/nikon11211/grpc-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExampleServiceClient struct {
	*grpc.ClientConn
}

func main() {
	cfg := &grpcclient.Config{
		Service:            "example-service",
		Address:            "localhost:50051",
		MaxCallRecvMsgSize: 4 * 1024 * 1024,
		RequestTimeout:     30,
	}

	logger := &CustomLogger{}

	constructor := func(conn *grpc.ClientConn) interface{} {
		return &ExampleServiceClient{ClientConn: conn}
	}

	client, err := grpcclient.New(
		cfg,
		logger,
		constructor,
		grpcclient.WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create gRPC client: %v", err))
	}
	defer client.Close()

	logger.Info("gRPC client started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gRPC client...")
}

type CustomLogger struct{}

func (c *CustomLogger) DebugF(format string, args ...any) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}

func (c *CustomLogger) Debug(msg string) {
	fmt.Printf("[DEBUG] %s\n", msg)
}

func (c *CustomLogger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (c *CustomLogger) Warn(msg string) {
	fmt.Printf("[WARN] %s\n", msg)
}

func (c *CustomLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}
