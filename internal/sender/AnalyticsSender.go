package sender

import (
	"context"
	ingester_v1 "ingester/pkg/pb/analytics"
)

type AnalyticsSender interface {
	Send(ctx context.Context, call *ingester_v1.ApiCall) error
}
