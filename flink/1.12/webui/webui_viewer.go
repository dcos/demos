
package main

import (
	"fmt"
  	"os"
	"os/signal"
  	"time"

	"github.com/icza/gowut/gwu"
	"github.com/Shopify/sarama"
)

func main() {

	var messages []string

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{"broker.kafka.l4lb.thisdcos.directory:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	topic := "fraud"
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(time.Second * 3)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Println("Detected Fraud: ")
				messg := "Potential Fraud: " + string(msg.Key) + string(msg.Value)
				messages = append(messages, messg)
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			case <- timer.C:
				doneCh <- struct{}{}
			}
		}
	}()
<-doneCh

////////////// WEBUI

	// Get host
	host := os.Getenv("HOST")
	port := os.Getenv("PORT0")

	// Create and build a window
	win := gwu.NewWindow("main", "Fraud Display")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	for _, message := range messages {
		win.Add(gwu.NewLabel(message))
	}



	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("fraud",  host + ":" + port)
	server.SetText("Fraud Display")
	server.AddWin(win)
	server.Start("") // Also opens windows list in browser


	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

}
