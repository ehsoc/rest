# resource
Resource is an experimental Web Resource abstraction for composing REST APIs in Go (Golang).
 
- **Rapid prototyping**.
- **Web Server generation (http.Handler)**
- **REST API Specification generation (Open-API v2)**

## Base components:
- RestAPI (Like a Resource, but it is the root node)
- Resource. Each resource is a node.

Code example:

```go
api := resource.RestAPI{}
	api.BasePath = "/v2"
	api.Host = "localhost"
	api.Resource("car", func(r *resource.Resource){
		r.Resource("findMatch", func(r *resource.Resource){
		}
		r.Resource("{carId}", func(r *resource.Resource){
		}
	})
	r.Resource("ping", func(r *resource.Resource){
	}
	r.Resource("user", func(r *resource.Resource){
		r.Resource("signOut", func(r *resource.Resource){
		}
	}
	
```
Diagram representation of the above code:

```
                                     +-----------+
                                     |           |
                          +----------+  RestAPI  +-----------+
                          |          |   "v2"    |           |
                          |          +-----+-----+           |
                          |                |                 |
                   +------v-----+    +-----v-----+     +-----v-----+
                   |            |    |           |     |           |
       +-----------+  Resource  |    |  Resource |     |  Resource |
       |           |   "car"    |    |  "ping"   |     |  "user"   |
       |           +-----+------+    +-----------+     +-----+-----+
       |                 |                                   |
+------v------+    +-----v------+                      +-----v------+
|             |    |            |                      |            |
|  Resource   |    |  Resource  |                      |  Resource  |
| "findMatch" |    |  "{carId}" |                      |  "signOut" |
+-------------+    +------------+                      +------------+

```


## RestAPI methods:
- GenerateServer(restAPI RestAPI) http.Handler
- GenerateAPISpec(w io.Writer, restAPI RestAPI)

## Resource main components:
- Methods: A collection of HTTP methods (Method)
- Method: A Method represents an HTTP method with an HTTP Handler.
- MethodOperation: Describes an Operation and Responses, one for Operation success, and another for failure.
- HTTPContentTypeSelector: Describes the available content types of for request and responses. Contains a default Content-Type negotiator that you can replace with your own implementation.
- Operation: Represents a logical operation function.`Operation` is an interface defined by an `Execute` method.

	- Execute method:
		- 	Inputs: Input
		- 	Output: `body` interface{}, `success` bool, and `err` error .
	

			1. `body` (interface{}): is the body that is going to be sent to the client.
			2. `success` (bool):  false value means that the operation result was not the expected, but it is not an API error nor a client error. This will 		trigger the `successResponse` (argument passed in the `NewMethodOperation` function) if `success` return value is true. If false, will return 		`failResponse` (argument passed in the `NewMethodOperation` function).

			3.  `err` (error): The `err`(error) is meant to indicate an internal server error when `err`!=nil, like a database failure or other API error. T		he `err`!=nil will trigger a 500 error.
	
- Parameters: For specification, validation, and get helper functions.

Example:
```go
getCarOperation := func(i res.Input) (body interface{}, success bool, err error) {
	carId, err := i.GetURIParam("carId")
	if err != nil {
		log.Println("error getting parameter: ", err)
		return Car{}, false, err
	}
	if carId == "error" {
		//Internal error trying to get the car. This will trigger a response code 500
		return nil, false, errors.New("Internal error")
	}
	if carId != "101" {
		//Car not found, success is false, no error. This will trigger a response code 404
		return nil, false, nil
	}
	//Car found, success true, error nil. This will trigger a response code 200
	return Car{carId, "Foo"}, true, nil
}
getCar := res.NewMethodOperation(
	res.OperationFunc(getCarOperation),
	res.NewResponse(200).WithOperationResultBody(Car{})).
	WithFailResponse(res.NewResponse(404))

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
