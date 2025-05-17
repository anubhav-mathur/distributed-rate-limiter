package api

import (
	"context"
	"github.com/anubhav-mathur/distributed-rate-limiter/internal/limiter"
	pb "github.com/anubhav-mathur/distributed-rate-limiter/proto"
)

type RateLimiterServer struct {
	pb.UnimplementedRateLimiterServer
	bucketMap map[string]*limiter.TokenBucket
}

func NewRateLimiterServer() *RateLimiterServer {
	return &RateLimiterServer{
		bucketMap: make(map[string]*limiter.TokenBucket),
	}
}

func (s *RateLimiterServer) getBucket(userID string) *limiter.TokenBucket {
	if _, exists := s.bucketMap[userID]; !exists {
		// 5 tokens total, refill 1 token every 2 seconds = 5 per 10 seconds
		s.bucketMap[userID] = limiter.NewTokenBucket(5, 5)
	}
	return s.bucketMap[userID]
}

func (s *RateLimiterServer) AllowRequest(ctx context.Context, req *pb.RequestInput) (*pb.RequestOutput, error) {
	bucket := s.getBucket(req.UserId)

	if bucket.Allow() {
		return &pb.RequestOutput{Allowed: true, Reason: "Request allowed"}, nil
	}
	return &pb.RequestOutput{Allowed: false, Reason: "Rate limit exceeded"}, nil
}

func (s *RateLimiterServer) GetUsage(ctx context.Context, req *pb.UsageInput) (*pb.UsageOutput, error) {
	bucket := s.getBucket(req.UserId)
	used, allowed := bucket.Usage()
	return &pb.UsageOutput{
		RequestsUsed:    int32(used),
		RequestsAllowed: int32(allowed),
	}, nil
}
