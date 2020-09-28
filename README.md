# resource
Resource is an experimental Web Resource abstraction for composing a REST API in Go (Golang).

- **Rapid prototyping**.
- **Web Server generation (http.Handler)**
- **REST API Specification generation (Open-API v2)**

## It is based on 2 main components:
1. Resource. Each resource is a node, and each resource node contains a collection of resources.
2. RestAPI (Like a Resource, but it is the root) 

## RestAPI main functions :
1. GenerateServer(restAPI RestAPI) http.Handler
2. GenerateAPISpec(w io.Writer, restAPI RestAPI)


**Resource components:**
- Methods: A collection of HTTP methods (Method)
- Method: A Method represents an HTTP method with a HTTP Handler.
- MethodOperation: Describes an Operation and two Responses, one for Operation success, and another for failure.
- HTTPContentTypeSelector: Describes the available content types of expected request and responses. Contains a default Content-Type negotiator that you can replace with your own implementation.
- Operation: Represents a logical operation function.
	- 	Inputs: *http.Request, and encdec.Decoder as inputs.
	- 	Output: interface{}, and error (error on failure, nil if success)
- Properties: Currently only for documenting purposes, because do not enforce the use of these properties nor affect any functionality, you roll your own request parsing.

Example:
```
	getCarOperation := func(r *http.Request, decoder encdec.Decoder) (interface{}, error) {
		return Car{"1", "Mazda"}, nil
	}
	getCar := res.NewMethodOperation(
		res.OperationFunc(getCarOperation),
		res.NewResponse(200),
		res.NewResponse(404),
		true,
	)
	ct := res.NewHTTPContentTypeSelector()
	ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)

	root := res.RestAPI{}
	root.BasePath = "/v1"
	root.Version = "v1"
	root.Title = "My simple car API"
	root.Resource("car", func(r *res.Resource) {
		carIdParam := res.NewURIParameter("carId", reflect.String)
		r.Resource("{carId}", func(r *res.Resource) {
			r.Get(getCar, ct).WithParameter(carIdParam)
		})
	})
	server := root.GenerateServer(chigenerator.ChiGenerator{})
	root.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})
	http.ListenAndServe(":8080", server)
```
