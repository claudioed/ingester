package svc

import (
	"context"
	"fmt"
	"ingester/internal/repository"
	ingester_v1 "ingester/pkg/pb/analytics"
)

type CollectorService struct {
	repo repository.ApiCallRepo
}

func (cs *CollectorService) Compute(ctx context.Context, req *ingester_v1.ApiCall) (res *ingester_v1.DataCollected, err error) {
	id, err := cs.repo.Add(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("insert: request: %w", err)
	}
	res = &ingester_v1.DataCollected{Uuid: *id}
	return res, nil
}

func NewCollectorService(repo repository.ApiCallRepo) *CollectorService {
	return &CollectorService{repo: repo}
}
