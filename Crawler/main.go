package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/swkkd/budget-google/crawler/htmlParser"
	"github.com/swkkd/budget-google/crawler/searchEngine"
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
		fmt.Println("Consumer created successfully")
	}
	defer consumer.Close()
	err = consumer.Subscribe("api-to-index", nil)
	if err != nil {
		panic(err)
	}

	//connect to elastic search
	searchEngine.ConnectToES()

	totalCount := 0

	for {

		msg, err := consumer.ReadMessage(1 * 100)
		if err != nil {
			// Errors are informational and automatically handled by the consumer
			continue
		}
		recordValue := msg.Value

		totalCount += 1
		fmt.Printf("Consumed record with value %s... total count: %v\n", recordValue, totalCount)
		htmlParser.Parser(recordValue)
	}

}
