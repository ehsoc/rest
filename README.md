# rest
Rest is an experimental Web Resource abstraction for composing REST APIs in Go (Golang).
 
- **Rapid prototyping**.
- **Web Server generation (http.Handler)**
- **REST API Specification generation (Open-API v2)**

## Base components:
- API (Like a Resource, but it is the root resource node)
- Resource. Each resource is a node in a URL path.

**Code example:**

```go
api := rest.API{}
api.BasePath = "/v2"
api.Host = "localhost"

api.Resource("car", func(r *rest.Resource) {
	r.Resource("findMatch", func(r *rest.Resource) {
	})
	r.Resource("{carId}", func(r *rest.Resource) {
	})
})

api.Resource("ping", func(r *rest.Resource) {
})

api.Resource("user", func(r *rest.Resource) {
	r.Resource("signOut", func(r *rest.Resource) {
	})
})
	
```

**Diagram of the above code:**

```
                                     +-----------+
                                     |  API      |
                          +----------+   "/2"    +-----------+
                          |          |           |           |
                          |          +-----+-----+           |
                          |                |                 |
                   +------v-----+    +-----v------+    +-----v------+
                   |  Resource  |    |  Resource  |    |  Resource  |
       +-----------+   "car"    |    |   "ping"   |    |  "user"    |
       |           |            |    |            |    |            |
       |           +-----+------+    +------------+    +-----+------+
       |                 |                                   |
+------v-----+     +-----v------+                      +-----v------+
|  Resource  |     |  Resource  |                      |  Resource  |
| "findMatch"|     | "{carId}"  |                      |  "signOut" |
|            |     |            |                      |            |
+------------+     +------------+                      +------------+

```


## API methods:
- GenerateServer(a API) http.Handler
- GenerateAPISpec(w io.Writer, api API)

## Resource
- Methods: A collection of HTTP methods (`Method`)
- Resources: Collection of child resources.

### Method
A `Method` represents an HTTP method with an HTTP Handler. A default handler will make sense of the method specification and return the appropriate HTTP response. The specification elements to be managed by the default handler are: Content negotiation, validation, and operation. 

- MethodOperation: Describes an `Operation` and responses (`Response` for success and failure).
- ContentTypes: Describes the available content-type encoder/decoders for request and responses. 
- Negotiator: Interface for content negotiation. A default implementation will be set when you create a Method.
- Parameters: The parameters expected to be sent by the client. The main purpose of the declaration of parameters is for API specification generation.
- Handler: The http.Handler of the method.  The default handler will be set when you create a new Method.

\* You can override the default handler if necessary, with the `Handler` property (method.Handler = MyHandler)
  
### Operation
Represents a logical operation upon a resource, like delete, list, create, ping, etc. `Operation` is an interface defined by an `Execute` method.

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
                                        |                      |                   +                 +
                               +--------+-------+      +-------+--------+
                               | MethodOperation|      |  ContentTypes  |
                           +---+                +---+  |                |
                           |   |                |   |  |                |
Your operation method      |   +--------+-------+   |  +----------------+
goes here                  |            |           |
      +                    |            |           |
      |             +------+----+ +-----+-----+ +---+-------+
      |             | Operation | | Response  | |  Response |
      +-----------> |           | | success   | |  fail     |
                    |           | |           | |           |
                    +-----------+ +-----------+ +-----------+


```

## Code example:
```go
getCarOperation := func(i rest.Input) (body interface{}, success bool, err error) {
    carID, err := i.GetURIParam("carID")

    if err != nil {
        log.Println("error getting parameter: ", err)
        return Car{}, false, err
    }

    if carID == "error" {
        // Internal error trying to get the car. This will trigger a response code 500
        return nil, false, errors.New("Internal error")
    }

    if carID != "101" {
        // Car not found, success is false, no error. This will trigger a response code 404
        return nil, false, nil
    }

    // Car found, success true, error nil. This will trigger a response code 200
    return Car{carID, "Foo"}, true, nil
}
// Responses
successResponse := rest.NewResponse(200).WithOperationResultBody(Car{})
failResponse := rest.NewResponse(404)
// Method Operation
getCar := rest.NewMethodOperation(rest.OperationFunc(getCarOperation), successResponse).WithFailResponse(failResponse)
// ContentTypes
ct := rest.NewContentTypes()
ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)

api := rest.NewAPI("/v1", "localhost", "My simple car API", "v1")
api.Resource("car", func(r *rest.Resource) {
    carIDParam := rest.NewURIParameter("carID", reflect.String)
    r.Resource("{carID}", func(r *rest.Resource) {
        r.Get(getCar, ct).WithParameter(carIDParam)
    })
})
// Generating OpenAPI v2 specification to standard output
api.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})

// Generating server routes
server := api.GenerateServer(chigenerator.ChiGenerator{})
```
