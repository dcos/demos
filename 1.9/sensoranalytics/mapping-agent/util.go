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
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/minio/minio-go"
	"os"
	"strconv"
)

func about() {
	fmt.Printf("\nThis is the mapping agent in version %s\n", onversion)
}

func frommsg(raw string) TrafficData {
	td := TrafficData{}
	buf := bytes.NewBufferString(raw)
	if err := json.NewDecoder(buf).Decode(&td); err != nil {
		log.WithFields(log.Fields{"func": "frommsg"}).Error(err)
	}
	return td
}

func syncstaticdata() MetaTrafficData {

	endpoint := ""
	accessKeyID := ""
	secretAccessKey := ""

	if ep := os.Getenv("PUBLIC_AGENT_IP"); ep != "" {
		endpoint = ep
	} else {
		log.WithFields(log.Fields{"func": "syncstaticdata"}).Error("Need $PUBLIC_AGENT_IP environment variable set")
		os.Exit(1)
	}

	if ak := os.Getenv("ACCESS_KEY_ID"); ak != "" {
		accessKeyID = ak
	} else {
		log.WithFields(log.Fields{"func": "syncstaticdata"}).Error("Need $ACCESS_KEY_ID environment variable set")
		os.Exit(1)
	}

	if sk := os.Getenv("SECRET_ACCESS_KEY"); sk != "" {
		secretAccessKey = sk
	} else {
		log.WithFields(log.Fields{"func": "syncstaticdata"}).Error("Need $SECRET_ACCESS_KEY environment variable set")
		os.Exit(1)
	}

	useSSL := false
	bucket := "aarhus"
	object := "route_metrics_data.json"
	mtd := MetaTrafficData{}

	log.WithFields(log.Fields{"func": "syncstaticdata"}).Info(fmt.Sprintf("Trying to retrieve %s/%s from Minio", bucket, object))
	if mc, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL); err != nil {
		log.WithFields(log.Fields{"func": "syncstaticdata"}).Fatal(fmt.Sprintf("%s ", err))
	} else {
		exists, err := mc.BucketExists(bucket)
		if err != nil || !exists {
			log.WithFields(log.Fields{"func": "syncstaticdata"}).Fatal(fmt.Sprintf("%s", err))
		} else {
			if err := mc.FGetObject(bucket, object, "./rmd.json"); err != nil {
				log.WithFields(log.Fields{"func": "syncstaticdata"}).Fatal(fmt.Sprintf("%s", err))
			} else {
				log.WithFields(log.Fields{"func": "syncstaticdata"}).Info(fmt.Sprintf("Retrieved route and metrics from Minio"))
				f, _ := os.Open("./rmd.json")
				defer os.Remove("./rmd.json")
				if err := json.NewDecoder(f).Decode(&mtd); err != nil {
					log.WithFields(log.Fields{"func": "frommsg"}).Fatal(err)
				}
				log.WithFields(log.Fields{"func": "syncstaticdata"}).Info(fmt.Sprintf("%+v", mtd))
			}
		}
	}
	return mtd
}

// Lookup joins a record ID from a traffic datapoint with
// the static route and metrics data set and creates a geo-coded
// marker for consumption in the OSM overlay from it
func lookup(id int) GeoMarker {
	gm := GeoMarker{}
	for _, record := range rmd.Result.Records {
		if record.ID == id {
			gm.Lat, _ = strconv.ParseFloat(record.Lat, 64)
			gm.Lng, _ = strconv.ParseFloat(record.Lng, 64)
			gm.Label = record.Name
			break
		}
	}
	return gm
}
