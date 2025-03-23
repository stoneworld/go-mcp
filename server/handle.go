package server

import (
	"context"
	"encoding/json"
	"fmt"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

func (server *Server) handleRequestWithInitialize(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithPing(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return protocol.NewPingResponse(), nil
}

func (server *Server) handleRequestWithListPrompts(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithGetPrompt(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListResources(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithReadResource(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithReadResourceTemplates(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithSubscribeResourceChange(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithUnSubscribeResourceChange(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListTools(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	request := &protocol.ListToolsRequest{}
	if err := parse(rawParams, request); err != nil {
		return nil, err
	}
	return server.listToolsHandler(ctx, request)
}

func (server *Server) handleRequestWithCallTool(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithCompleteRequest(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithSetLogLevel(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleNotifyWithCancelled(ctx context.Context, rawParams json.RawMessage) error {
	param := &protocol.CancelledNotification{}
	if err := parse(rawParams, param); err != nil {
		return err
	}
	return server.cancelledNotifyHandler(ctx, param)
}

func parse(rawParams json.RawMessage, request interface{}) error {
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return fmt.Errorf("JsonUnmarshal: rawParams=%s, err=%w", rawParams, err)
	}
	return nil
}
