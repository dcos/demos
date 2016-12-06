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
	"github.com/influxdata/influxdb/client/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION    string = "0.1.0"
	INFLUX_API string = "http://influxdb.marathon.l4lb.thisdcos.directory:8086"
)

var (
	version bool
	cities  = []string{}
	// FQDN/IP + port of a Kafka broker:
	broker string
	// into which InfluxDB database to ingest transactions:
	targetdb string
	// how many seconds to wait between ingesting transactions:
	ingestwaitsec time.Duration
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
}

func write(c client.Client, transaction string) {
	if bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  targetdb,
		Precision: "s", // second resultion
	}); err != nil {
		log.WithFields(log.Fields{"func": "write"}).Error(err)
	} else {
		log.WithFields(log.Fields{"func": "write"}).Info("Connected to ", fmt.Sprintf("%#v", c))
		t := strings.Split(transaction, " ")
		log.WithFields(log.Fields{"func": "write"}).Info("Extracted ", fmt.Sprintf("%#v", t))
		tags := map[string]string{
			"source": t[0],
			"target": t[1],
		}
		amount := 0
		if a, err := strconv.Atoi(t[2]); err == nil {
			amount = a
		}
		log.WithFields(log.Fields{"func": "write"}).Info(fmt.Sprintf("source=%s target=%s amount=%d", t[0], t[1], amount))
		fields := map[string]interface{}{
			"amount": amount,
		}
		if pt, err := client.NewPoint("transaction", tags, fields, time.Now()); err != nil {
			log.WithFields(log.Fields{"func": "write"}).Error(err)
		} else {
			bp.AddPoint(pt)
			log.WithFields(log.Fields{"func": "write"}).Info(fmt.Sprintf("Added %#v", pt))
		}
		if err := c.Write(bp); err != nil {
			log.WithFields(log.Fields{"func": "write"}).Error(err)
		} else {
			log.WithFields(log.Fields{"func": "write"}).Info(fmt.Sprintf("Wrote %#v", bp))
		}
	}
}

func ingest2Influx(transaction string) {
	log.WithFields(log.Fields{"func": "ingest2Influx"}).Info("Trying to ingest ", transaction)
	if c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     INFLUX_API,
		Username: "root",
		Password: "root",
	}); err != nil {
		log.WithFields(log.Fields{"func": "ingest2Influx"}).Error(err)
	} else {
		write(c, transaction)
	}
}

func consume(topic string) {
	var consumer sarama.Consumer
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
			ingest2Influx(string(msg.Value))
			time.Sleep(ingestwaitsec * time.Second)
			log.WithFields(log.Fields{"func": "consume"}).Info(fmt.Sprintf("%s: %#v", msg.Topic, string(msg.Value)))
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
		go consume(city)
	}

	for {
	}
}
