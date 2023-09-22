package kafka

import (
	"companies-service/config"
	"context"

	"github.com/segmentio/kafka-go"
)

// NewKafkaConn create new kafka connection
func NewKafkaConn(ctx context.Context, kafkaCfg *config.Kafka) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", kafkaCfg.Brokers[0])
}
