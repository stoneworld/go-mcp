package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"

	"go-mcp/pkg"
)

type stdioClientTransport struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader

	receiver ClientReceiver
}

func NewStdioClientTransport(command string, args ...string) (ClientTransport, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	return &stdioClientTransport{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}, nil
}

func (t *stdioClientTransport) Start() error {
	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go func() {
		for {
			t.receive()
		}
	}()
	return nil
}

func (t *stdioClientTransport) Send(ctx context.Context, msg Message) error {
	if _, err := t.stdin.Write(msg); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	return nil
}

func (t *stdioClientTransport) SetReceiver(receiver ClientReceiver) {
	t.receiver = receiver
}

func (t *stdioClientTransport) Close() error {
	if err := t.stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}
	return t.cmd.Wait()
}

func (t *stdioClientTransport) receive() {
	defer pkg.Recover()

	line, err := t.stdout.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			// TODO: 使用logger打印
			fmt.Errorf("stdioClientTransport serverReceive Error reading response: %v\n", err)
		}
		return
	}

	t.receiver.Receive(context.Background(), line)

}
