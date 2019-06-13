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

var a ArtistsJson

type DataSource struct {
	Data ArtistsJson
}

func (ds DataSource) fromFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("cannot open file", err.Error())
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("cannot read file", err.Error())
	}
	err = json.Unmarshal(bs, &ds.Data)
	if err != nil {
		panic(err)
	}
}

func (ds DataSource) getPage(page int, itemsPerPage int) ArtistsJson {
	first := itemsPerPage * (page - 1)
	last := itemsPerPage * (page)
	count := len(ds.Data)

	if first > count {
		first = count

	}

	if last > count {
		last = count
	}

	return ds.Data[first:last]
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

type Foo struct {
	Bar []int
}

func (f Foo) init() {
	f.Bar = append(f.Bar, 1)
}

func main() {
	var ds DataSource
	ds.fromFile("./response.json")
	h := newMyHandler(ds)
	var f Foo
	f.init()
	f.init()
	f.init()

	log.Println("serve on http://localhost:9000")
	http.ListenAndServe(":9000", h)
}
