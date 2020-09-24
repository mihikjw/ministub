package logger

import (
	"bytes"
	"testing"
)

// TestNewStdLogger1 ensures the constructor sets up all fields correctly
func TestNewStdLogger1(t *testing.T) {
	result := NewStdLogger()

	if result.stdErr == nil {
		t.Error("StdLogger.stdErr Is Nil On Create")
	}
	if result.stdOut == nil {
		t.Error("StdLogger.stdOut Is Nil On Create")
	}
}

// TestInfo1 ensures the given log message is output correctly - if the *log.Logger objects are internal are nil, we want to fail so don't test these
func TestInfo1(t *testing.T) {
	client := new(StdLogger)
	outBuffer := new(bytes.Buffer)
	client.stdOut = createOutput(outBuffer, "INFO: ")

	client.Info("Test Message")
	result := outBuffer.Bytes()

	if result != nil {
		if len(outBuffer.Bytes()) != 39 {
			t.Errorf("Info Output Is Unexpected Size: %d", len(result))
		}
	} else {
		t.Errorf("outBuffer.Bytes() Is Nil")
	}
}

// TestError1 ensures the given error log message is output correctly - if the *log.Logger objects are internal are nil, we want to fail so don't test these
func TestError1(t *testing.T) {
	client := new(StdLogger)
	outBuffer := new(bytes.Buffer)
	client.stdErr = createOutput(outBuffer, "ERROR: ")

	client.Error("Test Error Message")
	result := outBuffer.Bytes()

	if result != nil {
		if len(result) != 46 {
			t.Errorf("Error Output Is Unexpected Size: %d", len(result))
		}
	} else {
		t.Errorf("outBuffer.Bytes() Is Nil")
	}
}
