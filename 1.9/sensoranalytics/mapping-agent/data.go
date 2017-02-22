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

// MetaTrafficData is the payload from the
// Open Data Aarhus data API static route and metrics dataset
type MetaTrafficData struct {
	Result MetaTrafficDataResult `json:"result"`
}

// MetaTrafficDataResult holds the result
type MetaTrafficDataResult struct {
	Fields  []Field      `json:"fields"`
	Records []MetaRecord `json:"records"`
}

// MetaRecord is the actual data record
type MetaRecord struct {
	Lat      string `json:"POINT_1_LAT"`
	Lng      string `json:"POINT_1_LNG"`
	Name     string `json:"REPORT_NAME"`
	RecordID int    `json:"_id"`
	ID       int    `json:"REPORT_ID"`
}

// GeoMarker carries the data for an OSM overlay marker
type GeoMarker struct {
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Label        string  `json:"label"`
	TimeStamp    string  `json:"timestamp"`
	VehicleCount int     `json:"numvehicles"`
}
