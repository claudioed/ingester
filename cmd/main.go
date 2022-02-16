package main

import (
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"ingester/internal/infra/cloudevents"
	"ingester/internal/infra/elastic"
	"ingester/internal/infra/kafka"
	"ingester/internal/repository"
	sender2 "ingester/internal/sender"
	"ingester/internal/svc"
	ingester_v1 "ingester/pkg/pb/analytics"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	listener net.Listener
	server   *grpc.Server
	logger   *zap.Logger
)

func main() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
	initListener()

	escli, errES := elastic.NewESClient()
	if errES != nil {
		logger.Panic("Failed to connect ES", zap.Error(errES))
	}

	// Elastic Dependencies
	elastic := elastic.NewElasticSearch(escli)
	errIdx := elastic.CreateIndex("api_calls")
	if errIdx != nil {
		logger.Panic("Failed to create index ES", zap.Error(errES))
	}

	// Kafka + CE Dependencies
	sender, errSen := kafka.NewSender()
	if errSen != nil {
		logger.Panic("Failed to create sender", zap.Error(errES))
	}
	cec, errCe := cloudevents.NewCloudEventClientWithKafka(sender)
	if errCe != nil {
		logger.Panic("Failed to cloud events client", zap.Error(errES))
	}
	ks := sender2.NewKafkaAnalyticsSender(cec, logger)

	repo := repository.NewElasticApiCallRepository(elastic)
	collectorHandler := svc.NewCollectorService(repo, logger, ks)
	healthHandler := svc.NewHealthService()

	// gRPC Server
	server = grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	ingester_v1.RegisterCollectServer(server, collectorHandler)
	grpc_health_v1.RegisterHealthServer(server, healthHandler)

	logger.Info("Handlers registered")
	grpc_prometheus.Register(server)
	http.Handle("/metrics", promhttp.Handler())

	go signalsListener(server)
	logger.Info("Starting gRPC server...")
	if err := server.Serve(listener); err != nil {
		logger.Panic("Failed to start gRPC server", zap.Error(err))
	}
}

func initListener() {
	var err error
	addr := "localhost:50051"
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		logger.Panic("Failed to listen",
			zap.String("address", addr),
			zap.Error(err),
		)
	}
	logger.Info("Started listening...", zap.String("address", addr))
	return
}

func signalsListener(server *grpc.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	_ = <-sigs

	logger.Info("Gracefully stopping server...")
	server.GracefulStop()
}
