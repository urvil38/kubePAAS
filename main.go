// Copyright © 2018 Urvil Patel <patelurvil38@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/urvil38/kubepaas/cmd"
	"go.opencensus.io/trace"
)

func main() {
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "kubepaas",
	})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}
	defer sd.Flush()
	trace.RegisterExporter(sd)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1.0)})
	
	cmd.Execute()
}
