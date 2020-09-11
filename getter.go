package resource

import "net/http"

type Getter interface {
	Get(r *http.Request) string
}

// The GetterFunc type is an adapter to allow the use of
// ordinary functions as Getter. If f is a function
// with the appropriate signature, GetterFunc(f) is a
// Getter that calls f.
type GetterFunc func(r *http.Request) string

//Get calls gf(r)
func (f GetterFunc) Get(r *http.Request) string {
	return f(r)
}
