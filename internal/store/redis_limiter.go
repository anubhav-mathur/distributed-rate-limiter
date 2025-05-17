package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client      *redis.Client
	capacity    int           // Max tokens
	fillInterval time.Duration // Time between token refills
}

func NewRedisLimiter(addr string, capacity int, refillRate int) *RedisLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisLimiter{
		client:      client,
		capacity:    capacity,
		fillInterval: time.Second * time.Duration(10) / time.Duration(capacity),
	}
}

func (rl *RedisLimiter) Allow(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("rl:%s", userID)
	now := time.Now().Unix()

	// Use Lua script to perform atomic rate limiting
	script := redis.NewScript(`
		local tokens_key = KEYS[1]
		local now = tonumber(ARGV[1])
		local capacity = tonumber(ARGV[2])
		local fill_interval = tonumber(ARGV[3])

		local data = redis.call("HMGET", tokens_key, "tokens", "last_refill")
		local tokens = tonumber(data[1]) or capacity
		local last_refill = tonumber(data[2]) or now

		local elapsed = now - last_refill
		local refill_tokens = math.floor(elapsed / fill_interval)
		if refill_tokens > 0 then
			tokens = math.min(capacity, tokens + refill_tokens)
			last_refill = now
		end

		if tokens > 0 then
			tokens = tokens - 1
			redis.call("HMSET", tokens_key, "tokens", tokens, "last_refill", last_refill)
			redis.call("EXPIRE", tokens_key, 60)
			return 1
		else
			redis.call("HMSET", tokens_key, "tokens", tokens, "last_refill", last_refill)
			redis.call("EXPIRE", tokens_key, 60)
			return 0
		end
	`)

	res, err := script.Run(ctx, rl.client, []string{key},
		now,
		rl.capacity,
		int64(rl.fillInterval.Seconds()),
	).Result()

	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}
