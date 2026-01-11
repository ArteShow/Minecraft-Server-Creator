package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"path"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func (d *DockerService) UploadToVolume(
	name string,
	targetPath string,
	fileName string,
	data []byte,
) error {
	ctx := context.Background()

	targetPath = strings.TrimPrefix(targetPath, "/")
	fullPath := path.Join(targetPath, fileName)

	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine",
			Cmd:   []string{"sleep", "20"},
		},
		&container.HostConfig{
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
		Name: fullPath,
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
		"/data",
		buf,
		container.CopyToContainerOptions{},
	)
}
