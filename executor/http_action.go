package executor

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpAction struct {
	URL            string
	ExpectedStatus int
	ExpectedBody   string
	Timeout        int
	PreCommand     ActionExecutor
	PostCommand    ActionExecutor
}

func (h *HttpAction) Execute() error {
	if h.PreCommand != nil {
		h.PreCommand.Execute()
	}

	client := &http.Client{
		Timeout: time.Duration(h.Timeout) * time.Second,
	}
	resp, err := client.Get(h.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != h.ExpectedStatus {
		return fmt.Errorf("expected status %d but got %d", h.ExpectedStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(body), h.ExpectedBody) {
		return fmt.Errorf("expected body to contain '%s' but got '%s'", h.ExpectedBody, string(body))
	}

	if h.PostCommand != nil {
		h.PostCommand.Execute()
	}

	return nil
}
