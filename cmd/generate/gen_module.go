package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/morehao/golib/codegen"
	"github.com/morehao/golib/gast"
	"github.com/morehao/golib/gutils"
)

func genModule() error {
	moduleGenCfg := cfg.Module

	// 使用工具函数复制嵌入的模板文件到临时目录
	tplDir, err := CopyEmbeddedTemplatesToTempDir(TemplatesFS, "template/module")
	if err != nil {
		return err
	}
	// 清理临时目录
	defer os.RemoveAll(tplDir)

	rootDir := filepath.Join(workDir, moduleGenCfg.InternalAppRootDir)
	analysisCfg := &codegen.ModuleCfg{
		CommonConfig: codegen.CommonConfig{
			TplDir:            tplDir,
			PackageName:       moduleGenCfg.PackageName,
			RootDir:           rootDir,
			LayerParentDirMap: cfg.LayerParentDirMap,
			LayerNameMap:      cfg.LayerNameMap,
			LayerPrefixMap:    cfg.LayerPrefixMap,
			TplFuncMap: template.FuncMap{
				TplFuncIsSysField: IsSysField,
			},
		},
		TableName: moduleGenCfg.TableName,
	}
	gen := codegen.NewGenerator()
	analysisRes, analysisErr := gen.AnalysisModuleTpl(MysqlClient, analysisCfg)
	if analysisErr != nil {
		return fmt.Errorf("analysis module tpl error: %v", analysisErr)
	}

	var genParamsList []codegen.GenParamsItem
	for _, v := range analysisRes.TplAnalysisList {
		var modelFields []ModelField
		for _, field := range v.ModelFields {
			modelFields = append(modelFields, ModelField{
				FieldName:          gutils.ReplaceIdToID(field.FieldName),
				FieldLowerCaseName: gutils.SnakeToLowerCamel(field.FieldName),
				FieldType:          field.FieldType,
				ColumnName:         field.ColumnName,
				ColumnType:         field.ColumnType,
				Comment:            field.Comment,
				IsPrimaryKey:       field.ColumnKey == codegen.ColumnKeyPRI,
			})
		}

		genParamsList = append(genParamsList, codegen.GenParamsItem{
			TargetDir:      v.TargetDir,
			TargetFileName: v.TargetFilename,
			Template:       v.Template,
			ExtraParams: ModuleExtraParams{
				PackageName:            analysisRes.PackageName,
				ProjectRootDir:         moduleGenCfg.ProjectRootDir,
				TableName:              analysisRes.TableName,
				Description:            moduleGenCfg.Description,
				StructName:             analysisRes.StructName,
				ReceiverTypeName:       gutils.FirstLetterToLower(analysisRes.StructName),
				ReceiverTypePascalName: analysisRes.StructName,
				ApiDocTag:              moduleGenCfg.ApiDocTag,
				ApiGroup:               moduleGenCfg.ApiGroup,
				ApiPrefix:              strings.TrimSuffix(moduleGenCfg.ApiPrefix, "/"),
				Template:               v.Template,
				ModelFields:            modelFields,
			},
		})

	}
	genParams := &codegen.GenParams{
		ParamsList: genParamsList,
	}
	if err := gen.Gen(genParams); err != nil {
		return err
	}

	// 注册路由
	routerCallContent := fmt.Sprintf("%sRouter(routerGroup)", gutils.FirstLetterToLower(analysisRes.StructName))
	routerEnterFilepath := filepath.Join(rootDir, "/router/enter.go")
	if err := gast.AddContentToFunc(routerEnterFilepath, "RegisterRouter", routerCallContent); err != nil {
		return fmt.Errorf("appendContentToFunc error: %v", err)
	}
	return nil
}

type ModuleExtraParams struct {
	ServiceName            string
	ProjectRootDir         string
	PackageName            string
	PackagePascalName      string
	TableName              string
	Description            string
	StructName             string
	ReceiverTypeName       string
	ReceiverTypePascalName string
	ApiGroup               string
	ApiPrefix              string
	ApiDocTag              string
	Template               *template.Template
	ModelFields            []ModelField
}
