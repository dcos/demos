package main

import (
	// "fmt"
	"github.com/Shopify/sarama"
	// "os"
	"log"
)

func main() {

	producer, err := sarama.NewSyncProducer([]string{"10.0.3.178:9398"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	msg := &sarama.ProducerMessage{Topic: "fintrans", Value: sarama.StringEncoder("123")}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}
}
