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
	prompts := make([]protocol.Prompt, 0)
	server.prompts.Range(func(key string, entry *promptEntry) bool {
		prompts = append(prompts, *entry.prompt)
		return true
	})

	return &protocol.ListPromptsResult{
		Prompts: prompts,
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

	entry, ok := server.prompts.Load(request.Name)
	if !ok {
		return nil, fmt.Errorf("missing prompt, promptName=%s", request.Name)
	}
	return entry.handler(request)
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
	resources := make([]protocol.Resource, 0)
	server.resources.Range(func(key string, entry *resourceEntry) bool {
		resources = append(resources, *entry.resource)
		return true
	})

	return &protocol.ListResourcesResult{
		Resources: resources,
	}, nil
}

func (server *Server) handleRequestWithListResourceTemplates(rawParams json.RawMessage) (*protocol.ListResourceTemplatesResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.ListResourceTemplatesRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	// TODO: list resource template with cursor
	templates := make([]protocol.ResourceTemplate, 0)
	server.resourceTemplates.Range(func(key string, entry *resourceTemplateEntry) bool {
		templates = append(templates, *entry.resourceTemplate)
		return true
	})

	return &protocol.ListResourceTemplatesResult{
		ResourceTemplates: templates,
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

	entry, ok := server.resources.Load(request.URI)
	if !ok {
		return nil, fmt.Errorf("missing resource read handler, uri=%s", request.URI)
	}
	return entry.handler(request)
}

func (server *Server) handleRequestWithSubscribeResourceChange(sessionID string, rawParams json.RawMessage) (*protocol.SubscribeResult, error) {
	if server.capabilities.Resources == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.SubscribeRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	s, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return nil, pkg.ErrLackSession
	}
	s.subscribedResources.Set(request.URI, struct{}{})
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

	s, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return nil, pkg.ErrLackSession
	}
	s.subscribedResources.Remove(request.URI)
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
	tools := make([]*protocol.Tool, 0)
	server.tools.Range(func(key string, entry *toolEntry) bool {
		tools = append(tools, entry.tool)
		return true
	})

	return &protocol.ListToolsResult{Tools: tools}, nil
}

func (server *Server) handleRequestWithCallTool(rawParams json.RawMessage) (*protocol.CallToolResult, error) {
	if server.capabilities.Tools == nil {
		return nil, pkg.ErrServerNotSupport
	}

	var request *protocol.CallToolRequest
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return nil, err
	}

	entry, ok := server.tools.Load(request.Name)
	if !ok {
		return nil, fmt.Errorf("missing tool, toolName=%s", request.Name)
	}

	return entry.handler(request)
}

func (server *Server) handleNotifyWithInitialized(sessionID string, rawParams json.RawMessage) error {
	param := &protocol.InitializedNotification{}
	if err := pkg.JsonUnmarshal(rawParams, param); err != nil {
		return err
	}

	s, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return pkg.ErrLackSession
	}

	if !s.receiveInitRequest.Load() {
		return fmt.Errorf("the server has not received the client's initialization request")
	}
	s.ready.Store(true)
	return nil
}
