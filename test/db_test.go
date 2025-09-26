package meq_test

import (
	"testing"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/meq"
	"github.com/itsLeonB/meq/internal"
	"github.com/itsLeonB/meq/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewAsynqDB(t *testing.T) {
	logger := &testutil.MockLogger{}
	opts := asynq.RedisClientOpt{Addr: "localhost:6379"}

	db := meq.NewAsynqDB(logger, opts)

	assert.NotNil(t, db)
	assert.IsType(t, &internal.AsynqDB{}, db)
}
