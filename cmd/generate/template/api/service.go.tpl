package svc{{.PackageName}}

import (
	"github.com/gin-gonic/gin"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/internal/dto/dto{{.PackageName}}"
)

func {{.FunctionName}}(ctx *gin.Context, req *dto{{.PackageName}}.{{.FunctionName}}Request) (res *dto{{.PackageName}}.{{.FunctionName}}Response, err error) {
	res = &dto{{.PackageName}}.{{.FunctionName}}Response
	return res, nil
}
