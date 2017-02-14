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
	endpoint := "34.250.247.12"
	accessKeyID, secretAccessKey := "F3QE89J9WPSC49CMKCCG", "2/parG/rllluCLMgHeJggJfY9Pje4Go8VqOWEqI9"
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
