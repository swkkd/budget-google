package main

import (
	"github.com/swkkd/budget-google/APISearchRequest/searchEngine"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

//Template is a custom html/template renderer for Echo framework
type Template struct {
	templates *template.Template
}

//main start webserver
func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseFiles("html/search.html")),
	}
	e.Renderer = t

	//Controller
	e.GET("/", search)
	e.GET("/search", search)

	e.Logger.Fatal(e.Start(":9002"))
}

////helloWorld simple hello world webpage
//func helloWorld(c echo.Context) error {
//	return c.String(http.StatusOK, "Hello World")
//}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func search(c echo.Context) error {
	search := c.QueryParam("search")
	r := searchEngine.Search(search)

	return c.Render(http.StatusOK, "search.html", r)
}
