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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{client: cli}, nil
}

func (d *DockerService) CreateVolume(name string) error {
	_, err := d.client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name: "mc_" + name,
		Labels: map[string]string{
			"project": "minecraft-server",
		},
	})
	return err
}

func (d *DockerService) VolumeExists(name string) (bool, error) {
	vols, err := d.client.VolumeList(context.Background(), volume.ListOptions{Filters: filters.Args{}})
	if err != nil {
		return false, err
	}

	for _, v := range vols.Volumes {
		if v.Name == "mc_"+name {
			return true, nil
		}
	}
	return false, nil
}

func (d *DockerService) ListVolumes() ([]string, error) {
	args := filters.NewArgs()
	args.Add("label", "project=minecraft-server")

	vols, err := d.client.VolumeList(context.Background(), volume.ListOptions{Filters: args})
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, v := range vols.Volumes {
		names = append(names, v.Name)
	}
	return names, nil
}

func (d *DockerService) DeleteVolume(name string) error {
	return d.client.VolumeRemove(context.Background(), "mc_"+name, true)
}

func (d *DockerService) UploadToVolume(volumeName, targetPath, fileName string, data []byte) error {
	ctx := context.Background()

	resp, err := d.client.ContainerCreate(ctx, &container.Config{
		Image: "alpine:latest",
		Cmd:   []string{"sleep", "60"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: targetPath,
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}
	containerID := resp.ID

	defer func() {
		_ = d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	}()

	if err := d.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

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
	tw.Close()

	return d.client.CopyToContainer(ctx, containerID, targetPath, buf, container.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
}
