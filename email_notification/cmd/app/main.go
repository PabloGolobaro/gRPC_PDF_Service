package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

const topic = "notifications"

var wg sync.WaitGroup

func main() {
	conn1, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()

	partitions, err := conn1.ReadPartitions()
	if err != nil {
		log.Fatal(err)
	}
	m := map[string]struct{}{}
	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		if k == topic {
			wg.Add(1)
			log.Println(k)
			log.Println("Start notifications!")
			go ReadFromTopic(k)
		}

	}

	wg.Wait()
}
func ReadFromTopic(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Topic:   topic,
		Brokers: []string{"kafka:9092"},
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		log.Printf("message at offset %d: %s = %s\n", m.Offset,
			string(m.Key), string(m.Value))

	}
	if err := r.Close(); err != nil {
		log.Println("failed to close reader:", err)
	}
}
