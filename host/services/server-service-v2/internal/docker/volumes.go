package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"

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

func (d *DockerService) GetFileFromVolume(
	serverID string,
	targetPath string,
	fileName string,
) ([]byte, error) {
	if serverID == "" {
		return nil, errors.New("serverID is required")
	}
	if targetPath == "" {
		return nil, errors.New("targetPath is required")
	}
	if fileName == "" {
		return nil, errors.New("fileName is required")
	}

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
					Source: volumeName(serverID),
					Target: "/data",
				},
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return nil, err
	}

	containerID := resp.ID
	defer func() {
		_ = d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	}()

	if err := d.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, err
	}

	filePath := path.Join(targetPath, fileName)
	reader, _, err := d.client.CopyFromContainer(ctx, containerID, filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := readFileFromTar(reader, fileName)
	if err != nil {
		return nil, fmt.Errorf("read %s from volume %s: %w", filePath, volumeName(serverID), err)
	}

	return data, nil
}

func (d *DockerService) DeleteFileFromVolume(
	serverID string,
	targetPath string,
	fileName string,
) error {
	if serverID == "" {
		return errors.New("serverID is required")
	}
	if targetPath == "" {
		return errors.New("targetPath is required")
	}
	if fileName == "" {
		return errors.New("fileName is required")
	}

	ctx := context.Background()
	filePath := path.Join(targetPath, fileName)

	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine:3.19",
			Cmd:   []string{"rm", "-f", filePath},
		},
		&container.HostConfig{
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volumeName(serverID),
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
		_ = d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	}()

	if err := d.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := d.client.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case waitErr := <-errCh:
		if waitErr != nil {
			return waitErr
		}
	case status := <-statusCh:
		if status.Error != nil && status.Error.Message != "" {
			return fmt.Errorf("delete %s from volume %s: %s", filePath, volumeName(serverID), status.Error.Message)
		}
		if status.StatusCode != 0 {
			return fmt.Errorf("delete %s from volume %s failed with exit code %d", filePath, volumeName(serverID), status.StatusCode)
		}
	}

	return nil
}

func readFileFromTar(reader io.Reader, fileName string) ([]byte, error) {
	tr := tar.NewReader(reader)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("file %s not found in tar stream", fileName)
			}
			return nil, err
		}

		if hdr.FileInfo().IsDir() {
			continue
		}

		_, currentFile := path.Split(hdr.Name)
		if currentFile != fileName {
			continue
		}

		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
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
