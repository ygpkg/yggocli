package model

import (
	"time"

	"gorm.io/gorm"
)

// {{.StructName}}Entity {{.Description}}表结构体
type {{.StructName}}Entity struct {
{{- range .ModelFields}}
	{{- if .IsPrimaryKey}}
	{{.FieldName}} uint64 `gorm:"column:{{.ColumnName}};comment:{{.Comment}};primaryKey"`
	{{- else if eq .FieldName "DeletedAt"}}
	{{.FieldName}} gorm.DeletedAt `gorm:"column:{{.ColumnName}};comment:{{.Comment}}"`
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

func (l {{.StructName}}EntityList) ToMap() map[uint64]{{.StructName}}Entity {
	m := make(map[uint64]{{.StructName}}Entity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}