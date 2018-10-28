package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/euforia/cchain/pkg/dockerfile"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	testMultiStageDockerfile = "multi-stage.Dockerfile"
)

func Test_Docker_Build(t *testing.T) {
	_, err := os.Stat("/var/run/docker.sock")
	if err != nil {
		t.Skipf("Skipping: docker file descriptor: %v", err)
	}

	dkr, _ := NewDocker()

	req := &BuildRequest{
		ContextDir: "./test-fixtures",
		BuildOpts: &types.ImageBuildOptions{
			Dockerfile: testMultiStageDockerfile,
			NoCache:    true,
		},
	}

	resp, err := dkr.Build(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Context.Close()

	// defer resp.RawLog.Close()
	blr := NewStructedBuildLogReader(resp.RawLog, os.Stdout)
	err = blr.Read()
	if err != nil {
		t.Fatal(err)
	}

	dfPath := filepath.Join("./test-fixtures", testMultiStageDockerfile)
	df, err := dockerfile.NewDockerfile(dfPath)
	if err != nil {
		t.Fatal(err)
	}

	stages := df.Stages()
	c := len(stages)
	var i int
	for _, s := range stages {
		c += len(s.Commands)

		step := blr.Steps[i]
		assert.Contains(t, step.Cmd, s.BaseName)
		i++
		for _, cmd := range s.Commands {
			s := fmt.Sprint(cmd)
			step = blr.Steps[i]
			assert.Equal(t, s, step.Cmd)
			assert.Equal(t, i+1, step.Index)
			i++
		}
	}
	assert.Equal(t, c, len(blr.Steps))
}
