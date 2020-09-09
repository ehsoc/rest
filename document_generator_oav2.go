package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
)

type OpenAPIV2SpecGenerator struct {
	rest spec.Swagger
}

func (o *OpenAPIV2SpecGenerator) resolveResource(basePath string, apiResource Resource) {
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
			docMethod.Consumes = method.getMediaTypes()
		}
		//Parameters
		//Sorting parameters map for a consistent order in Marshaling
		pKeys := make([]string, 0)
		for key := range method.Parameters {
			pKeys = append(pKeys, key)
		}
		sort.Strings(pKeys)
		for _, key := range pKeys {
			parameter := method.Parameters[key]
			specParam := &spec.Parameter{}
			switch parameter.HTTPType {
			case BodyParameter:
				schema := o.toRef(parameter.Body)
				specParam = spec.BodyParam(parameter.Name, schema)
			case URIParameter:
				specParam = spec.PathParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			case HeaderParameter:
				specParam = spec.HeaderParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			}
			specParam.Description = parameter.Description
			specParam.Required = parameter.Required
			docMethod.AddParam(specParam)
		}

		docMethod.Produces = method.getMediaTypes()
		//Responses
		for _, response := range method.Responses {
			res := spec.NewResponse()
			if response.Body != nil {
				if reflect.TypeOf(response.Body) == reflect.TypeOf(method.methodOperation.entity) {
					res.Schema = o.toRef(response.Body)
				} else {
					res.Schema = toSchema(response.Body)
				}

			}
			res.Description = http.StatusText(response.Code)
			docMethod.RespondsWith(response.Code, res)
		}
		docMethod.Responses.Default = nil

		switch httpMethod {
		case http.MethodPost:
			pathItem.Post = docMethod
		case http.MethodGet:
			pathItem.Get = docMethod
		case http.MethodDelete:
			pathItem.Delete = docMethod
		}

	}
	if o.rest.Paths == nil {
		o.rest.Paths = &spec.Paths{}
		o.rest.Paths.Paths = make(map[string]spec.PathItem)
	}
	newBasePath := path.Join(basePath, apiResource.Path)
	o.rest.Paths.Paths[path.Join(basePath, apiResource.Path)] = pathItem
	for _, apiResource := range apiResource.Resources {
		o.resolveResource(newBasePath, apiResource)
	}
}

func typedParam(param *spec.Parameter, tpe reflect.Kind) {
	schema, err := simpleTypesToSchema(tpe)
	if err != nil {
		fmt.Println(err)
	}
	if len(schema.Type[0]) > 0 {
		param.Typed(schema.Type[0], schema.Format)
	}
}

func (o *OpenAPIV2SpecGenerator) GenerateAPISpec(w io.Writer, restApi RestAPI) {
	o.rest.BasePath = restApi.BasePath
	o.rest.Host = restApi.Host
	o.rest.ID = restApi.ID
	for _, apiResource := range restApi.Resources {
		o.resolveResource("", apiResource)
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
	case reflect.Array, reflect.Slice:
		schema = spec.ArrayProperty(toSchema(reflect.New(val.Type().Elem()).Interface()))
	case reflect.Struct:
		structName := val.Type().Name()
		schema = spec.RefSchema("#/definitions/" + structName)
		o.AddDefinition(structName, toSchema(v))
	default:
		schema, _ = simpleTypesToSchema(val.Kind())
	}
	return schema
}

func toSchema(v interface{}) *spec.Schema {
	val := getValue(v)
	schema := &spec.Schema{}
	switch val.Kind() {
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
	default:
		schema, _ = simpleTypesToSchema(val.Kind())
	}
	return schema
}

func simpleTypesToSchema(kind reflect.Kind) (*spec.Schema, error) {
	schema := &spec.Schema{}
	switch kind {
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
	case reflect.Array, reflect.Slice, reflect.Struct:
		return nil, errors.New("kind is a complex type, use toSchema function instead")
	}
	return schema, nil
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
