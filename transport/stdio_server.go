package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

type stdioServerTransport struct {
	receiver Receiver
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
	// TODO implement me
	panic("implement me")
}

func (t *stdioServerTransport) SetReceiver(receiver Receiver) {
	t.receiver = receiver
}

func (t *stdioServerTransport) Close() error {
	// t. cancel()
	return nil
}

func (t *stdioServerTransport) receive(ctx context.Context) {
	reader := bufio.NewReader(t.stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		readChan <- line

		line, err := t.stdout.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("Error reading response: %v\n", err)
			}
			return nil, nil
		}
		if err := t.receiver.Receive(ctx, line); err != nil {
			logs.Errorf("stdioClientTransport receive line=%s err=%+v", line, err)
		}
	}
}
