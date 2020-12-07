package petstore

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/ehsoc/rest"
	"github.com/ehsoc/rest/encdec"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

var notFoundResponse = rest.NewResponse(404)

func GeneratePetStore() rest.API {
	ct := rest.NewContentTypes()
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)
	ct.Add("application/xml", encdec.XMLEncoderDecoder{}, false)
	// POST
	create := rest.NewMethodOperation(rest.OperationFunc(operationCreate), rest.NewResponse(201)).WithFailResponse(rest.NewResponse(400))
	petScopes := map[string]string{"write:pets": "modify pets in your account", "read:pets": "read your pets"}
	petSO := rest.SecurityOperation{
		Authenticator: rest.AuthenticatorFunc(func(i rest.Input) rest.AuthError {
			return nil
		}),
		FailedAuthenticationResponse: rest.NewResponse(401),
		FailedAuthorizationResponse:  rest.NewResponse(403),
	}
	petAuthScheme := rest.NewOAuth2SecurityScheme("petstore_auth", petSO).
		WithImplicitOAuth2Flow("localhost:5050/oauth/dialog", petScopes)

	api := rest.API{}
	api.BasePath = "/v2"
	api.Host = "localhost"
	api.Resource("pet", func(r *rest.Resource) {
		r.Post(create, ct).
			WithRequestBody("Pet object that needs to be added to the store", Pet{}).
			WithSummary("Add a new pet to the store").
			WithSecurity(petAuthScheme)

		// PUT
		update := rest.NewMethodOperation(rest.OperationFunc(operationUpdate), rest.NewResponse(200)).WithFailResponse(rest.NewResponse(404).WithDescription("Pet not found"))
		r.Put(update, ct).
			WithRequestBody("Pet object that needs to be added to the store", Pet{}).
			WithSummary("Update an existing pet").
			WithValidation(rest.Validation{
				Validator: rest.ValidatorFunc(func(input rest.Input) error {
					pet := Pet{}
					body, _ := input.GetBody()
					respBody := new(bytes.Buffer)
					cBody := io.TeeReader(body, respBody)
					err := input.BodyDecoder.Decode(cBody, &pet)
					if err != nil {
						return err
					}
					input.Request.Body = ioutil.NopCloser(respBody)
					return nil
				}),
				Response: rest.NewResponse(400).WithDescription("Invalid ID supplied")})

		// Uri Parameters declaration, so it is available to all anonymous resources functions
		petIDURIParam := rest.NewURIParameter("petId", reflect.Int64).WithDescription("ID of pet to return").WithExample(1)
		// SubResource
		// New Resource with URIParam Resource GET By ID {petId} = /pet/{petId}
		r.Resource("{petId}", func(r *rest.Resource) {
			ct := rest.NewContentTypes()
			ct.AddEncoder("application/json", encdec.JSONEncoder{}, true)
			ct.AddEncoder("application/xml", encdec.XMLEncoder{}, false)
			getByID := rest.NewMethodOperation(rest.OperationFunc(operationGetPetByID), rest.NewResponse(200).WithOperationResultBody(Pet{})).WithFailResponse(notFoundResponse)
			petAPIKeySO := rest.SecurityOperation{
				Authenticator: rest.AuthenticatorFunc(func(i rest.Input) rest.AuthError {
					return nil
				}),
				FailedAuthenticationResponse: rest.NewResponse(401),
				FailedAuthorizationResponse:  rest.NewResponse(403),
			}
			apiKeyScheme := rest.NewAPIKeySecurityScheme("api_key", rest.NewHeaderParameter("api_key", reflect.String), petAPIKeySO)

			r.Get(getByID, ct).
				WithSummary("Find pet by ID").
				WithDescription("Returns a single pet").
				WithParameter(petIDURIParam).
				WithSecurity(apiKeyScheme)
			// Delete
			deleteByID := rest.NewMethodOperation(rest.OperationFunc(operationDeletePet), rest.NewResponse(200)).WithFailResponse(notFoundResponse)
			r.Delete(deleteByID, ct).
				WithSummary("Deletes a pet").
				WithParameter(
					petIDURIParam.WithDescription("Pet id to delete").
						WithValidation(
							rest.Validation{
								Validator: rest.ValidatorFunc(func(i rest.Input) error {
									petID, _ := i.GetURIParam("petId")
									_, err := getInt64Id(petID)
									if err != nil {
										return err
									}
									return nil
								}),
								Response: rest.NewResponse(400).WithDescription("Invalid ID supplied")})).
				WithParameter(rest.NewHeaderParameter("api_key", reflect.String))
			r.Resource("uploadImage", func(r *rest.Resource) {
				// Upload image resource under URIParameter Resource
				uploadImage := rest.NewMethodOperation(rest.OperationFunc(operationUploadImage), rest.NewResponse(200).WithBody(APIResponse{200, "OK", "image created"}).WithDescription("successful operation"))
				ct := rest.NewContentTypes()
				ct.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
				ct.AddDecoder("multipart/form-data", encdec.XMLEncoderDecoder{}, true)
				r.Post(uploadImage, ct).
					WithParameter(petIDURIParam.WithDescription("ID of pet to update")).
					WithParameter(rest.NewFormDataParameter("additionalMetadata", reflect.String, encdec.JSONDecoder{}).WithDescription("Additional data to pass to server")).
					WithParameter(rest.NewFileParameter("file").WithDescription("file to upload")).
					WithParameter(rest.NewFormDataParameter("jsonPetData", reflect.Struct, encdec.JSONDecoder{}).WithDescription("json format data").WithBody(Pet{})).
					WithSummary("uploads an image")
			})
		})
		// Resource /pet/findByStatus
		r.Resource("findByStatus", func(r *rest.Resource) {
			ct := rest.NewContentTypes()
			ct.AddEncoder("application/json", encdec.JSONEncoderDecoder{}, true)
			ct.AddEncoder("application/xml", encdec.XMLEncoderDecoder{}, false)
			// GET
			findByStatus := rest.NewMethodOperation(rest.OperationFunc(operationFindByStatus), rest.NewResponse(200).WithOperationResultBody([]Pet{}).WithDescription("successful operation")).WithFailResponse(rest.NewResponse(400).WithDescription("Invalid status value"))
			statusParam := rest.NewQueryArrayParameter("status", []interface{}{"available", "pending", "sold"}).AsRequired().WithDescription("Status values that need to be considered for filter")
			statusParam.CollectionFormat = "multi"
			petBasicAuthSO := rest.SecurityOperation{
				Authenticator: rest.AuthenticatorFunc(func(i rest.Input) rest.AuthError {
					return nil
				}),
				FailedAuthenticationResponse: rest.NewResponse(401),
				FailedAuthorizationResponse:  rest.NewResponse(403),
			}
			basicSecurity := rest.NewSecurityScheme("basicSecurity", rest.BasicSecurityType, petBasicAuthSO)
			r.Get(findByStatus, ct).
				WithSummary("Finds Pets by status").
				WithDescription("Multiple status values can be provided with comma separated strings").
				WithParameter(statusParam).
				WithSecurity(basicSecurity)
		})
	})

	return api
}
