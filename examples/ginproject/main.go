package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/neatframework/router/register/examples/ginproject/controller"
	"github.com/neatframework/router/register/examples/ginproject/facade"
	"github.com/neatframework/router/register/ginregister"
	"github.com/neatframework/router/register/openapi"
)

func main() {
	engine := gin.Default()
	ginregister.SetGinEngine(engine)
	ginregister.Service((*facade.Example)(nil), controller.ExampleController)

	openapiGen()

	engine.Run(":9205")
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
