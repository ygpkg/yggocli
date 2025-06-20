package generate

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/morehao/golib/codegen"
	"github.com/morehao/golib/gutils"
)

const (
	modelLayerParentDir = "models"
	modelLayerSuffix    = "type"

	nullableDefaultDesc = "not null"
	fieldDefaultKeyword = "default"
)

func genModel() error {
	modelGenCfg := cfg.Model

	// 使用工具函数复制嵌入的模板文件到临时目录
	tplDir, getTplErr := CopyEmbeddedTemplatesToTempDir(TemplatesFS, "template/model")
	if getTplErr != nil {
		return getTplErr
	}
	// 清理临时目录
	defer os.RemoveAll(tplDir)
	
	layerParentDirMap := map[codegen.LayerName]string{
		codegen.LayerNameModel: modelLayerParentDir,
		codegen.LayerNameDao:   modelLayerParentDir,
	}

	modelLayerName := codegen.LayerName(fmt.Sprintf("%s%s", cfg.appInfo.AppName, modelLayerSuffix))
	layerNameMap := map[codegen.LayerName]codegen.LayerName{
		codegen.LayerNameModel: modelLayerName,
		codegen.LayerNameDao:   codegen.LayerName(""),
	}

	layerPrefixMap := map[codegen.LayerName]codegen.LayerPrefix{
		codegen.LayerNameDao: codegen.LayerPrefix(""),
	}

	daoLayerName := cfg.appInfo.AppName
	if modelGenCfg.DaoLayerName != "" {
		daoLayerName = modelGenCfg.DaoLayerName
	}

	analysisCfg := &codegen.ModuleCfg{
		CommonConfig: codegen.CommonConfig{
			PackageName:       daoLayerName,
			TplDir:            tplDir,
			RootDir:           workDir,
			LayerParentDirMap: layerParentDirMap,
			LayerNameMap:      layerNameMap,
			LayerPrefixMap:    layerPrefixMap,
			TplFuncMap: template.FuncMap{
				TplFuncIsBuiltInField:      IsBuiltInField,
				TplFuncIsSysField:          IsSysField,
				TplFuncIsDefaultModelLayer: IsDefaultModelLayer,
				TplFuncIsDefaultDaoLayer:   IsDefaultDaoLayer,
			},
		},
		TableName: modelGenCfg.TableName,
	}
	gen := codegen.NewGenerator()
	analysisRes, analysisErr := gen.AnalysisModuleTpl(MysqlClient, analysisCfg)
	if analysisErr != nil {
		return fmt.Errorf("analysis model tpl error: %v", analysisErr)
	}

	var genParamsList []codegen.GenParamsItem
	for _, v := range analysisRes.TplAnalysisList {
		var modelFields []ModelField
		for _, field := range v.ModelFields {
			nullableDesc := nullableDefaultDesc
			if field.IsNullable {
				nullableDesc = ""
			}
			defaultValue := fmt.Sprintf("%s %s", fieldDefaultKeyword, field.DefaultValue)
			if field.DefaultValue == "" {
				defaultValue = ""
			}

			modelFields = append(modelFields, ModelField{
				FieldName:          gutils.ReplaceIdToID(field.FieldName),
				FieldLowerCaseName: gutils.SnakeToLowerCamel(field.FieldName),
				FieldType:          field.FieldType,
				ColumnName:         field.ColumnName,
				ColumnType:         field.ColumnType,
				NullableDesc:       nullableDesc,
				DefaultValue:       defaultValue,
				Comment:            field.Comment,
				IsPrimaryKey:       field.ColumnKey == codegen.ColumnKeyPRI,
			})
		}

		genParamsList = append(genParamsList, codegen.GenParamsItem{
			TargetDir:      v.TargetDir,
			TargetFileName: v.TargetFilename,
			Template:       v.Template,
			ExtraParams: ModelExtraParams{
				AppInfo: AppInfo{
					ProjectName:      cfg.appInfo.ProjectName,
					AppPathInProject: cfg.appInfo.AppPathInProject,
					AppName:          cfg.appInfo.AppName,
				},
				PackageName:  analysisRes.PackageName,
				TableName:    analysisRes.TableName,
				DaoLayerName: daoLayerName,
				Description:  modelGenCfg.Description,
				StructName:   analysisRes.StructName,
				Template:     v.Template,
				ModelFields:  modelFields,
			},
		})

	}
	genParams := &codegen.GenParams{
		ParamsList: genParamsList,
	}
	if err := gen.Gen(genParams); err != nil {
		return err
	}
	if err := addTableName(filepath.Join(workDir, modelLayerParentDir, string(modelLayerName), "db.go"), fmt.Sprintf("TableName%s", analysisRes.StructName), fmt.Sprintf("%q", analysisRes.TableName)); err != nil {
		return err
	}
	return nil
}

// addTableName 向指定文件第一个 const 块末尾添加一条常量定义，
// 并用 go/format 格式化源码，避免格式混乱。
func addTableName(filename, newConstName, newConstValue string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// 查找const块的结束位置（最后一个常量定义）
	constStartIdx := -1
	lastConstIdx := -1
	inConstBlock := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 检测const块开始
		if strings.HasPrefix(trimmed, "const (") {
			constStartIdx = i
			inConstBlock = true
			continue
		}

		// 检测const块结束
		if inConstBlock && trimmed == ")" {
			break
		}

		// 在const块中查找最后一个常量定义
		if inConstBlock && (strings.Contains(trimmed, "=") ||
			(strings.Contains(trimmed, `"`) && !strings.HasPrefix(trimmed, "//"))) {
			lastConstIdx = i
		}
	}

	if constStartIdx == -1 || lastConstIdx == -1 {
		return fmt.Errorf("could not find const block or last constant")
	}

	// 构造新的常量定义
	newConstLine := fmt.Sprintf("\t%s = %s", newConstName, newConstValue)

	// 插入新常量
	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:lastConstIdx+1]...)
	newLines = append(newLines, newConstLine)
	newLines = append(newLines, lines[lastConstIdx+1:]...)

	// 写回文件
	newContent := strings.Join(newLines, "\n")

	// 使用go/format格式化整个文件
	formatted, err := format.Source([]byte(newContent))
	if err != nil {
		return fmt.Errorf("failed to format source: %w", err)
	}

	return os.WriteFile(filename, formatted, 0644)
}

type ModelField struct {
	FieldName          string // 字段名称
	FieldLowerCaseName string // 字段名称小驼峰
	FieldType          string // 字段数据类型，如int、string
	ColumnName         string // 列名
	ColumnType         string // 列数据类型，如varchar(255)
	Comment            string // 字段注释
	NullableDesc       string // 是否允许为空描述，如 NOT NULL
	DefaultValue       string // 默认值,如 DEFAULT 0
	IsPrimaryKey       bool   // 是否是主键
}

type ModelExtraParams struct {
	AppInfo
	PackageName    string
	ModelLayerName string
	DaoLayerName   string
	TableName      string
	Description    string
	StructName     string
	Template       *template.Template
	ModelFields    []ModelField
}
