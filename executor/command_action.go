package executor

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CommandAction struct {
	Command        string
	ExpectedOutput string
	Timeout        int
	PreCommand     ActionExecutor
	PostCommand    ActionExecutor
}

func (c *CommandAction) Execute() error {
	if c.PreCommand != nil {
		c.PreCommand.Execute()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", c.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if !strings.Contains(string(output), c.ExpectedOutput) {
		return fmt.Errorf("expected output to contain '%s' but got '%s'", c.ExpectedOutput, string(output))
	}

	if c.PostCommand != nil {
		c.PostCommand.Execute()
	}

	return nil
}
