package transport

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
)

type sseClientTransport struct {
}

func NewSSEClientTransport(command string, args ...string) (Transport, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	return &sseClientTransport{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}, nil
}

func (c *sseClientTransport) Start() error {
	// TODO implement me
	panic("implement me")
}

func (c *sseClientTransport) Send(ctx context.Context, msg Message) error {
	// TODO implement me
	panic("implement me")
}

func (c *sseClientTransport) Receive() (Message, error) {
	// TODO implement me
	panic("implement me")
}

func (c *sseClientTransport) Close() error {
	// TODO implement me
	panic("implement me")
}
