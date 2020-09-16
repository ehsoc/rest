package oaiv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/ehsoc/resource"
	"github.com/go-openapi/spec"
)

type OpenAPIV2SpecGenerator struct {
	swagger spec.Swagger
}

func (o *OpenAPIV2SpecGenerator) resolveResource(basePath string, apiResource *resource.Resource) {
	pathItem := spec.PathItem{}
	for httpMethod, method := range apiResource.Methods {
		docMethod := spec.NewOperation("")
		docMethod.Description = method.Description
		docMethod.Summary = method.Summary
		if method.RequestBody.Body != nil {
			param := spec.BodyParam("body", o.toSchema(method.RequestBody.Body)).AsRequired()
			param.SimpleSchema = spec.SimpleSchema{}
			param.Description = method.RequestBody.Description
			docMethod.AddParam(param)
		}
		//Parameters
		//Sorting parameters map for a consistent order in Marshaling
		pKeys := make([]string, 0)
		//URI params will go first
		pURIKeys := make([]string, 0)
		pHeaderKeys := make([]string, 0)
		for key, p := range method.Parameters {
			if p.HTTPType == resource.URIParameter {
				pURIKeys = append(pURIKeys, key)
				continue
			}
			if p.HTTPType == resource.HeaderParameter {
				pHeaderKeys = append(pHeaderKeys, key)
				continue
			}
			pKeys = append(pKeys, key)
		}
		sort.Strings(pHeaderKeys)
		sort.Strings(pURIKeys)
		sort.Strings(pKeys)
		//Append two slices, uri params and the rest
		pHeaderKeys = append(pHeaderKeys, pURIKeys...)
		pKeys = append(pHeaderKeys, pKeys...)
		for _, key := range pKeys {
			parameter := method.Parameters[key]
			specParam := &spec.Parameter{}
			switch parameter.HTTPType {
			case resource.QueryParameter:
				specParam = spec.QueryParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			case resource.URIParameter:
				specParam = spec.PathParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			case resource.HeaderParameter:
				specParam = spec.HeaderParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			case resource.FormDataParameter:
				specParam = spec.FormDataParam(parameter.Name)
				//In case of multipart-form with schema type. This is not supported by OAI V2, this is a workaround.
				if parameter.Body != nil {
					parameter.Type = reflect.String
				}
				typedParam(specParam, parameter.Type)
			case resource.FileParameter:
				specParam = spec.FileParam(parameter.Name)
			}
			specParam.Description = parameter.Description
			specParam.Required = parameter.Required
			docMethod.AddParam(specParam)
		}
		docMethod.Consumes = method.GetDecoderMediaTypes()
		docMethod.Produces = method.GetEncoderMediaTypes()

		//Responses
		for _, response := range method.Responses {
			res := spec.NewResponse()
			if response.Body != nil {
				//res.Schema = o.toRef(response.Body)
				res.Schema = o.toSchema(response.Body)
			}

			//If response.Description is empty we will set a default response base on the status code
			if response.Description != "" {
				res.Description = response.Description
			} else {
				res.Description = http.StatusText(response.Code)
			}

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
	if o.swagger.Paths == nil {
		o.swagger.Paths = &spec.Paths{}
		o.swagger.Paths.Paths = make(map[string]spec.PathItem)
	}

	newBasePath := path.Join(basePath, apiResource.Path)
	o.swagger.Paths.Paths[path.Join(basePath, apiResource.Path)] = pathItem
	for _, apiResource := range apiResource.Resources {
		o.resolveResource(newBasePath, apiResource)
	}
}

func typedParam(param *spec.Parameter, tpe reflect.Kind) {
	schema, err := simpleTypesToSchema(tpe)
	if err != nil {
		log.Println("Warning on processing parameter", param.Name, ":", err)
	}
	if schema != nil && len(schema.Type[0]) > 0 {
		param.Typed(schema.Type[0], schema.Format)
	}
}

func (o *OpenAPIV2SpecGenerator) GenerateAPISpec(w io.Writer, restApi resource.RestAPI) {
	o.swagger.Swagger = "2.0"
	o.swagger.BasePath = restApi.BasePath
	o.swagger.Host = restApi.Host
	o.swagger.ID = restApi.ID
	for _, apiResource := range restApi.Resources {
		o.resolveResource("", apiResource)
	}
	json.NewEncoder(w).Encode(o.swagger)
}

func (o *OpenAPIV2SpecGenerator) AddDefinition(name string, schema *spec.Schema) {
	if o.swagger.Definitions == nil {
		o.swagger.Definitions = make(spec.Definitions)
	}
	o.swagger.Definitions[name] = *schema
}

func (o *OpenAPIV2SpecGenerator) toSchema(v interface{}) *spec.Schema {
	val := getValue(v)
	schema := &spec.Schema{}
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Struct {
		}
		schema = spec.ArrayProperty(o.toSchema(reflect.New(val.Type().Elem()).Interface()))
	case reflect.Struct:
		structName := val.Type().Name()
		refSchema := &spec.Schema{}
		refSchema = refSchema.Typed("object", "")
		refSchema.Description = fmt.Sprintf("A %s object.", structName)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			refSchema.SetProperty(getFieldName(field), *o.toSchema(val.Field(i).Interface()))
		}
		o.AddDefinition(structName, refSchema)
		schema = spec.RefSchema("#/definitions/" + structName)
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
