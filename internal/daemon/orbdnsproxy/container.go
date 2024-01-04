package orbdnsproxy

import (
	"context"
	"fmt"

	"github.com/dchest/uniuri"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type Container struct {
	logger     *zap.Logger
	docker     *client.Client
	configFile string
	hostname   string
	id         string
}

func NewContainer(logger *zap.Logger, docker *client.Client, configFile string, hostname string) (*Container, error) {
	return &Container{
		logger:     logger,
		docker:     docker,
		configFile: configFile,
		hostname:   hostname,
	}, nil
}

func (c *Container) CreateAndStart() error {
	c.logger.Debug("Creating container", zap.String("hostname", c.hostname))

	name := fmt.Sprintf("dotlocal-%s-%s", uniuri.NewLen(6), c.hostname)
	res, err := c.docker.ContainerCreate(context.Background(), &container.Config{
		Image: "nginx:1.24.0-alpine",
		Labels: map[string]string{
			"dev.orbstack.domains": c.hostname,
			"managed-dotlocal":     "true",
		},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: c.configFile,
				Target: "/etc/nginx/conf.d/default.conf",
			},
		},
	}, nil, nil, name)
	if err != nil {
		return err
	}
	c.id = res.ID

	err = c.docker.ContainerStart(context.Background(), c.id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	c.logger.Info("Started container", zap.String("hostname", c.hostname), zap.String("id", c.id))

	return nil
}

func (c *Container) Remove() error {
	if c.id == "" {
		return nil
	}

	c.logger.Debug("Removing container", zap.String("hostname", c.hostname), zap.String("id", c.id))
	err := c.docker.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	c.logger.Info("Removed container", zap.String("id", c.id))
	return nil
}
