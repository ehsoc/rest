package openapiv2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ehsoc/resource"
	"github.com/go-openapi/spec"
)

type OpenAPIV2SpecGenerator struct {
}

func (o OpenAPIV2SpecGenerator) GenerateAPISpec(w io.writer, r resource.RestAPI) string {
	rest := spec.Swagger{}
	rest.BasePath = r.BasePath
	rest.Host = r.Host
	rest.ID = r.ID
	for _, apiResource := range r.Resources {
		pathItem = spec.PathItem{}
		for key, method := range r.methods {
			docMethod := spec.NewOperation()
			docMethod.Description = method.Description
			docMethod.Summary = method.Summary
			if method.methodOperation.entityOnRequestBody {
				param := spec.BodyParam("body", ToSchema(method.methodOperation.entity)).AsRequired()
				docMethod.AddParam(param)
			}
			for _, response := range method.Responses {
				res := spec.NewResponse()
				res.Schema = ToSchema(response.Body)
				docMethod.RespondsWith(response.Code, res)
			}
			doc.Methods[key] = docMethod
		}
		rest.Paths.Paths[apiResource.Path] = pathItem
	}
	json.NewEncoder(w).Encode(rest)
}

// func getSchemaFromType(typ string, items *spec.Schema) *spec.Schema {
// 	schema := &spec.Schema{}
// 	switch typ {
// 	case "int", "int8", "int16", "int32", "uintptr":
// 		schema = spec.Int32Property()
// 	case "int64":
// 		schema = spec.Int64Property()
// 	case "string":
// 		schema = spec.StringProperty()
// 	case "bool":
// 		schema = spec.BoolProperty()
// 	case "float32":
// 		schema = spec.Float32Property()
// 	case "float64":
// 		schema = spec.Float64Property()
// 	case "array":
// 		schema = spec.ArrayProperty(items)
// 	}
// 	return schema
// }

// func StructToSchema(v interface{}) *spec.Schema {
// 	t := reflect.TypeOf(v)
// 	if t.Kind() != reflect.Struct {
// 		return &spec.Schema{}
// 	}
// 	val := getValue(v)
// 	schema := (&spec.Schema{}).Typed("object", "")
// 	schema.Description = fmt.Sprintf("A %s object.", val.Type().Name())
// 	walk(v, func(name string, fieldSchema spec.Schema) {
// 		schema.SetProperty(name, fieldSchema)
// 	})

// 	return schema
// }

func ToSchema(v interface{}) *spec.Schema {
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
		schema = spec.ArrayProperty(ToSchema(reflect.New(val.Type().Elem()).Interface()))
	case reflect.Struct:
		schema = schema.Typed("object", "")
		structName := val.Type().Name()
		schema.Description = fmt.Sprintf("A %s object.", structName)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			schema.SetProperty(getFieldName(field), *ToSchema(val.Field(i).Interface()))
		}
	}
	return schema
}

// func walk(x interface{}, fn func(name string, fieldSchema spec.Schema)) {
// 	val := getValue(x)
// 	walkValue := func(value reflect.Value) {
// 		walk(value.Interface(), fn)
// 	}

// 	switch val.Kind() {
// 	default:
// 		fn(val.Type().Name(), spec.Schema{})
// 	case reflect.Struct:
// 		for i := 0; i < val.NumField(); i++ {
// 			field := val.Type().Field(i)
// 			fieldName := getFieldName(field)
// 			subSchema := getSchema(field)
// 			fn(fieldName, *subSchema)
// 			walkValue(val.Field(i))
// 		}
// 	case reflect.Slice, reflect.Array:
// 		walkValue(val.Elem())
// 	case reflect.Map:
// 		for _, key := range val.MapKeys() {
// 			walkValue(val.MapIndex(key))
// 		}
// 	case reflect.Chan:
// 		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
// 			walk(v.Interface(), fn)
// 		}
// 	}
// }

// func getSchema(field reflect.StructField) *spec.Schema {
// 	if field.Type.Name() != "" {
// 		return getSchemaFromType(field.Type.Name(), nil)
// 	}
// 	switch field.Type.Kind() {
// 	case reflect.Slice, reflect.Array:
// 		sc := getSchemaFromType(field.Type.Elem().Name(), nil)
// 		return getSchemaFromType("array", sc)
// 	}
// 	return nil
// }

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
