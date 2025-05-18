package generate

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	TplFuncIsSysField          = "isSysField"
	TplFuncIsDefaultModelLayer = "isDefaultModelLayer"
	TplFuncIsDefaultDaoLayer   = "isDefaultDaoLayer"
)

func IsSysField(name string) bool {
	sysFieldMap := map[string]struct{}{
		"Id":        {},
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

// GetAppInfo 获取项目模块路径（包含项目根目录 + apps + 当前模块名）
// 例如输入：/Users/morehao/xxx/go-gin-web/apps/demo
// 返回：/Users/morehao/xxx/go-gin-web/apps/demo
func GetAppInfo(workDir string) (*AppInfo, error) {
	// 检查是否存在 internal 目录（确认是一个应用模块）
	internalPath := filepath.Join(workDir, "internal")
	info, err := os.Stat(internalPath)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("invalid module: %s does not contain internal/ directory", workDir)
	}

	// 向上追溯，找到 apps 目录
	appsDir := filepath.Dir(workDir)
	if filepath.Base(appsDir) != "apps" {
		return nil, fmt.Errorf("invalid structure: %s is not under apps/", workDir)
	}

	// 找到项目根（apps 的上一级）
	projectPath := filepath.Dir(appsDir)
	projectName := filepath.Base(projectPath)

	// 拼接最终路径：projectName/apps/appName
	appName := filepath.Base(workDir)
	projectAppPath := filepath.Join(projectName, "apps", appName)
	return &AppInfo{
		ProjectAppRelativePath: projectAppPath,
		ProjectName:            projectName,
		AppName:                appName,
	}, nil
}
