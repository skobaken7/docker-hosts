package main

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

type listners struct {
	dockerClient *client.Client
}

func newListeners(dockerClient *client.Client) *listners {
	return &listners{dockerClient: dockerClient}
}

func (l *listners) receive(event events.Message) error {
	if event.Type == events.NetworkEventType {
		if event.Action == "connect" {
			return l.networkConnected(event)
		}
		if event.Action == "disconnect" {
			return l.networkDisconnected(event)
		}
	}
	return nil
}

func (l *listners) init() error {
	containers, err := l.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list containers in init")
	}

	for _, container := range containers {
		for _, network := range container.NetworkSettings.Networks {
			if err := l.addHost(container.ID, network.NetworkID); err != nil {
				return errors.Wrapf(err, "failed to add hosts in init; networkID=%s, containerID=%s", network.NetworkID, container.ID)
			}
		}
	}
	return nil
}

func (l *listners) networkConnected(event events.Message) error {
	networkID := event.Actor.ID
	containerID := event.Actor.Attributes["container"]

	return l.addHost(containerID, networkID)
}

func (l *listners) addHost(containerID, networkID string) error {
	container, err := l.dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return errors.Wrapf(err, "failed to get container attributes; containerID=%s", containerID)
	}

	for networkName, network := range container.NetworkSettings.Networks {
		if network.NetworkID == networkID {
			if len(network.Aliases) == 0 {
				return nil
			}

			fqdn := fmt.Sprintf("%s.%s.docker.internal", network.Aliases[0], networkName)
			if err := addHost(containerID, networkID, network.IPAddress, fqdn); err != nil {
				return errors.Wrapf(err, "failed to add hosts; fqdn=%s", fqdn)
			}
		}
	}

	return errors.Wrapf(err, "failed to get network attributes; networkID=%s, containerID=%s", networkID, containerID)
}

func (l *listners) networkDisconnected(event events.Message) error {
	networkID := event.Actor.ID
	containerID := event.Actor.Attributes["container"]

	if err := removeHostByID(containerID, networkID); err != nil {
		return errors.Wrapf(err, "failed to remove hosts; containerID=%s, networkID=%s", containerID, networkID)
	}

	return nil
}
