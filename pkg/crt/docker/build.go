package docker

import (
	"github.com/docker/docker/pkg/archive"
	"golang.org/x/net/context"
)

// Build builds an image with the request params
func (orch *Docker) Build(ctx context.Context, req *BuildRequest) (*BuildResponse, error) {
	err := req.SetIgnores()
	if err != nil {
		return nil, err
	}

	rdc, err := archive.TarWithOptions(req.ContextDir, req.TarOpts)
	if err != nil {
		return nil, err
	}
	defer rdc.Close()

	resp, err := orch.cli.ImageBuild(ctx, rdc, *req.BuildOpts)
	if err != nil {
		return nil, err
	}

	return &BuildResponse{
		RawLog:  resp.Body,
		Context: rdc,
	}, nil
}
