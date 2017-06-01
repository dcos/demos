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
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
)

const (
	onversion     string = "0.1.0"
	maxAccountNum int    = 1000
	maxAmount     int    = 10000
)

var (
	version bool
	cities  = []string{}
	// FQDN/IP + port of a Kafka broker:
	broker string
	// how many seconds to wait between generating a message (default is 2):
	genwaitsec time.Duration
	// the Kafka producer:
	producer sarama.SyncProducer
)

func about() {
	fmt.Printf("\nThis is the fintrans generator in version %s\n", onversion)
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

	genwaitsec = 2
	if gw := os.Getenv("GEN_WAIT_SEC"); gw != "" {
		if gwi, err := strconv.Atoi(gw); err == nil {
			genwaitsec = time.Duration(gwi)
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

	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		city := cities[r.Intn(len(cities))]
		source := r.Intn(maxAccountNum)
		target := r.Intn(maxAccountNum)
		amount := r.Intn(maxAmount)
		if source != target {
			rawmsg := fmt.Sprintf("%d %d %d", source, target, amount)
			msg := &sarama.ProducerMessage{Topic: string(city), Value: sarama.StringEncoder(rawmsg)}
			if _, _, err := producer.SendMessage(msg); err != nil {
				log.Error("Failed to send message ", err)
			} else {
				log.Info(fmt.Sprintf("%#v", msg))
			}
		}
		time.Sleep(genwaitsec * time.Second)
	}
}
