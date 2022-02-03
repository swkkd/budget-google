package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/elastic/go-elasticsearch/v8/esapi"
	"html/template"
	"log"
	"strings"
)

type ESSingleResponse struct {
	response *ESResponse
}

type ESResponse struct {
	Url           interface{}
	ContentOfPage template.HTML
}

func Search(searchQuery string) []ESResponse {
	log.SetFlags(0)

	var (
		r map[string]interface{}
	)

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
			"nested": map[string]interface{}{
				"path": "parseData",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": map[string]interface{}{
							"match": map[string]interface{}{
								"parseData.contentOfPage": searchQuery,
							},
						},
					},
				},
				"inner_hits": map[string]interface{}{
					"highlight": map[string]interface{}{
						"fields": map[string]interface{}{
							"parseData.contentOfPage": map[string]interface{}{
								"pre_tags":  "<mark>",
								"post_tags": "</mark>",
							},
						},
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("my-index-000001"),
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

	var responses []ESResponse
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		//log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"].(map[string]interface{})["URL"])
		hitsContent := hit.(map[string]interface{})["inner_hits"].(map[string]interface{})["parseData"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["highlight"].(map[string]interface{})["parseData.contentOfPage"].([]interface{})
		singleResponse := ESResponse{
			Url:           hit.(map[string]interface{})["_source"].(map[string]interface{})["parseData"].([]interface{})[0].(map[string]interface{})["url"],
			ContentOfPage: template.HTML(fmt.Sprintf("%v", hitsContent)),
		}
		responses = append(responses, singleResponse)

	}
	log.Printf("%v", responses)

	log.Println(strings.Repeat("=", 37))
	return responses

}
