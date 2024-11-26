package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"search-engine-client/src/elasticsearch"
	"search-engine-client/src/structs"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/home.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Print("Template executing error: ", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	searchInput := r.Form.Get("input")

	log.Print("Querying database for: ", searchInput)

	pages := elasticsearch.SearchContent(searchInput)

	searchResult := structs.SearchResult{
		Input: searchInput,
		Pages: pages,
	}

	jsonData, err := json.Marshal(searchResult)
	if err != nil {
		log.Print("JSON executing error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	elasticsearch.NewElasticSearchClient()
	exists := elasticsearch.ExistsIndex(elasticsearch.IndexName)

	if !exists {
		elasticsearch.CreateIndex(elasticsearch.IndexName)
	}

	mux := mux.NewRouter()

	mux.HandleFunc("/", homeHandler).Methods("GET")
	mux.HandleFunc("/search", searchHandler).Methods("GET")
	http.ListenAndServe(":80", mux)
}
