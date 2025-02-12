package openapi

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/neatframework/router/register"
	"github.com/neatframework/router/register/internal/stringx"
	"github.com/neatframework/router/register/openapi/tags"
)

var resolved = map[string]bool{}

var Info = &openapi3.Info{
	TermsOfService: "github.com/neatframework/router/register",
	Title:          "Gofr",
	Description:    "Gofr API Documentation",
	Version:        "1.0",
}

func Print() string {
	t := &openapi3.T{
		Extensions: nil,
		OpenAPI:    "3.0.1",
		Components: new(openapi3.Components),
		Info:       Info,
		Paths:      openapi3.NewPaths(),
	}
	t.Components.Schemas = make(map[string]*openapi3.SchemaRef)
	t.Components.Schemas["any"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: "Can be anything: string, number, array, object, etc., including `null`",
		},
	}

	t.Components.Schemas["null"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: "null type",
		},
	}

	for _, route := range register.Routes {
		ResolveRoute(t, route)
	}

	content, err := t.MarshalJSON()
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})

	_ = json.Unmarshal(content, &m)
	content, err = json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(content)
}

func ResolveRoute(t *openapi3.T, route register.OpenRoute) {
	operation := new(openapi3.Operation)

	key := route.HttpMethod() + "." + route.HttpPath()

	if _, ok := resolved[strings.ToLower(key)]; ok {
		return
	} else {
		resolved[strings.ToLower(key)] = true
	}

	resolveRouteRequest(t, operation, route)
	resolveRouteResponse(operation, route)
}

type Response[T any] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
}

// fillRefValue 递归填充ref的值
func fillRefValue(ref *openapi3.SchemaRef, typ reflect.Type) {
	for name, ref := range ref.Value.Properties {

		if field, ok := findField(typ, name); ok {
			ref.Value.Title = field.Tag.Get("title")
			if field.Type.Kind() == reflect.Ptr {
				ref.Value.Nullable = true
			}

			if ref.Value.Type.Is("array") {
				fillRefValue(ref.Value.Items, field.Type.Elem())
			}

			if ref.Value.Type.Is("object") {
				fillRefValue(ref, field.Type)
			}
		}
	}
}

func deepClone(ref *openapi3.SchemaRef) *openapi3.SchemaRef {

	marshalJSON, err := ref.MarshalJSON()
	if err != nil {
		panic(err)
	}

	schemaClone := &openapi3.SchemaRef{}
	err = schemaClone.UnmarshalJSON(marshalJSON)
	if err != nil {
		panic(err)
	}

	return schemaClone
}

func resolveRouteResponse(operation *openapi3.Operation, route register.OpenRoute) {
	scheam, err := openapi3gen.NewSchemaRefForValue(&Response[any]{}, openapi3.Schemas{})
	if err != nil {
		panic(err)
	}

	data, err := scheam.MarshalJSON()
	if err != nil {
		panic(err)
	}

	desc := ""

	if route.ResponseType() != nil {
		var dataTypeString string

		if route.ResponseType().Kind() == reflect.Map {
			// nothing
			dataTypeString = `{"type":"object","additionalProperties":{"$ref":"#/components/schemas/any"}}`
		} else {
			value, err := openapi3gen.NewSchemaRefForValue(reflect.New(route.ResponseType()).Interface(), openapi3.Schemas{})

			if err != nil {
				panic(err)
			}

			// 通过 NewSchemaRefForValue 生成的schema，相同的类型会引用 同一个SchemaRef对象，
			// 包括 string, int 基础类型，这样修改一个string对象里面的值，就相当于修改所有string对象
			// 使用 json 重新生成对象，就不是相互引用关系
			schemaClone := deepClone(value)

			fillRefValue(schemaClone, route.ResponseType())

			if route.ResponseType().Kind() == reflect.Ptr {
				schemaClone.Value.Nullable = true
			}

			marshalJSON, err := schemaClone.MarshalJSON()
			if err != nil {
				panic(err)
			}

			dataTypeString = string(marshalJSON)
		}

		err = scheam.UnmarshalJSON([]byte(strings.ReplaceAll(string(data), "\"data\":{}", "\"data\":"+string(dataTypeString))))
		if err != nil {
			panic(err)
		}

		desc = "响应描述"
	}

	operation.AddResponse(200, &openapi3.Response{
		Description: &desc,
		Content: map[string]*openapi3.MediaType{
			"application/json": {
				Schema: scheam,
			},
		},
		Links: nil,
	})
}

