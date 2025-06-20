package svc{{.PackageName}}

import (
	"github.com/gin-gonic/gin"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/internal/dto/dto{{.PackageName}}"
)

func {{.FunctionName}}(ctx *gin.Context, req *dto{{.PackageName}}.{{.FunctionName}}Request, resp *dto{{.PackageName}}.{{.FunctionName}}Response) error {
	return nil
}
