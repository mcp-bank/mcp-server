package middleware

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcp-bank/mcp-server/internal/metrics"
)

func Metrics() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			toolName := request.Params.Name
			metrics.RecordToolCall(toolName)
			start := time.Now()
			defer func() {
				metrics.RecordToolDuration(toolName, time.Since(start))
			}()
			result, err := next(ctx, request)
			if err != nil {
				metrics.RecordToolCallError(toolName)
				return nil, err
			}
			return result, nil
		}
	}
}
