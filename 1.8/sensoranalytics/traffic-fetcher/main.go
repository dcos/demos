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
	"strconv"
	"time"
)

const (
	onversion string = "0.1.0"
	dataurl   string = "https://www.odaa.dk/api/action/datastore_search?resource_id=b3eeb0ff-c8a8-4824-99d6-e0a3747c8b0d&limit=10"
)

var (
	version bool
	// FQDN/IP + port of a Kafka broker:
	broker string
	// how many seconds to wait between generating a message (default is 2):
	genwaitsec time.Duration
	// the HTTP client used to communicate with the Open Data Aarhus API:
	httpc http.Client
	// the Kafka producer:
	producer sarama.SyncProducer
)

// TrafficData is the payload from the
// Open Data Aarhus data API
type TrafficData struct {
	Result TrafficDataResult `json:"result"`
}

// TrafficDataResult holds the query result
type TrafficDataResult struct {
	Fields  []Field  `json:"fields"`
	Records []Record `json:"records"`
}

// Field is the schema pair
type Field struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Record is the actual data record
type Record struct {
	Status       string `json:"status"`
	AvgMTime     int    `json:"avgMeasuredTime"`
	TimeStamp    string `json:"TIMESTAMP"`
	MMTime       int    `json:"medianMeasuredTime"`
	AvgSpeed     int    `json:"avgSpeed"`
	VehicleCount int    `json:"vehicleCount"`
	RecordID     int    `json:"_id"`
	ID           int    `json:"REPORT_ID"`
}

func about() {
	fmt.Printf("\nThis is the traffic data fetcher in version %s\n", onversion)
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

	genwaitsec = 10
	if gw := os.Getenv("GEN_WAIT_SEC"); gw != "" {
		if gwi, err := strconv.Atoi(gw); err == nil {
			genwaitsec = time.Duration(gwi)
		}
	}

	httpc = http.Client{Timeout: 10 * time.Second}
}

func pulldata(target interface{}) error {
	r, err := httpc.Get(dataurl)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
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
		d := TrafficData{}
		pulldata(&d)
		log.Info(fmt.Sprintf("%#v", d))

		// msg := &sarama.ProducerMessage{Topic: "trafficdata", Value: sarama.StringEncoder(rawmsg)}
		// if _, _, err := producer.SendMessage(msg); err != nil {
		// 	log.Error("Failed to send message ", err)
		// } else {
		// 	log.Info(fmt.Sprintf("%#v", msg))
		// }
		time.Sleep(genwaitsec * time.Second)
	}
}
