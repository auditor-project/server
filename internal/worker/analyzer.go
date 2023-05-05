package worker

import (
	"context"
	"encoding/json"

	"auditor.z9fr.xyz/server/internal/proto"
	"github.com/hibiken/asynq"
)

const (
	FAILED_TO_ENQUEUE = "failed to enqueue task"
)

func (distributor *RedisTaskDistributor) DistributeTaskAuditStartRequest(
	ctx context.Context,
	payload *proto.AuditStartRequest,
	opts ...asynq.Option,
) error {
	var info *asynq.TaskInfo
	requestId := ctx.Value("requestId").(string)
	option := asynq.TaskID(requestId)
	queue := asynq.Queue(QueueAnalysis)

	opts = append(opts, option)
	opts = append(opts, queue)

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		distributor.logger.Errorf(err.Error(), "error", "failed to marshal task payload", "taskId", requestId)
		return err
	}

	task := asynq.NewTask(TaskStartAnalyser, jsonPayload, opts...)

	info, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		distributor.logger.Errorw(err.Error(), "error", FAILED_TO_ENQUEUE, "taskId", requestId)
		return err
	}

	if err != nil {
		distributor.logger.Errorw(err.Error(), "error", "failed to enqueue task", "taskId", requestId)
		return err
	}

	distributor.logger.Infow("task enque success", "info", info, "taskId", requestId)
	return nil
}

func (processor *RedisTaskProcessor) ProcessStartAuditorAnalysis(ctx context.Context, task *asynq.Task) error {
	var payload *proto.AuditStartRequest
	requestId, _ := asynq.GetTaskID(ctx)
	processor.logger.Infow("starting to process task", "taskid", requestId)

	if err := json.Unmarshal([]byte(task.Payload()), &payload); err != nil {
		processor.logger.Errorw("failed to unmarshal task payload", "error", err, "taskId", requestId)
		return err
	}

	ok, err := processor.analyzer.InitiateAnalyzer()

	if err != nil {
		processor.logger.Errorw("failed to send email", "data", payload, "error", err, "taskId", requestId)
		return err
	}

	processor.logger.Infow("task completed", "response", ok, "taskId", requestId)
	return nil
}
