package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testServiceFile = "test-fixtures/service.yml"
	testDBFile      = "test-fixtures/db.yml"
)

func Test_Compose(t *testing.T) {
	c, err := NewCompose(".", nil, testServiceFile, testDBFile)
	if err != nil {
		t.Fatal(err)
	}

	config := c.Config()
	assert.NotNil(t, config)

	assert.Equal(t, "3.0", config.Version)

	have := make(map[string]struct{})
	for _, svc := range config.Services {
		have[svc.Name] = struct{}{}
	}

	for _, n := range []string{"db", "app", "ui"} {
		assert.NotNil(t, have[n])
	}
}
