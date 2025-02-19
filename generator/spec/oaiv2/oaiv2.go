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
	"time"

	"github.com/ehsoc/rest"
	"github.com/go-openapi/spec"
)

type OpenAPIV2SpecGenerator struct {
	swagger spec.Swagger
}

func (o *OpenAPIV2SpecGenerator) resolveResource(basePath string, apiResource rest.Resource) {
	pathItem := spec.PathItem{}
	for _, method := range apiResource.Methods() {
		specMethod := spec.NewOperation("")
		specMethod.Description = method.Description
		specMethod.Summary = method.Summary

		if method.RequestBody.Body != nil {
			param := spec.BodyParam("body", o.toSchema(method.RequestBody.Body)).AsRequired()
			param.SimpleSchema = spec.SimpleSchema{}
			param.Description = method.RequestBody.Description
			specMethod.AddParam(param)
		}
		// Parameters
		// Sorting parameters map for a consistent order in Marshaling
		pKeys := make([]rest.Parameter, 0)
		// URI params will go first
		pURIKeys := make([]rest.Parameter, 0)
		pHeaderKeys := make([]rest.Parameter, 0)

		for _, p := range method.Parameters() {
			if p.HTTPType == rest.URIParameter {
				pURIKeys = append(pURIKeys, p)
				continue
			}

			if p.HTTPType == rest.HeaderParameter {
				pHeaderKeys = append(pHeaderKeys, p)
				continue
			}

			pKeys = append(pKeys, p)
		}

		sort.Slice(pHeaderKeys, func(i, j int) bool {
			return pHeaderKeys[i].Name < pHeaderKeys[j].Name
		})

		sort.Slice(pURIKeys, func(i, j int) bool {
			return pURIKeys[i].Name < pURIKeys[j].Name
		})

		sort.Slice(pKeys, func(i, j int) bool {
			return pKeys[i].Name < pKeys[j].Name
		})

		// Append two slices, uri params and the rest
		pHeaderKeys = append(pHeaderKeys, pURIKeys...)
		pKeys = append(pHeaderKeys, pKeys...)
		for _, parameter := range pKeys {
			specParam := &spec.Parameter{}

			switch parameter.HTTPType {
			case rest.QueryParameter:
				specParam = spec.QueryParam(parameter.Name)
				if parameter.Type == reflect.Array {
					specParam.Type = "array"
					specParam.Items = spec.NewItems()
					if parameter.EnumValues != nil && len(parameter.EnumValues) > 0 {
						specParam.Items.WithEnum(parameter.EnumValues...).WithDefault(parameter.EnumValues[0])
						specParam.CollectionFormat = parameter.CollectionFormat

						ssch := o.toSchema(parameter.EnumValues[0])
						if len(ssch.SchemaProps.Type) > 0 {
							specParam.Items.Type = ssch.SchemaProps.Type[0]
						}
					}
				} else {
					typedParam(specParam, parameter.Type)
				}
			case rest.URIParameter:
				specParam = spec.PathParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			case rest.HeaderParameter:
				specParam = spec.HeaderParam(parameter.Name)
				typedParam(specParam, parameter.Type)
			case rest.FormDataParameter:
				specParam = spec.FormDataParam(parameter.Name)
				// In case of multipart-form with schema type. This is not supported by OAI V2, this is a workaround.
				if parameter.Body != nil {
					parameter.Type = reflect.String
				}
				typedParam(specParam, parameter.Type)
			case rest.FileParameter:
				specParam = spec.FileParam(parameter.Name)
			}

			specParam.Description = parameter.Description
			specParam.Required = parameter.Required
			// Example on parameters is not allowed, so a extension is set.
			if parameter.Example != nil {
				specParam.AddExtension("x-example", parameter.Example)
			}

			specMethod.AddParam(specParam)
		}
		specMethod.Consumes = method.GetDecoderMediaTypes()
		specMethod.Produces = method.GetEncoderMediaTypes()
		// Security
		for _, security := range method.SecurityCollection {
			secSchemes := map[string][]string{}
			for _, securityScheme := range security.SecuritySchemes {
				switch securityScheme.Type {
				case rest.BasicSecurityType:
					secScheme := spec.BasicAuth()
					// Add to secSchemes map
					secSchemes[securityScheme.Name] = []string{}
					o.addSecurityDefinition(securityScheme.Name, secScheme)
				case rest.APIKeySecurityType:
					secParam := convertParameter(securityScheme.Parameter)
					secScheme := spec.APIKeyAuth(securityScheme.Name, secParam.In)
					// Add to secSchemes map
					secSchemes[securityScheme.Name] = []string{}
					o.addSecurityDefinition(secScheme.Name, secScheme)
				case rest.OAuth2SecurityType:
					if securityScheme.OAuth2Flows != nil {
						// OpenAPI v2 doesn't support multiple flows, so will create a oauth scheme per flow
						for k, flow := range securityScheme.OAuth2Flows {
							secScheme := getOAuth2SecScheme(flow)
							if k != 0 {
								securityScheme.Name += "_" + secScheme.Flow
							}
							scopes := []string{}
							for scp := range secScheme.Scopes {
								scopes = append(scopes, scp)
							}
							sort.Strings(scopes)
							// Add to secSchemes map
							secSchemes[securityScheme.Name] = scopes
							o.addSecurityDefinition(securityScheme.Name, secScheme)
						}
					}
				}
			}
			specMethod.Security = append(specMethod.Security, secSchemes)
		}
		// Responses
		for _, response := range method.Responses() {
			res := spec.NewResponse()
			// Body() returns an interfaces so can be nil
			if response.Body() != nil {
				res.Schema = o.toSchema(response.Body())
			}
			// If response.Description is empty we will set a default response base on the status code
			if response.Description() != "" {
				res.Description = response.Description()
			} else {
				res.Description = http.StatusText(response.Code())
			}
			specMethod.RespondsWith(response.Code(), res)
			specMethod.Responses.Default = nil
		}

		switch method.HTTPMethod {
		case http.MethodPost:
			pathItem.Post = specMethod
		case http.MethodPut:
			pathItem.Put = specMethod
		case http.MethodGet:
			pathItem.Get = specMethod
		case http.MethodDelete:
			pathItem.Delete = specMethod
		}
	}
	if o.swagger.Paths == nil {
		o.swagger.Paths = &spec.Paths{}
		o.swagger.Paths.Paths = make(map[string]spec.PathItem)
	}

	newBasePath := path.Join(basePath, apiResource.Path())
	// Only Paths with methods should be in the Paths map
	if len(apiResource.Methods()) > 0 {
		o.swagger.Paths.Paths[newBasePath] = pathItem
	}

	for _, apiResource := range apiResource.Resources() {
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

func convertParameter(parameter rest.Parameter) *spec.Parameter {
	specParam := &spec.Parameter{}

	switch parameter.HTTPType {
	case rest.QueryParameter:
		specParam = spec.QueryParam(parameter.Name)

		if parameter.Type == reflect.Array {
			specParam.Type = "array"
			specParam.Items = spec.NewItems()

			if parameter.EnumValues != nil && len(parameter.EnumValues) > 0 {
				specParam.Items.WithEnum(parameter.EnumValues...).WithDefault(parameter.EnumValues[0])
				specParam.CollectionFormat = parameter.CollectionFormat
			}
		} else {
			typedParam(specParam, parameter.Type)
		}

	case rest.URIParameter:
		specParam = spec.PathParam(parameter.Name)
		typedParam(specParam, parameter.Type)
	case rest.HeaderParameter:
		specParam = spec.HeaderParam(parameter.Name)
		typedParam(specParam, parameter.Type)
	case rest.FormDataParameter:
		specParam = spec.FormDataParam(parameter.Name)
		// In case of multipart-form with schema type. This is not supported by OAI V2, this is a workaround.
		if parameter.Body != nil {
			parameter.Type = reflect.String
		}
		typedParam(specParam, parameter.Type)
	case rest.FileParameter:
		specParam = spec.FileParam(parameter.Name)
	}

	specParam.Description = parameter.Description
	specParam.Required = parameter.Required
	// Example on parameters is not allowed, so a extension is set.
	if parameter.Example != nil {
		specParam.AddExtension("x-example", parameter.Example)
	}

	return specParam
}

func (o *OpenAPIV2SpecGenerator) GenerateAPISpec(w io.Writer, api rest.API) {
	o.swagger.Swagger = "2.0"
	o.swagger.BasePath = api.BasePath
	o.swagger.Host = api.Host
	o.swagger.ID = api.ID
	info := &spec.Info{}
	info.Description = api.Description
	info.Title = api.Title
	info.Version = api.Version
	o.swagger.Info = info

	for _, apiResource := range api.Resources() {
		o.resolveResource("/", apiResource)
	}

	e := json.NewEncoder(w)

	e.SetIndent(" ", "  ")
	e.Encode(o.swagger)
}

func getOAuth2SecScheme(flow rest.OAuth2Flow) *spec.SecurityScheme {
	secScheme := getBaseFlow(flow)
	secScheme.Scopes = flow.Scopes

	return secScheme
}

func getBaseFlow(flow rest.OAuth2Flow) *spec.SecurityScheme {
	switch flow.Name {
	case rest.FlowImplicitType:
		return spec.OAuth2Implicit(flow.AuthorizationURL)
	case rest.FlowPasswordType:
		return spec.OAuth2Password(flow.TokenURL)
	case rest.FlowAuthCodeType:
		return spec.OAuth2AccessToken(flow.AuthorizationURL, flow.TokenURL)
	case rest.FlowClientCredentialType:
		return spec.OAuth2Application(flow.TokenURL)
	default:
		return &spec.SecurityScheme{SecuritySchemeProps: spec.SecuritySchemeProps{
			Type:             "oauth2",
			Flow:             flow.Name,
			TokenURL:         flow.TokenURL,
			AuthorizationURL: flow.AuthorizationURL,
		}}
	}
}

func (o *OpenAPIV2SpecGenerator) addDefinition(name string, schema *spec.Schema) {
	if o.swagger.Definitions == nil {
		o.swagger.Definitions = make(spec.Definitions)
	}

	o.swagger.Definitions[name] = *schema
}

func (o *OpenAPIV2SpecGenerator) addSecurityDefinition(name string, schema *spec.SecurityScheme) {
	if o.swagger.SecurityDefinitions == nil {
		o.swagger.SecurityDefinitions = make(spec.SecurityDefinitions)
	}

	o.swagger.SecurityDefinitions[name] = schema
}

func (o *OpenAPIV2SpecGenerator) toSchema(v interface{}) *spec.Schema {
	val := getValue(v)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		return spec.ArrayProperty(o.toSchema(reflect.New(val.Type().Elem()).Interface()))
	case reflect.Struct:
		if _, ok := v.(time.Time); ok {
			return spec.DateTimeProperty()
		}
		structName := val.Type().Name()
		refSchema := &spec.Schema{}
		refSchema = refSchema.Typed("object", "")
		refSchema.Description = fmt.Sprintf("A %s object.", structName)

		for i := 0; i < val.NumField(); i++ {
			// Avoiding panic on unexported fields
			if val.Field(i).CanInterface() {
				field := val.Type().Field(i)
				refSchema.SetProperty(getFieldName(field), *o.toSchema(val.Field(i).Interface()))
			}
		}
		o.addDefinition(structName, refSchema)

		return spec.RefSchema("#/definitions/" + structName)
	default:
		schema, _ := simpleTypesToSchema(val.Kind())
		return schema
	}
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
