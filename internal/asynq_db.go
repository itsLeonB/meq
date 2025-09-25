package internal

import (
	"errors"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/ezutil/v2"
)

type AsynqDB struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewAsynqDB(logger ezutil.Logger, opts asynq.RedisConnOpt) *AsynqDB {
	if opts == nil {
		logger.Warn("opts is nil, will not initialize asynq db")
		return nil
	}
	return &AsynqDB{
		Client:    asynq.NewClient(opts),
		Inspector: asynq.NewInspector(opts),
	}
}

func (d *AsynqDB) Ping() error {
	var errs error

	if err := d.Client.Ping(); err != nil {
		errs = errors.Join(errs, err)
	}
	if _, err := d.Inspector.Queues(); err != nil {
		errs = errors.Join(errs, err)
	}

	return errs
}

func (d *AsynqDB) Shutdown() error {
	var errs error

	if err := d.Client.Close(); err != nil {
		errs = errors.Join(errs, err)
	}
	if err := d.Inspector.Close(); err != nil {
		errs = errors.Join(errs, err)
	}

	return errs
}
