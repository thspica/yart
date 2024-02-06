package executor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	expectedBody = "Hello, World!"
	invalidBody  = "Goodbye, World!"
	notFoundBody = "Not Found"
)

func TestHttpAction_Execute_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	}))
	defer server.Close()

	action := HttpAction{
		URL:            server.URL,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   expectedBody,
		Timeout:        2,
	}

	err := action.Execute()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHttpAction_Execute_StatusCodeMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(notFoundBody))
	}))
	defer server.Close()

	action := HttpAction{
		URL:            server.URL,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   expectedBody,
		Timeout:        2,
	}

	err := action.Execute()
	if err == nil {
		t.Errorf("expected an error, got nil")
	} else if err.Error() != "expected status 200 but got 404" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHttpAction_Execute_BodyMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(invalidBody))
	}))
	defer server.Close()

	action := HttpAction{
		URL:            server.URL,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   expectedBody,
		Timeout:        2,
	}

	err := action.Execute()
	if err == nil {
		t.Errorf("expected an error, got nil")
	} else if !strings.Contains(err.Error(), "expected body to contain 'Hello, World!'") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHttpAction_Execute_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	}))
	defer server.Close()

	action := HttpAction{
		URL:            server.URL,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   expectedBody,
		Timeout:        1,
	}

	err := action.Execute()
	if err == nil {
		t.Errorf("expected a timeout error, got nil")
	} else if !strings.Contains(err.Error(), "Client.Timeout exceeded") {
		t.Errorf("unexpected error message: %v", err)
	}
}
