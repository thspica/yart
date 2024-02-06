package executor

import (
	"strings"
	"testing"
)

const (
	successCommand = "echo 'success'"
	failCommand    = "false"
)

func TestCommandExecutor_Execute_Success(t *testing.T) {
	executor := CommandExecutor{
		Command:      "echo 'Hello World'",
		CheckCommand: successCommand,
		Timeout:      2,
		CheckRetry:   3,
	}

	err := executor.Execute()
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestCommandExecutor_Execute_CheckCommandFailure(t *testing.T) {
	executor := CommandExecutor{
		Command:      "echo 'Hello World'",
		CheckCommand: failCommand,
		Timeout:      2,
		CheckRetry:   3,
	}

	err := executor.Execute()
	if err == nil {
		t.Errorf("expected an error, but got nil")
	} else if !strings.Contains(err.Error(), "check command failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCommandExecutor_Execute_Timeout(t *testing.T) {
	executor := CommandExecutor{
		Command:      "sleep 5",
		CheckCommand: successCommand,
		Timeout:      1,
		CheckRetry:   1,
	}

	err := executor.Execute()
	if err == nil {
		t.Errorf("expected a timeout error, but got nil")
	} else if !strings.Contains(err.Error(), "signal: killed") {
		t.Errorf("unexpected error message: %v", err)
	}
}