func resolveTag(t *openapi3.T, route register.OpenRoute) *openapi3.Tag {
	var tag string

	tag = route.Object().Name()
	tag = t.Info.Title + "/" + tag

	for _, t2 := range t.Tags {
		if t2.Name == tag {
			return t2
		}
	}

	tagE := &openapi3.Tag{
		Extensions: map[string]interface{}{},
		Name:       tag,
	}
	t.Tags = append(t.Tags, tagE)

	return tagE
}

func simpleTypeMap(t reflect.Type) string {
	kind := t.Kind()
	if kind == reflect.Ptr {
		kind = t.Elem().Kind()
	}

	var typ string
	switch kind {
	case reflect.Bool:
		typ = "boolean"

	case reflect.Int:
		typ = "integer"
	case reflect.Int8:
		typ = "integer"
	case reflect.Int16:
		typ = "integer"
	case reflect.Int32:
		typ = "integer"
	case reflect.Int64:
		typ = "long"
	case reflect.Uint:
		typ = "integer"
	case reflect.Uint8:
		typ = "integer"
	case reflect.Uint16:
		typ = "integer"
	case reflect.Uint32:
		typ = "integer"
	case reflect.Uint64:
		typ = "long"
	case reflect.Float32:
		typ = "float"
	case reflect.Float64:
		typ = "double"
	case reflect.String:
		typ = "string"
	case reflect.Slice:
		typ = "array"
	default:
		panic("error simple type" + kind.String())
	}
	return typ
}

func resolveRouteRequest(t *openapi3.T, operation *openapi3.Operation, route register.OpenRoute) {
	path := route.HttpPath()

	field, _ := route.Object().FieldByName(stringx.UcFirst(route.MethodName()))

	if summary, ok := field.Tag.Lookup(tags.TagSummary); ok {
		operation.Summary = summary
	}

	t.Paths.Set(path, &openapi3.PathItem{
		Summary:     "",
		Description: "",
	})
	t.Paths.Find(path).SetOperation(route.HttpMethod(), operation)

	tag := resolveTag(t, route)

	operation.OperationID = route.HttpMethod() + "_" + route.HttpPath()
	operation.Tags = append(operation.Tags, tag.Name)

	if route.RequestType() != nil {
		if LikeGet(route.HttpMethod()) {
			requestType := route.RequestType()
			if requestType.Kind() == reflect.Map {
				// nothing
			} else {
				if requestType.Kind() == reflect.Ptr {
					requestType = requestType.Elem()
				}

				for i := 0; i < requestType.NumField(); i++ {
					parameter := new(openapi3.Parameter)
					parameter.Extensions = map[string]interface{}{}

					field := requestType.Field(i)

					parameter.Name = fieldName(field)
					parameter.Description = field.Tag.Get("title")
					parameter.Example = field.Tag.Get("examples")
					parameter.In = "query"

					parameter.Extensions["x-lang-type"] = simpleTypeMap(field.Type)

					if field.Tag.Get("required") == "true" {
						parameter.Required = true
					}
					operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
						Value: parameter,
					})
				}
			}
		} else {
			scheam, err := openapi3gen.NewSchemaRefForValue(reflect.New(route.RequestType()).Interface(), openapi3.Schemas{})

			if err != nil {
				panic(err)
			}

			// 通过 NewSchemaRefForValue 生成的schema，相同的类型会引用 同一个SchemaRef对象，
			// 包括 string, int 基础类型，这样修改一个string对象里面的值，就相当于修改所有string对象
			// 使用 json 重新生成对象，就不是相互引用关系
			schemaClone := deepClone(scheam)
			fillRefValue(schemaClone, route.RequestType())
			if route.RequestType().Kind() == reflect.Ptr {
				schemaClone.Value.Nullable = true
			}

			operation.RequestBody = &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Extensions:  nil,
					Description: "请求说明",
					Required:    false,
					Content: map[string]*openapi3.MediaType{
						"application/json": {
							Schema: schemaClone,
						},
					},
				},
			}
		}
	}
}

func findField(typ reflect.Type, field string) (reflect.StructField, bool) {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		if fieldName(typ.Field(i)) == field {
			return typ.Field(i), true
		}
	}
	return reflect.StructField{}, false
}

func fieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	name := ""
	if jsonTag != "-" {
		for i, part := range strings.Split(jsonTag, ",") {
			if i == 0 {
				if part != "" {
					name = part
				}
			}
		}
	}

	if name == "" {
		name = field.Name
	}

	return name
}

func LikeGet(method string) bool {
	return method != http.MethodPost
}
