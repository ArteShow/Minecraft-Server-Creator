package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

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

func (ds *DockerService) RemoveContainer(containerID string) error {
	ctx := context.Background()

	return ds.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: true,
	})
}

func (ds *DockerService) StopContainer(containerID string) error {
	ctx := context.Background()
	return ds.client.ContainerKill(ctx, containerID, "SIGINT")
}

func (ds *DockerService) StartServerContainer(
	serverID string,
	image string,
	port int,
) (string, error) {

	ctx := context.Background()

	resp, err := ds.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			WorkingDir: "/data",
			Tty:        true,

			Cmd: []string{
				"sh",
				"-c",
				`
				while [ ! -f server.jar ]; do
					echo "waiting for server.jar..."
					sleep 1
				done

				echo "starting minecraft server"
				java -Xms1G -Xmx2G -jar server.jar nogui
				`,
			},

			ExposedPorts: nat.PortSet{
				"25565/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volumeName(serverID),
					Target: "/data",
				},
			},
			PortBindings: nat.PortMap{
				"25565/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprint(port),
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		nil,
		nil,
		"mc_container_"+serverID,
	)
	if err != nil {
		return "", err
	}

	if err := ds.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
