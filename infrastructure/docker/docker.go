package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/toolfordev/local-cli/application/models"
)

type Docker struct {
	client *client.Client
}

func (docker *Docker) Init() (err error) {
	docker.client, err = client.NewClientWithOpts(client.FromEnv)
	return
}

func (docker *Docker) NetworkCreate(name string) (err error) {
	_, err = docker.client.NetworkCreate(
		context.Background(),
		name,
		types.NetworkCreate{
			CheckDuplicate: true,
		},
	)
	return
}

func (docker *Docker) NetworkRemove(name string) (err error) {
	err = docker.client.NetworkRemove(
		context.Background(),
		name,
	)
	return
}

func (docker *Docker) NetworkConnect(networkName, containerName string) (err error) {
	err = docker.client.NetworkConnect(
		context.Background(),
		networkName,
		containerName,
		nil,
	)
	return
}

func (docker *Docker) ImagePull(name string) (err error) {
	_, err = docker.client.ImagePull(context.Background(), name, types.ImagePullOptions{})
	return
}

func (docker *Docker) ImageIsReady(name string) (isReady bool, err error) {
	images, err := docker.client.ImageList(context.Background(), types.ImageListOptions{Filters: filters.NewArgs(filters.KeyValuePair{Key: "reference", Value: name})})
	isReady = len(images) == 1
	return
}

func (docker *Docker) ContainerCreate(application models.ApplicationConfig) (err error) {
	getEnv := func() []string {
		envs := make([]string, len(application.EnvironmentVariables))
		for i, environmentVariable := range application.EnvironmentVariables {
			envs[i] = fmt.Sprintf("%v=%v", environmentVariable.Name, environmentVariable.Value)
		}
		return envs
	}

	portBindings := make(nat.PortMap, len(application.Ports))
	for _, port := range application.Ports {
		portBindings[nat.Port(fmt.Sprintf("%v/%v", port.Application, port.Protocol))] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: port.External,
			},
		}
	}

	_, err = docker.client.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:    application.Image,
			Env:      getEnv(),
			Hostname: fmt.Sprintf("%v.local.toolfor.dev", application.ContainerName),
		},
		&container.HostConfig{
			RestartPolicy: container.RestartPolicy{Name: "always"},
			PortBindings:  portBindings,
		},
		nil,
		nil,
		application.ContainerName,
	)
	return
}

func (docker *Docker) ContainerStart(application models.ApplicationConfig) (err error) {
	err = docker.client.ContainerStart(context.Background(), application.ContainerName, types.ContainerStartOptions{})
	return
}

func (docker *Docker) ContainerRemove(application models.ApplicationConfig) (err error) {
	err = docker.client.ContainerRemove(
		context.Background(),
		application.ContainerName,
		types.ContainerRemoveOptions{
			Force: true,
		},
	)
	return
}
