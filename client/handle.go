package client

import (
	"context"
	"errors"

	"go-mcp/protocol"
)

func (client *Client) handleRequestWithInitialize(ctx context.Context, request *protocol.JSONRPCRequest) (protocol.Result, error) {
	return nil, nil
}

func (client *Client) handleRequestWithPing(ctx context.Context, request *protocol.JSONRPCRequest) (protocol.Result, error) {
	return nil, nil
}

func (client *Client) handleRequestWithListRoots(ctx context.Context, request *protocol.JSONRPCRequest) (protocol.Result, error) {
	return nil, nil
}

func (client *Client) handleRequestWithCreateMessagesSampling(ctx context.Context, request *protocol.JSONRPCRequest) (protocol.Result, error) {
	return nil, nil
}

func (client *Client) handleNotify(ctx context.Context, notify *protocol.JSONRPCNotification) error {
	if notify.Method == "" {
		return errors.New("notify method can't is \"\"")
	}

	// TODO: 使用client里定义一个 notifyMethod2handler 对通知进行处理
	handler, ok := client.notifyMethod2handler[notify.Method]
	if !ok {
		// 打印 warn/info 日志
		// 此处也可以向上抛error，在上层识别error统一打日志
		return nil
	}
	handler(notify.Params)
	return nil
}
