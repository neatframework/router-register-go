package register

import (
	"context"
	"reflect"
)

type OpenRoute interface {
	Object() reflect.Type
	Action() reflect.Value
	MethodName() string

	ResponseType() reflect.Type
	RequestType() reflect.Type
	HttpPath() string
	HttpMethod() string
}

var config = &Config{}

type Config struct {
	BaseUri string
}

func SetBaseUri(baseUri string) {
	config.BaseUri = baseUri
}

type Route struct {
	httpPath   string
	httpMethod string
	key        string

	requestType  reflect.Type
	responseType reflect.Type
	methodName   string
	object       reflect.Type
	action       reflect.Value
}

func (r Route) Action() reflect.Value {
	return r.action
}

func (r Route) RequestType() reflect.Type {
	return r.requestType
}

func (r Route) ResponseType() reflect.Type {
	return r.responseType
}

func (r Route) MethodName() string {
	return r.methodName
}

func (r Route) Object() reflect.Type {
	return r.object
}

func (r Route) HttpMethod() string {
	return r.httpMethod
}

func (r Route) HttpPath() string {
	return r.httpPath
}

func NewRoute(routeKey string, interfaceType reflect.Type, methodName string, action reflect.Value) *Route {
	var requestType reflect.Type

	if !action.Type().In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		panic("method first parameter must be context.Context: " + action.String())
	}

	if action.Type().NumIn() == 2 {
		requestType = action.Type().In(1)
	} else {
		requestType = nil
	}

	return &Route{
		key:          routeKey,
		object:       interfaceType,
		action:       action,
		methodName:   methodName,
		responseType: action.Type().Out(0),
		requestType:  requestType,
	}
}
