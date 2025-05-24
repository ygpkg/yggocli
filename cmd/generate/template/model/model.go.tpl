package model

import (
	"time"

	"gorm.io/gorm"
)

// {{.StructName}}Entity {{.Description}}表结构体
type {{.StructName}}Entity struct {
    gorm.Model
{{- range .ModelFields}}
    {{- if isBuiltInField .FieldName}}
        {{- continue}}
    {{- else}}
	{{.FieldName}} {{.FieldType}} `gorm:"column:{{.ColumnName}};comment:{{.Comment}}"`
	{{- end}}
{{- end}}
}

type {{.StructName}}EntityList []{{.StructName}}Entity

const TblName{{.StructName}} = "{{.TableName}}"

func ({{.StructName}}Entity ) TableName() string {
  return TblName{{.StructName}}
}

func (l {{.StructName}}EntityList) ToMap() map[uint]{{.StructName}}Entity {
	m := make(map[uint]{{.StructName}}Entity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}