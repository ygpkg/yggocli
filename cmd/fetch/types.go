package fetch

type BaseResponse struct {
	Code int    `json:"code"`
	Env  string `json:"env"`
}
type ResourceType string
type KnownowForest struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}
type FileItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	PublicUrl string `json:"public_url"`
	FileID    uint   `json:"file_id"`
	ForestID  uint   `json:"forest_id"`
}
type RequestBody struct {
	Cmd     string `json:"cmd"`
	Env     string `json:"env"`
	Request struct {
		ResourceIDs  []int  `json:"resource_ids"`
		ResourceType string `json:"resource_type"`
	} `json:"request"`
	Version string `json:"version"`
}

// --- 灵活的响应结构体 (使用 map[string]interface{}) ---

// GenericResource 承载所有类型都共有的字段，Meta 使用 map
type GenericResource struct {
	ID           uint                   `json:"id"`
	Meta         map[string]interface{} `json:"meta"` // 关键：灵活的 Meta
	ResourceType ResourceType           `json:"resource_type"`
}

// GenericEmbedResponse 包含通用资源列表
type GenericEmbedResponse struct {
	Data []*GenericResource `json:"data"`
}

// GenericApiResponse 响应的最外层
type GenericApiResponse struct {
	BaseResponse
	Response GenericEmbedResponse `json:"Response"`
}

type Processor interface {
	// Route 返回该处理器对应的 API 路由
	Route() string
	// PreProcessData 用于根据 flags 准备请求体
	PreProcessData(flags *FlagSet) RequestBody
	// Process 负责处理“通用”的响应数据，并执行其特定的业务逻辑
	Process(data []*GenericResource, baseDir string) error
}
