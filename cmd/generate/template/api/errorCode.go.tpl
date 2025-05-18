package errorCode

import "github.com/morehao/golib/gerror"

var {{.ReceiverTypePascalName}}{{.FunctionName}}Err = gerror.Error{
	Code: 100100,
	Msg:  "{{.Description}}失败",
}
