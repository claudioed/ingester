package sender

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/uuid"
	"go.uber.org/zap"
	ingester_v1 "ingester/pkg/pb/analytics"
)

type KafkaAnalyticsSender struct {
	c      client.Client
	logger *zap.Logger
}

func (kas *KafkaAnalyticsSender) Send(ctx context.Context, call *ingester_v1.ApiCall) error {
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
	if result := kas.c.Send(
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(call.TenantId)),
		evt,
	); cloudevents.IsUndelivered(result) {
		kas.logger.Error("failed to send: %v", zap.Bool("evt.sent", cloudevents.IsACK(result)))
	} else {
		kas.logger.Debug("sent: message", zap.String("evt.id", evt.ID()), zap.Bool("evt.sent", cloudevents.IsACK(result)))
	}

	return nil
}

func NewKafkaAnalyticsSender(c cloudevents.Client, logger *zap.Logger) *KafkaAnalyticsSender {
	return &KafkaAnalyticsSender{c: c, logger: logger}
}
