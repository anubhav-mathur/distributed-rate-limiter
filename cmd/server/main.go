package main

import (
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
	go func() {
		metrics.Init()
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics available at :2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRateLimiterServer(grpcServer, api.NewRateLimiterServer())

	log.Println("gRPC Rate Limiter server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
