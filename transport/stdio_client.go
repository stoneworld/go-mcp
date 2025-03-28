package transport

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"go-mcp/pkg"
)

const mcpMessageDelimiter = '\n'

type stdioClientTransport struct {
	cmd      *exec.Cmd
	receiver ClientReceiver
	reader   io.ReadCloser
	writer   io.WriteCloser

	logger pkg.Logger

	cancel          context.CancelFunc
	receiveShutDone chan struct{}
}

func NewStdioClientTransport(command string, env []string, args ...string) (ClientTransport, error) {
	cmd := exec.Command(command, args...)

	cmd.Env = append(os.Environ(), env...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	client := &stdioClientTransport{
		cmd:             cmd,
		reader:          stdout,
		writer:          stdin,
		logger:          pkg.DefaultLogger,
		receiveShutDone: make(chan struct{}),
	}

	return client, nil
}

func (t *stdioClientTransport) Start() error {
	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		defer pkg.Recover()

		t.receive(ctx)
		close(t.receiveShutDone)
	}()

	return nil
}

func (t *stdioClientTransport) Send(ctx context.Context, msg Message) error {
	_, err := t.writer.Write(append(msg, mcpMessageDelimiter))
	return err
}

func (t *stdioClientTransport) SetReceiver(receiver ClientReceiver) {
	t.receiver = receiver
}

func (t *stdioClientTransport) Close() error {
	t.cancel()

	if err := t.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	if err := t.reader.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	if err := t.cmd.Process.Kill(); err != nil {
		return err
	}

	<-t.receiveShutDone

	return nil
}

func (t *stdioClientTransport) receive(ctx context.Context) {
	s := bufio.NewScanner(t.reader)

	for s.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			if err := t.receiver.Receive(ctx, s.Bytes()); err != nil {
				t.logger.Errorf("receiver failed: %v", err)
				return
			}
		}
	}

	if err := s.Err(); err != nil {
		if !errors.Is(err, io.ErrClosedPipe) { // 单测时会保这个错，在这里屏蔽掉
			t.logger.Errorf("client receive unexpected error reading input: %v", err)
		}
		return
	}
}
