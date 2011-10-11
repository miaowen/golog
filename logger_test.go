package golog

import (
	"gomock.googlecode.com/hg/gomock"
	"testing"
)

func TestOutput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	message := &LogMessage{}

	mockLogOuter := NewMockLogOuter(mockCtrl)
	mockLogOuter.EXPECT().Output(message)

	logger := NewLogger(mockLogOuter, 0, nil)
	logger.Log(0, func() *LogMessage { return message })
}

func TestNoOutput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogOuter := NewMockLogOuter(mockCtrl)

	var called bool = false

	logger := NewLogger(mockLogOuter, 1, nil)
	// The message is logged below the filter level.
	logger.Log(0, func() *LogMessage { called = true; return nil })

	if called {
		t.Error("Closure evaluated even though no output produced")
	}
}