package demoapp

import (
	"os"
	"testing"

	"github.com/morehao/gcli/cmd/generate"
	"github.com/stretchr/testify/assert"
)

func TestLoadTemplates(t *testing.T) {
	// 测试模板目录是否存在
	dirs := []string{"template/module", "template/model", "template/api"}
	for _, dir := range dirs {
		entries, err := generate.TemplatesFS.ReadDir(dir)
		if err != nil {
			t.Errorf("Failed to read directory %s: %v", dir, err)
			continue
		}
		if len(entries) == 0 {
			t.Errorf("Directory %s is empty", dir)
		}
		t.Logf("Directory %s is not empty", dir)
	}
}

func TestGetAppRootDir(t *testing.T) {
	workDir, _ := os.Getwd()
	appDir, err := generate.GetAppInfo(workDir)
	assert.Nil(t, err)
	t.Log(appDir)
}

func TestGenerateCommand(t *testing.T) {
	t.Run("generate model code", func(t *testing.T) {
		_, err := executeCommand(generate.Cmd, "--mode", "model")
		if err != nil {
			t.Errorf("Failed to execute command with config: %v", err)
		}
	})
}
