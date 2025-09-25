package task

import (
	"time"

	"github.com/google/uuid"
)

type Task[T Message] struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Type      string    `json:"type"`
	Message   T         `json:"message"`
}

func (t Task[T]) IsZero() bool {
	return t.ID == uuid.Nil
}

func NewTask[T Message](source string, message T) Task[T] {
	return Task[T]{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Source:    source,
		Type:      message.Type(),
		Message:   message,
	}
}

type Message interface {
	Type() string
}
