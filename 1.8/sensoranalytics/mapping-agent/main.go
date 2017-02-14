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
	// "encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	onversion string = "0.1.0"
)

var (
	version bool
	// FQDN/IP + port of a Kafka broker:
	broker string
	// how many seconds to wait between generating a message (default is 2):
	pollwaitsec time.Duration
)

func about() {
	fmt.Printf("\nThis is the mapping agent in version %s\n", onversion)
}

func init() {

	flag.BoolVar(&version, "version", false, "Display version information")
	// flag.StringVar(&broker, "broker", "", "The FQDN or IP address and port of a Kafka broker. Example: broker-1.kafka.mesos:9382 or 10.0.3.178:9398")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [args]\n\n", os.Args[0])
		fmt.Println("Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	pollwaitsec = 10
	if pw := os.Getenv("POLL_WAIT_SEC"); pw != "" {
		if pwi, err := strconv.Atoi(pw); err == nil {
			pollwaitsec = time.Duration(pwi)
		}
	}

}

func main() {
	if version {
		about()
		os.Exit(0)
	}
	// if broker == "" {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }
	fileServer := http.FileServer(http.Dir("content/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.ListenAndServe(":8080", nil)

}
