package compose

import (
	"github.com/docker/cli/cli/compose/types"
)

// Compose holds a set of parsed composed file that belong together
type Compose struct {
	workdir string
	files   []string
	env     map[string]string
	config  *types.Config
}

// NewCompose parses the given docker compose files and returns a new Compose object given
// the settings
func NewCompose(workdir string, env map[string]string, files ...string) (*Compose, error) {

	finalEnv, err := buildEnvironment()
	if err != nil {
		return nil, err
	}

	if len(env) > 0 {
		for k, v := range env {
			finalEnv[k] = v
		}
	}

	config, err := parseDockerComposeFiles(finalEnv, workdir, files...)
	if err != nil {
		return nil, err
	}

	c := &Compose{
		workdir: workdir,
		files:   files,
		env:     finalEnv,
		config:  config,
	}

	return c, nil
}

// Config returns a docker-compose config
func (c *Compose) Config() *types.Config {
	return c.config
}

func parseDockerComposeFiles(env map[string]string, workdir string, files ...string) (*types.Config, error) {
	// Load first compose
	config, err := parseDockerComposeFile(env, workdir, files[0])
	if err != nil {
		return nil, err
	}

	// Load all others merging back to the config above to produce
	// one config
	for _, f := range files[1:] {
		var conf *types.Config
		conf, err = parseDockerComposeFile(env, workdir, f)
		if err != nil {
			break
		}

		config, err = mergeComposeObject(config, conf)
		if err != nil {
			break
		}
	}

	return config, err
}
