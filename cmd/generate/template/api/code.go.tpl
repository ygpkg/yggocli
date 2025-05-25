package code

import "github.com/morehao/golib/gerror"

var {{.StructName}}{{.FunctionName}}Err = gerror.Error{
	Code: 100100,
	Msg:  "{{.Description}}失败",
}
