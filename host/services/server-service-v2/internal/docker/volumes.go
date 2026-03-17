package docker

import (
	"archive/tar"
	"bytes"
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
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

func (d *DockerService) UploadToVolume(
	name string,
	targetPath string,
	fileName string,
	data []byte,
) error {
	ctx := context.Background()

	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine:3.19",
			Cmd:   []string{"sleep", "20"},
		},
		&container.HostConfig{
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volumeName(name),
					Target: "/data",
				},
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return err
	}

	containerID := resp.ID
	defer func() {
		_ = d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
			Force: true,
		})
	}()

	if err := d.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	hdr := &tar.Header{
		Name: fileName,
		Mode: 0644,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := tw.Write(data); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}

	return d.client.CopyToContainer(
		ctx,
		containerID,
		targetPath,
		buf,
		container.CopyToContainerOptions{},
	)
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
