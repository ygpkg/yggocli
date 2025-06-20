package dto{{.PackageName}}

import (
	"github.com/ygpkg/yg-go/apis/apiobj"
)

type {{.FunctionName}}Response struct {
	apiobj.BaseResponse
	Response {{.FunctionName}}EmbedResponse
}

type {{.FunctionName}}EmbedResponse struct {
}
