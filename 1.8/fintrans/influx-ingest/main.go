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
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	onversion     string = "0.1.0"
	prodInfluxAPI string = "http://influxdb.marathon.l4lb.thisdcos.directory:8086"
)

var (
	version bool
	wg      sync.WaitGroup
	cities  []string
	// FQDN/IP + port of a Kafka broker:
	broker string
	// the URL for the InfluxDB HTTP API:
	influxAPI string
	// into which InfluxDB database to ingest transactions:
	targetdb string
	// how many seconds to wait between ingesting transactions:
	ingestwaitsec time.Duration
	// the ingestion queue
	iqueue chan Transaction
)

// Transaction defines a single transaction:
// in which City it originated, the Source
// account the money came from as well as
// the Target account the money goes to and
// last but not least the Amount that was
// transferred in the transaction.
type Transaction struct {
	City   string
	Source string
	Target string
	Amount int
}

func about() {
	fmt.Printf("\nThis is the fintrans InfluxDB ingestion consumer in version %s\n", onversion)
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
	influxAPI = prodInfluxAPI
	if ia := os.Getenv("INFLUX_API"); ia != "" {
		influxAPI = ia
	}
	targetdb = "fintrans"
	if td := os.Getenv("INFLUX_TARGET_DB"); td != "" {
		targetdb = td
	}
	ingestwaitsec = 1
	if iw := os.Getenv("INGEST_WAIT_SEC"); iw != "" {
		if iwi, err := strconv.Atoi(iw); err == nil {
			ingestwaitsec = time.Duration(iwi)
		}
	}
	// creating the buffered channel holding up to 100 transactions:
	iqueue = make(chan Transaction, 100)
}

func ingest() {
	if c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxAPI,
		Username: "root",
		Password: "root",
	}); err != nil {
		log.WithFields(log.Fields{"func": "ingest"}).Error(err)
	} else {
		defer c.Close()

		log.WithFields(log.Fields{"func": "ingest"}).Info("Connected to ", fmt.Sprintf("%#v", c))
		for {
			t := <-iqueue
			log.WithFields(log.Fields{"func": "ingest"}).Info(fmt.Sprintf("Dequeued %#v", t))
			bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  targetdb,
				Precision: "s", // second resultion
			})
			log.WithFields(log.Fields{"func": "ingest"}).Info(fmt.Sprintf("Preparing batch %#v", bp))
			tags := map[string]string{
				"source": t.Source,
				"target": t.Target,
			}
			fields := map[string]interface{}{
				"amount": t.Amount,
			}
			pt, _ := client.NewPoint(t.City, tags, fields, time.Now())
			bp.AddPoint(pt)
			log.WithFields(log.Fields{"func": "ingest"}).Info(fmt.Sprintf("Added point %#v", pt))
			if err := c.Write(bp); err != nil {
				log.WithFields(log.Fields{"func": "ingest"}).Error("Could not ingest transaction: ", err.Error())
			} else {
				log.WithFields(log.Fields{"func": "ingest"}).Info(fmt.Sprintf("Ingested %#v", bp))
			}
			time.Sleep(ingestwaitsec * time.Second)
			log.WithFields(log.Fields{"func": "ingest"}).Info(fmt.Sprintf("Current queue length: %d", len(iqueue)))
		}
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
	go ingest()
	for _, city := range cities {
		go consume(city)
	}
	wg.Wait()
}
