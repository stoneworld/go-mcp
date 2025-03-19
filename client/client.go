package client

import (
	"context"
	"fmt"
	"sync/atomic"

	"go-mcp/protocol"
	"go-mcp/transport"
)

type Client struct {
	transport transport.Transport

	requestID atomic.Int64
}

func Init(t transport.Transport) (*Client, error) {
	client := &Client{
		transport: t,
	}
	if err := client.transport.Start(); err != nil {
		return nil, fmt.Errorf("init mcp client transpor start fail: %w", err)
	}

	return client, nil
}

// 请求

// 1. 请求构造
// 2. 发送请求 client.callServer(ctx)
// 3. 响应解析

func (client *Client) Ping(ctx context.Context) error {
	// client.call(ctx)
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

func (client *Client) ListTools(ctx context.Context) error {
	return nil
}

func (client *Client) CallTools(ctx context.Context) error {
	return nil
}

func (client *Client) RequestCompletions(ctx context.Context) error {
	return nil
}

func (client *Client) SetLogLevel(ctx context.Context) error {
	return nil
}

// 通知
// 1. 构造通知结构体
// 2. 发送通知 client.sendMsgWithNotification(ctx)

func (client *Client) SendNotification4Cancelled(ctx context.Context) error {
	return nil
}

func (client *Client) SendNotification4Progress(ctx context.Context) error {
	return nil
}

func (client *Client) SendNotification4RootListChanges(ctx context.Context) error {
	return nil
}

func (client *Client) Close() error {
	// TODO 还有一些其他处理操作也可以放在这里
	// client.transport.Close()
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
