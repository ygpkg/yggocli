package demoapp

import (
	"testing"

	"github.com/morehao/gocli/cmd/generate"
)

func TestGenerateCommand(t *testing.T) {
	t.Run("generate model code", func(t *testing.T) {
		_, err := generate.ExecuteCommand(generate.Cmd, "--mode", "model")
		if err != nil {
			t.Errorf("Failed to execute command with config: %v", err)
		}
	})
	t.Run("generate module code", func(t *testing.T) {
		_, err := generate.ExecuteCommand(generate.Cmd, "--mode", "module")
		if err != nil {
			t.Errorf("Failed to execute command with config: %v", err)
		}
	})
	t.Run("generate api code", func(t *testing.T) {
		_, err := generate.ExecuteCommand(generate.Cmd, "--mode", "api")
		if err != nil {
			t.Errorf("Failed to execute command with config: %v", err)
		}
	})
}
