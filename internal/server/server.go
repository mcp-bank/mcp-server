package server

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcp-bank/mcp-server/internal/tools"
)

type Server struct {
	McpServer *server.MCPServer
	tools     *tools.Service
}

func New(tools *tools.Service) *Server {
	return &Server{
		McpServer: server.NewMCPServer("mcp-bank/mcp-server", "v0.0.1", server.WithRecovery()), // TODO убрать хардкод
		tools:     tools,
	}
}

func (s *Server) RegisterTools() {
	getPortfolio := mcp.NewTool("get_portfolio",
		mcp.WithDescription("Get user portfolio"),
		mcp.WithString("user_id",
			mcp.Required(),
			mcp.Description("User ID")),
	)
	getStockPrice := mcp.NewTool("get_stock_price",
		mcp.WithDescription("Get user stock_price"),
		mcp.WithString("isin",
			mcp.Required(),
			mcp.Description("ISIN")),
	)
	getAccountBalance := mcp.NewTool("get_account_balance",
		mcp.WithDescription("Get user account_balance"),
		mcp.WithString("user_id",
			mcp.Required(),
			mcp.Description("User ID")),
	)
	getNews := mcp.NewTool("get_news",
		mcp.WithDescription("Get user news"),
		mcp.WithString("isin",
			mcp.Required(),
			mcp.Description("ISIN")),
	)

	s.McpServer.AddTool(getPortfolio, s.tools.HandleGetPortfolio)
	s.McpServer.AddTool(getStockPrice, s.tools.HandleGetStockPrice)
	s.McpServer.AddTool(getAccountBalance, s.tools.HandleGetAccountBalance)
	s.McpServer.AddTool(getNews, s.tools.HandleGetNews)
}
