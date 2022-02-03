package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/swkkd/budget-google/crawler/middleware"
	"log"
	"os"
)

var kafkaServer, kafkaTopic string

func init() {
	kafkaServer = readFromENV("KAFKA_BROKER", "localhost:29092")
	kafkaTopic = readFromENV("KAFKA_TOPIC", "api-to-index")

	fmt.Println("Kafka Broker - ", kafkaServer)
	fmt.Println("Kafka topic - ", kafkaTopic)
}
func main() {
	//connect to kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
		"group.id":          "GROUP-1",
		"auto.offset.reset": "smallest"})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		panic(err)
	} else {
		log.Println("KAFKA CONSUMER CREATED!")
	}
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Printf("consumer Close faulure! %v", err)
		}
	}(consumer)
	err = consumer.Subscribe(kafkaTopic, nil)
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

		body, err := middleware.HtmlToReadable(recordValue)
		if err != nil {
			// Errors are informational and automatically handled by the consumer
			continue
		}
		//fmt.Println(body)
		middleware.Insert(recordValue, body)

	}

}
func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

//todo prometheus metrics
//todo
