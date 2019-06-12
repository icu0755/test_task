package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type myHandler int

var response string

func (h myHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/api/posts":
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, response)
	}
}

func main() {
	var h myHandler

	f, err := os.Open("./response.json")
	if err != nil {
		log.Fatalln("cannot open file", err.Error())
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("cannot read file", err.Error())
	}

	response = string(bs)
	log.Println("serve on http://localhost:9000")
	http.ListenAndServe(":9000", h)
}
