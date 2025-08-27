// internal/messaging/kafka.go
package messaging

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
    writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *KafkaProducer {
    return &KafkaProducer{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(broker),
            Topic:    topic,
            Balancer: &kafka.LeastBytes{},
        },
    }
}

func (p *KafkaProducer) Publish(ctx context.Context, key string, event interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    return p.writer.WriteMessages(ctx,
        kafka.Message{
            Key:   []byte(key),
            Value: data,
        },
    )
}
