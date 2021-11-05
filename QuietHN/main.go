package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gophercises/quiet_hn/hn"
)

const numWorkers = 5

var siteCache = Cache{false, nil}

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	go tickerFunction()

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var sortedStories []item
		var clientErr int
		if !siteCache.Valid {
			sortedStories, clientErr = retrieveItems(numStories)
			if clientErr != 0 {
				http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
				return
			}
			siteCache.CachedStories = sortedStories
			siteCache.Valid = true
		}
		/*
			var client hn.Client
			ids, err := client.TopItems()
			if err != nil {
				http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
				return
			}
			stories := make(map[int]item)
			var totalStories = 0
			var lockNumStories sync.Mutex
			var lockStories sync.Mutex
			var toWorkers = workerAttr{
				Client:       client,
				TotalStories: &totalStories,
				NumStories:   numStories,
				Ids:          ids,
				Stories:      &stories,
				muNumStories: &lockNumStories,
				muStories:    &lockStories,
			}

			results := make(chan int, numWorkers)

			for i := 0; i < numWorkers; i++ {
				go worker(toWorkers, i, results)
			}

			for i := 0; i < numWorkers; i++ {
				<-results
			}

			keys := make([]int, 0)
			sortedStories := make([]item, 0)
			for k := range stories {
				keys = append(keys, k)
			}

			sort.Ints(keys)
			for _, v := range keys {
				sortedStories = append(sortedStories, stories[v])
			}*/

		data := templateData{
			Stories: siteCache.CachedStories[:30],
			Time:    time.Now().Sub(start),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

func worker(data workerAttr, id int, results chan<- int) {
	share := data.NumStories / numWorkers
	var startIdx int
	var endIdx int

	if id == 0 {
		startIdx = 0
	} else {
		startIdx = id * int(float32(share)*1.25)
	}

	if id == numWorkers-1 {
		endIdx = int(float32(data.NumStories) * 1.25)
	} else {
		endIdx = (id + 1) * int(float32(share)*1.25)
	}
	for i := startIdx; i < endIdx; i++ {
		hnItem, err := data.Client.GetItem(data.Ids[i])
		if err != nil {
			results <- 1
			return
		}
		item := parseHNItem(hnItem)
		if isStoryLink(item) {
			data.muStories.Lock()
			(*data.Stories)[data.Ids[i]] = item
			data.muStories.Unlock()

			data.muNumStories.Lock()
			if *data.TotalStories >= data.NumStories {
				data.muNumStories.Unlock()
				results <- 1
				return
			}
			*data.TotalStories += 1
			data.muNumStories.Unlock()
		}
	}
	results <- 1
}

func retrieveItems(numStories int) ([]item, int) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	stories := make(map[int]item)
	var totalStories = 0
	var lockNumStories sync.Mutex
	var lockStories sync.Mutex
	var toWorkers = workerAttr{
		Client:       client,
		TotalStories: &totalStories,
		NumStories:   numStories,
		Ids:          ids,
		Stories:      &stories,
		muNumStories: &lockNumStories,
		muStories:    &lockStories,
	}

	results := make(chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(toWorkers, i, results)
	}

	for i := 0; i < numWorkers; i++ {
		<-results
	}

	keys := make([]int, 0)
	sortedStories := make([]item, 0)
	for k := range stories {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, v := range keys {
		sortedStories = append(sortedStories, stories[v])
	}
	return sortedStories, 0
}

func tickerFunction() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		siteCache.Valid = false
	}
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type workerAttr struct {
	Client       hn.Client
	TotalStories *int
	NumStories   int
	Ids          []int
	Stories      *map[int]item
	muNumStories *sync.Mutex
	muStories    *sync.Mutex
}

type Cache struct {
	Valid         bool
	CachedStories []item
}
