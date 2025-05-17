package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/anubhav-mathur/distributed-rate-limiter/proto"
	"github.com/anubhav-mathur/distributed-rate-limiter/internal/api"
)

func main() {
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
