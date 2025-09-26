package internal

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/meq/task"
	"github.com/rotisserie/eris"
)

type AsynqTaskQueue[T task.Message] struct {
	logger    ezutil.Logger
	client    *asynq.Client
	inspector *asynq.Inspector
}

func NewAsynqTaskQueue[T task.Message](logger ezutil.Logger, db *AsynqDB) *AsynqTaskQueue[T] {
	if logger == nil {
		panic("logger cannot be nil")
	}

	tq := &AsynqTaskQueue[T]{logger: logger}

	if db != nil {
		tq.client = db.Client
		tq.inspector = db.Inspector
	}

	return tq
}

func (tq *AsynqTaskQueue[T]) Enqueue(ctx context.Context, source string, message T) error {
	t := task.NewTask(source, message)

	payload, err := json.Marshal(t)
	if err != nil {
		return eris.Wrap(err, "error marshaling task")
	}

	asynqTask := asynq.NewTask(t.Type, payload)

	info, err := tq.client.EnqueueContext(ctx, asynqTask, asynq.Queue(t.Type))
	if err != nil {
		return eris.Wrap(err, "error enqueuing task")
	}

	tq.logger.Infof("enqueued task: ID=%s, Queue=%s", info.ID, info.Queue)

	return nil
}

func (tq *AsynqTaskQueue[T]) GetAllPending(ctx context.Context) ([]task.Task[T], error) {
	pendingTasks, err := tq.inspector.ListPendingTasks(tq.taskType(), asynq.PageSize(1000))
	if err != nil {
		return nil, eris.Wrap(err, "error listing pending tasks")
	}
	return ezutil.MapSliceWithError(pendingTasks, tq.mapToTask)
}

func (tq *AsynqTaskQueue[T]) DeleteAll(ctx context.Context) error {
	if err := tq.inspector.DeleteQueue(tq.taskType(), true); err != nil {
		return eris.Wrap(err, "error deleting queue")
	}
	return nil
}

func (tq *AsynqTaskQueue[T]) taskType() string {
	var msg T
	return msg.Type()
}

func (tq *AsynqTaskQueue[T]) mapToTask(taskInfo *asynq.TaskInfo) (task.Task[T], error) {
	if taskInfo == nil {
		return task.Task[T]{}, eris.New("task info is nil")
	}

	var payload task.Task[T]
	if err := json.Unmarshal(taskInfo.Payload, &payload); err != nil {
		return task.Task[T]{}, eris.Wrap(err, "error unmarshal to task")
	}

	return payload, nil
}
