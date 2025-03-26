package transport

import (
	"bufio"
	"context"
	"io"
	"sync"

	"go-mcp/pkg"
)

type stdioClientTransport struct {
	receiver ClientReceiver
	reader   *bufio.Reader
	writer   io.Writer

	cancel context.CancelFunc
	done   chan struct{}
	once   sync.Once
}

func NewStdioClientTransport(in io.Reader, out io.Writer) ClientTransport {
	return &stdioClientTransport{
		reader: bufio.NewReader(in),
		writer: out,

		done: make(chan struct{}),
	}
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
