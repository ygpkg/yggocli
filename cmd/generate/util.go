package generate

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	TplFuncIsBuiltInField      = "isBuiltInField"
	TplFuncIsSysField          = "isSysField"
	TplFuncIsDefaultModelLayer = "isDefaultModelLayer"
	TplFuncIsDefaultDaoLayer   = "isDefaultDaoLayer"
)

func IsBuiltInField(name string) bool {
	buildInFieldMap := map[string]struct{}{
		"ID":        {},
		"CreatedAt": {},
		"UpdatedAt": {},
		"DeletedAt": {},
	}
	_, ok := buildInFieldMap[name]
	return ok
}

func IsSysField(name string) bool {
	sysFieldMap := map[string]struct{}{
		"ID":        {},
		"CreatedAt": {},
		"CreatedBy": {},
		"UpdatedAt": {},
		"UpdatedBy": {},
		"DeletedAt": {},
		"DeletedBy": {},
	}
	_, ok := sysFieldMap[name]
	return ok
}

func IsDefaultModelLayer(name string) bool {
	return name == "model"
}

func IsDefaultDaoLayer(name string) bool {
	return name == "dao"
}

// CopyEmbeddedTemplatesToTempDir 将嵌入的模板文件复制到临时目录，并返回该目录的路径。
func CopyEmbeddedTemplatesToTempDir(embeddedFS embed.FS, root string) (string, error) {
	// 创建一个临时目录来存放模板文件
	tempDir, err := os.MkdirTemp("", "codegen_templates")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// 将嵌入的模板文件复制到临时目录
	err = fs.WalkDir(embeddedFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			data, readErr := embeddedFS.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			// 保持目录结构
			relPath, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			targetPath := filepath.Join(tempDir, relPath)
			if mkDirErr := os.MkdirAll(filepath.Dir(targetPath), 0755); mkDirErr != nil {
				return mkDirErr
			}
			if writeErr := os.WriteFile(targetPath, data, 0644); writeErr != nil {
				return writeErr
			}
		}
		return nil
	})
	if err != nil {
		// 如果复制失败，清理临时目录
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to copy templates: %v", err)
	}

	return tempDir, nil
}

// GetAppInfo 应用模块路径信息
// 输入示例：/Users/morehao/xxx/go-gin-web/internal/apps/demo
func GetAppInfo(workDir string) (*AppInfo, error) {
	cleanPath := filepath.Clean(workDir)
	segments := strings.Split(cleanPath, string(filepath.Separator))

	// 查找 "apps/{appName}"
	var appsIndex = -1
	for i := 0; i < len(segments)-1; i++ {
		if segments[i] == "apps" {
			appsIndex = i
			break
		}
	}
	if appsIndex == -1 {
		return nil, fmt.Errorf("invalid structure: path does not contain /apps/{appName}")
	}

	appName := segments[appsIndex+1]

	// 解析项目名、app名、相对路径
	projectNameIndex := appsIndex - 1
	if projectNameIndex < 0 {
		return nil, fmt.Errorf("cannot determine project name from path: %s", workDir)
	}
	projectName := segments[projectNameIndex]

	appPathInProject := filepath.Join(projectName, "apps", appName)

	return &AppInfo{
		AppPathInProject: appPathInProject,
		ProjectName:      projectName,
		AppName:          appName,
	}, nil
}

// ExecuteCommand 执行命令并捕获输出
func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), err
}
