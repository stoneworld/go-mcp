package server

import (
	"encoding/json"
	"fmt"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

func (server *Server) handleRequestWithPing() (*protocol.PingResult, error) {
	return protocol.NewPingResponse(), nil
}

func (server *Server) handleRequestWithInitialize(sessionID string, rawParams json.RawMessage) (*protocol.InitializeResult, error) {
	var request *protocol.InitializeRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	if request.ProtocolVersion != protocol.Version {
		return nil, fmt.Errorf("protocol version not supported, supported version is %v", protocol.Version)
	}

	s := newSession()
	s.clientInfo = &request.ClientInfo
	s.clientCapabilities = &request.Capabilities
	s.receiveInitRequest.Store(true)

	server.sessionID2session.Store(sessionID, s)

	return &protocol.InitializeResult{
		ServerInfo:      *server.serverInfo,
		Capabilities:    *server.capabilities,
		ProtocolVersion: protocol.Version,
		Instructions:    server.instructions,
	}, nil
}

func (server *Server) handleRequestWithListPrompts(rawParams json.RawMessage) (*protocol.ListPromptsResult, error) {
	if server.capabilities.Prompts == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.ListPromptsRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	// TODO: list prompt with cursor
	return &protocol.ListPromptsResult{
		Prompts: server.prompts,
	}, nil
}

func (server *Server) handleRequestWithGetPrompt(rawParams json.RawMessage) (*protocol.GetPromptResult, error) {
	if server.capabilities.Prompts == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.GetPromptRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	// TODO: validate request's arguments
	handlerFunc, ok := server.promptHandlers[request.Name]
	if !ok {
		return nil, fmt.Errorf("missing prompt handler, promptName=%s", request.Name)
	}
	return handlerFunc(request)
}

func (server *Server) handleRequestWithListResources(rawParams json.RawMessage) (*protocol.ListResourcesResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.ListResourcesRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	// TODO: list resources with cursor
	return &protocol.ListResourcesResult{
		Resources: server.resources,
	}, nil
}

func (server *Server) handleRequestWithReadResource(rawParams json.RawMessage) (*protocol.ReadResourceResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.ReadResourceRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	handlerFunc, ok := server.resourceHandlers[request.URI]
	if !ok {
		return nil, fmt.Errorf("missing resource read handler, uri=%s", request.URI)
	}
	return handlerFunc(request)
}

func (server *Server) handleRequestWitListResourceTemplates(rawParams json.RawMessage) (*protocol.ListResourceTemplatesResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.ListResourceTemplatesRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	// TODO: list resource template with cursor
	return &protocol.ListResourceTemplatesResult{
		ResourceTemplates: server.resourceTemplates,
	}, nil
}

func (server *Server) handleRequestWithSubscribeResourceChange(sessionID string, rawParams json.RawMessage) (*protocol.SubscribeResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.SubscribeRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	value, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return nil, pkg.ErrLackSession
	}
	value.(*session).subscribedResources.Set(request.URI, struct{}{})
	return protocol.NewSubscribeResponse(), nil
}

func (server *Server) handleRequestWithUnSubscribeResourceChange(sessionID string, rawParams json.RawMessage) (*protocol.UnsubscribeResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.UnsubscribeRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	value, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return nil, pkg.ErrLackSession
	}
	value.(*session).subscribedResources.Remove(request.URI)
	return protocol.NewUnsubscribeResponse(), nil
}

func (server *Server) handleRequestWithListTools(rawParams json.RawMessage) (*protocol.ListToolsResult, error) {
	if server.capabilities.Tools == nil {
		return nil, pkg.ErrServerNotSupport
	}

	request := &protocol.ListToolsRequest{}
	if err := pkg.JsonUnmarshal(rawParams, request); err != nil {
		return nil, err
	}
	// TODO: 需要处理request.Cursor的翻页操作
	return &protocol.ListToolsResult{Tools: server.tools}, nil
}

func (server *Server) handleRequestWithCallTool(rawParams json.RawMessage) (*protocol.CallToolResult, error) {
	if server.capabilities.Tools == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.CallToolRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	// TODO: validate request params
	handlerFunc, ok := server.toolHandlers[request.Name]
	if !ok {
		return nil, fmt.Errorf("missing tool call handler, toolName=%s", request.Name)
	}

	return handlerFunc(request)
}

func (server *Server) handleNotifyWithInitialized(sessionID string, rawParams json.RawMessage) error {
	param := &protocol.InitializedNotification{}
	if err := pkg.JsonUnmarshal(rawParams, param); err != nil {
		return err
	}

	val, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return pkg.ErrLackSession
	}
	s := val.(*session)

	if !s.receiveInitRequest.Load() {
		return fmt.Errorf("the server has not received the client's initialization request")
	}
	s.ready.Store(true)
	return nil
}
