# gofr-router-register

go 轻量级路由注册器, 通过一套统一的接口开发标准, 适配不同的框架, 并提供一些通用的能力

- [x] 通过tag注册路由
  - [x] get_mapping
  - [x] post_mapping
  - [x] put_mapping
  - [x] any_mapping
- [x] 统一入参出参处理
- [x] 根据接口结构生成openapi文档
  - [x] 通过tag声明字段说明, 类型
  - [ ] 控制器说明

框架适配

- [x] gin
  - [x] controller 批量注册
  - [ ] handler 单路由注册

# example with gin

### register

```go
package main

import (
	"context"
	"github.com/neatframework/router/register/ginregister"
	"github.com/gin-gonic/gin"
)

type HelloRequest struct {
    Name string `json:"name" form:"name" title:"姓名"`
}

type HelloResponse struct {
    Message string `json:"message" title:"消息"`
}

type Example struct {
    Hello func(ctx context.Context, req HelloRequest) (HelloResponse, error) `get_mapping:"/example/hello"`
}

var exampleController = Example{
    Hello: Hello,
}

func Hello(ctx context.Context, req HelloRequest) (HelloResponse, error) {
    name := req.Name
    if name == "" {
        name = "world"
    }
    res := HelloResponse{
        Message: "hello, " + name,
    }
    
    return res, nil
}

func main() {
    engine := gin.Default()
    ginregister.SetGinEngine(engine)
    ginregister.Service((*Example)(nil), exampleController)
    
    engine.Run(":9205")
}
```

### openapi

```go
package main

import (
    "github.com/neatframework/router/register/register/openapi"
    "os"
)

func main() {
    openapiGen()
}

func openapiGen() {
    _ = os.Remove("openapi.json")
    file, err := os.OpenFile("openapi.json", os.O_CREATE|os.O_RDWR, 0777)
    if err != nil {
        panic(err)
    }
    
    defer file.Close()
    
    _, err = file.WriteString(openapi.Print())
    if err != nil {
        panic(err)
    }
}
```

更多例子请查看example目录