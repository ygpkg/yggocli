package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/internal/dto/dto{{.PackageName}}"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/internal/services/svc{{.PackageName}}"
	"github.com/openrpacloud/{{.ProjectName}}/pkgs/apis/errcode"
	"github.com/ygpkg/yg-go/logs"
	"github.com/ygpkg/yg-go/validate"
)

// {{.FunctionName}} {{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary {{.Description}}
// @Description {{.Description}}
// @Router /{{.AppName}}.{{.FunctionName}} [post]
// @Param request body dto{{.PackageName}}.{{.FunctionName}}Request true "request"
// @Success 200 {object} dto{{.PackageName}}.{{.FunctionName}}Response "response"
func {{.FunctionName}}(ctx *gin.Context, req *dto{{.PackageName}}.{{.FunctionName}}Request, resp *dto{{.PackageName}}.{{.FunctionName}}Response) {
	if err := validate.IsValidStruct(req, false); err != nil {
		resp.Code = errcode.ErrCode_BadRequest
		resp.Message = err.Error()
		logs.ErrorContextf(ctx, "[{{.FunctionName}}] validate.IsValidStruct failed, err:%v, req:%s", err, logs.JSON(req))
		return
	}
	// TODO: 需要手动注册路由 and 手动定义错误码
	if err := svc{{.PackageName}}.{{.FunctionName}}(ctx, req, resp); err != nil {
		logs.ErrorContextf(ctx, "[{{.FunctionName}}] svc{{.PackageName}}.{{.FunctionName}} failed, err:%v, req:%s", err, logs.JSON(req))
		resp.Code = errcode.ErrCode_{{.FunctionName}}
		resp.Message = errcode.GetMessage(errcode.ErrCode_{{.FunctionName}})
		return
	}
}