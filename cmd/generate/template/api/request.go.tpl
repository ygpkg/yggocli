package dto{{.PackageName}}

import (
	"github.com/ygpkg/yg-go/apis/apiobj"
)

type {{.FunctionName}}Request struct {
	apiobj.BaseRequest
	Request {{.FunctionName}}EmbedRequest
}

type {{.FunctionName}}EmbedRequest struct {
}