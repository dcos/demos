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
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"os"
)

const (
	VERSION string = "0.1.0"
)

var (
	version bool
	cities  = []string{}
	// FQDN/IP + port of a Kafka broker:
	broker string
)

func about() {
	fmt.Printf("\nThis is the fintrans InfluxDB ingestion util in version %s\n", VERSION)
}

func init() {
	cities = []string{
		"London",
		"NYC",
		"SF",
		"Moscow",
		"Tokyo",
	}
	flag.BoolVar(&version, "version", false, "Display version information")
	flag.StringVar(&broker, "broker", "", "The FQDN or IP address and port of a Kafka broker. Example: broker-1.kafka.mesos:9382 or 10.0.3.178:9398")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [args]\n\n", os.Args[0])
		fmt.Println("Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func ingest(topic string) {
	var consumer sarama.Consumer
	if c, err := sarama.NewConsumer([]string{broker}, nil); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		consumer = c
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()

	if partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		for {
			msg := <-partitionConsumer.Messages()
			log.Info(fmt.Sprintf("%s: %#v", msg.Topic, string(msg.Value)))
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

	for _, city := range cities {
		go ingest(city)
	}

	for {
	}
}
