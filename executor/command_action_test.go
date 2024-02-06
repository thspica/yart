package executor

import (
	"testing"
)

func TestCommandAction_Execute_Success(t *testing.T) {
	action := CommandAction{
		Command:        "echo 'Hello World'",
		ExpectedOutput: "Hello World",
		Timeout:        2,
	}

	err := action.Execute()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
func TestCommandAction_Execute_Failure(t *testing.T) {
	action := CommandAction{
		Command:        "echo 'Goodbye World'",
		ExpectedOutput: "Hello World",
		Timeout:        2,
	}

	err := action.Execute()
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestCommandAction_Execute_Timeout(t *testing.T) {
	action := CommandAction{
		Command:        "sleep 5",
		ExpectedOutput: "",
		Timeout:        1,
	}

	err := action.Execute()
	if err == nil {
		t.Errorf("expected a timeout error, got nil")
	} else if err.Error() != "signal: killed" {
		t.Errorf("expected context deadline exceeded error, got %v", err)
	}
}
