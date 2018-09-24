package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testFile = "test-fixtures/docker-compose.yml"
)

func Test_Compose(t *testing.T) {
	c, err := NewCompose(".", nil, testFile)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, c.Config())
}
