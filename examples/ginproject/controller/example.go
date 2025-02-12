package controller

import (
	"context"

	"github.com/neatframework/router/register/examples/ginproject/facade"
	"github.com/neatframework/router/register/examples/ginproject/facade/exampledto"
)

var ExampleController = facade.Example{
	Hello:    Hello,
	HelloPtr: HelloPtr,
}

func HelloPtr(ctx context.Context, req *exampledto.HelloRequest) (string, error) {
	name := req.Name
	if name == "" {
		name = "world"
	}

	return "hello, " + name, nil
}

func Hello(ctx context.Context, req exampledto.HelloRequest) (string, error) {
	name := req.Name
	if name == "" {
		name = "world"
	}

	return "hello, " + name, nil
}
