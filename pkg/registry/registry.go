package registry

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

// Registry implements a registry interface
type Registry interface {
	// ID as specified in the config
	ID() string
	// Registry domain ie. image name prefix e.g. 123456789012.ecr.xxx
	Domain() string
	// Initialize the registry provider
	Init(conf *Config) error
	// Create a new repository
	CreateRepo(string) (interface{}, error)
	// Get repo info
	GetRepo(string) (interface{}, error)
	// Get image manifest
	GetImageManifest(name, tag string) (interface{}, error)
	// Name of the image with the registry. Needed for deployments
	ImageName(string) string
	// Returns a docker AuthConfig
	GetAuthConfig() (types.AuthConfig, error)
}

// New returns a new registry based on the config.  It returns an error if an
// unsupported provider is supplied or fails to initialize the underlying
// registry provider
func New(conf *Config) (Registry, error) {
	var (
		reg Registry
		err error
	)

	switch conf.Provider {
	case "ecr":
		reg = &awsContainerRegistry{}

	case "docker":
		reg = &localDocker{}

	case "dockerhub":
		reg = &dockerHub{}

	default:
		err = fmt.Errorf("unsupported container registry: '%s'", conf.Provider)

	}

	if err != nil {
		return nil, err
	}

	err = reg.Init(conf)
	return reg, err
}
