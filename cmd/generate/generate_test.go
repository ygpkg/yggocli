package generate

import (
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	// 测试模板文件是否存在
	_, err := templatesFS.ReadFile("template/module.go.tmpl")
	if err != nil {
		t.Errorf("Failed to read module template: %v", err)
	}

	_, err = templatesFS.ReadFile("template/model.go.tmpl")
	if err != nil {
		t.Errorf("Failed to read model template: %v", err)
	}

	_, err = templatesFS.ReadFile("template/api.go.tmpl")
	if err != nil {
		t.Errorf("Failed to read api template: %v", err)
	}
}

// func TestGenerateModule(t *testing.T) {
// 	// 设置测试配置
// 	cfg = &Config{
// 		Mysql: "test:test@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
// 		CodeGen: CodeGen{
// 			Mode: "module",
// 			Module: ModuleConfig{
// 				InternalAppRootDir: "test/internal",
// 				ProjectRootDir:     "test",
// 				Description:        "test module",
// 				ApiDocTag:          "test",
// 				ApiGroup:           "test",
// 				ApiPrefix:          "/test",
// 				PackageName:        "test",
// 				TableName:          "test_table",
// 			},
// 		},
// 	}
//
// 	// 测试生成模块
// 	err := genModule()
// 	if err != nil {
// 		t.Errorf("Failed to generate module: %v", err)
// 	}
// }
