package dockerfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFile = "test-fixtures/dockerfile"

func Test_dockerfile(t *testing.T) {
	df, err := NewDockerfile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testFile, df.filename)
	assert.Equal(t, 4, len(df.BuildArgs()))
	assert.Equal(t, 1, len(df.EnvVars()))
	assert.Equal(t, 2, len(df.Stages()))
	assert.Equal(t, "alpine", df.BaseImage())
	assert.Equal(t, "artifact", df.Image())
	assert.Equal(t, 1, len(df.BuildBaseImages()))
	assert.Equal(t, 1, len(df.Ports()))
	assert.Equal(t, 2, len(df.Volumes()))
}
