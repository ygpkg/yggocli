package router

import (
	"{{.AppPathInProject}}/controller/ctr{{.PackageName}}"

	"github.com/gin-gonic/gin"
)

// {{.PackageName}}Router 初始化{{.Description}}路由信息
func {{.PackageName}}Router(routerGroup *gin.RouterGroup) {
	{{.PackageName}}Ctr := ctr{{.PackageName}}.New{{.StructName}}Ctr()
	{{.PackageName}}Group := routerGroup.Group("{{.PackageName}}")
	{
		{{.PackageName}}Group.POST("create", {{.PackageName}}Ctr.Create)   // 新建{{.Description}}
		{{.PackageName}}Group.POST("delete", {{.PackageName}}Ctr.Delete)   // 删除{{.Description}}
		{{.PackageName}}Group.POST("update", {{.PackageName}}Ctr.Update)   // 更新{{.Description}}
		{{.PackageName}}Group.GET("detail", {{.PackageName}}Ctr.Detail)    // 根据ID获取{{.Description}}
        {{.PackageName}}Group.GET("pageList", {{.PackageName}}Ctr.PageList)  // 获取{{.Description}}列表
	}
}
