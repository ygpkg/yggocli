package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/morehao/golib/codegen"
	"github.com/morehao/golib/gast"
	"github.com/morehao/golib/gutils"
)

func genModule() error {
	moduleGenCfg := cfg.Module

	// 使用工具函数复制嵌入的模板文件到临时目录
	tplDir, getTplErr := CopyEmbeddedTemplatesToTempDir(TemplatesFS, "template/module")
	if getTplErr != nil {
		return getTplErr
	}
	// 清理临时目录
	defer os.RemoveAll(tplDir)

	analysisCfg := &codegen.ModuleCfg{
		CommonConfig: codegen.CommonConfig{
			PackageName:       moduleGenCfg.PackageName,
			TplDir:            tplDir,
			RootDir:           workDir,
			LayerParentDirMap: cfg.LayerParentDirMap,
			LayerNameMap:      cfg.LayerNameMap,
			LayerPrefixMap:    cfg.LayerPrefixMap,
			TplFuncMap: template.FuncMap{
				TplFuncIsBuiltInField:      IsBuiltInField,
				TplFuncIsSysField:          IsSysField,
				TplFuncIsDefaultModelLayer: IsDefaultModelLayer,
				TplFuncIsDefaultDaoLayer:   IsDefaultDaoLayer,
			},
		},
		TableName: moduleGenCfg.TableName,
	}
	gen := codegen.NewGenerator()
	analysisRes, analysisErr := gen.AnalysisModuleTpl(MysqlClient, analysisCfg)
	if analysisErr != nil {
		return fmt.Errorf("analysis module tpl error: %v", analysisErr)
	}

	var modelLayerName, daoLayerName codegen.LayerName
	for _, v := range analysisRes.TplAnalysisList {
		if v.OriginLayerName == codegen.LayerNameModel {
			modelLayerName = v.LayerName
		}
		if v.OriginLayerName == codegen.LayerNameDao {
			daoLayerName = v.LayerName
		}
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
				AppInfo: AppInfo{
					ProjectName:      cfg.appInfo.ProjectName,
					AppPathInProject: cfg.appInfo.AppPathInProject,
					AppName:          cfg.appInfo.AppName,
				},
				PackageName:          analysisRes.PackageName,
				TableName:            analysisRes.TableName,
				ModelLayerName:       string(modelLayerName),
				DaoLayerName:         string(daoLayerName),
				Description:          moduleGenCfg.Description,
				StructName:           analysisRes.StructName,
				StructNameLowerCamel: gutils.FirstLetterToLower(analysisRes.StructName),
				Template:             v.Template,
				ModelFields:          modelFields,
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
	routerContent := fmt.Sprintf("%sRouter(routerGroup)", gutils.FirstLetterToLower(analysisRes.StructName))
	routerEnterFilepath := filepath.Join(workDir, "/router/enter.go")
	if err := gast.AddContentToFunc(routerEnterFilepath, "RegisterRouter", routerContent); err != nil {
		return fmt.Errorf("router appendContentToFunc error: %v", err)
	}

	// 注册错误码
	codeContent := fmt.Sprintf("registerError(%sErrorMsgMap)", gutils.FirstLetterToLower(analysisRes.StructName))
	codeEnterFilepath := filepath.Join(workDir, "/code/enter.go")
	if err := gast.AddContentToFunc(codeEnterFilepath, "init", codeContent); err != nil {
		return fmt.Errorf("code appendContentToFunc error: %v", err)
	}
	return nil
}

type ModuleExtraParams struct {
	AppInfo
	PackageName          string
	ModelLayerName       string
	DaoLayerName         string
	TableName            string
	Description          string
	StructName           string
	StructNameLowerCamel string // 结构体小写驼峰名
	Template             *template.Template
	ModelFields          []ModelField
}
