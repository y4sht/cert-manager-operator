// Copyright 2020 The Operator-SDK Authors
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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/openshift/cert-manager-operator/test/e2e"
	lib "github.com/openshift/cert-manager-operator/test/library"
	scapiv1alpha3 "github.com/operator-framework/api/pkg/apis/scorecard/v1alpha3"
	apimanifests "github.com/operator-framework/api/pkg/manifests"

	"k8s.io/klog/v2"

	"log"
)

var PodBundleRoot = "/bundle"

func main() {
	entrypoint := os.Args[1:]
	if len(entrypoint) == 2 {
		PodBundleRoot = entrypoint[1]
	}
	// Disable klog for this test
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()

	// Output logs to a buffer instead of stdout.
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	if len(entrypoint) == 0 {
		log.Fatalf("Test name argument is required")
	}
	fmt.Printf("%s\n", logBuf.String())

	// Read the pod's untar'd bundle from a well-known path.
	cfg, err := apimanifests.GetBundleFromDir(PodBundleRoot)
	if err != nil {
		log.Fatalf("Failed to read bundle")
		fmt.Printf("%s\n", logBuf.String())
	}

	var result scapiv1alpha3.TestResult
	tests := e2e.Tests{T: &testing.T{}}
	reflection, err := lib.Invoke(tests, entrypoint[0], cfg)
	if err != nil {
		log.Printf("Failed to invoke test: %v", err)
	}
	result = reflection.Interface().(scapiv1alpha3.TestResult)
	result.Log = logBuf.String()

	// Convert scapiv1alpha3.TestResult to json.
	prettyJSON, err := json.MarshalIndent(lib.WrapResult(result), "", "    ")
	if err != nil {
		log.Fatalf("Failed to generate json")
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
