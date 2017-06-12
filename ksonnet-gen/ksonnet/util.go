package ksonnet

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
)

const constructorName = "new"

var specialProperties = map[kubespec.PropertyName]kubespec.PropertyName{
	"apiVersion": "apiVersion",
	"metadata":   "metadata",
	"kind":       "kind",
}

func isSpecialProperty(pn kubespec.PropertyName) bool {
	_, ok := specialProperties[pn]
	return ok
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
