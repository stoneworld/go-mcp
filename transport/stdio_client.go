package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"go-mcp/pkg"
)

type stdioClientTransport struct {
	cmd      *exec.Cmd
	receiver ClientReceiver

	reader *bufio.Reader
	writer io.WriteCloser

	cancel context.CancelFunc
	done   chan struct{}
	once   sync.Once
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

	client := &stdioClientTransport{
		cmd:    cmd,
		reader: bufio.NewReader(stdout),
		writer: stdin,
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	return client, nil
}

func (t *stdioClientTransport) Start() error {
	t.once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		t.cancel = cancel

		go pkg.SafeRunGo(ctx, func() {
			t.receive(ctx)
			close(t.done)
		})
	})

	return nil
}

func (t *stdioClientTransport) Send(ctx context.Context, msg Message) error {
	msg = append(msg, mcpMessageDelimiter)

	_, err := t.writer.Write(msg)
	return err
}

func (t *stdioClientTransport) SetReceiver(receiver ClientReceiver) {
	t.receiver = receiver
}

func (t *stdioClientTransport) Close(ctx context.Context) error {
	if t.cancel == nil {
		return nil
	}

	t.cancel()

	if err := t.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return t.cmd.Wait()
}

func (t *stdioClientTransport) receive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := t.reader.ReadBytes(mcpMessageDelimiter)
			if err != nil {
				if err != io.EOF {
					// todo: handler error
				}
				return
			}

			t.receiver.Receive(ctx, line)
		}
	}
}
