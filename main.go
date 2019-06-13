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

type ArtistsJson []struct {
	Name string `json:"name"`
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

func (ds *DataSource) getPage(page int, itemsPerPage int) ArtistsJson {
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

		artists := h.ds.getPage(page, ItemsPerPage)

		data, _ := json.Marshal(artists)
		res.Write(data)
	}
}

func main() {
	var ds DataSource
	ds.fromFile("./response.json")
	h := newMyHandler(ds)

	log.Println("serve on http://localhost:9000")
	http.ListenAndServe(":9000", h)
}
