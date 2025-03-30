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

type Option func(*Server)

func WithLogger(logger pkg.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

type Server struct {
	transport transport.ServerTransport

	tools             []*protocol.Tool
	toolHandlers      map[string]ToolHandlerFunc
	prompts           []protocol.Prompt
	promptHandlers    map[string]PromptHandlerFunc
	resources         []protocol.Resource
	resourceHandlers  map[string]ResourceHandlerFunc
	resourceTemplates []protocol.ResourceTemplate

	// TODO：需要定期清理无效session
	sessionID2session *pkg.MemorySessionStore

	inShutdown   atomic.Bool // true when server is in shutdown
	inFlyRequest sync.WaitGroup

	ServerCapabilities *protocol.ServerCapabilities
	ServerInfo         *protocol.Implementation
	ServerInstructions string

	logger pkg.Logger
}

type session struct {
	requestID atomic.Int64

	reqID2respChan cmap.ConcurrentMap[string, chan *protocol.JSONRPCResponse]

	// cache client initialize reqeust info
	clientInfo         *protocol.Implementation
	clientCapabilities *protocol.ClientCapabilities

	// subscribed resources
	subscribedResources cmap.ConcurrentMap[string, struct{}]

	receiveInitRequest atomic.Bool
	ready              atomic.Bool
}

func newSession() *session {
	return &session{
		reqID2respChan:      cmap.New[chan *protocol.JSONRPCResponse](),
		subscribedResources: cmap.New[struct{}](),
	}
}

func NewServer(t transport.ServerTransport, initialize *protocol.InitializeResult, opts ...Option) (*Server, error) {
	server := &Server{
		transport:          t,
		logger:             pkg.DefaultLogger,
		sessionID2session:  pkg.NewMemorySessionStore(),
		ServerCapabilities: &initialize.Capabilities,
		ServerInfo:         &initialize.ServerInfo,
		ServerInstructions: initialize.Instructions,
		toolHandlers:       make(map[string]ToolHandlerFunc),
		promptHandlers:     make(map[string]PromptHandlerFunc),
		resourceHandlers:   make(map[string]ResourceHandlerFunc),
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

type ToolHandlerFunc func(*protocol.CallToolRequest) (*protocol.CallToolResult, error)

func (server *Server) AddTool(tool *protocol.Tool, toolHandler ToolHandlerFunc) {
	server.tools = append(server.tools, tool)
	server.toolHandlers[tool.Name] = toolHandler

}

type PromptHandlerFunc func(*protocol.GetPromptRequest) (*protocol.GetPromptResult, error)

func (server *Server) AddPrompt(prompt protocol.Prompt, promptHandler PromptHandlerFunc) {
	server.prompts = append(server.prompts, prompt)
	server.promptHandlers[prompt.Name] = promptHandler
}

type ResourceHandlerFunc func(*protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error)

func (server *Server) AddResource(resource protocol.Resource, resourceHandler ResourceHandlerFunc) {
	server.resources = append(server.resources, resource)
	server.resourceHandlers[resource.URI] = resourceHandler
}

func (server *Server) AddResourceTemplate(tmpl protocol.ResourceTemplate) {
	server.resourceTemplates = append(server.resourceTemplates, tmpl)
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
