package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"go-mcp/pkg"
)

const stdioSessionID = "stdio"

type stdioServerTransport struct {
	receiver ServerReceiver
	stdin    *bufio.Reader
	stdout   io.Writer

	cancel context.CancelFunc
}

func NewStdioServerTransport() (ServerTransport, error) {
	return &stdioServerTransport{
		stdin:  bufio.NewReader(os.Stdin),
		stdout: os.Stdout,
	}, nil
}

func (t *stdioServerTransport) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		for {
			t.receive(ctx)
		}
	}()

	return nil
}

func (t *stdioServerTransport) Send(ctx context.Context, sessionID string, msg Message) error {
	if _, err := t.stdout.Write(msg); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	return nil
}

func (t *stdioServerTransport) SetReceiver(receiver ServerReceiver) {
	t.receiver = receiver
}

func (t *stdioServerTransport) Close() error {
	// t. cancel()
	return nil
}

func (t *stdioServerTransport) receive(ctx context.Context) {
	defer pkg.Recover()

	line, err := t.stdin.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			// TODO: 记录日志
			return
		}
		return
	}
	t.receiver.Receive(ctx, stdioSessionID, line)
}
