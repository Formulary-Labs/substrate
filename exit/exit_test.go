package exit_test

import (
	"testing"

	"github.com/Formulary-Labs/substrate/exit"
)

func TestConstants(t *testing.T) {
	if exit.OK != 0 {
		t.Errorf("exit.OK must be 0, got %d", exit.OK)
	}
	if exit.Validation != 1 {
		t.Errorf("exit.Validation must be 1, got %d", exit.Validation)
	}
	if exit.ToolError != 2 {
		t.Errorf("exit.ToolError must be 2, got %d", exit.ToolError)
	}
}

func TestIfErr_nil(t *testing.T) {
	// Should not panic when err is nil.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("IfErr with nil panicked: %v", r)
		}
	}()
	exit.IfErr(nil)
}
