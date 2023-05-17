package worker

import (
	"context"

	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/redis"
	"auditor.z9fr.xyz/server/internal/service"
	"github.com/hibiken/asynq"
)

const (
	QueueAnalysis = "analysis-queue"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessStartAuditorAnalysis(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server   *asynq.Server
	logger   lib.Logger
	analyzer *service.AnalyzerService
}

func NewRedisTaskProcessor(redis redis.RedisConnection, logger lib.Logger, analyzer *service.AnalyzerService) TaskProcessor {
	logger.Debug("Init redis task processor")
	server := asynq.NewServer(
		redis.RedisClientOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueAnalysis: 5,
				QueueDefault:  5,
			},
			Logger: logger,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Errorw(err.Error(), task)
			}),
		},
	)

	return &RedisTaskProcessor{
		server:   server,
		logger:   logger,
		analyzer: analyzer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskStartAnalyser, processor.ProcessStartAuditorAnalysis)
	return processor.server.Start(mux)
}
