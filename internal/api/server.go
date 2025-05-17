package api

import (
	"context"
	"github.com/anubhav-mathur/distributed-rate-limiter/internal/store"
	"github.com/anubhav-mathur/distributed-rate-limiter/internal/metrics"
	pb "github.com/anubhav-mathur/distributed-rate-limiter/proto"
)

type RateLimiterServer struct {
	pb.UnimplementedRateLimiterServer
	limiter *store.RedisLimiter
}

func NewRateLimiterServer() *RateLimiterServer {
	limiter := store.NewRedisLimiter("localhost:6379", 5, 5)
	return &RateLimiterServer{limiter: limiter}
}

// func (s *RateLimiterServer) getBucket(userID string) *limiter.TokenBucket {
// 	if _, exists := s.bucketMap[userID]; !exists {
// 		// 5 tokens total, refill 1 token every 2 seconds = 5 per 10 seconds
// 		s.bucketMap[userID] = limiter.NewTokenBucket(5, 5)
// 	}
// 	return s.bucketMap[userID]
// }

func (s *RateLimiterServer) AllowRequest(ctx context.Context, req *pb.RequestInput) (*pb.RequestOutput, error) {
	ok, err := s.limiter.Allow(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if ok {
		metrics.RequestsTotal.WithLabelValues(req.UserId, "allowed").Inc()
		return &pb.RequestOutput{Allowed: true, Reason: "Request allowed"}, nil
	}
	metrics.RequestsTotal.WithLabelValues(req.UserId, "denied").Inc()
	return &pb.RequestOutput{Allowed: false, Reason: "Rate limit exceeded"}, nil
}

func (s *RateLimiterServer) GetUsage(ctx context.Context, req *pb.UsageInput) (*pb.UsageOutput, error) {
	used, allowed, err := s.limiter.Usage(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.UsageOutput{
		RequestsUsed:    int32(used),
		RequestsAllowed: int32(allowed),
	}, nil
}

