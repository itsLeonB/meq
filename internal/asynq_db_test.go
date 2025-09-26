package internal

import (
	"testing"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/meq/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewAsynqDB(t *testing.T) {
	tests := []struct {
		name      string
		opts      asynq.RedisConnOpt
		expectNil bool
	}{
		{
			name:      "with valid opts",
			opts:      asynq.RedisClientOpt{Addr: "localhost:6379"},
			expectNil: false,
		},
		{
			name:      "with nil opts",
			opts:      nil,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &testutil.MockLogger{}
			if tt.opts == nil {
				logger.On("Warn", "opts is nil, will not initialize asynq db")
			}

			result := NewAsynqDB(logger, tt.opts)

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Client)
				assert.NotNil(t, result.Inspector)
			}

			logger.AssertExpectations(t)
		})
	}
}

func TestAsynqDB_Ping(t *testing.T) {
	// This test requires actual asynq client/inspector instances
	// Testing with real Redis connection would be integration test
	t.Skip("Ping method requires real asynq client/inspector instances")
}

func TestAsynqDB_Shutdown(t *testing.T) {
	// This test requires actual asynq client/inspector instances
	// Testing with real Redis connection would be integration test
	t.Skip("Shutdown method requires real asynq client/inspector instances")
}
