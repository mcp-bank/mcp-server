package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mcp-bank/proto/gen/brokerv1"
)

func (s *Service) HandleGetStockPrice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var err error
	defer func() {
		if err != nil {
			slog.Error("HandleGetStockPrice failed:",
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
	portfolio, err := s.grpcClient.GetStockPrice(ctx, &brokerv1.GetStockPriceRequest{Isin: isin})
	if err != nil {
		err = fmt.Errorf("GetStockPrice: %w", err)
		return nil, err
	}
	bytes, err := json.Marshal(portfolio)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(bytes)), nil
}
