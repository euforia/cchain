package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/archive"
)

// ContainerCreateConfig is used to create a container
type ContainerCreateConfig struct {
	Container *container.Config
	Host      *container.HostConfig
	Network   *network.NetworkingConfig
}

// BuildRequest is a request to build an image
type BuildRequest struct {
	ContextDir string
	TarOpts    *archive.TarOptions
	BuildOpts  *types.ImageBuildOptions
	// Output     io.Writer
}

type BuildResponse struct {
	Context io.ReadCloser
	RawLog  io.ReadCloser
}

// SetIgnores reads the .dockerignore file if it exists and appends them to the
// exclude pattern for the tar options
func (req *BuildRequest) SetIgnores() error {
	ign, err := ParseIgnoresFile(req.ContextDir)
	if err != nil {
		return err
	}

	if req.TarOpts == nil {
		req.TarOpts = &archive.TarOptions{}
	}

	if len(ign) > 0 {
		req.TarOpts.ExcludePatterns = dedup(req.TarOpts.ExcludePatterns, ign)
	}

	return nil
}

// PushRequest is a container image push request
type PushRequest struct {
	Image   string
	Tag     string
	Output  io.Writer
	Options types.ImagePushOptions
}

// PullRequest ..
type PullRequest struct {
	Image   string
	Tag     string
	Output  io.Writer
	Options types.ImagePullOptions
}

// Ref returns the fully qualified image reference
func (req *PullRequest) Ref() string {
	return req.Image + ":" + req.Tag
}
