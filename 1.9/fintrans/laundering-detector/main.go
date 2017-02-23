// Copyright 2016 Mesosphere. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
)

const (
	VERSION         string = "0.1.0"
	MAX_ACCOUNT_NUM int    = 1000
)

var (
	version bool
	wg      sync.WaitGroup
	cities  []string
	// the aggregate transaction volume [source][target]:
	volume [MAX_ACCOUNT_NUM][MAX_ACCOUNT_NUM]int
	// FQDN/IP + port of a Kafka broker:
	broker string
	// the ingestion queue:
	iqueue chan Transaction
	// aggregate volume alert threshold:
	threshold int
	// defines how the output is rendered:
	prodout bool
)

type Transaction struct {
	City   string
	Source string
	Target string
	Amount int
}

func about() {
	fmt.Printf("\nThis is the fintrans money laundering detection consumer in version %s\n", VERSION)
}

func init() {
	cities = []string{
		"London",
		"NYC",
		"SF",
		"Moscow",
		"Tokyo",
	}
	// the commend line parameters:
	flag.BoolVar(&version, "version", false, "Display version information")
	flag.StringVar(&broker, "broker", "", "The FQDN or IP address and port of a Kafka broker. Example: broker-1.kafka.mesos:9382 or 10.0.3.178:9398")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [args]\n\n", os.Args[0])
		fmt.Println("Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// the optional environment variables:
	threshold = 10000
	if at := os.Getenv("ALERT_THRESHOLD"); at != "" {
		if ati, err := strconv.Atoi(at); err == nil {
			threshold = ati
		}
	}

	prodout = true
	if po := os.Getenv("PROD_OUTPUT"); po != "" {
		if strings.ToLower(po) == "false" {
			prodout = false
		}
	}

	// creating the aggregate transaction volume matrix with format [source][target]:
	volume = [MAX_ACCOUNT_NUM][MAX_ACCOUNT_NUM]int{}
	// creating the buffered channel holding up to 100 transactions:
	iqueue = make(chan Transaction, 100)
}

func detect() {
	for {
		t := <-iqueue
		log.WithFields(log.Fields{"func": "detect"}).Info(fmt.Sprintf("Dequeued %#v", t))
		source := -1
		if s, err := strconv.Atoi(t.Source); err == nil {
			source = s
		}
		target := -1
		if t, err := strconv.Atoi(t.Target); err == nil {
			target = t
		}
		if source != -1 && target != -1 { // we have a valid transaction
			volume[source][target] = volume[source][target] + t.Amount
			log.WithFields(log.Fields{"func": "detect"}).Info(fmt.Sprintf("%s -> %s totalling %d now", t.Source, t.Target, volume[source][target]))
			if volume[source][target] > threshold {
				if prodout {
					fmt.Println(fmt.Sprintf("POTENTIAL MONEY LAUNDERING: %s -> %s totalling %d now ", t.Source, t.Target, volume[source][target]))
				} else {
					fmt.Println(fmt.Sprintf("\x1b[31;1mPOTENTIAL MONEY LAUNDERING:\x1b[0m %s -> %s totalling %d now ", t.Source, t.Target, volume[source][target]))
				}
			}
		}
		log.WithFields(log.Fields{"func": "detect"}).Info(fmt.Sprintf("Current queue length: %d", len(iqueue)))
	}
}

func consume(topic string) {
	var consumer sarama.Consumer
	defer wg.Done()
	if c, err := sarama.NewConsumer([]string{broker}, nil); err != nil {
		log.WithFields(log.Fields{"func": "consume"}).Error(err)
		return
	} else {
		consumer = c
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.WithFields(log.Fields{"func": "consume"}).Error(err)
		}
	}()

	if partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest); err != nil {
		log.WithFields(log.Fields{"func": "consume"}).Error(err)
		return
	} else {
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.WithFields(log.Fields{"func": "consume"}).Error(err)
			}
		}()
		for {
			msg := <-partitionConsumer.Messages()
			traw := strings.Split(string(msg.Value), " ")
			amount := 0
			if a, err := strconv.Atoi(traw[2]); err == nil {
				amount = a
			}
			t := Transaction{City: msg.Topic, Source: traw[0], Target: traw[1], Amount: amount}
			iqueue <- t
			log.WithFields(log.Fields{"func": "consume"}).Info(fmt.Sprintf("Queued %#v", t))
		}
	}
}

func main() {
	if version {
		about()
		os.Exit(0)
	}
	if broker == "" {
		flag.Usage()
		os.Exit(1)
	}
	wg.Add(len(cities))
	go detect()
	for _, city := range cities {
		go consume(city)
	}
	wg.Wait()
}
