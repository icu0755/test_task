package main

import (
	"io"
	"net/http"
)

type myHandler int

func (h myHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/api/posts":
		io.WriteString(res, "Hello World")
	}
}

func main() {
	var h myHandler
	http.ListenAndServe(":9000", h)
}
