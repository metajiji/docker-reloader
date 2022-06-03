package docker

import (
	"bytes"
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// ExecResult represents a result returned from Exec()
type ExecResult struct {
	ExitCode  int
	outBuffer *bytes.Buffer
	errBuffer *bytes.Buffer
}

// Stdout returns stdout output of a command run by Exec()
func dockerStdout(res *ExecResult) string {
	data := res.outBuffer.String()
	return trimNewlines(data)
}

// Stderr returns stderr output of a command run by Exec()
func dockerStderr(res *ExecResult) string {
	data := res.errBuffer.String()
	return trimNewlines(data)
}

// Combined returns combined stdout and stderr output of a command run by Exec()
func dockerCombined(res *ExecResult) string {
	data := res.outBuffer.String() + res.errBuffer.String()
	return trimNewlines(data)
}

func trimNewlines(data string) string {
	data = strings.TrimRight(data, "\n") // Remove all newlines, will be added later in log
	return data
}

// Exec executes a command inside a container, returning the result
// containing stdout, stderr, and exit code. Note:
//  - this is a synchronous operation;
//  - cmd stdin is closed.
func dockerExec(ctx context.Context, cli client.APIClient, id string, cmd []string) (ExecResult, error) {
	// prepare exec
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	cresp, err := cli.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return ExecResult{}, err
	}
	execID := cresp.ID

	// run it, with stdout/stderr attached
	aresp, err := cli.ContainerExecAttach(ctx, execID, types.ExecStartCheck{})
	if err != nil {
		return ExecResult{}, err
	}
	defer aresp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, aresp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return ExecResult{}, err
		}
		break

	case <-ctx.Done():
		return ExecResult{}, ctx.Err()
	}

	// get the exit code
	iresp, err := cli.ContainerExecInspect(ctx, execID)
	if err != nil {
		return ExecResult{}, err
	}

	return ExecResult{ExitCode: iresp.ExitCode, outBuffer: &outBuf, errBuffer: &errBuf}, nil
}

func ExecCommand(dClient client.APIClient, container string, command []string) ExecResult {
	ctx := context.Background()

	d, err := dockerExec(ctx, dClient, container, command)
	if err != nil {
		log.Printf(`Command "%s" failed in container "%s" with error: %s`, command, container, err)
		return d
	}

	verdict := `executed successfully`
	if d.ExitCode > 0 {
		verdict = `failed`
	}

	log.Printf(`Command "%s" %s in container "%s" with exit code: %d and output: %s`, command, verdict, container, d.ExitCode, dockerCombined(&d))

	return d
}
