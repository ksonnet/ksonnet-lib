package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	k8s "github.com/ksonnet/ksonnet-lib/ksonnet-gen/k8sSwagger"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ksonnet"
)

var usage = "Usage: ksonnet-gen [path to k8s swagger-spec folder] [output dir]"

func main() {
	if len(os.Args) != 3 {
		log.Fatal(usage)
	}

	writeOutAppSpec("v1.json", "core.v1.libsonnet")
	writeOutAppSpec("apps_v1beta1.json", "apps.v1beta1.libsonnet")
	writeOutAppSpec("extensions_v1beta1.json", "extensions.v1beta1.libsonnet")
}

func writeOutAppSpec(sourceFilename, destFilename string) {
	filename := fmt.Sprintf("%s/%s", os.Args[1], sourceFilename)
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file at '%s'", filename)
	}

	appSpec, err := k8s.AppSpecFromJson(text)
	if err != nil {
		log.Fatalf(
			"Could not deserialize swagger spec at '%s':\n%v",
			filename,
			err)
	}

	bytes, err := ksonnet.Marshal(appSpec)
	if err != nil {
		log.Fatalf(
			"Failed to generate '%s' with error:\n%v\n", destFilename, err)
	}
	err = ioutil.WriteFile(os.Args[2]+"/"+destFilename, bytes, 0644)
	if err != nil {
		log.Fatalf(
			"Failed to generate '%s' with error:\n%v\n", destFilename, err)
	}
}

func init() {
	// Get rid of time in logs.
	log.SetFlags(0)
}
