package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mcpServerPkg "github.com/mark3labs/mcp-go/server"
	"github.com/mcp-bank/mcp-server/internal/broker"
	"github.com/mcp-bank/mcp-server/internal/cache"
	"github.com/mcp-bank/mcp-server/internal/messaging"
	"github.com/mcp-bank/mcp-server/internal/server"
	"github.com/mcp-bank/mcp-server/internal/tools"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			slog.Error("main:",
				"err", err,
			)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	grpcClient, err := broker.New()
	if err != nil {
		return
	}
	err = messaging.Init()
	if err != nil {
		return
	}
	kafka := messaging.New()
	rdb, err := cache.New()
	if err != nil {
		return
	}
	service := tools.New(grpcClient, rdb, kafka)
	mcpServer := server.New(service)
	mcpServer.RegisterTools()
	sseServer := mcpServerPkg.NewSSEServer(mcpServer.McpServer, mcpServerPkg.WithBaseURL("http://mcp-server:8080")) // TODO убрать хардкод
	go func() {
		err = sseServer.Start(":8080") // TODO убрать хардкод
		if err != nil {
			slog.Error("stopping sseServer (may be ok, if stopping with graceful shutdown)",
				"err", err)
			return
		}
	}()
	<-quit
	slog.Info("graceful shutdown")
	if err = kafka.GracefulShutdown(); err != nil {
		err = fmt.Errorf("cannot properly shutdown kafka %w", err)
	}
	if err = rdb.Close(); err != nil {
		err = fmt.Errorf("cannot properly shutdown redis %w", err)
	}
	if err = sseServer.Shutdown(context.Background()); err != nil {
		err = fmt.Errorf("cannot properly shutdown server %w", err)
	}
}
