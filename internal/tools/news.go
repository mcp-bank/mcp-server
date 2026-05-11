package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mcp-bank/mcp-server/internal/messaging"
	"github.com/mcp-bank/proto/gen/brokerv1"
)

func (s *Service) HandleGetNews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var err error
	defer func(kafka *messaging.Kafka, ctx context.Context, tool string) {
		err = kafka.PublishNotification(ctx, tool)
		if err != nil {
			slog.Error("HandleGetNews:",
				"err", err)
		}
	}(s.kafka, ctx, "HandleGetNews")
	defer func() {
		if err != nil {
			slog.Error("HandleGetNews:",
				"err", err)
		}
	}()

	rawIsin, exists := request.GetArguments()["isin"]
	if !exists {
		err = fmt.Errorf("isin is required")
		return nil, err
	}
	isin, ok := rawIsin.(string)
	if !ok {
		err = fmt.Errorf("isin must be string, got %T", rawIsin)
		return nil, err
	}
	get := s.rdb.Get(ctx, isin)
	if err = get.Err(); err != nil {
		var news *brokerv1.GetNewsResponse
		start := time.Now()
		news, err = s.grpcClient.GetNews(ctx, &brokerv1.GetNewsRequest{Isin: isin})
		if err != nil {
			err = fmt.Errorf("GetNews: %w", err)
			return nil, err
		}
		var bytes []byte
		bytes, err = json.Marshal(news)
		if err != nil {
			return nil, err
		}
		s.rdb.Set(ctx, isin, bytes, time.Hour*24)
		slog.Info("without redis",
			"start_time", start,
			"duration", time.Since(start),
		)
		return mcp.NewToolResultText(string(bytes)), nil
	}
	start := time.Now()
	news, err := get.Result()
	if err != nil {
		err = fmt.Errorf("get.Result: %w", err)
		return nil, err
	}
	slog.Info("with redis",
		"start_time", start,
		"duration", time.Since(start),
	)
	return mcp.NewToolResultText(news), nil
}
