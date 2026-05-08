package tools

import (
	"github.com/mcp-bank/proto/gen/brokerv1"
)

type Service struct {
	grpcClient brokerv1.BrokerServiceClient
}

func New(grpcClient brokerv1.BrokerServiceClient) *Service {
	return &Service{grpcClient: grpcClient}
}
