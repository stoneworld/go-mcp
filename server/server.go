package server

import (
	"context"
	"fmt"
	"sync/atomic"

	"go-mcp/pkg"
	"go-mcp/protocol"
	"go-mcp/transport"
)

type Server struct {
	transport transport.ServerTransport

	tools []protocol.Tool

	cancelledNotifyHandler func(ctx context.Context, notifyParam *protocol.CancelledNotification) error

	// TODO：需要定期清理无效session
	sessionID2session map[string]*session

	requestID atomic.Int64

	logger pkg.Logger
}

type session struct {
	reqID2respChan map[int]chan *protocol.JSONRPCResponse
	first          bool
	readyChan      chan struct{}
}

func NewServer(t transport.ServerTransport, opts ...Option) (*Server, error) {
	server := &Server{
		transport: t,
		logger:    &pkg.Log{},
	}
	t.SetReceiver(server)

	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

func (server *Server) Start() error {
	if err := server.transport.Start(); err != nil {
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

func (server *Server) Shutdown() error {
	// TODO 还有一些其他处理操作也可以放在这里
	// server.transport.Close()
	return nil
}
