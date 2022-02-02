package svc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
)

type HealthService struct {
}

func (h *HealthService) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Printf("✅ Server's status is %s", grpc_health_v1.HealthCheckResponse_SERVING)
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthService) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return nil
}

func NewHealthService() *HealthService {
	return &HealthService{}
}
