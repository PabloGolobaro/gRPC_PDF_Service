package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

const topic = "notifications"
const addr = "kafka:9092"

var wg sync.WaitGroup

type Effector func(ctx context.Context) (*kafka.Conn, error)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	retry := Retry(connectToKafka, 3, 4*time.Second)
	conn1, err := retry(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()
	wg.Add(1)

	log.Println("Start notifications!")
	go ReadFromTopic(topic)

	wg.Wait()
}
func ReadFromTopic(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Topic:   topic,
		Brokers: []string{addr},
		GroupID: "1",
	})
	log.Println(r.Offset())
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		log.Printf("message at offset %d: %s = %s\n", m.Offset,
			string(m.Key), string(m.Value))
		log.Println(r.Offset())
	}
	if err := r.Close(); err != nil {
		log.Println("failed to close reader:", err)
	}
}
func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (*kafka.Conn, error) {
		for r := 0; ; r++ {
			responce, err := effector(ctx)
			if err == nil || r >= retries {
				return responce, err
			}
			log.Printf("Attempt %d failde; retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}

func connectToKafka(ctx context.Context) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", addr)

}
