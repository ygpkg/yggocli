package errorCode

import "github.com/morehao/golib/gerror"

const (
    {{.StructName}}CreateErr      = 100100
    {{.StructName}}DeleteErr      = 100101
    {{.StructName}}UpdateErr      = 100102
    {{.StructName}}GetDetailErr   = 100103
    {{.StructName}}GetPageListErr = 100104
    {{.StructName}}NotExistErr    = 100105
)

var {{.StructName}}ErrMsgMap = gerror.CodeMsgMap{
    {{.StructName}}CreateErr:      "创建{{.Description}}失败",
    {{.StructName}}DeleteErr:      "删除{{.Description}}失败",
    {{.StructName}}UpdateErr:      "修改{{.Description}}失败",
    {{.StructName}}GetDetailErr:   "查看{{.Description}}失败",
    {{.StructName}}GetPageListErr: "查看{{.Description}}列表失败",
    {{.StructName}}NotExistErr:    "{{.Description}}不存在",
}