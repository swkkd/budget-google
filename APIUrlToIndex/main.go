package main

import (
	_ "encoding/json"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
)

const topic string = "api-to-index"

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

type urlToIndex struct {
	url string
}

type Controller struct {
	producer *Producer
}

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func (p *Producer) Send(s []byte) error {

	delivery_channel := make(chan kafka.Event)
	p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic,
			Partition: kafka.PartitionAny},
		Value: s,
	}, delivery_channel)
	r := <-delivery_channel
	m := r.(*kafka.Message)

	return m.TopicPartition.Error
}

func (p *Producer) Close() {
	p.producer.Flush(1 * 1000)
	p.producer.Close()
}

func NewProducer(topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	return &Producer{p, topic}, nil
}

func NewController(p *Producer) (c *Controller) {
	return &Controller{producer: p}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//todo const topic
func main() {
	p, err := NewProducer(topic)
	if err != nil {
		panic(err)
	}

	defer p.Close()
	controller := NewController(p)

	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("static/index.html")),
	}
	e.Renderer = renderer

	e.GET("/urltoparse", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
	e.POST("/urltoparse", controller.sendUrlToIndex)
	e.Logger.Fatal(e.Start(":9001"))
}

// get JSON (id and URL) from POST request and send it to kafka?
func (co *Controller) sendUrlToIndex(c echo.Context) error {
	var urls urlToIndex
	urls.url = c.FormValue("url")
	log.Printf("URL: %v", urls.url)

	err := co.producer.Send([]byte(urls.url))

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "index.html", nil)
}