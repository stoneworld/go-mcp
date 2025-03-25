package transport

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-mcp/pkg"
)

type SSEClientTransportOption func(*SSEClientTransport)

func WithSSEClientOptionReceiveTimeout(timeout time.Duration) SSEClientTransportOption {
	return func(o *SSEClientTransport) {
		o.receiveTimeout = timeout
	}
}

func WithSSEClientOptionHTTPClient(client *http.Client) SSEClientTransportOption {
	return func(o *SSEClientTransport) {
		o.client = client
	}
}

func WithSSEClientOptionLogger(log pkg.Logger) SSEClientTransportOption {
	return func(o *SSEClientTransport) {
		o.log = log
	}
}

type SSEClientTransport struct {
	ctx    context.Context
	cancel context.CancelFunc

	serverURL string

	endpointChan    chan struct{}
	messageEndpoint *url.URL
	receiver        ClientReceiver

	// options
	log            pkg.Logger
	receiveTimeout time.Duration
	client         *http.Client
}

func NewSSEClientTransport(parent context.Context, serverURL string, opts ...SSEClientTransportOption) (ClientTransport, error) {
	ctx, cancel := context.WithCancel(parent)

	x := &SSEClientTransport{
		ctx:             ctx,
		cancel:          cancel,
		serverURL:       serverURL,
		endpointChan:    make(chan struct{}, 1),
		messageEndpoint: nil,
		receiver:        nil,
		log:             pkg.DefaultLogger,
		receiveTimeout:  time.Second * 30,
		client:          http.DefaultClient,
	}

	for _, opt := range opts {
		opt(x)
	}

	return x, nil
}

func (x *SSEClientTransport) Start() error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	req, err = http.NewRequest(http.MethodGet, x.serverURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	if resp, err = x.client.Do(req); err != nil {

		return fmt.Errorf("failed to connect to SSE stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	go x.readSSE(resp.Body)

	// Wait for the endpoint to be received

	select {
	case <-x.endpointChan:
		// Endpoint received, proceed
	case <-time.After(30 * time.Second): // Add a timeout
		return fmt.Errorf("timeout waiting for endpoint")
	}

	return nil
}

// readSSE continuously reads the SSE stream and processes events.
// It runs until the connection is closed or an error occurs.
func (x *SSEClientTransport) readSSE(reader io.ReadCloser) {
	defer pkg.Recover()

	defer func() {
		_ = reader.Close()
	}()

	br := bufio.NewReader(reader)
	var event, data string

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Process any pending event before exit
				if event != "" && data != "" {
					x.handleSSEEvent(event, data)
				}
				break
			}
			select {
			case <-x.ctx.Done():
				return
			default:
				fmt.Printf("SSE stream error: %v\n", err)
				return
			}
		}

		// Remove only newline markers
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			// Empty line means end of event
			if event != "" && data != "" {
				x.handleSSEEvent(event, data)
				event = ""
				data = ""
			}
			continue
		}

		if strings.HasPrefix(line, "event:") {
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			data = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}
}

// handleSSEEvent processes SSE events based on their type.
// Handles 'endpoint' events for connection setup and 'message' events for JSON-RPC communication.
func (x *SSEClientTransport) handleSSEEvent(event, data string) {
	switch event {
	case "endpoint":
		endpoint, err := url.Parse(data)
		if err != nil {
			fmt.Printf("Error parsing endpoint URL: %v\n", err)
			return
		}
		x.log.Debugf("Received endpoint: %s", endpoint.String())
		x.messageEndpoint = endpoint
		close(x.endpointChan)

	case "message":
		ctx, cancel := context.WithTimeout(x.ctx, x.receiveTimeout)
		defer cancel()
		x.receiver.Receive(ctx, []byte(data))
	}
}

func (x *SSEClientTransport) Send(ctx context.Context, msg Message) error {
	x.log.Debugf("Sending message: %s to %s", msg, x.messageEndpoint.String())

	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, x.messageEndpoint.String(), bytes.NewReader(msg))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if resp, err = x.client.Do(req); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (x *SSEClientTransport) SetReceiver(receiver ClientReceiver) {
	x.receiver = receiver
}

func (x *SSEClientTransport) Close() error {
	x.cancel()
	return nil
}
