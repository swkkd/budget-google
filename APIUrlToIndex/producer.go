package main

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Controller struct {
	producer *Producer
}

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func (p *Producer) Send(s []byte) error {

	deliveryChannel := make(chan kafka.Event)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic,
			Partition: kafka.PartitionAny},
		Value: s,
	}, deliveryChannel)

	if err != nil {
		return err
	}

	r := <-deliveryChannel
	m := r.(*kafka.Message)

	return m.TopicPartition.Error
}

func (p *Producer) Close() {
	p.producer.Flush(1 * 1000)
	p.producer.Close()
}

func NewProducer(topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	if err != nil {
		return nil, err
	}
	return &Producer{p, topic}, nil
}

func NewController(p *Producer) (c *Controller) {
	return &Controller{producer: p}
}
