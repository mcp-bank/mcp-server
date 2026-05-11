package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	notificationsWriter *kafka.Writer
}

func New() *Kafka {
	return &Kafka{notificationsWriter: &kafka.Writer{Addr: kafka.TCP("message-broker:9092"), Topic: "notifications"}}
}

func Init() error {
	var conn *kafka.Conn
	var err error
	defer func() {
		if conn != nil {
			if err = conn.Close(); err != nil {
				slog.Error("error closing connection to kafka in messaging.Init")
			}
		}
	}()
	for i := 0; i < 10; i++ {
		conn, err = kafka.Dial("tcp", "message-broker:9092")
		if err == nil {
			slog.Info("connection to kafka succeed")
			break
		}
		slog.Warn(fmt.Sprintf("trying connect to kafka, attempt $%v", i+1))
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		slog.Error("cannot connect to kafka",
			"err", err)
		return err
	}
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             "notifications",
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *Kafka) GracefulShutdown() error {
	err := k.notificationsWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (k *Kafka) PublishNotification(ctx context.Context, tool string) error {
	message := kafka.Message{
		Value: []byte(tool),
	}
	err := k.notificationsWriter.WriteMessages(ctx, message)
	if err != nil {
		err = fmt.Errorf("PublishNotification: %w", err)
		return err
	}
	return nil
}
