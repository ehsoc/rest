package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
)

type OpenAPIV2SpecGenerator struct {
	rest spec.Swagger
}

func (o *OpenAPIV2SpecGenerator) GenerateAPISpec(w io.Writer, restApi RestAPI) {
	o.rest.BasePath = restApi.BasePath
	o.rest.Host = restApi.Host
	o.rest.ID = restApi.ID
	for _, apiResource := range restApi.Resources {
		pathItem := spec.PathItem{}
		for httpMethod, method := range apiResource.methods {
			docMethod := spec.NewOperation("")
			docMethod.Description = method.Description
			docMethod.Summary = method.Summary
			if method.methodOperation.entityOnRequestBody {
				param := spec.BodyParam("body", o.toRef(method.Request.GetBody())).AsRequired()
				param.SimpleSchema = spec.SimpleSchema{}
				param.Description = method.Request.Description
				docMethod.AddParam(param)
			}
			docMethod.Consumes = method.getMediaTypes()
			docMethod.Produces = method.getMediaTypes()
			for _, response := range method.Responses {
				res := spec.NewResponse()
				if response.Body != nil {
					res.Schema = toSchema(response.Body)
				}
				res.Description = http.StatusText(response.Code)
				docMethod.RespondsWith(response.Code, res)
			}
			docMethod.Responses.Default = nil
			if httpMethod == "POST" {
				pathItem.Post = docMethod
			}
		}
		if o.rest.Paths == nil {
			o.rest.Paths = &spec.Paths{}
			o.rest.Paths.Paths = make(map[string]spec.PathItem)
		}
		o.rest.Paths.Paths[apiResource.Path] = pathItem
	}
	json.NewEncoder(w).Encode(o.rest)
}

func (o *OpenAPIV2SpecGenerator) AddDefinition(name string, schema *spec.Schema) {
	if o.rest.Definitions != nil {

	}
	o.rest.Definitions = make(spec.Definitions)
	o.rest.Definitions[name] = *schema
}

func (o *OpenAPIV2SpecGenerator) toRef(v interface{}) *spec.Schema {
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
		structName := val.Type().Name()
		schema = spec.RefSchema("#/definitions/" + structName)
		o.AddDefinition(structName, toSchema(v))
	}
	return schema
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
