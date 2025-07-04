package obj{{.PackageName}}

type {{.StructName}}BaseInfo struct {
{{- range .ModelFields}}
{{- if isSysField .FieldName}}
    {{- continue}}
{{- end}}
{{- if eq .FieldType "time.Time"}}
    {{.FieldName}} int64 `json:"{{.FieldLowerCaseName}}" form:"{{.FieldLowerCaseName}}"` // {{.Comment}}
{{- else}}
    {{.FieldName}} {{.FieldType}} `json:"{{.FieldLowerCaseName}}" form:"{{.FieldLowerCaseName}}"` // {{.Comment}}
{{- end}}
{{- end}}
}
