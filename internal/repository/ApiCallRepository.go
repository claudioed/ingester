package repository

import (
	"context"
	ingester_v1 "ingester/pkg/pb/analytics"
)

type ApiCallRepo interface {
	Add(ctx context.Context, call *ingester_v1.ApiCall) (*string, error)
}
