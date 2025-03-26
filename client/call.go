package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

// 请求

// 1. 请求构造
// 2. 发送请求 client.callServer(ctx)
// 3. 响应解析

func (client *Client) initialization(ctx context.Context) error {
	// client.callServer(ctx)
	// client.SendNotification4Initialized(ctx)
	return nil
}

func (client *Client) Ping(ctx context.Context) error {
	// client.callServer(ctx)
	return nil
}

func (client *Client) ListPrompts(ctx context.Context) error {
	return nil
}

func (client *Client) GetPrompt(ctx context.Context) error {
	return nil
}

func (client *Client) ListResources(ctx context.Context) error {
	return nil
}

func (client *Client) ReadResource(ctx context.Context) error {
	return nil
}

func (client *Client) ListResourceTemplates(ctx context.Context) error {
	return nil
}

func (client *Client) SubscribeResourceChange(ctx context.Context) error {
	return nil
}

func (client *Client) UnSubscribeResourceChange(ctx context.Context) error {
	return nil
}

func (client *Client) ListTools(ctx context.Context) error {
	return nil
}

func (client *Client) CallTool(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	result := &protocol.CallToolResult{}

	if err := client.callAndParse(ctx, protocol.ToolsCall, &request, &result); err != nil {
		return nil, fmt.Errorf("CallTool: %w", err)
	}
	return result, nil
}

func (client *Client) CompleteRequest(ctx context.Context) error {
	return nil
}

func (client *Client) SetLogLevel(ctx context.Context) error {
	return nil
}

// 通知
// 1. 构造通知结构体
// 2. 发送通知 client.sendMsgWithNotification(ctx)

func (client *Client) sendNotification4Initialized(ctx context.Context, notify *protocol.InitializedNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationInitialized, notify)
}

func (client *Client) SendNotification4Cancelled(ctx context.Context, notify *protocol.CancelledNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationCancelled, notify)
}

func (client *Client) SendNotification4Progress(ctx context.Context, notify *protocol.ProgressNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationProgress, notify)
}

func (client *Client) SendNotification4RootListChanges(ctx context.Context, notify *protocol.RootsListChangedNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationRootsListChanged, notify)
}

func (client *Client) callAndParse(ctx context.Context, method protocol.Method, request protocol.ClientRequest, result protocol.ServerResponse) error {
	rawResult, err := client.callServer(ctx, method, request)
	if err != nil {
		return err
	}

	if err := pkg.JsonUnmarshal(rawResult, &result); err != nil {
		return fmt.Errorf("JsonUnmarshal: rawResult=%s, err=%w", rawResult, err)
	}
	return nil
}

// 负责request和response的拼接
func (client *Client) callServer(ctx context.Context, method protocol.Method, params protocol.ClientRequest) (json.RawMessage, error) {
	requestID := strconv.FormatInt(client.requestID.Add(1), 10)
	// 发送请求
	if err := client.sendMsgWithRequest(ctx, requestID, method, params); err != nil {
		return nil, fmt.Errorf("callServer: %w", err)
	}

	respChan := make(chan *protocol.JSONRPCResponse)

	client.reqID2respChan.Set(requestID, respChan)

	select {
	case <-ctx.Done():
		client.reqID2respChan.Remove(requestID)
		return nil, ctx.Err()
	case response := <-respChan:
		if err := response.Error; err != nil {
			return nil, pkg.NewResponseError(err.Code, err.Message, err.Data)
		}
		return response.RawResult, nil
	}
}
