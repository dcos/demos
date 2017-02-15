// Copyright 2017 Mesosphere. All Rights Reserved.
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
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"sync"
)

const (
	onversion string = "0.1.0"
)

var (
	version bool
	wg      sync.WaitGroup
	// FQDN/IP + port of a Kafka broker:
	broker string
	// the data point ingestion queue:
	tqueue chan TrafficData
	rmd    MetaTrafficData
)

func init() {
	flag.BoolVar(&version, "version", false, "Display version information")
	flag.StringVar(&broker, "broker", "", "The FQDN or IP address and port of a Kafka broker. Example: broker-1.kafka.mesos:9382 or 10.0.3.178:9398")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [args]\n\n", os.Args[0])
		fmt.Println("Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// creating the buffered channel holding up to 10 traffic data points:
	tqueue = make(chan TrafficData, 10)
}

func servecontent() {
	fileServer := http.FileServer(http.Dir("content/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		t := <-tqueue
		log.WithFields(log.Fields{"func": "servecontent"}).Info(fmt.Sprintf("Serving data: %v records", len(t.Result.Records)))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		// join with route and metrics data and create geomap overlay data
		gmlist := []GeoMarker{}
		for _, record := range t.Result.Records {
			gm := lookup(record.ID)
			gm.TimeStamp = record.TimeStamp
			gm.VehicleCount = record.VehicleCount
			gmlist = append(gmlist, gm)
			log.WithFields(log.Fields{"func": "servecontent"}).Info(fmt.Sprintf("%v", gm))
		}
		json.NewEncoder(w).Encode(gmlist)

	})
	log.WithFields(log.Fields{"func": "servecontent"}).Info("Starting app server")
	http.ListenAndServe(":8080", nil)
}

func consume(topic string) {
	var consumer sarama.Consumer
	var partitionConsumer sarama.PartitionConsumer
	var err error
	defer wg.Done()
	if consumer, err = sarama.NewConsumer([]string{broker}, nil); err != nil {
		log.WithFields(log.Fields{"func": "consume"}).Error(err)
		return
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			log.WithFields(log.Fields{"func": "consume"}).Error(err)
		}
	}()

	if partitionConsumer, err = consumer.ConsumePartition(topic, 0, sarama.OffsetNewest); err != nil {
		log.WithFields(log.Fields{"func": "consume"}).Error(err)
		return
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.WithFields(log.Fields{"func": "consume"}).Error(err)
		}
	}()

	for {
		msg := <-partitionConsumer.Messages()
		traw := string(msg.Value)
		t := frommsg(traw)
		tqueue <- t
		log.WithFields(log.Fields{"func": "consume"}).Debug(fmt.Sprintf("%#v", t))
		log.WithFields(log.Fields{"func": "servecontent"}).Info("Received data from Kafka and queued it")
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

	// pull static route and metrics data from Minio:
	rmd = syncstaticdata()

	// serve static content (HTML page with OSM overlay) from /static
	// and traffic data in JSON from /data endpoint
	go servecontent()

	// kick of consuming data from Kafka:
	wg.Add(1)
	go consume("trafficdata")
	wg.Wait()
}
