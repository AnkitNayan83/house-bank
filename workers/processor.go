package workers

import (
	"context"

	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	QueueueCritical = "critical"
	QueueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskPorcessor(redisOptions *asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOptions, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			QueueueCritical: 10,
			QueueueDefault:  5,
		},
	})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessSendVerifyEmail)

	return processor.server.Start(mux)
}
