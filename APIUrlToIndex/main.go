package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swkkd/budget-google/APIUrlToIndex/middleware"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

var kafkaServer, kafkaTopic string

type urlToIndex struct {
	url string
}

func init() {
	kafkaServer = readFromENV("KAFKA_BROKER", "localhost:29092")
	kafkaTopic = readFromENV("KAFKA_TOPIC", "api-to-index")

	fmt.Println("Kafka Broker - ", kafkaServer)
	fmt.Println("Kafka topic - ", kafkaTopic)

	err := prometheus.Register(middleware.TotalRequests)
	if err != nil {
		log.Printf("Failed to register prometheus data with error: %s", err)
	}
}

func main() {
	p, err := NewProducer(kafkaTopic)
	if err != nil {
		panic(err)
	}

	defer p.Close()
	controller := NewController(p)
	router := mux.NewRouter()
	router.Use(middleware.PrometheusMiddleware)
	router.HandleFunc("/", controller.sendUrlToIndex)

	router.Path("/metrics").Handler(promhttp.Handler())

	http.Handle("/", router)

	err = http.ListenAndServe(":9001", router)
	log.Fatal(err)
}
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//todo check if the url is already parsed
//todo some sort of visualization what urls are already in DB
//todo check if parsed data is up-to-date

func (co *Controller) sendUrlToIndex(w http.ResponseWriter, r *http.Request) {
	var urls urlToIndex
	tmpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if r.Method == "POST" {
		urls.url = r.FormValue("url")
	}
	if IsUrl(urls.url) == true {
		log.Printf("URL: %v", urls.url)

		err := co.producer.Send([]byte(urls.url))
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Printf("%s IS NOT VALID URL!", urls.url)
	}
}

func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

//todo return the response to the html page if url added successfully!
