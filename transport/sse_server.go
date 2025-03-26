package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go-mcp/pkg"

	"github.com/google/uuid"
)

type SSEServerTransportOption func(*sseServerTransport)

func WithSSEServerTransportOptionLogger(logger pkg.Logger) SSEServerTransportOption {
	return func(t *sseServerTransport) {
		t.logger = logger
	}
}

func WithSSEServerTransportOptionSSEPath(ssePath string) SSEServerTransportOption {
	return func(t *sseServerTransport) {
		t.ssePath = ssePath
	}
}

func WithSSEServerTransportOptionMessagePath(messagePath string) SSEServerTransportOption {
	return func(t *sseServerTransport) {
		t.messagePath = messagePath
	}
}

func WithSSEServerTransportOptionURLPrefix(urlPrefix string) SSEServerTransportOption {
	return func(t *sseServerTransport) {
		t.urlPrefix = urlPrefix
	}
}

type SSEServerTransportAndHandlerOption func(*sseServerTransport)

func WithSSEServerTransportAndHandlerOptionLogger(logger pkg.Logger) SSEServerTransportAndHandlerOption {
	return func(t *sseServerTransport) {
		t.logger = logger
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
func (h *SSEHandler) HandleSSE() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.transport.handleSSE(w, r)
	})
}

// HandleMessage processes incoming JSON-RPC messages from clients and sends responses
// back through both the SSE connection and HTTP response.
func (h *SSEHandler) HandleMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.transport.handleMessage(w, r)
	})
}

// NewSSEServerTransport 返回会启动 http 的 transport
func NewSSEServerTransport(addr string, opts ...SSEServerTransportOption) (ServerTransport, error) {
	ctx, cancel := context.WithCancel(context.Background())

	t := &sseServerTransport{
		ctx:          ctx,
		cancel:       cancel,
		sessionStore: pkg.NewMemorySessionStore(),
		logger:       pkg.DefaultLogger,
		ssePath:      "/sse",
		messagePath:  "/message",
		urlPrefix:    "http://" + addr,
	}
	for _, opt := range opts {
		opt(t)
	}
	var err error
	t.messageEndpointFullURL, err = url.JoinPath(t.urlPrefix + t.messagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to join URL: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(t.ssePath, t.handleSSE)
	mux.HandleFunc(t.messagePath, t.handleMessage)

	t.httpSvr = &http.Server{
		Addr:        addr,
		Handler:     mux,
		IdleTimeout: time.Minute,
	}

	return t, nil
}

// NewSSEServerTransportAndHandler 返回不启动 http 的 transport + Handler 让用户自行启动
func NewSSEServerTransportAndHandler(messageEndpointFullURL string, opts ...SSEServerTransportAndHandlerOption) (ServerTransport, *SSEHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	t := &sseServerTransport{
		ctx:                    ctx,
		cancel:                 cancel,
		messageEndpointFullURL: messageEndpointFullURL,
		sessionStore:           pkg.NewMemorySessionStore(),
		logger:                 pkg.DefaultLogger,
	}
	for _, opt := range opts {
		opt(t)
	}

	return t, &SSEHandler{transport: t}, nil
}

func (t *sseServerTransport) Run() error {
	if t.httpSvr == nil {
		<-t.ctx.Done()
		return nil
	}

	if err := t.httpSvr.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (t *sseServerTransport) Send(ctx context.Context, sessionID string, msg Message) error {
	conn, ok := t.sessionStore.Load(sessionID)
	if !ok {
		return pkg.NewLackSessionError(sessionID)
	}
	c, ok := conn.(chan []byte)
	if !ok {
		return fmt.Errorf("sessionID=%s invalid connection type: %T", sessionID, conn)
	}

	select {
	case c <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-t.ctx.Done():
		return fmt.Errorf("transport is shutting down")
	}
}

func (t *sseServerTransport) SetReceiver(receiver ServerReceiver) {
	t.receiver = receiver
}

// handleSSE handles incoming SSE connections from clients and sends messages to them.
func (t *sseServerTransport) handleSSE(w http.ResponseWriter, r *http.Request) {
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
	t.sessionStore.Store(sessionID, sessionChan)
	defer t.sessionStore.Delete(sessionID)

	uri := fmt.Sprintf("%s?sessionID=%s", t.messageEndpointFullURL, sessionID)
	// Send the initial endpoint event
	_, _ = fmt.Fprintf(w, "event: endpoint\ndata: %s\n\n", uri)
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.ctx.Done():
			// server closed
			return
		case msg := <-sessionChan:
			_, err := fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
			if err != nil {
				t.logger.Errorf("Failed to write message: %v", err)
				continue
			}
			flusher.Flush()
		}
	}
}

// handleMessage processes incoming JSON-RPC messages from clients and sends responses
// back through both the SSE connection and HTTP response.
func (t *sseServerTransport) handleMessage(w http.ResponseWriter, r *http.Request) {
	defer pkg.RecoverWithFunc(func(r any) {
		t.writeError(w, http.StatusInternalServerError, "Internal server error")
	})

	if r.Method != http.MethodPost {
		t.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		t.writeError(w, http.StatusBadRequest, "Missing session ID")
		return
	}

	_, ok := t.sessionStore.Load(sessionID)
	if !ok {
		t.writeError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	ctx := r.Context()
	// Parse message as raw JSON
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		t.writeError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	t.receiver.Receive(ctx, sessionID, bs)
	// Process message through MCPServer

	// For notifications, just send 202 Accepted with no body
	t.logger.Debugf("Received message: %s", string(bs))
	// ref: https://github.com/encode/httpx/blob/master/httpx/_status_codes.py#L8
	// in official httpx, 2xx is success
	w.WriteHeader(http.StatusAccepted)
}

// writeError writes a JSON-RPC error response with the given error details.
func (t *sseServerTransport) writeError(w http.ResponseWriter, code int, message string) {
	t.logger.Errorf("sseServerTransport writeError: code: %d, message: %s", code, message)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	if _, err := w.Write([]byte(message)); err != nil {
		t.logger.Errorf("sseServerTransport writeError: %+v", err)
	}
}

func (t *sseServerTransport) Shutdown(ctx context.Context) error {
	if t.httpSvr == nil {
		t.logger.Warnf("shutdown sse server without httpSvr")
		return nil
	}

	t.httpSvr.RegisterOnShutdown(func() {
		<-ctx.Done()
		t.cancel()
	})

	if err := t.httpSvr.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}
