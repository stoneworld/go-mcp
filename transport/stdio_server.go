package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

type stdioServerTransport struct {
	receiver receiver
	stdin    io.Reader
	stdout   io.Writer

	cancel context.CancelFunc
}

func NewStdioServerTransport() (Transport, error) {
	return &stdioServerTransport{
		stdin:  os.Stdin,
		stdout: os.Stdout,
	}, nil
}

func (t *stdioServerTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go t.receive(ctx)

	return nil
}

func (t *stdioServerTransport) Send(ctx context.Context, msg Message) error {
	if _, err := t.stdout.Write(msg); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	return nil
}

func (t *stdioServerTransport) SetReceiver(receiver receiver) {
	t.receiver = receiver
}

func (t *stdioServerTransport) Close() error {
	// t. cancel()
	return nil
}

func (t *stdioServerTransport) receive(ctx context.Context) {
	reader := bufio.NewReader(t.stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				// TODO: 记录日志
				return
			}
			return
		}
		t.receiver.Receive(ctx, line)
	}
}
