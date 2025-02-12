package exampledto

type HelloRequest struct {
	Name string `json:"name" form:"name" title:"姓名"`
}
