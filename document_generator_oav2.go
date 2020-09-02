package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
)

type OpenAPIV2SpecGenerator struct {
}

func (o OpenAPIV2SpecGenerator) GenerateAPISpec(w io.Writer, restApi RestAPI) {
	rest := spec.Swagger{}
	rest.BasePath = restApi.BasePath
	rest.Host = restApi.Host
	rest.ID = restApi.ID
	for _, apiResource := range restApi.Resources {
		pathItem := spec.PathItem{}
		for httpMethod, method := range apiResource.methods {
			docMethod := spec.NewOperation("")
			docMethod.Description = method.Description
			docMethod.Summary = method.Summary
			if method.methodOperation.entityOnRequestBody {
				param := spec.BodyParam("body", toSchema(method.methodOperation.entity)).AsRequired()
				docMethod.AddParam(param)
			}
			for _, response := range method.Responses {
				res := spec.NewResponse()
				res.Schema = toSchema(response.Body)
				docMethod.RespondsWith(response.Code, res)
			}
			if httpMethod == "POST" {
				pathItem.Post = docMethod
			}

		}
		rest.Paths.Paths[apiResource.Path] = pathItem
	}
	json.NewEncoder(w).Encode(rest)
}

func toSchema(v interface{}) *spec.Schema {
	val := getValue(v)
	schema := &spec.Schema{}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		schema = spec.Int32Property()
	case reflect.Int64:
		schema = spec.Int64Property()
	case reflect.String:
		schema = spec.StringProperty()
	case reflect.Bool:
		schema = spec.BoolProperty()
	case reflect.Float32:
		schema = spec.Float32Property()
	case reflect.Float64:
		schema = spec.Float64Property()
	case reflect.Array, reflect.Slice:
		schema = spec.ArrayProperty(toSchema(reflect.New(val.Type().Elem()).Interface()))
	case reflect.Struct:
		schema = schema.Typed("object", "")
		structName := val.Type().Name()
		schema.Description = fmt.Sprintf("A %s object.", structName)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			schema.SetProperty(getFieldName(field), *toSchema(val.Field(i).Interface()))
		}
	}
	return schema
}

func getFieldName(field reflect.StructField) string {
	fieldName := field.Name
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
			return jsonTag[:commaIdx]
		}
		return jsonTag
	}
	return fieldName
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}
