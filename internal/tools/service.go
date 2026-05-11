package tools

import (
	"github.com/mcp-bank/mcp-server/internal/messaging"
	"github.com/mcp-bank/proto/gen/brokerv1"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	grpcClient brokerv1.BrokerServiceClient
	rdb        *redis.Client
	kafka      *messaging.Kafka
}

func New(
	grpcClient brokerv1.BrokerServiceClient,
	rdb *redis.Client,
	kafka *messaging.Kafka,
) *Service {
	return &Service{
		grpcClient: grpcClient,
		rdb:        rdb,
		kafka:      kafka,
	}
}
