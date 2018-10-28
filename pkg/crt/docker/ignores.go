package docker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DockerIgnoresFile is the docker ignore filename
const DockerIgnoresFile = ".dockerignore"

// ParseIgnoresFile reads and parses the ignores file from the directory
func ParseIgnoresFile(dir string) ([]string, error) {
	fpath := filepath.Join(dir, DockerIgnoresFile)
	_, err := os.Stat(fpath)
	if err != nil {
		return []string{}, nil
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}

func dedup(ign ...[]string) []string {
	m := make(map[string]struct{})
	for _, sl := range ign {
		for _, ign := range sl {
			m[ign] = struct{}{}
		}
	}

	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
