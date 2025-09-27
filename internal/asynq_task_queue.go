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
	queueName string
}

func NewAsynqTaskQueue[T task.Message](logger ezutil.Logger, db *AsynqDB) *AsynqTaskQueue[T] {
	if logger == nil {
		panic("logger cannot be nil")
	}

	var msg T

	tq := &AsynqTaskQueue[T]{
		logger:    logger,
		queueName: msg.Type(),
	}

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
	pendingTasks, err := tq.inspector.ListPendingTasks(tq.queueName, asynq.PageSize(1000))
	if err != nil {
		return nil, eris.Wrap(err, "error listing pending tasks")
	}
	return ezutil.MapSliceWithError(pendingTasks, tq.mapToTask)
}

func (tq *AsynqTaskQueue[T]) DeleteAll(ctx context.Context) error {
	if err := tq.inspector.DeleteQueue(tq.queueName, true); err != nil {
		return eris.Wrap(err, "error deleting queue")
	}
	return nil
}

func (tq *AsynqTaskQueue[T]) GetOldest(ctx context.Context) (task.Task[T], string, error) {
	tasks, err := tq.inspector.ListPendingTasks(tq.queueName, asynq.PageSize(1))
	if err != nil {
		return task.Task[T]{}, "", eris.Wrap(err, "error listing pending tasks")
	}

	if len(tasks) < 1 {
		return task.Task[T]{}, "", nil
	}

	msg, err := tq.mapToTask(tasks[0])
	if err != nil {
		return task.Task[T]{}, "", err
	}

	return msg, tasks[0].ID, nil
}

func (tq *AsynqTaskQueue[T]) Delete(ctx context.Context, id string) error {
	if err := tq.inspector.DeleteTask(tq.queueName, id); err != nil {
		return eris.Wrapf(err, "error deleting task id %s", id)
	}
	return nil
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
