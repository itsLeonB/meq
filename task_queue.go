package meq

import (
	"context"
	"fmt"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/meq/internal"
	"github.com/itsLeonB/meq/task"
)

type TaskQueue[T task.Message] interface {
	Enqueue(ctx context.Context, source string, message T) error
	GetAllPending(ctx context.Context) ([]task.Task[T], error)
	DeleteAll(ctx context.Context) error
	GetOldest(ctx context.Context) (task.Task[T], string, error)
	Delete(ctx context.Context, id string) error
}

func NewTaskQueue[T task.Message](logger ezutil.Logger, db DB) TaskQueue[T] {
	switch d := db.(type) {
	case *internal.AsynqDB:
		return internal.NewAsynqTaskQueue[T](logger, d)
	default:
		panic(fmt.Sprintf("unsupported db type: %T", d))
	}
}
