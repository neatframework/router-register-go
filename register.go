package register

import (
	"log"
	"path"
	"reflect"
	"strings"

	"github.com/neatframework/router/register/internal/stringx"
)

var Routes = map[string]*Route{}

func Service(handlerType interface{}, impl interface{}) []OpenRoute {
	ht := reflect.TypeOf(handlerType).Elem()

	if !reflect.TypeOf(impl).AssignableTo(ht) {
		log.Fatalf("Server.Service found the handler of type %v that does not satisfy %v", reflect.TypeOf(impl), ht)
	}

	st := reflect.ValueOf(impl)

	var serviceRoutes []OpenRoute

	for i := 0; i < ht.NumField(); i++ {
		field := ht.Field(i)

		httpTagMap := map[string]string{
			TagGetMapping:  "GET",
			TagPostMapping: "POST",
			TagPutMapping:  "PUT",
			TagAnyMapping:  "ANY",
		}

		routeKey := st.Type().PkgPath() + "." + stringx.SnakeCase(st.Type().Name()) + "@" + stringx.LcFirst(field.Name)

		route := NewRoute(routeKey, st.Type(), field.Name, st.FieldByName(field.Name))
		for t, m := range httpTagMap {
			if value, ok := field.Tag.Lookup(t); ok {
				route.httpMethod = m
				route.httpPath = "/" + strings.TrimLeft(path.Join(config.BaseUri, value), "/")
			}
		}

		serviceRoutes = append(serviceRoutes, route)
		Routes[routeKey] = route
	}

	return serviceRoutes
}
