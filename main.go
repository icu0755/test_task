package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const ItemsPerPage = 5

type ArtistsJson []struct {
	Name string `json:"name"`
}

type PostsResponse struct {
	Items ArtistsJson `json:"items"`
	Pages int         `json:"pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataSource struct {
	Artists ArtistsJson
	Count   int
}

func (ds *DataSource) fromFile(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("cannot open file", err.Error())
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("cannot read file", err.Error())
	}
	err = json.Unmarshal(bs, &ds.Artists)
	if err != nil {
		panic(err)
	}

	ds.Count = len(ds.Artists)
}

func (ds *DataSource) getPageItems(page int, itemsPerPage int) ArtistsJson {
	first := itemsPerPage * (page - 1)
	last := itemsPerPage * (page)

	if first > ds.Count {
		first = ds.Count

	}

	if last > ds.Count {
		last = ds.Count
	}

	return ds.Artists[first:last]
}

func (ds *DataSource) getPages(itemsPerPage int) int {
	pages := len(ds.Artists) / itemsPerPage
	if ds.Count%itemsPerPage != 0 {
		pages += 1
	}
	return pages
}

type myHandler struct {
	ds DataSource
}

func newMyHandler(ds DataSource) *myHandler {
	return &myHandler{ds: ds}
}

func (h myHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var data []byte

	switch req.URL.Path {
	case "/api/posts":
		page := h.getPage(req)

		res.Header().Set("Content-Type", "application/json")

		if simulateError() {
			res.WriteHeader(http.StatusInternalServerError)
			response := ErrorResponse{
				"Error occured. Please try later.",
			}
			data, _ = json.Marshal(response)
		} else {
			response := PostsResponse{
				h.ds.getPageItems(page, ItemsPerPage),
				h.ds.getPages(ItemsPerPage),
			}
			data, _ = json.Marshal(response)
		}

		res.Write(data)
	}
}

func (h myHandler) getPage(req *http.Request) int {
	var page int
	pages, ok := req.URL.Query()["page"]
	if !ok {
		page = 1
	} else {
		page, _ = strconv.Atoi(pages[0])
	}
	return page
}

func simulateError() bool {
	return (rand.Int() % 2) == 0
}

func main() {
	var ds DataSource
	ds.fromFile("./response.json")
	h := newMyHandler(ds)

	log.Println("serve on http://localhost:9000")
	http.ListenAndServe(":9000", h)
}
