package dockerfile

import (
	"io"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Dockerfile profiles helper functions to inspect a dockerfile
type Dockerfile struct {
	filename string
	stages   []instructions.Stage
}

// NewDockerfile parses the given dockerfiler and returns a Dockerfile
// object
func NewDockerfile(filename string) (*Dockerfile, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	df, err := ParseDockerfile(fh)
	if err == nil {
		df.filename = filename
	}

	return df, err
}

// ParseDockerfile parses a dockerfile into an object from the reader
func ParseDockerfile(rd io.Reader) (*Dockerfile, error) {
	p, err := parser.Parse(rd)
	if err == nil {
		var stages []instructions.Stage
		stages, _, err = instructions.Parse(p.AST)
		if err == nil {
			return &Dockerfile{stages: stages}, nil
		}
	}
	return nil, err
}

// Stages returns all stages in the dockerfile
func (d *Dockerfile) Stages() []instructions.Stage {
	return d.stages
}

// BuildArgs returns all build args
func (d *Dockerfile) BuildArgs() map[string]string {
	out := make(map[string]string)

	for _, s := range d.stages {
		for _, c := range s.Commands {
			name := c.Name()
			if name == "arg" {
				cmd := c.(*instructions.ArgCommand)
				if cmd.Value == nil {
					out[cmd.Key] = ""
				} else {
					out[cmd.Key] = *cmd.Value
				}

			}
		}
	}

	return out
}

// Ports returns the exposed ports. These are all expose statements
// in the last stage
func (d *Dockerfile) Ports() []string {

	ports := make([]string, 0)

	last := d.stages[len(d.stages)-1]
	for _, cmd := range last.Commands {
		name := cmd.Name()
		if name == "expose" {
			expose := cmd.(*instructions.ExposeCommand)
			ports = append(ports, expose.Ports...)
		}
	}

	return ports
}

// Volumes returns a list of volumes for the final stage ie. artifact
func (d *Dockerfile) Volumes() []string {

	vols := make([]string, 0)

	last := d.stages[len(d.stages)-1]
	for _, cmd := range last.Commands {
		name := cmd.Name()
		if name == "volume" {
			vol := cmd.(*instructions.VolumeCommand)
			vols = append(vols, vol.Volumes...)
		}
	}

	return vols
}

// EnvVars returns a map of env variables in the final stage ie. artifact
func (d *Dockerfile) EnvVars() map[string]string {
	out := make(map[string]string)

	last := d.stages[len(d.stages)-1]
	for _, cmd := range last.Commands {
		name := cmd.Name()
		if name == "env" {
			env := cmd.(*instructions.EnvCommand)
			for _, kv := range env.Env {
				out[kv.Key] = kv.Value
			}
		}
	}

	return out
}

// BaseImage returns the images use to build the final output artifact
func (d *Dockerfile) BaseImage() string {
	last := d.stages[len(d.stages)-1]
	return last.BaseName
}

// Image returns the name of the image for the last stage.  This is the
// key after AS
func (d *Dockerfile) Image() string {
	last := d.stages[len(d.stages)-1]
	return last.Name
}

// BuildBaseImages returns all base images for all but the
// last stage
func (d *Dockerfile) BuildBaseImages() []string {
	stages := d.stages[:len(d.stages)-1]
	out := make([]string, len(stages))
	for i, stage := range stages {
		out[i] = stage.BaseName
	}
	return out
}
