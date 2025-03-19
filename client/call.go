package client

import (
	"context"

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

func (client *Client) CallTool(ctx context.Context) error {
	return nil
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

func (client *Client) SendNotification4Initialized(ctx context.Context) error {
	return nil
}

func (client *Client) SendNotification4Cancelled(ctx context.Context) error {
	return nil
}

func (client *Client) SendNotification4Progress(ctx context.Context) error {
	return nil
}

func (client *Client) SendNotification4RootListChanges(ctx context.Context) error {
	return nil
}

// 负责request和response的拼接
func (client *Client) callServer(ctx context.Context, method protocol.Method, params interface{}) ([]byte, error) {
	requestID := client.requestID.Add(1)
	// 发送请求
	if err := client.sendMsgWithRequest(ctx, requestID, method, params); err != nil {
		return nil, err
	}

	// TODO：
	// 通过chan阻塞等待response
	// 使用ctx进行超时控制
	return nil, nil
}
