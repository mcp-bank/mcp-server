package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mcp-bank/proto/gen/brokerv1"
)

func (s *Service) HandleGetAccountBalance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var err error
	defer func() {
		if err != nil {
			slog.Error("HandleGetAccountBalance:",
				"err", err)
		}
	}()

	rawUserID, exists := request.GetArguments()["user_id"]
	if !exists {
		err = fmt.Errorf("user_id is required")
		return nil, err
	}
	userID, ok := rawUserID.(string)
	if !ok {
		err = fmt.Errorf("user_id must be string, got %T", rawUserID)
		return nil, err
	}
	portfolio, err := s.grpcClient.GetAccountBalance(ctx, &brokerv1.GetAccountBalanceRequest{Uuid: userID})
	if err != nil {
		err = fmt.Errorf("GetAccountBalance: %w", err)
		return nil, err
	}
	bytes, err := json.Marshal(portfolio)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(bytes)), nil
}
