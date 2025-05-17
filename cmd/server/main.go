package main

import (
	"os"
	"log"
	"net"
	"net/http"

	"github.com/anubhav-mathur/distributed-rate-limiter/internal/api"
	"github.com/anubhav-mathur/distributed-rate-limiter/internal/metrics"
	pb "github.com/anubhav-mathur/distributed-rate-limiter/proto"

	"google.golang.org/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "2112"
	}

	go func() {
		metrics.Init()
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics available at :%s/metrics", metricsPort)
		log.Fatal(http.ListenAndServe(":"+metricsPort, nil))
	}()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRateLimiterServer(grpcServer, api.NewRateLimiterServer())

	log.Println("gRPC Rate Limiter server is running on port %s...\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
