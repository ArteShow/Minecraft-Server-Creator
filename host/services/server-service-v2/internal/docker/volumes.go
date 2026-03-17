package docker

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerService struct {
	client *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &DockerService{client: cli}, nil
}

func volumeName(name string) string {
	return "mc_" + name
}

func (d *DockerService) CreateVolume(name string) error {
	_, err := d.client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name: volumeName(name),
		Labels: map[string]string{
			"project": "minecraft-server",
		},
	})
	return err
}

func (d *DockerService) VolumeExists(name string) (bool, error) {
	vols, err := d.client.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, v := range vols.Volumes {
		if v.Name == volumeName(name) {
			return true, nil
		}
	}
	return false, nil
}

func (d *DockerService) ListVolumes() ([]string, error) {
	args := filters.NewArgs()
	args.Add("label", "project=minecraft-server")

	vols, err := d.client.VolumeList(context.Background(), volume.ListOptions{
		Filters: args,
	})
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(vols.Volumes))
	for _, v := range vols.Volumes {
		names = append(names, v.Name)
	}
	return names, nil
}

func (d *DockerService) DeleteVolume(name string) error {
	return d.client.VolumeRemove(
		context.Background(),
		volumeName(name),
		true,
	)
}
