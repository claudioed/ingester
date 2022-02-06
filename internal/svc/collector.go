package svc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"ingester/internal/repository"
	ingester_v1 "ingester/pkg/pb/analytics"
)

type CollectorService struct {
	repo   repository.ApiCallRepo
	logger *zap.Logger
}

func (cs *CollectorService) Compute(ctx context.Context, req *ingester_v1.ApiCall) (res *ingester_v1.DataCollected, err error) {
	cs.logger.Info("Receiving Data...")
	id, err := cs.repo.Add(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("insert: request: %w", err)
	}
	res = &ingester_v1.DataCollected{Uuid: *id}
	cs.logger.Info("Data processed")
	return res, nil
}

func NewCollectorService(repo repository.ApiCallRepo, logger *zap.Logger) *CollectorService {
	return &CollectorService{repo: repo, logger: logger}
}
