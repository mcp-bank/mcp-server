package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mcp-bank/proto/gen/brokerv1"
)

func (s *Service) HandleGetNews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var err error
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
		return mcp.NewToolResultText(string(bytes)), nil
	}
	news, err := get.Result()
	if err != nil {
		err = fmt.Errorf("get.Result: %w", err)
		return nil, err
	}
	return mcp.NewToolResultText(news), nil
}
