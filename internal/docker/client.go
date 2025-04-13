package docker

import (
	"github.com/docker/docker/client"
	"github.com/trinhminhtriet/docker-activity/pkg/error"
)

// Client wraps the Docker client.
type Client struct {
	*client.Client
}

// NewClient creates a Docker client.
func NewClient() (*Client, error.Error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, error.WrapIO(err)
	}
	return &Client{Client: cli}, nil
}
