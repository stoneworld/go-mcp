package server

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go-mcp/pkg"
	"go-mcp/protocol"
	"go-mcp/transport"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type Server struct {
	transport transport.ServerTransport

	tools []*protocol.Tool

	cancelledNotifyHandler func(ctx context.Context, notifyParam *protocol.CancelledNotification) error

	// TODO：需要定期清理无效session
	sessionID2session *pkg.MemorySessionStore

	inShutdown   atomic.Bool // true when server is in shutdown
	inFlyRequest sync.WaitGroup

	// The result requirements
	protocolVersion string
	capabilities    protocol.ServerCapabilities
	serverInfo      protocol.Implementation

	logger pkg.Logger
}

type session struct {
	requestID atomic.Int64

	reqID2respChan cmap.ConcurrentMap[string, chan *protocol.JSONRPCResponse]

	// cache client initialize reqeust info
	clientInitializeRequest *protocol.InitializeRequest

	first     bool
	readyChan chan struct{}
}

func NewServer(t transport.ServerTransport, opts ...Option) (*Server, error) {
	server := &Server{
		transport:         t,
		logger:            pkg.DefaultLogger,
		sessionID2session: pkg.NewMemorySessionStore(),
		protocolVersion:   protocol.PROTOCOL_VERSION,
		capabilities: protocol.ServerCapabilities{
			Prompts: &protocol.PromptsCapability{
				ListChanged: true,
			},
			Resources: &protocol.ResourcesCapability{
				Subscribe:   true,
				ListChanged: true,
			},
			Tools: &protocol.ToolsCapability{
				ListChanged: true,
			},
		},
	}
	t.SetReceiver(server)

	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

func (server *Server) Start() error {
	if err := server.transport.Run(); err != nil {
		return fmt.Errorf("init mcp server transpor start fail: %w", err)
	}
	return nil
}

type Option func(*Server)

func WithCancelNotifyHandler(handler func(ctx context.Context, notifyParam *protocol.CancelledNotification) error) Option {
	return func(s *Server) {
		s.cancelledNotifyHandler = handler
	}
}

func WithLogger(logger pkg.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithProtocolVersion(protocolVersion string) Option {
	return func(s *Server) {
		s.protocolVersion = protocolVersion
	}
}

func WithCapabilities(capabilities protocol.ServerCapabilities) Option {
	return func(s *Server) {
		s.capabilities = capabilities
	}
}

func WithInfo(serverInfo protocol.Implementation) Option {
	return func(s *Server) {
		s.serverInfo = serverInfo
	}
}

func (server *Server) AddTool(tool *protocol.Tool) {
	server.tools = append(server.tools, tool)
}

func (server *Server) Shutdown(userCtx context.Context) error {
	server.inShutdown.Store(true)

	serverCtx, cancel := context.WithCancel(userCtx)
	defer cancel()

	go func() {
		defer pkg.Recover()

		server.inFlyRequest.Wait()
		cancel()
	}()

	return server.transport.Shutdown(userCtx, serverCtx)
}
