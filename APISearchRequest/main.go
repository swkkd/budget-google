package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swkkd/budget-google/APISearchRequest/handlers"
	"github.com/swkkd/budget-google/APISearchRequest/middleware"
	"log"
	"net/http"
)

func init() {
	err := prometheus.Register(middleware.TotalRequests)
	if err != nil {
		log.Printf("Failed to register prometheus data with error: %s", err)
	}
}

//main start webserver
func main() {
	router := mux.NewRouter()
	router.Use(middleware.PrometheusMiddleware)
	router.HandleFunc("/", handlers.SearchHandler).Queries("search", "{search}")

	router.Path("/metrics").Handler(promhttp.Handler())

	http.Handle("/", router)

	err := http.ListenAndServe(":9002", router)
	log.Fatal(err)
}
