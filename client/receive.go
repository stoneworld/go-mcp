package client

import (
	"context"

	"go-mcp/pkg"
	"go-mcp/protocol"

	"github.com/tidwall/gjson"
)

// 对来自客户端的 message(request、response、notification)进行接收处理
// 如果是 request、notification 路由到对应的handler处理，如果是 response 则传递给对应 reqID 的 chan

func (client *Client) Receive(ctx context.Context, msg []byte) {
	if !gjson.GetBytes(msg, "id").Exists() {
		notify := &protocol.JSONRPCNotification{}
		if err := pkg.JsonUnmarshal(msg, &notify); err != nil {
			// 打印日志
			return
		}
		if err := client.receiveNotify(ctx, notify); err != nil {
			// TODO: 打印日志
			return
		}
		return
	}

	// 判断 request和response
	if !gjson.GetBytes(msg, "method").Exists() {
		resp := &protocol.JSONRPCResponse{}
		if err := pkg.JsonUnmarshal(msg, &resp); err != nil {
			// 打印日志
			return
		}
		if err := client.receiveResponse(ctx, resp); err != nil {
			// TODO: 打印日志
			return
		}
		return
	}

	req := &protocol.JSONRPCRequest{}
	if err := pkg.JsonUnmarshal(msg, &req); err != nil {
		// 打印日志
		return
	}
	if err := client.receiveRequest(ctx, req); err != nil {
		// TODO: 打印日志
		return
	}
	return
}

func (client *Client) receiveRequest(ctx context.Context, request *protocol.JSONRPCRequest) *protocol.JSONRPCResponse {
	if !request.IsValid() {
		// return protocol.NewJSONRPCErrorResponse(request.ID,)
	}

	var (
		result protocol.ClientResponse
		err    error
	)

	switch request.Method {
	case protocol.Ping:
		result, err = client.handleRequestWithPing(ctx, request.RawParams)
	case protocol.RootsList:
		result, err = client.handleRequestWithListRoots(ctx, request.RawParams)
	case protocol.SamplingCreateMessage:
		result, err = client.handleRequestWithCreateMessagesSampling(ctx, request.RawParams)
	default:
		// return protocol.NewJSONRPCErrorResponse(request.ID)
	}

	if err != nil {
		// return &protocol.NewJSONRPCErrorResponse(request.ID, ,err.Error())
	}

	return protocol.NewJSONRPCSuccessResponse(request.ID, result)
}

func (client *Client) receiveNotify(ctx context.Context, notify *protocol.JSONRPCNotification) error {
	return client.handleNotify(ctx, notify)
}

func (client *Client) receiveResponse(ctx context.Context, response *protocol.JSONRPCResponse) error {
	// 通过 client.reqID2respChan[response.ID] 将 response 传回发送 request 的协程
	return nil
}
