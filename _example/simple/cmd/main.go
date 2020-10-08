package main

import (
	"net/http"

	"github.com/ehsoc/resource/_example/simple"
)

func main() {
	server := simple.GenerateServer()
	http.ListenAndServe(":8080", server)
}
