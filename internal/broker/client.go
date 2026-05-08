package broker

import (
	"fmt"

	"github.com/mcp-bank/proto/gen/brokerv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New() (brokerv1.BrokerServiceClient, error) {
	client, err := grpc.NewClient("broker-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) // TODO убрать хардкод
	if err != nil {
		err = fmt.Errorf("new: %w", err)
		return nil, err
	}
	return brokerv1.NewBrokerServiceClient(client), nil
}
