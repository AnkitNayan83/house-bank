package workers

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOptions *asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOptions)
	return &RedisTaskDistributor{
		client: client,
	}
}
