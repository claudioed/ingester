package sender

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/uuid"
	ingester_v1 "ingester/pkg/pb/analytics"
	"log"
)

type KafkaAnalyticsSender struct {
	c client.Client
}

func (receiver *KafkaAnalyticsSender) Send(ctx context.Context, call *ingester_v1.ApiCall) error {
	evt := cloudevents.NewEvent()
	evt.SetID(uuid.New().String())
	evt.SetType("tech.claudioed.analytics.http.api")
	evt.SetSource("https://tech.claudioed/ingester")

	m := jsonpb.Marshaler{}
	js, err := m.MarshalToString(call)
	if err != nil {
		return err
	}
	_ = evt.SetData(cloudevents.ApplicationJSON, js)
	if result := receiver.c.Send(
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(call.TenantId)),
		evt,
	); cloudevents.IsUndelivered(result) {
		log.Printf("failed to send: %v", result)
	} else {
		log.Printf("sent: %s, accepted: %t", evt.ID(), cloudevents.IsACK(result))
	}

	return nil
}

func NewKafkaAnalyticsSender(c cloudevents.Client) *KafkaAnalyticsSender {
	return &KafkaAnalyticsSender{c: c}
}
