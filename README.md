# resource
Resource is an experimental Web Resource abstraction for composing REST APIs in Go (Golang).
 
- **Rapid prototyping**.
- **Web Server generation (http.Handler)**
- **REST API Specification generation (Open-API v2)**

## Base components:
- RestAPI (Like a Resource, but it is the root resource node)
- Resource. Each resource is a node in a URL path.

**Code example:**

```go
api := resource.RestAPI{}
api.BasePath = "/v2"
api.Host = "localhost"

api.Resource("car", func(r *resource.Resource) {
	r.Resource("findMatch", func(r *resource.Resource) {
	})
	r.Resource("{carId}", func(r *resource.Resource) {
	})
})

api.Resource("ping", func(r *resource.Resource) {
})

api.Resource("user", func(r *resource.Resource) {
	r.Resource("signOut", func(r *resource.Resource) {
	})
})
	
```

**Diagram of the above code:**

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

## Resource
- Methods: A collection of HTTP methods (`Method`)
- Resources: Collection of child resources.

### Method
A `Method` represents an HTTP method with an HTTP Handler.
- MethodOperation: Describes an `Operation` and responses (`Response` for success and failure).
- Renderers: Describes the available renderers for request and responses. 
- Negotiator: Interface responsable for content negotiation. A default implementation will be set when you create a Method.
- Parameters: The parameters expected to be sent by the client. The main purpose of the declaration of parameters is for API specification generation.
  
### Operation
Represents a logical operation upon a resource, like delete, list, create, etc. `Operation` is an interface defined by an `Execute` method.

#### Execute method
- 	Inputs: `Input` type
- 	Output: `body` interface{}, `success` bool, and `err` error .


	1. `body` (interface{}): Is the body that is going to be send to the client.(Optional)
	2. `success` (bool): If the value is true, it will trigger the `successResponse` (argument passed in the `NewMethodOperation` function). If the value is false, it will trigger the `failResponse` (argument passed in the `NewMethodOperation` function). False means that the most positive operation output didn't happened, but is not either an API error or a client error.

	3.  `err` (error): The `err`(error) is meant to indicate an API error, or any internal server error, like a database failure, i/o error, etc. The `err`!=nil will always trigger a 500 code error.

### Resource main components diagram:

```
                                                                   +-----------+
                                                                   | Resource  |
                                                               +---+ "{carId}" +---+
                                                               |   |           |   |
                                                               |   +-----------+   |
                                                               |                   |
                                                         +-----+-----+        +----+------+
                                                         |  Method   |        |  Method   |
                                        +----------------+   "GET"   |        | "DELETE"  +----------+
                                        |                |           |        |           |          |
                                        |                +-----+-----+        +----+------+          |
                                        |                      |                   |                 |
                               +--------+-------+      +-------+--------+
                               | MethodOperation|      |    Renderers   |
                        +------+                +---+  |                |
                        |      |                |   |  |                |
Your operation method   |      +------+---------+   |  +----------------+
goes here               |             |             |
    +                   |             |             |
    |             +-----+-----+  +----+------+   +--+--------+
    |             | Operation |  | Response  |   |  Response |
    +-----------> |           |  | success   |   |  fail     |
                  |           |  |           |   |           |
                  +-----------+  +-----------+   +-----------+




```

## Code example:
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
successResponse := res.NewResponse(200).WithOperationResultBody(Car{})
failResponse := res.NewResponse(404)
getCar := res.NewMethodOperation(res.OperationFunc(getCarOperation), successResponse).WithFailResponse(failResponse)

ct := res.NewRenderers()
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
