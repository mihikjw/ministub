package logger

import (
	"reflect"
	"testing"
)

// TestNewLogger1 tests the correct logger is created from the given args
func TestNewLogger1(t *testing.T) {
	supportedLogTypes := []string{
		"std",
	}

	for _, mode := range supportedLogTypes {
		anyValidType := false
		result := NewLogger(mode)

		if _, valid := result.(*StdLogger); valid {
			anyValidType = true
			break
		}

		if !anyValidType {
			t.Errorf("%s Is Not Supported Logger Mode", mode)
		}
	}

}

func TestNewLogger2(t *testing.T) {
	result := NewLogger("")

	if result != nil {
		t.Errorf("Unsupported Logger Mode Did Not Return 'nil', Result Type: %s", reflect.TypeOf(result))
	}
}
