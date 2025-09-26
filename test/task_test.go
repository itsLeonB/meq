package meq_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/meq/task"
	"github.com/itsLeonB/meq/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTask_IsZero(t *testing.T) {
	tests := []struct {
		name string
		task task.Task[testutil.TestMessage]
		want bool
	}{
		{
			name: "zero task",
			task: task.Task[testutil.TestMessage]{},
			want: true,
		},
		{
			name: "non-zero task",
			task: task.Task[testutil.TestMessage]{
				ID: uuid.New(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.task.IsZero())
		})
	}
}

func TestNewTask(t *testing.T) {
	source := "test-source"
	message := testutil.TestMessage{Content: "test"}

	result := task.NewTask(source, message)

	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.WithinDuration(t, time.Now(), result.Timestamp, time.Second)
	assert.Equal(t, source, result.Source)
	assert.Equal(t, "test_message", result.Type)
	assert.Equal(t, message, result.Message)
}

func TestTestMessage_Type(t *testing.T) {
	msg := testutil.TestMessage{Content: "test"}
	assert.Equal(t, "test_message", msg.Type())
}
