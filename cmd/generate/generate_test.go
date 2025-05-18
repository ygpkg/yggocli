package generate

import (
	"bytes"
	"testing"

	"github.com/morehao/golib/gutils"
	"github.com/spf13/cobra"
)

// 测试模板文件加载
func TestLoadTemplates(t *testing.T) {
	// 测试模板目录是否存在
	dirs := []string{"template/module", "template/model", "template/api"}
	for _, dir := range dirs {
		entries, err := TemplatesFS.ReadDir(dir)
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

// 辅助函数：执行命令并捕获输出
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), err
}

// 测试配置加载
func TestConfigLoading(t *testing.T) {
	// 执行命令
	_, err := executeCommand(Cmd, "--mode", "model")
	if err != nil {
		t.Errorf("Failed to execute command with config: %v", err)
	}
	t.Log(gutils.ToJsonString(cfg))
}

// 测试命令基本功能
func TestGenerateCommand(t *testing.T) {

	t.Run("generate model code", func(t *testing.T) {
		_, err := executeCommand(Cmd, "--mode", "model")
		if err != nil {
			t.Errorf("Failed to execute command with config: %v", err)
		}
	})
}
