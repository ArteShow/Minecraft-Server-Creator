package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func (ds *DockerService) RemoveContainer(containerID string) error {
	ctx := context.Background()

	return ds.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: true,
	})
}

func (ds *DockerService) StopContainer(containerID string) error {
	ctx := context.Background()

	if err := ds.sendMinecraftStopCommand(ctx, containerID); err == nil {
		if err := ds.waitContainerStopped(ctx, containerID, 30*time.Second); err == nil {
			return nil
		}
	}

	timeout := 30
	if err := ds.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return ds.client.ContainerKill(ctx, containerID, "SIGINT")
	}

	return nil
}

func (ds *DockerService) sendMinecraftStopCommand(ctx context.Context, containerID string) error {
	attached, err := ds.client.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: false,
		Stderr: false,
		Logs:   false,
	})
	if err != nil {
		return err
	}
	defer attached.Close()

	if _, err := io.WriteString(attached.Conn, "stop\n"); err != nil {
		return err
	}

	if c, ok := attached.Conn.(interface{ CloseWrite() error }); ok {
		_ = c.CloseWrite()
	}

	return nil
}

func (ds *DockerService) waitContainerStopped(ctx context.Context, containerID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		inspect, err := ds.client.ContainerInspect(ctx, containerID)
		if err != nil {
			return err
		}
		if !inspect.State.Running {
			return nil
		}

		time.Sleep(300 * time.Millisecond)
	}

	return errors.New("timeout waiting for container to stop gracefully")
}

func (ds *DockerService) StartServerContainer(
	serverID string,
	image string,
	hostPort int,
	containerPort int,
) (string, error) {

	ctx := context.Background()
	port := nat.Port(fmt.Sprintf("%d/tcp", containerPort))

	resp, err := ds.client.ContainerCreate(
		ctx,
		&container.Config{
			Image:      image,
			WorkingDir: "/data",
			Tty:        true,
			OpenStdin:  true,
			AttachStdin: true,
			StdinOnce:  false,

			Cmd: []string{
				"sh",
				"-c",
				`
				while [ ! -f server.jar ]; do
					echo "waiting for server.jar..."
					sleep 1
				done

				echo "starting minecraft server"
				exec java -Xms1G -Xmx2G -jar server.jar nogui
				`,
			},

			ExposedPorts: nat.PortSet{
				port: struct{}{},
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
				port: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprint(hostPort),
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "no",
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
