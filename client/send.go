package client

import (
	"context"
	"fmt"

	"go-mcp/protocol"

	"github.com/bytedance/sonic"
)

func (client *Client) sendMsgWithRequest(ctx context.Context, requestID protocol.RequestID, method protocol.Method, params interface{}) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	req := protocol.NewJSONRPCRequest(requestID, method, params)

	message, err := sonic.Marshal(req)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendRequest: transport send: %w", err)
	}
	return nil
}

func (client *Client) sendMsgWithResponse(ctx context.Context, requestID protocol.RequestID, result interface{}) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	resp := protocol.NewJSONRPCSuccessResponse(requestID, result)

	message, err := sonic.Marshal(resp)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}

func (client *Client) sendMsgWithNotification(ctx context.Context, method protocol.Method, params interface{}) error {
	notify := protocol.NewJSONRPCNotification(method, params)

	message, err := sonic.Marshal(notify)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendNotification: transport send: %w", err)
	}
	return nil
}

func (client *Client) sendMsgWithError(ctx context.Context, requestID protocol.RequestID, code int, msg string) error {
	if requestID == nil {
		return fmt.Errorf("requestID can't is nil")
	}

	resp := protocol.NewJSONRPCErrorResponse(requestID, code, msg)

	message, err := sonic.Marshal(resp)
	if err != nil {
		return err
	}

	if err := client.transport.Send(ctx, message); err != nil {
		return fmt.Errorf("sendResponse: transport send: %w", err)
	}
	return nil
}
