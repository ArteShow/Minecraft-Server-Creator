package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
)

func (ds *DockerService) ExecInContainer(
	containerID string,
	cmd []string,
) error {

	ctx := context.Background()

	execResp, err := ds.client.ContainerExecCreate(
		ctx,
		containerID,
		container.ExecOptions{
			Cmd:          cmd,
			AttachStdout: true,
			AttachStderr: true,
		},
	)
	if err != nil {
		return err
	}

	attach, err := ds.client.ContainerExecAttach(
		ctx,
		execResp.ID,
		container.ExecAttachOptions{},
	)
	if err != nil {
		return err
	}
	defer attach.Close()

	_, err = io.Copy(io.Discard, attach.Reader)
	return err
}