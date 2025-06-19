package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
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
	addTableName(filepath.Join(workDir, modelLayerParentDir, string(modelLayerName), "db.go"), fmt.Sprintf("TableName%s", analysisRes.StructName), fmt.Sprintf("%q", analysisRes.TableName))
	return nil
}

// addTableName 将一个新的 const 常量添加到目标文件中第一个 const 块的末尾
func addTableName(filename, newConstName, newConstValue string) error {
	fset := token.NewFileSet()
	node, parseFileErr := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if parseFileErr != nil {
		return fmt.Errorf("failed to parse file: %v", parseFileErr)
	}

	inserted := false
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		// 构建新的 const ValueSpec
		newSpec := &ast.ValueSpec{
			Names:  []*ast.Ident{ast.NewIdent(newConstName)},
			Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: newConstValue}},
		}

		// 添加到 const block 的末尾
		genDecl.Specs = append(genDecl.Specs, newSpec)

		inserted = true
		break // 只处理第一个 const 块
	}

	if !inserted {
		return fmt.Errorf("no const block found to insert into")
	}

	// 打开文件用于覆盖写入
	outFile, rewriteErr := os.Create(filename)
	if rewriteErr != nil {
		return fmt.Errorf("failed to open file for writing: %v", rewriteErr)
	}
	defer outFile.Close()

	// 输出格式化 AST 到文件
	if err := printer.Fprint(outFile, fset, node); err != nil {
		return fmt.Errorf("failed to write modified file: %v", err)
	}
	return nil
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
