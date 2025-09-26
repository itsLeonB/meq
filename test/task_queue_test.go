package meq_test

import (
	"testing"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/meq"
	"github.com/itsLeonB/meq/internal"
	"github.com/itsLeonB/meq/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskQueue(t *testing.T) {
	tests := []struct {
		name   string
		db     meq.DB
		panics bool
	}{
		{
			name: "with AsynqDB",
			db:   internal.NewAsynqDB(&testutil.MockLogger{}, asynq.RedisClientOpt{Addr: "localhost:6379"}),
		},
		{
			name:   "with unsupported DB type",
			db:     &mockDB{},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &testutil.MockLogger{}

			if tt.panics {
				assert.Panics(t, func() {
					meq.NewTaskQueue[testutil.TestMessage](logger, tt.db)
				})
			} else {
				tq := meq.NewTaskQueue[testutil.TestMessage](logger, tt.db)
				assert.NotNil(t, tq)
			}
		})
	}
}

type mockDB struct{}

func (m *mockDB) Ping() error   { return nil }
func (m *mockDB) Shutdown() error { return nil }
