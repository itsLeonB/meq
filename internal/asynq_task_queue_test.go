package internal

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/itsLeonB/meq/task"
	"github.com/itsLeonB/meq/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewAsynqTaskQueue_NilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "logger cannot be nil", r)
		} else {
			t.Error("Expected panic but none occurred")
		}
	}()

	NewAsynqTaskQueue[testutil.TestMessage](nil, &AsynqDB{})
}

func TestNewAsynqTaskQueue_ValidCases(t *testing.T) {
	tests := []struct {
		name   string
		logger *testutil.MockLogger
		db     *AsynqDB
	}{
		{
			name:   "with valid logger and db",
			logger: &testutil.MockLogger{},
			db:     &AsynqDB{},
		},
		{
			name:   "with valid logger and nil db",
			logger: &testutil.MockLogger{},
			db:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tq := NewAsynqTaskQueue[testutil.TestMessage](tt.logger, tt.db)
			assert.NotNil(t, tq)
		})
	}
}

func TestAsynqTaskQueue_Enqueue(t *testing.T) {
	// This test requires actual asynq client
	// Testing with real Redis connection would be integration test
	t.Skip("Enqueue method requires real asynq client")
}

func TestAsynqTaskQueue_GetAllPending(t *testing.T) {
	// This test requires actual asynq inspector
	// Testing with real Redis connection would be integration test
	t.Skip("GetAllPending method requires real asynq inspector")
}

func TestAsynqTaskQueue_DeleteAll(t *testing.T) {
	// This test requires actual asynq inspector
	// Testing with real Redis connection would be integration test
	t.Skip("DeleteAll method requires real asynq inspector")
}

func TestAsynqTaskQueue_taskType(t *testing.T) {
	logger := &testutil.MockLogger{}
	tq := NewAsynqTaskQueue[testutil.TestMessage](logger, nil)
	
	result := tq.taskType()
	assert.Equal(t, "test_message", result)
}

func TestAsynqTaskQueue_mapToTask(t *testing.T) {
	logger := &testutil.MockLogger{}
	tq := NewAsynqTaskQueue[testutil.TestMessage](logger, nil)

	tests := []struct {
		name     string
		taskInfo *asynq.TaskInfo
		wantErr  bool
	}{
		{
			name:     "nil task info",
			taskInfo: nil,
			wantErr:  true,
		},
		{
			name: "invalid payload",
			taskInfo: &asynq.TaskInfo{
				ID:      "test-id",
				Payload: []byte("invalid json"),
			},
			wantErr: true,
		},
		{
			name: "valid payload",
			taskInfo: &asynq.TaskInfo{
				ID:      "test-id",
				Payload: createValidTaskPayload(t),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tq.mapToTask(tt.taskInfo)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, result.IsZero())
			} else {
				assert.NoError(t, err)
				assert.False(t, result.IsZero())
			}
		})
	}
}

func createValidTaskPayload(t *testing.T) []byte {
	testTask := task.Task[testutil.TestMessage]{
		ID:      uuid.New(),
		Source:  "test-source",
		Type:    "test_message",
		Message: testutil.TestMessage{Content: "test"},
	}
	payload, err := json.Marshal(testTask)
	assert.NoError(t, err)
	return payload
}
