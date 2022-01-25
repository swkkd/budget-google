package main

import (
	"github.com/gorilla/mux"
	"github.com/swkkd/budget-google/APISearchRequest/searchEngine"
	"html/template"
	"log"
	"net/http"
)

//Template is a custom html/template renderer for Echo framework

//main start webserver
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", search).Queries("search", "{search}")
	http.Handle("/", r)

	http.ListenAndServe(":9002", r)
}

////helloWorld simple hello world webpage
//func helloWorld(c echo.Context) error {
//	return c.String(http.StatusOK, "Hello World")
//}

func search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]
	req := searchEngine.Search(search)
	log.Printf("RENDERING: %#v", req)
	//fmt.Fprintf(w, "<html>"+html.UnescapeString(req[0].ContentOfPage))
	tmpl, err := template.ParseFiles("html/search.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := tmpl.Execute(w, req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}
