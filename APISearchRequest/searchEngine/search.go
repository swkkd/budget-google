package searchEngine

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
)

type ESSingleResponse struct {
	response *ESResponses
}

type ESResponses struct {
	Url           interface{}
	ContentOfPage interface{}
}

func Search(searchQuery string) []ESResponses {
	log.SetFlags(0)

	var (
		r map[string]interface{}
		//wg sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.1.75:9200"},
		// ...
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	// ---------------------------

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"contentOfPage": searchQuery,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("budget-google"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	//log.Printf("RESPONSE: ------------ %s", r)
	// Print the ID and document source for each hit.
	var responses []ESResponses
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		//log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"].(map[string]interface{})["URL"])
		singleResponse := ESResponses{
			Url:           hit.(map[string]interface{})["_source"].(map[string]interface{})["URL"],
			ContentOfPage: hit.(map[string]interface{})["_source"].(map[string]interface{})["contentOfPage"],
		}

		responses = append(responses, singleResponse)

	}
	log.Println(responses)

	log.Println(strings.Repeat("=", 37))
	return responses
}
