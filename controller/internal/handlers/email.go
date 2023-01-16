package handlers

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func WriteToTopic(topic string, message string) {
	w := kafka.Writer{
		Topic: topic,
		Addr:  kafka.TCP("kafka:9092"),
	}

	if err := w.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte("push"),
			Value: []byte(message),
		},
	); err != nil {
		log.Fatal(err)
	}
	log.Println("Message was sent to kafka!")
	if err := w.Close(); err != nil {
		fmt.Println("failed to close writer: ", err)
	}
	log.Println("Writer closed!")
}
