package executor

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"yart/internal/config"
)

type CommandExecutor struct {
	Command      string
	CheckCommand string
	Timeout      int
	CheckRetry   int
}

func (c *CommandExecutor) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", c.Command)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if c.CheckCommand != "" {
		backoff := config.InitialBackoff
		for i := 0; i < c.CheckRetry; i++ {
			checkCtx, checkCancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
			defer checkCancel()

			checkCmd := exec.CommandContext(checkCtx, "sh", "-c", c.CheckCommand)
			output, checkErr := checkCmd.CombinedOutput()
			if checkErr == nil && strings.TrimSpace(string(output)) != "" {
				return nil
			}

			time.Sleep(time.Duration(backoff) * time.Second)

			backoff = backoff * 2
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}

		return fmt.Errorf("check command failed after %d retries", c.CheckRetry)
	}

	return nil
}
