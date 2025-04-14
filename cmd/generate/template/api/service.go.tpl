package svc{{.PackagePascalName}}

import (
    "{{.ProjectRootDir}}/internal/app/dto/dto{{.PackagePascalName}}"

    "github.com/gin-gonic/gin"
)

{{if not .TargetFileExist}}
type {{.ReceiverTypePascalName}}Svc interface {
    {{.FunctionName}}(c *gin.Context, req *dto{{.PackagePascalName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Req) (*dto{{.PackagePascalName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Resp, error)
}

type {{.ReceiverTypeName}}Svc struct {
}

var _ {{.ReceiverTypePascalName}}Svc = (*{{.ReceiverTypeName}}Svc)(nil)

func New{{.ReceiverTypePascalName}}Svc() {{.ReceiverTypePascalName}}Svc {
    return &{{.ReceiverTypeName}}Svc{
    }
}
{{end}}
func (svc *{{.ReceiverTypeName}}Svc) {{.FunctionName}}(c *gin.Context, req *dto{{.PackagePascalName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Req) (*dto{{.PackagePascalName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Resp, error) {
    return &dto{{.PackagePascalName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Resp{}, nil
}
