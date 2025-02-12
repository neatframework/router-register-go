package facade

import (
	"context"

	"github.com/neatframework/router/register/examples/ginproject/facade/exampledto"
)

type Example struct {
	Hello    func(ctx context.Context, req exampledto.HelloRequest) (string, error)  `get_mapping:"/example/hello"`
	HelloPtr func(ctx context.Context, req *exampledto.HelloRequest) (string, error) `get_mapping:"/example/hello_ptr"`
}
