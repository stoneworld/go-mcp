package server

import (
	"context"
	"encoding/json"

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
	return nil, nil
}

func (server *Server) handleRequestWithGetPrompt(ctx context.Context, rawParams json.RawMessage) (*protocol.GetPromptResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListResources(ctx context.Context, rawParams json.RawMessage) (*protocol.ListResourcesResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWithReadResource(ctx context.Context, rawParams json.RawMessage) (*protocol.ReadResourceResult, error) {
	return nil, nil
}

func (server *Server) handleRequestWitListResourceTemplates(ctx context.Context, rawParams json.RawMessage) (*protocol.ListResourceTemplatesResult, error) {
	return nil, nil
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
