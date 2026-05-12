package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mcpServerPkg "github.com/mark3labs/mcp-go/server"
	"github.com/mcp-bank/mcp-server/internal/broker"
	"github.com/mcp-bank/mcp-server/internal/cache"
	"github.com/mcp-bank/mcp-server/internal/messaging"
	"github.com/mcp-bank/mcp-server/internal/metrics"
	"github.com/mcp-bank/mcp-server/internal/middleware"
	"github.com/mcp-bank/mcp-server/internal/server"
	"github.com/mcp-bank/mcp-server/internal/tools"
	"github.com/redis/go-redis/v9"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	grpcClient, err := broker.New()
	if err != nil {
		slog.Error("error creating broker",
			"err", err)
		return
	}
	err = messaging.Init()
	if err != nil {
		slog.Error("error initializing broker",
			"err", err)
		return
	}
	kafka := messaging.New()
	rdb, err := cache.New()
	if err != nil {
		slog.Error("error creating caching",
			"err", err)
		return
	}
	service := tools.New(grpcClient, rdb, kafka)
	mcpServer := server.New(service)
	mcpServer.RegisterTools()
	mcpServer.McpServer.Use(middleware.Metrics())
	sseServer := mcpServerPkg.NewSSEServer(mcpServer.McpServer, mcpServerPkg.WithBaseURL("http://mcp-server:8080")) // TODO убрать хардкод
	metricsServer := metrics.NewServer()
	go func() {
		err = sseServer.Start(":8080") // TODO убрать хардкод
		if err != nil {
			slog.Error("stopping sseServer (may be ok, if stopping with graceful shutdown)",
				"err", err)
			return
		}
	}()
	go func() {
		err = metricsServer.ListenAndServe()
		if err != nil {
			slog.Error("stopping metrics server (may be ok, if stopping with graceful shutdown)",
				"err", err)
			return
		}
	}()
	<-quit
	gracefulShutdown(kafka, rdb, metricsServer, sseServer)
}

func gracefulShutdown(kafka *messaging.Kafka, rdb *redis.Client, metricsServer *http.Server, sseServer *mcpServerPkg.SSEServer) {
	slog.Info("graceful shutdown")
	if err := kafka.GracefulShutdown(); err != nil {
		slog.Error("cannot properly shutdown kafka",
			"err", err)
	}
	if err := rdb.Close(); err != nil {
		slog.Error("cannot properly shutdown redis",
			"err", err)
	}
	if err := metricsServer.Shutdown(context.Background()); err != nil {
		slog.Error("cannot properly shutdown server",
			"err", err)
	}
	if err := sseServer.Shutdown(context.Background()); err != nil {
		slog.Error("cannot properly shutdown server",
			"err", err)
	}
}
