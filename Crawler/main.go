package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/swkkd/budget-google/crawler/htmlParser"
	"github.com/swkkd/budget-google/crawler/middleware"
)

func main() {
	//connect to kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "GROUP-1",
		"auto.offset.reset": "smallest"})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		panic(err)
	} else {
		log.Println("KAFKA CONSUMER CREATED!")
	}
	defer consumer.Close()
	err = consumer.Subscribe("api-to-index", nil)
	if err != nil {
		panic(err)
	}

	//connect to elastic search
	middleware.ConnectToES()
	totalCount := 0

	for {

		msg, err := consumer.ReadMessage(1 * 100)
		if err != nil {
			// Errors are informational and automatically handled by the consumer
			continue
		}
		recordValue := string(msg.Value)
		totalCount += 1
		fmt.Printf("Consumed record with value %s... total count: %v\n", recordValue, totalCount)
		// html, err := htmlParser.Parser(recordValue)
		// if err != nil {
		// 	log.Printf("ERROR PARSING: %s", err)
		// }
		//url, body := htmlParser.HTMLToReadable(string(recordValue))
		// log.Printf("INSERTING INTO DB: --[URL: %s]-- \n --[%s]--", recordValue, html)
		body := htmlParser.HtmlToReadable(recordValue)
		fmt.Println(body)
		middleware.Insert(recordValue, body)
		//middleware.ConnectToES(string(html), string(url))
		//fmt.Printf("URL: %s %s", string(url), string(body))
	}

}
