package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/morehao/golib/codegen"
	"github.com/morehao/golib/gutils"
)

const (
	controllerLayerName = "apis"
	serviceLayerName    = "services"

	internalDirName = "internal"
)

func genApi() error {
	apiGenCfg := cfg.Api

	// 使用工具函数复制嵌入的模板文件到临时目录
	tplDir, getTplErr := CopyEmbeddedTemplatesToTempDir(TemplatesFS, "template/api")
	if getTplErr != nil {
		return getTplErr
	}
	// 清理临时目录
	defer os.RemoveAll(tplDir)

	layerParentDirMap := map[codegen.LayerName]string{
		// codegen.LayerNameController: internalDirName,
		// codegen.LayerNameService: internalDirName,
		// codegen.LayerNameDto:        internalDirName,
	}

	layerNameMap := map[codegen.LayerName]codegen.LayerName{
		codegen.LayerNameController: controllerLayerName,
		codegen.LayerNameService:    serviceLayerName,
	}

	layerPrefixMap := map[codegen.LayerName]codegen.LayerPrefix{
		codegen.LayerNameController: codegen.LayerPrefix(""),
	}

	analysisCfg := &codegen.ApiCfg{
		CommonConfig: codegen.CommonConfig{
			PackageName:       apiGenCfg.PackageName,
			TplDir:            tplDir,
			RootDir:           filepath.Join(workDir, internalDirName),
			LayerParentDirMap: layerParentDirMap,
			LayerNameMap:      layerNameMap,
			LayerPrefixMap:    layerPrefixMap,
		},
		TargetFilename: apiGenCfg.TargetFilename,
	}
	gen := codegen.NewGenerator()
	analysisRes, analysisErr := gen.AnalysisApiTpl(analysisCfg)
	if analysisErr != nil {
		return fmt.Errorf("analysis api tpl error: %v", analysisErr)
	}

	structName := gutils.SnakeToPascal(gutils.TrimFileExtension(apiGenCfg.TargetFilename))
	structNameLowerCamel := gutils.FirstLetterToLower(structName)
	functionName := gutils.FirstLetterToUpper(apiGenCfg.FunctionName)
	functionNameLowerCamel := gutils.FirstLetterToLower(apiGenCfg.FunctionName)
	var genParamsList []codegen.GenParamsItem

	for _, v := range analysisRes.TplAnalysisList {
		param := codegen.GenParamsItem{
			TargetDir:      v.TargetDir,
			TargetFileName: v.TargetFilename,
			Template:       v.Template,
			ExtraParams: ApiExtraParams{
				AppInfo: AppInfo{
					ProjectName:      cfg.appInfo.ProjectName,
					AppPathInProject: cfg.appInfo.AppPathInProject,
					AppName:          cfg.appInfo.AppName,
				},
				PackageName:            analysisRes.PackageName,
				TargetFileExist:        v.TargetFileExist,
				Description:            apiGenCfg.Description,
				StructName:             structName,
				StructNameLowerCamel:   structNameLowerCamel,
				FunctionName:           functionName,
				FunctionNameLowerCamel: functionNameLowerCamel,
				HttpMethod:             apiGenCfg.HttpMethod,
				ApiDocTag:              apiGenCfg.ApiDocTag,
				Template:               v.Template,
			},
		}
		switch v.LayerName {
		case controllerLayerName:
			param.TargetDir = filepath.Dir(v.TargetDir)
		case serviceLayerName:
			last := filepath.Base(v.TargetDir)
			dir := filepath.Dir(v.TargetDir)
			secondLast := filepath.Base(dir)
			appDir := filepath.Dir(filepath.Dir(filepath.Dir(v.TargetDir)))
			param.TargetDir = filepath.Join(appDir, secondLast, last)
		}
		genParamsList = append(genParamsList, param)

	}
	genParams := &codegen.GenParams{
		ParamsList: genParamsList,
	}
	if err := gen.Gen(genParams); err != nil {
		return err
	}
	// routerCallContent := fmt.Sprintf(`routerGroup.%s("/%s", %sCtr.%s) // %s`, apiGenCfg.HttpMethod, functionNameLowerCamel, structNameLowerCamel, functionName, apiGenCfg.Description)
	// routerEnterFilepath := filepath.Join(workDir, fmt.Sprintf("/router/%s.go", gutils.TrimFileExtension(apiGenCfg.PackageName)))
	// if err := gast.AddContentToFuncWithLineNumber(routerEnterFilepath, fmt.Sprintf("%sRouter", structNameLowerCamel), routerCallContent, -2); err != nil {
	// 	return fmt.Errorf("appendContentToFunc error: %v", err)
	// }
	return nil
}

type ApiExtraParams struct {
	AppInfo
	PackageName            string
	Description            string
	TargetFileExist        bool
	HttpMethod             string
	StructName             string
	StructNameLowerCamel   string
	FunctionName           string
	FunctionNameLowerCamel string
	ApiDocTag              string
	Template               *template.Template
}
