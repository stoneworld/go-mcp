package server

import (
	"context"
	"fmt"

	"go-mcp/protocol"

	"github.com/bytedance/sonic"
)

func (server *Server) sendMsgWithRequest(ctx context.Context, sessionID string, requestID protocol.RequestID, method protocol.Method, params protocol.ServerRequest) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	req := protocol.NewJSONRPCRequest(requestID, method, params)

	message, err := sonic.Marshal(req)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendRequest: transport send: %w", err)
	}
	return nil
}

func (server *Server) sendMsgWithResponse(ctx context.Context, sessionID string, response *protocol.JSONRPCResponse) error {
	message, err := sonic.Marshal(response)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}

func (server *Server) sendMsgWithNotification(ctx context.Context, sessionID string, method protocol.Method, params protocol.ServerNotify) error {
	notify := protocol.NewJSONRPCNotification(method, params)

	message, err := sonic.Marshal(notify)
	if err != nil {
		return err
	}

	if err := server.transport.Send(ctx, sessionID, message); err != nil {
		return fmt.Errorf("sendNotification: transport send: %w", err)
	}
	return nil
}
