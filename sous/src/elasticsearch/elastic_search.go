package elasticsearch

// Elastic search client
import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"search-engine-client/src/structs"

	elastic "github.com/olivere/elastic/v7"
)

const (
	IndexName    = "recipes"
	IndexMapping = `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"title": {
					"type":"text"
				},
				"description": {
					"type":"text"
				},
				"body": {
					"type":"text"
				},
				"url": {
					"type":"text"
				}
			}
		}
	}`
)

var client *elastic.Client

// NewElasticSearchClient returns an elastic seach client
func NewElasticSearchClient() *elastic.Client {
	var err error
	connected := false
	retries := 0

	// Custom retry strategy for docker-compose initialization
	for connected == false {
		// Create a new elastic client
		client, err = elastic.NewClient(
			elastic.SetURL("http://192.168.1.78:9200"), elastic.SetSniff(false))
		if err != nil {
			// log.Fatal(err)
			if retries == 5 {
				log.Fatal(err)
			}
			fmt.Println("Elasticsearch isn't ready for connection", 5-retries, "less")
			retries++
			time.Sleep(3 * time.Second)
		} else {
			connected = true
		}
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.1.78:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return client
}

// ExistsIndex checks if the given index exists or not
func ExistsIndex(i string) bool {
	// Check if index exists
	exists, err := client.IndexExists(i).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return exists
}

// CreateIndex creates a new index
func CreateIndex(i string) {
	createIndex, err := client.CreateIndex(IndexName).
		Body(IndexMapping).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}
	if !createIndex.Acknowledged {
		log.Println("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}
}

// SearchContent returns the results for a given query
func SearchContent(input string) []structs.Page {
	pages := []structs.Page{}

	ctx := context.Background()
	// Search for a page in the database using multi match query
	q := elastic.NewMultiMatchQuery(input, "title", "description", "body", "url").
		Type("most_fields").
		Fuzziness("2")
	result, err := client.Search().
		Index(IndexName).
		Query(q).
		From(0).Size(50).
		Sort("_score", false).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var ttyp structs.Page
	for _, page := range result.Each(reflect.TypeOf(ttyp)) {
		p := page.(structs.Page)
		pages = append(pages, p)
	}

	return pages
}
