package internal_test

import (
	"testing"

	"github.com/itsLeonB/meq/internal"
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
	
	internal.NewAsynqTaskQueue[testutil.TestMessage](nil, &internal.AsynqDB{})
}

func TestNewAsynqTaskQueue_ValidCases(t *testing.T) {
	tests := []struct {
		name   string
		logger *testutil.MockLogger
		db     *internal.AsynqDB
	}{
		{
			name:   "with valid logger and db",
			logger: &testutil.MockLogger{},
			db:     &internal.AsynqDB{},
		},
		{
			name:   "with valid logger and nil db",
			logger: &testutil.MockLogger{},
			db:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tq := internal.NewAsynqTaskQueue[testutil.TestMessage](tt.logger, tt.db)
			assert.NotNil(t, tq)
		})
	}
}
