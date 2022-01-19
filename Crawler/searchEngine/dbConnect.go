package searchEngine

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

type UrlContent struct {
	ID   int
	Url  string
	Body string
}

func ConnectToES() {
	//if config is not specified it uses default port to connect to! use
	//.NewClient(cfg) to connect to specific port!
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()
	log.Println(res)

}
