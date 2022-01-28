package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swkkd/budget-google/APISearchRequest/handlers"
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_search_request_http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		totalRequests.WithLabelValues(path).Inc()
	})
}

func init() {
	prometheus.Register(totalRequests)
}

//main start webserver
func main() {
	router := mux.NewRouter()
	router.Use(prometheusMiddleware)
	router.HandleFunc("/", handlers.SearchHandler).Queries("search", "{search}")

	router.Path("/metrics").Handler(promhttp.Handler())

	http.Handle("/", router)

	err := http.ListenAndServe(":9002", router)
	log.Fatal(err)
}
