package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"
	"time"
)

// ESData todo add timestamp
type ESData struct {
	ParseData []ParseData `json:"parseData"`
	ParseDate time.Time   `json:"parseDate"`
}

// ParseData todo in the elasticsearch parsedata shouldn't be nested field
type ParseData struct {
	URL           string   `json:"url"`
	ContentOfPage []string `json:"contentOfPage"`
}

/*todo
create database if it doesnt exist
should i do it here in code?
*/

// ConnectToES connects to ES and create new client
//basically it only checks if database is UP
func ConnectToES() {
	var (
		r map[string]interface{}
	)
	//if config is not specified it uses default port to connect to! use
	//.NewClient(prometheus) to connect to specific port!
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
	log.Println(strings.Repeat("~", 37))
	log.Println("SUCCESSFULLY CONNECTED TO ELASTIC SEARCH")
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
}

// Insert insert data into elastic search
func Insert(url string, body []string) {

	log.SetFlags(0)

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.1.75:9200"},
		// ...
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	//todo receive timestamp at the other side
	doc := &ESData{
		ParseDate: time.Now(),
		ParseData: []ParseData{
			{
				ContentOfPage: body,
				URL:           url,
			},
		},
	}

	payload, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}

	req := esapi.IndexRequest{
		Index: "my-index-000001",
		//i dont specify the document id so ES can generate one
		//DocumentID: strconv.Itoa(i + 1),
		Body:    bytes.NewReader(payload),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID", res)
		//log.Printf("[%s] Error indexing document ID", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}
