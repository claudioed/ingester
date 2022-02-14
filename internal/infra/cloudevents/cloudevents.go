package cloudevents

import (
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func NewCloudEventClientWithKafka(sender *kafka_sarama.Sender) (cloudevents.Client, error) {
	c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}
	return c, nil
}
