package ksonnet

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
)

var specialProperties = map[kubespec.PropertyName]kubespec.PropertyName{
	"apiVersion": "apiVersion",
	"metadata":   "metadata",
	"kind":       "kind",
}

func isSpecialProperty(pn kubespec.PropertyName) bool {
	_, ok := specialProperties[pn]
	return ok
}

func toJsonnetName(ok kubespec.ObjectKind) kubespec.ObjectKind {
	if len(ok) == 0 {
		log.Fatalf("Can't lowercase first letter of 0-rune string")
	}
	kindString := string(ok)

	upper := strings.ToLower(kindString[:1])
	return kubespec.ObjectKind(upper + kindString[1:])
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
