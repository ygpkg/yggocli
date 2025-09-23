package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/internal/dto/dto{{.PackageName}}"
	"github.com/openrpacloud/{{.ProjectName}}/apps/{{.AppName}}/services/svc{{.PackageName}}"
	"github.com/ygpkg/yg-go/apis/errcode"
	"github.com/ygpkg/yg-go/logs"
)

// {{.FunctionName}} {{.Description}}
// @Tags {{.ApiDocTag}}
// @Summary {{.Description}}
// @Description {{.Description}}
// @Router /{{.AppName}}.{{.FunctionName}} [post]
// @Param request body dto{{.PackageName}}.{{.FunctionName}}Request true "request"
// @Success 200 {object} dto{{.PackageName}}.{{.FunctionName}}Response "response"
func {{.FunctionName}}(ctx *gin.Context, req *dto{{.PackageName}}.{{.FunctionName}}Request, resp *dto{{.PackageName}}.{{.FunctionName}}Response) {
	if req.Validity(resp); resp.Code != 0 {
		logs.ErrorContextf(ctx, "[{{.FunctionName}}] request invalid, req: %s, error message: %v", logs.JSON(req), resp.Message)
		return
	}

	// TODO: 需要手动注册路由和修改 Message 的值
	res, err := svc{{.PackageName}}.{{.FunctionName}}(ctx, req)
	if err != nil {
		logs.ErrorContextf(ctx, "[{{.FunctionName}}] svc{{.PackageName}}.{{.FunctionName}} failed, err: %v", err)
		resp.Code = errcode.ErrCode_InternalError
		resp.Message = errcode.GetMessage(errcode.ErrCode_InternalError)
		return
	}
	resp.Code = res.Code
	resp.Message = res.Message
	resp.Response = res.Response
}