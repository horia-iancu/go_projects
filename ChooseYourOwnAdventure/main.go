package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func parseJson(filename string, storyData *map[string]GopherObj) {
	fileData, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(fileData, &storyData); err != nil {
		log.Fatal(err)
	}
}

var templ *template.Template

type handler map[string]GopherObj

type GopherObj struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string
}

func NewHandler(m map[string]GopherObj) handler {
	return handler(m)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	err := templ.Execute(w, h[path])
	if err != nil {
		log.Fatal(err)
	}
}

func initTemplate() {
	templ = template.Must(template.ParseFiles("layout.html"))
}

func main() {
	var storyData map[string]GopherObj
	parseJson("gopher.json", &storyData)

	initTemplate()

	h := NewHandler(storyData)
	http.ListenAndServe(":8080", h)
}
