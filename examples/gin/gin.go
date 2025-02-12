package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neatframework/router/register/ginregister"
)

func old() {
	engine := gin.Default()
	// 注册一个路由和处理函数
	engine.Any("/example/hello", WebRoot)
	// 绑定端口，然后启动应用
	engine.Run(":9206")
}

func WebRoot(context *gin.Context) {
	context.String(http.StatusOK, "hello, world")
}

func main() {
	engine := gin.Default()
	ginregister.SetGinEngine(engine)
	ginregister.Service((*Example)(nil),
		exampleController)

	// for compare
	go func() {
		old()
	}()

	engine.Run(":9205")
}

type Example struct {
	Hello func(ctx context.Context) (string, error) `get_mapping:"/example/hello"`
}

var exampleController = Example{
	Hello: func(ctx context.Context) (string, error) {
		return "hello, world", nil
	},
}
