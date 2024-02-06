package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"yart/executor"
	"yart/internal/config"

	"gopkg.in/yaml.v3"
)

type Scenario struct {
	Name     string   `yaml:"name"`
	Startup  Command  `yaml:"startup"`
	Tests    []Action `yaml:"tests"`
	Shutdown Command  `yaml:"shutdown"`
}

type Command struct {
	Command      string `yaml:"command"`
	CheckCommand string `yaml:"check_command"`
	Timeout      int    `yaml:"timeout"`
	Retry        int    `yaml:"retry"`
}

type Action struct {
	Type           string   `yaml:"type"`
	URL            string   `yaml:"url,omitempty"`
	Command        string   `yaml:"command,omitempty"`
	ExpectedStatus int      `yaml:"expected_status,omitempty"`
	ExpectedOutput string   `yaml:"expected_output,omitempty"`
	Timeout        int      `yaml:"timeout"`
	PreCommand     *Command `yaml:"pre_command"`
	PostCommand    *Command `yaml:"post_command"`
}

func CreateCommand(command Command) (executor.ActionExecutor, error) {
	return &executor.CommandExecutor{
		Command:      command.Command,
		CheckCommand: command.CheckCommand,
		Timeout: func() int {
			if command.Timeout == 0 {
				return config.DefaultTimeout
			}
			return command.Timeout
		}(),
		CheckRetry: func() int {
			if command.Retry == 0 {
				return config.DefaultRetries
			}
			return command.Retry
		}(),
	}, nil
}

func CreateAction(action Action) (executor.ActionExecutor, error) {
	pre_command, _ := func() (executor.ActionExecutor, error) {
		if action.PreCommand == nil {
			return nil, fmt.Errorf("Command not defined")
		}
		return CreateCommand(*action.PreCommand)
	}()
	post_command, _ := func() (executor.ActionExecutor, error) {
		if action.PostCommand == nil {
			return nil, fmt.Errorf("Command not defined")
		}
		return CreateCommand(*action.PostCommand)
	}()

	switch action.Type {
	case "http":
		return &executor.HttpAction{
			URL:            action.URL,
			ExpectedStatus: action.ExpectedStatus,
			ExpectedBody:   action.ExpectedOutput,
			Timeout:        action.Timeout,
			PreCommand:     pre_command,
			PostCommand:    post_command,
		}, nil
	case "command":
		return &executor.CommandAction{
			Command:        action.Command,
			ExpectedOutput: action.ExpectedOutput,
			Timeout:        action.Timeout,
			PreCommand:     pre_command,
			PostCommand:    post_command,
		}, nil
	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func main() {
	flag.Parse()
	fileName := flag.Arg(0)

	_, err := os.Stat(fileName)
	if err != nil {
		log.Fatalf("file not found: %s\n", err)
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("failed to read file: %s\n", err)
	}

	s := &Scenario{}
	err = yaml.Unmarshal([]byte(fileContent), &s)
	if err != nil {
		log.Fatalf("failed to unmarshal yaml: %s\n", err)
	}

	startupExecutor, err := CreateCommand(s.Startup)
	if err != nil {
		log.Fatalf("failed to create startup executor: %s\n", err)
	}
	log.Println("Executing startup...")
	err = startupExecutor.Execute()
	if err != nil {
		log.Fatalf("Startup failed: %s\n", err)
	}
	log.Println("Startup completed successfully.")

	for i, test := range s.Tests {
		testExecutor, err := CreateAction(test)
		if err != nil {
			log.Fatalf("failed to create test executor: %s\n", err)
		}
		log.Printf("Executing test %d: %s\n", i+1, test.Type)
		err = testExecutor.Execute()
		if err != nil {
			log.Printf("Test %d failed: %s\n", i+1, err)
		} else {
			log.Printf("Test %d passed successfully.\n", i+1)
		}
	}

	shutdownExecutor, err := CreateCommand(s.Shutdown)
	if err != nil {
		log.Fatalf("failed to create shutdown executor: %s\n", err)
	}
	log.Println("Executing shutdown...")
	err = shutdownExecutor.Execute()
	if err != nil {
		log.Fatalf("Shutdown failed: %s\n", err)
	}
	log.Println("Shutdown completed successfully.")
}
