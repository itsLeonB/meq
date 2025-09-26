package meq

import (
	"github.com/hibiken/asynq"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/meq/internal"
)

type DB interface {
	Ping() error
	Shutdown() error
}

func NewAsynqDB(logger ezutil.Logger, opts asynq.RedisConnOpt) DB {
	return internal.NewAsynqDB(logger, opts)
}
