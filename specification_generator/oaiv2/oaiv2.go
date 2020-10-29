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

func (o *OpenAPIV2SpecGenerator) resolveResource(basePath string, apiResource resource.Resource) {
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
		//Parameters
		//Sorting parameters map for a consistent order in Marshaling
		pKeys := make([]resource.Parameter, 0)
		//URI params will go first
		pURIKeys := make([]resource.Parameter, 0)
		pHeaderKeys := make([]resource.Parameter, 0)
		for _, p := range method.GetParameters() {
			if p.HTTPType == resource.URIParameter {
				pURIKeys = append(pURIKeys, p)
				continue
			}
			if p.HTTPType == resource.HeaderParameter {
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
		//Append two slices, uri params and the rest
		pHeaderKeys = append(pHeaderKeys, pURIKeys...)
		pKeys = append(pHeaderKeys, pKeys...)
		for _, parameter := range pKeys {
			specParam := &spec.Parameter{}
			switch parameter.HTTPType {
			case resource.QueryParameter:
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
			//Example on parameters is not allowed, so a extension is set.
			if parameter.Example != nil {
				specParam.AddExtension("x-example", parameter.Example)
			}
			specMethod.AddParam(specParam)
		}
		specMethod.Consumes = method.GetDecoderMediaTypes()
		specMethod.Produces = method.GetEncoderMediaTypes()
		//Security
		for _, security := range method.SecuritySchemes {
			switch security.Type {
			case resource.BasicSecurityType:
				secScheme := spec.BasicAuth()
				specMethod.SecuredWith(security.Name, []string{}...)
				o.addSecurityDefinition(security.Name, secScheme)
			case resource.ApiKeySecurityType:
				params := security.Parameters.GetParameters()
				if len(params) > 0 {
					secParam := convertParameter(params[0])
					secScheme := spec.APIKeyAuth(security.Name, secParam.In)
					specMethod.SecuredWith(secScheme.Name, []string{}...)
					o.addSecurityDefinition(secScheme.Name, secScheme)
				}
			case resource.OAuth2SecurityType:
				if security.OAuth2Flows != nil {
					//OpenAPI v2 doesn't support multiple flows, so will create a oauth scheme per flow
					for k, flow := range security.OAuth2Flows {
						secScheme := getOAuth2SecScheme(flow)
						if k != 0 {
							security.Name += "_" + secScheme.Flow
						}
						scopes := []string{}
						for scp := range secScheme.Scopes {
							scopes = append(scopes, scp)
						}
						sort.Strings(scopes)
						specMethod.SecuredWith(security.Name, scopes...)
						o.addSecurityDefinition(security.Name, secScheme)
					}
				}
			}
		}
		//Responses
		for _, response := range method.GetResponses() {
			res := spec.NewResponse()
			//Body()  returns an interfaces so can be nil
			if response.Body() != nil {
				//res.Schema = o.toRef(response.Body)
				res.Schema = o.toSchema(response.Body())
			}
			//If response.Description is empty we will set a default response base on the status code
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
	//Only Paths with methods should be in the Paths map
	if len(apiResource.Methods()) > 0 {
		o.swagger.Paths.Paths[newBasePath] = pathItem
	}

	for _, apiResource := range apiResource.GetResources() {
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

func convertParameter(parameter resource.Parameter) *spec.Parameter {
	specParam := &spec.Parameter{}
	switch parameter.HTTPType {
	case resource.QueryParameter:
		specParam = spec.QueryParam(parameter.Name)
		if parameter.Type == reflect.Array {
			specParam.Type = "array"
			specParam.Items = spec.NewItems()
			if parameter.EnumValues != nil && len(parameter.EnumValues) > 0 {
				specParam.Items.WithEnum(parameter.EnumValues...).WithDefault(parameter.EnumValues[0])
				specParam.CollectionFormat = parameter.CollectionFormat
				// ssch := o.toSchema(parameter.EnumValues[0])
				// if len(ssch.SchemaProps.Type) > 0 {
				// 	specParam.Items.Type = ssch.SchemaProps.Type[0]
				// }
			}
		} else {
			typedParam(specParam, parameter.Type)
		}
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
	//Example on parameters is not allowed, so a extension is set.
	if parameter.Example != nil {
		specParam.AddExtension("x-example", parameter.Example)
	}
	return specParam
}

func (o *OpenAPIV2SpecGenerator) GenerateAPISpec(w io.Writer, restApi resource.RestAPI) {
	o.swagger.Swagger = "2.0"
	o.swagger.BasePath = restApi.BasePath
	o.swagger.Host = restApi.Host
	o.swagger.ID = restApi.ID
	info := &spec.Info{}
	info.Description = restApi.Description
	info.Title = restApi.Title
	info.Version = restApi.Version
	o.swagger.Info = info
	for _, apiResource := range restApi.GetResources() {
		o.resolveResource("/", apiResource)
	}
	e := json.NewEncoder(w)
	e.SetIndent(" ", "  ")
	e.Encode(o.swagger)
}

func getOAuth2SecScheme(flow resource.OAuth2Flow) *spec.SecurityScheme {
	secScheme := &spec.SecurityScheme{}
	switch flow.Name {
	case resource.FlowImplicitType:
		secScheme = spec.OAuth2Implicit(flow.AuthorizationURL)
	case resource.FlowPasswordType:
		secScheme = spec.OAuth2Password(flow.TokenURL)
	case resource.FlowAuthCodeType:
		secScheme = spec.OAuth2AccessToken(flow.AuthorizationURL, flow.TokenURL)
	case resource.FlowClientCredentialType:
		secScheme = spec.OAuth2Application(flow.TokenURL)
	default:
		secScheme = &spec.SecurityScheme{SecuritySchemeProps: spec.SecuritySchemeProps{
			Type:             "oauth2",
			Flow:             flow.Name,
			TokenURL:         flow.TokenURL,
			AuthorizationURL: flow.AuthorizationURL,
		}}
	}

	secScheme.Scopes = flow.Scopes
	return secScheme
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
		o.addDefinition(structName, refSchema)
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
