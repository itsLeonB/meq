package meq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/meq/internal"
	"github.com/itsLeonB/meq/task"
)

type TaskQueue[T task.Message] interface {
	Enqueue(ctx context.Context, source string, message T) error
	GetAllPending(ctx context.Context) ([]task.Task[T], error)
	DeleteAll(ctx context.Context) error
}

func NewAsynqTaskQueue[T task.Message](
	logger ezutil.Logger,
	client *asynq.Client,
	inspector *asynq.Inspector,
) TaskQueue[T] {
	return internal.NewAsynqTaskQueue[T](
		logger,
		client,
		inspector,
	)
}
