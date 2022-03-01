package svc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"ingester/internal/repository"
	"ingester/internal/sender"
	ingester_v1 "ingester/pkg/pb/analytics"
)

type CollectorService struct {
	repo   repository.ApiCallRepo
	sender sender.AnalyticsSender
	logger *zap.Logger
}

func (cs *CollectorService) Compute(ctx context.Context, req *ingester_v1.ApiCall) (res *ingester_v1.DataCollected, err error) {
	cs.logger.Debug("Receiving Data...")
	id, err := cs.repo.Add(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("insert: request: %w", err)
	}

	if err := cs.sender.Send(context.Background(), req); err != nil {
		return nil, fmt.Errorf("kafka: request: %w", err)
	}
	res = &ingester_v1.DataCollected{Uuid: *id}
	cs.logger.Debug("Data processed")
	return res, nil
}

func NewCollectorService(repo repository.ApiCallRepo, logger *zap.Logger, analyticsSender sender.AnalyticsSender) *CollectorService {
	return &CollectorService{repo: repo, logger: logger, sender: analyticsSender}
}
