package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const ItemsPerPage = 5

type myHandler int
type ArtistsJson []struct {
	Name string `json:"name"`
}

var response string
var a ArtistsJson

func (h myHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var page int

	switch req.URL.Path {
	case "/api/posts":
		pages, ok := req.URL.Query()["page"]
		if !ok {
			page = 1
		} else {
			page, _ = strconv.Atoi(pages[0])
		}

		fmt.Println("page: ", page)

		res.Header().Set("Content-Type", "application/json")

		artists := getPage(page, ItemsPerPage)

		data, _ := json.Marshal(artists)
		res.Write(data)
	}
}

func getPage(page int, itemsPerPage int) ArtistsJson {
	return a[itemsPerPage*(page-1) : itemsPerPage*(page)]
}

func readResponse() {
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
	err = json.Unmarshal(bs, &a)
	if err != nil {
		panic(err)
	}
}

func main() {
	var h myHandler

	readResponse()
	log.Println("serve on http://localhost:9000")
	http.ListenAndServe(":9000", h)
}
