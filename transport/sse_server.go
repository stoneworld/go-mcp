package transport

import (
	"context"
	"errors"
	"fmt"
	"go-mcp/pkg"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type SSEServerTransportOption func(*sseServerTransport)

func WithSSEServerTransportOptionLogger(logger pkg.Logger) SSEServerTransportOption {
	return func(o *sseServerTransport) {
		o.logger = logger
	}
}

func WithSSEServerTransportOptionSSEPath(ssePath string) SSEServerTransportOption {
	return func(transport *sseServerTransport) {
		transport.ssePath = ssePath
	}
}

func WithSSEServerTransportOptionMessagePath(messagePath string) SSEServerTransportOption {
	return func(transport *sseServerTransport) {
		transport.messagePath = messagePath
	}
}

func WithSSEServerTransportOptionURLPrefix(urlPrefix string) SSEServerTransportOption {
	return func(transport *sseServerTransport) {
		transport.urlPrefix = urlPrefix
	}
}

type SSEServerTransportAndHandlerOption func(*sseServerTransport)

func WithSSEServerTransportAndHandlerOptionLogger(logger pkg.Logger) SSEServerTransportAndHandlerOption {
	return func(o *sseServerTransport) {
		o.logger = logger
	}
}

type sseServerTransport struct {
	ctx    context.Context
	cancel context.CancelFunc

	httpSvr *http.Server

	messageEndpointFullURL string // 自动生成

	sessionStore pkg.SessionStore

	receiver ServerReceiver

	// options
	logger      pkg.Logger
	ssePath     string
	messagePath string
	urlPrefix   string
}

type SSEHandler struct {
	transport *sseServerTransport
}

// HandleSSE handles incoming SSE connections from clients and sends messages to them.
func (x *SSEHandler) HandleSSE() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x.transport.handleSSE(w, r)
	})
}

// HandleMessage processes incoming JSON-RPC messages from clients and sends responses
// back through both the SSE connection and HTTP response.
func (x *SSEHandler) HandleMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x.transport.handleMessage(w, r)
	})
}

// NewSSEServerTransport 返回会启动 http 的 transport
func NewSSEServerTransport(addr string, opts ...SSEServerTransportOption) (ServerTransport, error) {
	ctx, cancel := context.WithCancel(context.Background())

	x := &sseServerTransport{
		ctx:                    ctx,
		cancel:                 cancel,
		httpSvr:                nil,
		messageEndpointFullURL: "",
		sessionStore:           pkg.DefaultSessionStore,
		receiver:               nil,
		logger:                 pkg.DefaultLogger,
		ssePath:                "/sse",
		messagePath:            "/message",
		urlPrefix:              "http://" + addr,
	}
	for _, opt := range opts {
		opt(x)
	}
	var err error
	x.messageEndpointFullURL, err = url.JoinPath(x.urlPrefix + x.messagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to join URL: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(x.ssePath, x.handleSSE)
	mux.HandleFunc(x.messagePath, x.handleMessage)
	x.httpSvr = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return x, nil
}

// NewSSEServerTransportAndHandler 返回不启动 http 的 transport + Handler 让用户自行启动
func NewSSEServerTransportAndHandler(messageEndpointFullURL string, opts ...SSEServerTransportAndHandlerOption) (ServerTransport, *SSEHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	x := &sseServerTransport{
		ctx:                    ctx,
		cancel:                 cancel,
		httpSvr:                nil,
		messageEndpointFullURL: messageEndpointFullURL,
		sessionStore:           pkg.DefaultSessionStore,
		receiver:               nil,
		logger:                 pkg.DefaultLogger,
		ssePath:                "",
		messagePath:            "",
		urlPrefix:              "",
	}
	for _, opt := range opts {
		opt(x)
	}

	return x, &SSEHandler{
		transport: x,
	}, nil
}

func (x *sseServerTransport) Run() error {
	if x.httpSvr == nil {
		<-x.ctx.Done()
		return fmt.Errorf("httpSvr is nil, please init before Run()")
	}
	err := x.httpSvr.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (x *sseServerTransport) Send(ctx context.Context, sessionID string, msg Message) error {
	conn, ok := x.sessionStore.Load(sessionID)
	if !ok {
		return nil
	}
	c, ok := conn.(chan []byte)
	if !ok {
		return nil
	}
	select {
	case c <- msg:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (x *sseServerTransport) SetReceiver(receiver ServerReceiver) {
	x.receiver = receiver
}

// handleSSE handles incoming SSE connections from clients and sends messages to them.
func (x *sseServerTransport) handleSSE(w http.ResponseWriter, r *http.Request) {
	defer pkg.Recover()

	//nolint:govet // Ignore error since we're just logging
	ctx := r.Context()
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create flush-supporting writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	// Create an SSE connection
	sessionChan := make(chan []byte, 64)
	sessionID := uuid.New().String()
	x.sessionStore.Store(sessionID, sessionChan)
	defer x.sessionStore.Delete(sessionID)
	uri := fmt.Sprintf("%s?sessionID=%s", x.messageEndpointFullURL, sessionID)
	// Send the initial endpoint event
	_, _ = fmt.Fprintf(w, "event: endpoint\ndata: %s\n\n", uri)
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-x.ctx.Done():
			// server closed
			return
		case msg := <-sessionChan:
			_, err := fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
			if err != nil {
				x.logger.Errorf("Failed to write message: %v", err)
				continue
			}
			flusher.Flush()
		}
	}
}

// handleMessage processes incoming JSON-RPC messages from clients and sends responses
// back through both the SSE connection and HTTP response.
func (x *sseServerTransport) handleMessage(w http.ResponseWriter, r *http.Request) {
	defer pkg.RecoverWithFunc(func(r any) {
		x.writeError(w, http.StatusInternalServerError, "Internal server error")
	})

	if r.Method != http.MethodPost {
		x.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		x.writeError(w, http.StatusBadRequest, "Missing session ID")
		return
	}

	_, ok := x.sessionStore.Load(sessionID)
	if !ok {
		x.writeError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	ctx := r.Context()
	// Parse message as raw JSON
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		x.writeError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	x.receiver.Receive(ctx, sessionID, bs)
	// Process message through MCPServer

	// For notifications, just send 202 Accepted with no body
	x.logger.Debugf("Received message: %s", string(bs))
	// ref: https://github.com/encode/httpx/blob/master/httpx/_status_codes.py#L8
	// in official httpx, 2xx is success
	w.WriteHeader(http.StatusAccepted)
}

// writeError writes a JSON-RPC error response with the given error details.
func (x *sseServerTransport) writeError(w http.ResponseWriter, code int, message string) {
	x.logger.Errorf("sseServerTransport writeError: %d %s", code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write([]byte(message)); err != nil {
		x.logger.Errorf("sseServerTransport writeError: %+v", err)
	}
}

func (x *sseServerTransport) Shutdown(ctx context.Context) error {
	x.cancel()

	if x.httpSvr == nil {
		x.logger.Warnf("shutdown sse server without httpSvr")
		return nil
	}

	if err := x.httpSvr.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}
