package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	conf := &Config{
		Provider: "ecr",
	}
	reg, err := New(conf)
	assert.Nil(t, err)

	conf.Config = map[string]interface{}{"region": "us-west-2"}
	reg, err = New(conf)
	assert.Nil(t, err)
	r := reg.(*awsContainerRegistry)
	assert.Equal(t, "us-west-2", *r.sess.Config.Region)

	conf.Provider = "unsupported"
	_, err = New(conf)
	assert.Contains(t, err.Error(), "unsupported")
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Registry_ECR(t *testing.T) {
	conf := &Config{
		Provider: "ecr",
		Config: map[string]interface{}{
			"region": "us-west-2",
		},
	}
	// 	for k, v := range ecrCreds {
	// 		conf.Config[k] = v
	// 	}

	treg, _ := New(conf)
	reg := treg.(*awsContainerRegistry)
	// assert.NotEmpty(t, os.Getenv("AWS_ACCESS_KEY_ID"))
	// assert.NotEmpty(t, os.Getenv("AWS_SECRET_ACCESS_KEY"))

	_, err := reg.GetRepo("does-not-exist")
	assert.NotNil(t, err)

	// 	repoName := "test-repo/test-sub"
	// 	_, err = reg.CreateRepo(repoName)
	// 	assert.Nil(t, err)

	// 	_, err = reg.GetImageManifest("test-sub", "notfound")
	// 	assert.NotNil(t, err)

	// 	_, err = reg.GetRepo(repoName)
	// 	assert.Nil(t, err)
	// 	_, err = reg.DeleteRepo(repoName)
	// 	assert.Nil(t, err)

	// 	_, err = reg.GetRepo(repoName)
	// 	assert.NotNil(t, err)
	// 	// nr := nrepo.(*ecr.Repository)
}
