package server

import (
	"fmt"
	"sync/atomic"

	"go-mcp/protocol"
	"go-mcp/transport"
)

type Server struct {
	transport transport.Transport

	reqID2respChan map[int]chan *protocol.JSONRPCResponse

	notifyMethod2handler map[protocol.Method]func(notifyParam interface{})

	requestID atomic.Int64
}

func NewServer(t transport.Transport, opts ...Option) (*Server, error) {
	server := &Server{
		transport: t,
	}
	t.SetReceiver(server)

	for _, opt := range opts {
		opt(server)
	}

	if err := server.transport.Start(); err != nil {
		return nil, fmt.Errorf("init mcp server transpor start fail: %w", err)
	}

	return server, nil
}

type Option func(*Server)

func WithNotifyHandler(notifyMethod2handler map[protocol.Method]func(notifyParam interface{})) Option {
	return func(s *Server) {
		s.notifyMethod2handler = notifyMethod2handler
	}
}

func (server *Server) Close() error {
	// TODO 还有一些其他处理操作也可以放在这里
	// server.transport.Close()
	return nil
}
