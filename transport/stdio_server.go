package transport

import (
	"bufio"
	"context"
	"io"
	"os"
	"sync"

	"go-mcp/pkg"
)

const stdioSessionID = "stdio"

type stdioServerTransport struct {
	receiver ServerReceiver
	reader   *bufio.Reader
	writer   io.Writer

	cancel context.CancelFunc
	done   chan struct{}
	once   sync.Once
}

func NewStdioServerTransport() ServerTransport {
	return &stdioServerTransport{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,

		done: make(chan struct{}),
	}
}

func (t *stdioServerTransport) Run() error {
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

func (t *stdioServerTransport) Send(ctx context.Context, sessionID string, msg Message) error {
	msg = append(msg, mcpMessageDelimiter)

	_, err := t.writer.Write(msg)
	return err
}

func (t *stdioServerTransport) SetReceiver(receiver ServerReceiver) {
	t.receiver = receiver
}

func (t *stdioServerTransport) Shutdown(ctx context.Context) error {
	if t.cancel == nil {
		return nil
	}

	t.cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), defaultStdioTransportCloseTimeout)
	defer cancelTimeout()

	select {
	case <-t.done:
		return nil
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *stdioServerTransport) receive(ctx context.Context) {
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

			t.receiver.Receive(ctx, stdioSessionID, line)
		}
	}
}
