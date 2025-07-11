package {{.ModelLayerName}}

import (
	"gorm.io/gorm"
)

// {{.StructName}} {{.Description}}表结构体
type {{.StructName}} struct {
    gorm.Model
{{- range .ModelFields}}
    {{- if isBuiltInField .FieldName}}
        {{- continue}}
    {{- else}}
	{{.FieldName}} {{.FieldType}} `gorm:"column:{{.ColumnName}};type:{{.ColumnType}};{{.NullableDesc}};{{.DefaultValue}};comment:{{.Comment}}"`
	{{- end}}
{{- end}}
}

type {{.StructName}}List []{{.StructName}}

func ({{.StructName}} ) TableName() string {
  return TableName{{.StructName}}
}

func (l {{.StructName}}List) ToMap() map[uint]{{.StructName}} {
	m := make(map[uint]{{.StructName}})
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}