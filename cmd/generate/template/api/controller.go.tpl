package ctr{{.PackageName}}

import (
    "{{.AppPathInProject}}/dto/dto{{.PackageName}}"
    "{{.AppPathInProject}}/service/svc{{.PackageName}}"

    "github.com/gin-gonic/gin"
    "github.com/morehao/golib/gcontext/gincontext"
)
{{if not .TargetFileExist}}
type {{.ReceiverTypePascalName}}Ctr interface {
	{{.FunctionName}}(c *gin.Context)
}

type {{.ReceiverTypeName}}Ctr struct {
	{{.ReceiverTypeName}}Svc svc{{.PackageName}}.{{.ReceiverTypePascalName}}Svc
}

var _ {{.ReceiverTypePascalName}}Ctr = (*{{.ReceiverTypeName}}Ctr)(nil)

func New{{.ReceiverTypePascalName}}Ctr() {{.ReceiverTypePascalName}}Ctr {
	return &{{.ReceiverTypeName}}Ctr{
		{{.ReceiverTypeName}}Svc: svc{{.PackageName}}.New{{.ReceiverTypePascalName}}Svc(),
	}
}
{{end}}
{{if eq .HttpMethod "POST"}}
// {{.FunctionName}} {{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary {{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackageName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Req true "{{.Description}}"
// @Success 200 {object} dto.DefaultRender{data=dto{{.PackageName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Resp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router {{.ApiPrefix}}/{{.ApiSuffix}} [post]
func (ctr *{{.ReceiverTypeName}}Ctr) {{.FunctionName}}(c *gin.Context) {
	var req dto{{.PackageName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Req
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.{{.ReceiverTypeName}}Svc.{{.FunctionName}}(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
{{else if eq .HttpMethod "GET"}}
// {{.FunctionName}} {{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary {{.Description}}
// @accept application/json
// @Produce application/json
// @Param req query dto{{.PackageName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Req true "{{.Description}}"
// @Success 200 {object} dto.DefaultRender{data=dto{{.PackageName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Resp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router {{.ApiPrefix}}/{{.ApiSuffix}} [get]
func (ctr *{{.ReceiverTypeName}}Ctr){{.FunctionName}}(c *gin.Context) {
	var req dto{{.PackageName}}.{{.ReceiverTypePascalName}}{{.FunctionName}}Req
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.{{.ReceiverTypeName}}Svc.{{.FunctionName}}(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
{{end}}
