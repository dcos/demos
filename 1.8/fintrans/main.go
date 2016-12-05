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
	"strconv"
)

const (
	VERSION string = "0.1.0"
)

var (
	DEBUG    bool
	version  bool
	broker   string
	producer sarama.SyncProducer
)

func about() {
	fmt.Printf("\nThis is the fintrans generator in version %s\n", VERSION)
}

func init() {
	flag.BoolVar(&version, "version", false, "Display version information")
	flag.StringVar(&broker, "broker", "", "The FQDN or IP address and port of a Kafka broker. Example: broker-1.kafka.mesos:9382 or 10.0.3.178:9398")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [args]\n\n", os.Args[0])
		fmt.Println("Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	DEBUG = false
	if envd := os.Getenv("DEBUG"); envd != "" {
		if d, err := strconv.ParseBool(envd); err == nil {
			DEBUG = d
		}
	}
}

func main() {
	if version || broker == "" {
		about()
		os.Exit(0)
	}
	if broker == "" {
		flag.Usage()
		os.Exit(1)
	}

	if p, err := sarama.NewSyncProducer([]string{broker}, nil); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		producer = p
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()

	msg := &sarama.ProducerMessage{Topic: "fintrans", Value: sarama.StringEncoder("123")}
	if partition, offset, err := producer.SendMessage(msg); err != nil {
		log.Error("Failed to send message ", err)
		os.Exit(1)
	} else {
		log.Info("Message sent to partition ", partition, " at offset ", offset)
	}
}
