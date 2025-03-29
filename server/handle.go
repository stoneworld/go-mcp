package server

import (
	"context"
	"encoding/json"
	"fmt"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

func (server *Server) handleRequestWithInitialize(ctx context.Context, rawParams json.RawMessage) (*protocol.InitializeResult, error) {
	var req protocol.InitializeRequest
	if err := pkg.JsonUnmarshal(rawParams, &req); err != nil {
		return nil, err
	}

	// TODO: validate client initialize request

	// cache client information to session
	sessionID, _ := getSessionIDFromCtx(ctx)
	value, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return nil, pkg.NewLackSessionError(sessionID)
	}
	session := value.(*session)
	session.clientInitializeRequest = &req

	result := protocol.InitializeResult{
		ProtocolVersion: server.protocolVersion,
		Capabilities:    server.capabilities,
		ServerInfo:      server.serverInfo,
	}
	return &result, nil
}

func (server *Server) handleRequestWithPing(ctx context.Context, rawParams json.RawMessage) (*protocol.PingResult, error) {
	return protocol.NewPingResponse(), nil
}

func (server *Server) handleRequestWithListPrompts(ctx context.Context, rawParams json.RawMessage) (*protocol.ListPromptsResult, error) {
	var req protocol.ListPromptsRequest
	if err := pkg.JsonUnmarshal(rawParams, &req); err != nil {
		return nil, err
	}

	// TODO: list prompt with cursor
	return &protocol.ListPromptsResult{
		Prompts: server.prompts,
	}, nil
}

func (server *Server) handleRequestWithGetPrompt(ctx context.Context, rawParams json.RawMessage) (*protocol.GetPromptResult, error) {
	var req protocol.GetPromptRequest
	if err := pkg.JsonUnmarshal(rawParams, &req); err != nil {
		return nil, err
	}

	// TODO: validate request's arguments
	handleFunc, ok := server.promptHandlers[req.Name]
	if !ok {
		return nil, fmt.Errorf("missing prompt handler, promptName=%s", req.Name)
	}
	return handleFunc(req), nil
}

func (server *Server) handleRequestWithListResources(ctx context.Context, rawParams json.RawMessage) (*protocol.ListResourcesResult, error) {
	var req protocol.ListResourcesRequest
	if err := pkg.JsonUnmarshal(rawParams, &req); err != nil {
		return nil, err
	}

	// TODO: list resources with cursor
	return &protocol.ListResourcesResult{
		Resources: server.resources,
	}, nil
}

func (server *Server) handleRequestWithReadResource(ctx context.Context, rawParams json.RawMessage) (*protocol.ReadResourceResult, error) {
	var req protocol.ReadResourceRequest
	if err := pkg.JsonUnmarshal(rawParams, &req); err != nil {
		return nil, err
	}

	handleFunc, ok := server.resourceHandlers[req.URI]
	if !ok {
		return nil, fmt.Errorf("missing resource read handler, uri=%s", req.URI)
	}
	return handleFunc(req), nil
}

func (server *Server) handleRequestWitListResourceTemplates(ctx context.Context, rawParams json.RawMessage) (*protocol.ListResourceTemplatesResult, error) {
	var req protocol.ListResourceTemplatesRequest
	if err := pkg.JsonUnmarshal(rawParams, &req); err != nil {
		return nil, err
	}

	// TODO: list resource template with cursor
	return &protocol.ListResourceTemplatesResult{
		ResourceTemplates: server.resourceTemplates,
	}, nil
}

func (server *Server) handleRequestWithSubscribeResourceChange(ctx context.Context, rawParams json.RawMessage) (*protocol.SubscribeResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWithUnSubscribeResourceChange(ctx context.Context, rawParams json.RawMessage) (*protocol.UnsubscribeResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListTools(ctx context.Context, rawParams json.RawMessage) (*protocol.ListToolsResult, error) {
	request := &protocol.ListToolsRequest{}
	if err := pkg.JsonUnmarshal(rawParams, request); err != nil {
		return nil, err
	}
	// TODO: 需要处理request.Cursor的翻页操作
	return &protocol.ListToolsResult{Tools: server.tools}, nil
}

func (server *Server) handleRequestWithCallTool(ctx context.Context, rawParams json.RawMessage) (*protocol.CallToolResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWithCompleteRequest(ctx context.Context, rawParams json.RawMessage) (*protocol.CompleteResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWithSetLogLevel(ctx context.Context, rawParams json.RawMessage) (*protocol.SetLoggingLevelResult, error) {
	return nil, nil
}

func (server *Server) handleNotifyWithCancelled(ctx context.Context, rawParams json.RawMessage) error {
	param := &protocol.CancelledNotification{}
	if err := pkg.JsonUnmarshal(rawParams, param); err != nil {
		return err
	}
	return server.cancelledNotifyHandler(ctx, param)
}
