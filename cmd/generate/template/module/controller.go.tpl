package ctr{{.PackageName}}

import (
    "{{.AppPathInProject}}/dto/dto{{.PackageName}}"
    "{{.AppPathInProject}}/service/svc{{.PackageName}}"

    "github.com/gin-gonic/gin"
    "github.com/morehao/golib/gcontext/gincontext"
)

type {{.StructName}}Ctr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
}

type {{.PackageName}}Ctr struct {
	{{.PackageName}}Svc svc{{.PackageName}}.{{.StructName}}Svc
}

var _ {{.StructName}}Ctr = (*{{.PackageName}}Ctr)(nil)

func New{{.StructName}}Ctr() {{.StructName}}Ctr {
	return &{{.PackageName}}Ctr{
		{{.PackageName}}Svc: svc{{.PackageName}}.New{{.StructName}}Svc(),
	}
}


// Create 创建{{.Description}}
// @Tags {{.Description}}
// @Summary 创建{{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackageName}}.{{.StructName}}CreateReq true "创建{{.Description}}"
// @Success 200 {object} gincontext.DtoRender{data=dto{{.PackageName}}.{{.StructName}}CreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /{{.AppName}}/{{.PackageName}}/create [post]
func (ctr *{{.PackageName}}Ctr) Create(ctx *gin.Context) {
	var req dto{{.PackageName}}.{{.StructName}}CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.{{.PackageName}}Svc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Delete 删除{{.Description}}
// @Tags {{.Description}}
// @Summary 删除{{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackageName}}.{{.StructName}}DeleteReq true "删除{{.Description}}"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /{{.AppName}}/{{.PackageName}}/delete [post]
func (ctr *{{.PackageName}}Ctr) Delete(ctx *gin.Context) {
	var req dto{{.PackageName}}.{{.StructName}}DeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.{{.PackageName}}Svc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

// Update 修改{{.Description}}
// @Tags {{.Description}}
// @Summary 修改{{.Description}}
// @accept application/json
// @Produce application/json
// @Param req body dto{{.PackageName}}.{{.StructName}}UpdateReq true "修改{{.Description}}"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /{{.AppName}}/{{.PackageName}}/update [post]
func (ctr *{{.PackageName}}Ctr) Update(ctx *gin.Context) {
	var req dto{{.PackageName}}.{{.StructName}}UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.{{.PackageName}}Svc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

// Detail {{.Description}}详情
// @Tags {{.Description}}
// @Summary {{.Description}}详情
// @accept application/json
// @Produce application/json
// @Param req query dto{{.PackageName}}.{{.StructName}}DetailReq true "{{.Description}}详情"
// @Success 200 {object} gincontext.DtoRender{data=dto{{.PackageName}}.{{.StructName}}DetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /{{.AppName}}/{{.PackageName}}/detail [get]
func (ctr *{{.PackageName}}Ctr) Detail(ctx *gin.Context) {
	var req dto{{.PackageName}}.{{.StructName}}DetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.{{.PackageName}}Svc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// PageList {{.Description}}列表
// @Tags {{.Description}}
// @Summary {{.Description}}列表分页
// @accept application/json
// @Produce application/json
// @Param req query dto{{.PackageName}}.{{.StructName}}PageListReq true "{{.Description}}列表"
// @Success 200 {object} gincontext.DtoRender{data=dto{{.PackageName}}.{{.StructName}}PageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /{{.AppName}}/{{.PackageName}}/pageList [get]
func (ctr *{{.PackageName}}Ctr) PageList(ctx *gin.Context) {
	var req dto{{.PackageName}}.{{.StructName}}PageListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.{{.PackageName}}Svc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
