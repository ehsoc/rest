# rest
Rest is an experimental web resource abstraction for composing REST APIs in Go (Golang).

- **Rapid prototyping**.
- **Web Server generation (http.Handler)**
- **REST API Specification generation (Open-API v2)**

## Base components:
- API. Is the root of all resources. It contains the server information, and can generate the server handler, and the API specification.
- Resource. Each resource is a node in a URL path, and contains the method and other resources.

**Code example:**

```go
api := rest.NewAPI("/api/v1", "localhost", "My simple car API", "v1")

api.Resource("car", func(r *rest.Resource) {
	r.Resource("findMatch", func(r *rest.Resource) {
	})
	carID := rest.NewURIParameter("carID", reflect.String)
	r.ResourceP(carID, func(r *rest.Resource) {
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
                                     |   API     |
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
- GenerateServer(g rest.ServerGenerator) http.Handler
- GenerateSpec(w io.Writer, api rest.API)

## Example:
```go
api := rest.NewAPI("/v1", "localhost", "My simple car API", "v1")
// Generating OpenAPI v2 specification to standard output
api.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})
// Generating server handler
server := api.GenerateServer(chigenerator.ChiGenerator{})
```

## Resource
A Resource can contain:
- Methods: A collection of HTTP methods.
- Resources: Collection of child resources.

## Example:
```go
api.Resource("user", func(r *rest.Resource) {
    // Method
    r.Get(getUser, ct)
    // Resource
    r.Resource("logout", func(r *rest.Resource) {
        // Method
        r.Get(logOut, ct)
    })
})
```

## ResourceP (URI parameter Resource)
ResourceP creates a new URI parameter resource node.
The first argument must be a Parameter of URIParameter type. Use NewURIParameter to create one.

## Example:
```go
carID := rest.NewURIParameter("carID", reflect.String)
r.ResourceP(carID, func(r *rest.Resource) {
    r.Get(getCar, ct).WithParameter(carID)
})
```
In the example the route to get a car will be `/car/{carID}`, where `{carID}` is the variable part.

### Method
A `Method` represents an HTTP method with an HTTP Handler. A default handler will make sense of the method specification and return the appropriate HTTP response. The specification elements to be managed by the default handler are: **Content negotiation, security, validation, and operation.**

- MethodOperation: Describes an `Operation` and responses (`Response` for success and failure).
- ContentTypes: Describes the available content-types and encoder/decoders for request and responses. 
- Negotiator: Interface for content negotiation. A default implementation will be set when you create a Method.
- SecurityCollection: Is the security definition.
- Parameters: The parameters expected to be sent by the client. The main purpose of the declaration of parameters is for API specification generation.
- Handler: The http.Handler of the method.  The default handler will be set when you create a new Method.

### Operation
Represents a logical operation upon a resource, like delete, list, create, ping, etc. `Operation` is an interface defined by an `Execute` method.

#### Execute method
- 	Inputs: `Input` type
- 	Output: `body` interface{}, `success` bool, and `err` error .


	1. `body` (interface{}): Is the body that is going to be send to the client.(Optional)
	2. `success` (bool): If the value is true, it will trigger the `successResponse` (argument passed in the `NewMethodOperation` function). If the value is false, it will trigger the `failResponse` (set it with `WithFailResponse` method). False means that the most positive operation output didn't happened, but is not an API nor a client error.
	3.  `err` (error): The `err`(error) is meant to indicate an API error, or any internal server error, like a database failure, i/o error, etc. The `err`!=nil will always trigger a 500 code error.

### Method:

```
                               +----------------+
                               |                |
                               |     Method     |
                               |                |
                               +--------+-------+
                                        |
                               +--------+-------+
                               |                |
                           +---+ MethodOperation+---+
                           |   |                |   |
Your operation method      |   +--------+-------+   |
goes here                  |            |           |
      +                    |            |           |
      |             +------+----+ +-----+-----+ +---+-------+
      |             |           | | Response  | |  Response |
      +-----------> | Operation | | success   | |  fail     |
                    |           | |           | |           |
                    +-----------+ +-----------+ +-----------+

```

## Full code example:
```go
// Responses
successResponse := rest.NewResponse(200).WithOperationResultBody(Car{})
failResponse := rest.NewResponse(404)

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
        // Car not found, success is false, no error. This will trigger the `failResponse` (response code 404)
        return nil, false, nil
    }

    // Car found, success true, error nil. This will trigger the `successResponse` (response code 200)
    return Car{carID, "Foo"}, true, nil
}
// Method Operation
getCar := rest.NewMethodOperation(rest.OperationFunc(getCarOperation), successResponse).WithFailResponse(failResponse)
// ContentTypes
ct := rest.NewContentTypes()
ct.Add("application/json", encdec.JSONEncoderDecoder{}, true)

api := rest.NewAPI("/v1", "localhost", "My simple car API", "v1")
api.Resource("car", func(r *rest.Resource) {
    carID := rest.NewURIParameter("carID", reflect.String)
    r.ResourceP(carID, func(r *rest.Resource) {
        r.Get(getCar, ct).WithParameter(carID)
    })
})
// Generating OpenAPI v2 specification to standard output
api.GenerateSpec(os.Stdout, &oaiv2.OpenAPIV2SpecGenerator{})

// Generating server routes
server := api.GenerateServer(chigenerator.ChiGenerator{})
http.ListenAndServe(":8080", server)

```
