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
	"testing"
)

func TestCanParseODAData(t *testing.T) {
	d := TrafficData{
		Result: TrafficDataResult{
			Fields: []Field{
				Field{
					Type: "text",
					ID:   "status",
				},
				Field{
					Type: "int4",
					ID:   "avgMeasuredTime",
				},
				Field{
					Type: "timestamp",
					ID:   "TIMESTAMP",
				},
				Field{
					Type: "int4",
					ID:   "medianMeasuredTime",
				},
				Field{
					Type: "int4",
					ID:   "avgSpeed",
				},
				Field{
					Type: "int4",
					ID:   "vehicleCount",
				},
				Field{
					Type: "int4",
					ID:   "_id",
				},
				Field{
					Type: "int4",
					ID:   "REPORT_ID",
				},
			},
			Records: []Record{
				Record{
					Status:       "OK",
					AvgMTime:     1,
					TimeStamp:    "2017-01-13T00:00:00",
					MMTime:       1,
					AvgSpeed:     1,
					VehicleCount: 1,
					RecordID:     1,
					ID:           1,
				},
			},
		},
	}

	expected := `{"result":{"fields":[{"type":"text","id":"status"},{"type":"int4","id":"avgMeasuredTime"},{"type":"timestamp","id":"TIMESTAMP"},{"type":"int4","id":"medianMeasuredTime"},{"type":"int4","id":"avgSpeed"},{"type":"int4","id":"vehicleCount"},{"type":"int4","id":"_id"},{"type":"int4","id":"REPORT_ID"}],"records":[{"status":"OK","avgMeasuredTime":1,"TIMESTAMP":"2017-01-13T00:00:00","medianMeasuredTime":1,"avgSpeed":1,"vehicleCount":1,"_id":1,"REPORT_ID":1}]}}`
	actual := tomsg(d)
	if actual != expected {
		t.Error("Test failed:\n", expected, "\n\n", actual)
	}
}
