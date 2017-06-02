package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
)

var usage = "Usage: ksonnet-gen [path to k8s OpenAPI swagger.json] [output dir]"

func main() {
	if len(os.Args) != 3 {
		log.Fatal(usage)
	}

	swaggerPath := os.Args[1]
	text, err := ioutil.ReadFile(swaggerPath)
	if err != nil {
		log.Fatalf("Could not read file at '%s':\n%v", swaggerPath, err)
	}

	// Deserialize the API object.
	s := kubespec.APISpec{}
	err = json.Unmarshal(text, &s)
	if err != nil {
		log.Fatalf("Could not deserialize schema:\n%v", err)
	}

	// Print names of definitions.
	for defName := range s.Definitions {
		fmt.Println(defName)
	}
}

func init() {
	// Get rid of time in logs.
	log.SetFlags(0)
}
