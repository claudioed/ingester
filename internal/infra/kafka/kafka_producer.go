package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	"os"
)

func NewSender() (*kafka_sarama.Sender, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	sender, err := kafka_sarama.NewSender([]string{os.Getenv("KAFKA_ADDRESSES")}, saramaConfig, os.Getenv("KAFKA_TOPIC"))
	if err != nil {
		return nil, err
	}
	return sender, nil
}
