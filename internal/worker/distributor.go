package worker

import (
	"context"

	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/proto"
	"auditor.z9fr.xyz/server/internal/redis"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskAuditStartRequest(
		ctx context.Context,
		payload *proto.AuditStartRequest,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
	logger lib.Logger
}

func NewRedisTaskDistributor(redis redis.RedisConnection, logger lib.Logger) TaskDistributor {
	client := asynq.NewClient(redis.RedisClientOpt)
	return &RedisTaskDistributor{
		client: client,
		logger: logger,
	}
}
