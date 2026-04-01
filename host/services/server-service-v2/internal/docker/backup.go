package docker

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func (d *DockerService) resolveServerVolume(serverID string) (string, error) {
	ctx := context.Background()
	containerName := "mc_container_" + serverID

	inspect, err := d.client.ContainerInspect(ctx, containerName)
	if err == nil {
		for _, m := range inspect.Mounts {
			if strings.EqualFold(m.Destination, "/data") && m.Type == mount.TypeVolume {
				if m.Name == "" {
					return "", fmt.Errorf("server container %s has /data mount without a named volume", containerName)
				}
				return m.Name, nil
			}
		}
		return "", fmt.Errorf("server container %s has no /data volume mount", containerName)
	}

	// Fallback for cases where container is not currently inspectable but volume still exists.
	fallback := volumeName(serverID)
	exists, existsErr := d.VolumeExists(serverID)
	if existsErr != nil {
		return "", existsErr
	}
	if !exists {
		return "", fmt.Errorf("server volume %s not found; refusing to create backup in a new empty volume", fallback)
	}

	return fallback, nil
}

func (d *DockerService) CreateBackup(serverID string, backupName string) error {
	if serverID == "" || backupName == "" {
		return errors.New("serverID and backupName are required")
	}

	ctx := context.Background()
	serverVolume, err := d.resolveServerVolume(serverID)
	if err != nil {
		return err
	}

	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine:3.19",
			Cmd: []string{
				"sh", "-c",
				"tar --exclude='*.tar.gz' -czf " + path.Join("/data", backupName+".tar.gz") + " -C /data .",
			},
		},
		&container.HostConfig{
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: serverVolume,
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
	defer d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})

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
		if status.StatusCode != 0 {
			return fmt.Errorf("backup container failed with exit code %d", status.StatusCode)
		}
	}

	return nil
}

func (d *DockerService) DownloadBackup(serverID string, backupName string) ([]byte, error) {
	if serverID == "" || backupName == "" {
		return nil, errors.New("serverID and backupName are required")
	}

	ctx := context.Background()
	filePath := path.Join("/data", backupName+".tar.gz")
	serverVolume, err := d.resolveServerVolume(serverID)
	if err != nil {
		return nil, err
	}

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
					Source: serverVolume,
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
	defer d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})

	if err := d.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, err
	}

	reader, _, err := d.client.CopyFromContainer(ctx, containerID, filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := readFileFromTar(reader, backupName+".tar.gz")
	if err != nil {
		return nil, fmt.Errorf("failed to read backup %s: %w", backupName, err)
	}

	return data, nil
}

func (d *DockerService) DeleteBackup(serverID string, backupName string) error {
	if serverID == "" || backupName == "" {
		return errors.New("serverID and backupName are required")
	}

	ctx := context.Background()
	filePath := path.Join("/data", backupName+".tar.gz")
	serverVolume, err := d.resolveServerVolume(serverID)
	if err != nil {
		return err
	}

	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine:3.19",
			Cmd:   []string{"rm", filePath},
		},
		&container.HostConfig{
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: serverVolume,
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
	defer d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})

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
		if status.StatusCode != 0 {
			return fmt.Errorf("delete backup failed with exit code %d", status.StatusCode)
		}
	}

	return nil
}
