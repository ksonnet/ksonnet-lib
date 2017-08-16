package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ksonnet"
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
	s.Text = text
	s.FilePath = filepath.Dir(swaggerPath)

	// Emit Jsonnet code.
	ksonnetLibSHA := getSHARevision(".")
	k8sSHA := getSHARevision(s.FilePath)
	jsonnetBytes, err := ksonnet.Emit(&s, &ksonnetLibSHA, &k8sSHA)
	if err != nil {
		log.Fatalf("Could not write ksonnet library:\n%v", err)
	}

	// Write out.
	outfile := fmt.Sprintf("%s/%s", os.Args[2], "k8s.libsonnet")
	err = ioutil.WriteFile(outfile, jsonnetBytes, 0644)
	if err != nil {
		log.Fatalf("Could not write `kube.libsonnet`:\n%v", err)
	}
}

func getSHARevision(dir string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could get working directory:\n%v", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		log.Fatalf("Could cd to directory of repository at '%s':\n%v", dir, err)
	}

	sha, err := exec.Command("sh", "-c", "git rev-parse HEAD").Output()
	if err != nil {
		log.Fatalf("Could not find SHA of HEAD:\n%v", err)
	}

	err = os.Chdir(cwd)
	if err != nil {
		log.Fatalf("Could cd back to current directory '%s':\n%v", cwd, err)
	}

	return strings.TrimSpace(string(sha))
}

func init() {
	// Get rid of time in logs.
	log.SetFlags(0)
}
