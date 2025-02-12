package ginregister

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/neatframework/router/register"
)

var e *gin.Engine

type BaseResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

var renderError = func(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

var renderResponse = func(ctx *gin.Context, data interface{}) {
	var baseResponse BaseResponse
	baseResponse.Code = http.StatusOK
	baseResponse.Data = data

	ctx.JSON(http.StatusOK, baseResponse)
}

func SetGinEngine(engine *gin.Engine) {
	e = engine
}

func SetRenderError(f func(ctx *gin.Context, err error)) {
	renderError = f
}

func SetRenderResponse(f func(ctx *gin.Context, data interface{})) {
	renderResponse = f
}

func Service(handlerType interface{}, impl interface{}) {
	routes := register.Service(handlerType, impl)
	for _, route := range routes {
		e.Handle(route.HttpMethod(), route.HttpPath(), func(ctx *gin.Context) {
			var parameters []reflect.Value
			parameters = []reflect.Value{
				reflect.ValueOf(ctx),
			}

			if route.RequestType() != nil {
				reqTyp := route.RequestType()
				if reqTyp.Kind() == reflect.Ptr {
					reqTyp = reqTyp.Elem()
				}

				object := reflect.New(reqTyp).Interface()

				var err error

				if route.HttpMethod() == http.MethodPost {
					err = ctx.Bind(object)
				} else {
					err = ctx.BindQuery(object)
				}

				if err != nil {
					renderError(ctx, err)
					return
				}

				if object == nil {
					parameters = append(parameters, reflect.New(route.RequestType().Elem()))
				} else {
					if route.RequestType().Kind() == reflect.Ptr {
						parameters = append(parameters, reflect.ValueOf(object))
					} else {
						parameters = append(parameters, reflect.Indirect(reflect.ValueOf(object)))
					}
				}
			}

			// 调用接口
			res := route.Action().Call(parameters)

			if err, ok := res[1].Interface().(error); ok {
				renderError(ctx, err)
			} else {
				renderResponse(ctx, res[0].Interface())
			}
		})
	}
}
