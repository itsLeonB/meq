package task

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestMessage struct {
	Content string `json:"content"`
}

func (tm TestMessage) Type() string {
	return "test_message"
}

func TestTask_IsZero(t *testing.T) {
	tests := []struct {
		name string
		task Task[TestMessage]
		want bool
	}{
		{
			name: "zero task",
			task: Task[TestMessage]{},
			want: true,
		},
		{
			name: "non-zero task",
			task: Task[TestMessage]{
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
	message := TestMessage{Content: "test"}

	result := NewTask(source, message)

	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.WithinDuration(t, time.Now(), result.Timestamp, time.Second)
	assert.Equal(t, source, result.Source)
	assert.Equal(t, "test_message", result.Type)
	assert.Equal(t, message, result.Message)
}

func TestTestMessage_Type(t *testing.T) {
	msg := TestMessage{Content: "test"}
	assert.Equal(t, "test_message", msg.Type())
}
