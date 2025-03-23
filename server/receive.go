package server

import (
	"context"

	"go-mcp/pkg"
	"go-mcp/protocol"

	"github.com/tidwall/gjson"
)

// 对来自客户端的 message(request、response、notification)进行接收处理
// 如果是 request、notification 路由到对应的handler处理，如果是 response 则传递给对应 reqID 的 chan

func (server *Server) Receive(ctx context.Context, sessionID string, msg []byte) {
	ctx = setSessionIDToCtx(ctx, sessionID)

	if !gjson.GetBytes(msg, "id").Exists() {
		notify := &protocol.JSONRPCNotification{}
		if err := pkg.JsonUnmarshal(msg, &notify); err != nil {
			// 打印日志
			return
		}
		if err := server.receiveNotify(ctx, sessionID, notify); err != nil {
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
		if err := server.receiveResponse(ctx, sessionID, resp); err != nil {
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
	if err := server.receiveRequest(ctx, sessionID, req); err != nil {
		// TODO: 打印日志
		return
	}
	return
}

func (server *Server) receiveRequest(ctx context.Context, sessionID string, request *protocol.JSONRPCRequest) *protocol.JSONRPCResponse {
	if !request.IsValid() {
		// return protocol.NewJSONRPCErrorResponse(request.ID,)
	}

	var (
		result protocol.ServerResponse
		err    error
	)

	switch request.Method {
	case protocol.Ping:
		result, err = server.handleRequestWithPing(ctx, request.RawParams)
	case protocol.Initialize:
		result, err = server.handleRequestWithInitialize(ctx, request.RawParams)
	case protocol.PromptsList:
		result, err = server.handleRequestWithListPrompts(ctx, request.RawParams)
	case protocol.PromptsGet:
		result, err = server.handleRequestWithGetPrompt(ctx, request.RawParams)
	case protocol.ResourcesList:
		result, err = server.handleRequestWithListResources(ctx, request.RawParams)
	case protocol.ResourcesRead:
		result, err = server.handleRequestWithReadResource(ctx, request.RawParams)
	case protocol.ResourceListTemplates:
		result, err = server.handleRequestWithListPrompts(ctx, request.RawParams)
	case protocol.ResourcesSubscribe:
		result, err = server.handleRequestWithSubscribeResourceChange(ctx, request.RawParams)
	case protocol.ResourcesUnsubscribe:
		result, err = server.handleRequestWithUnSubscribeResourceChange(ctx, request.RawParams)
	case protocol.ToolsList:
		result, err = server.handleRequestWithListTools(ctx, request.RawParams)
	case protocol.ToolsCall:
		result, err = server.handleRequestWithCallTool(ctx, request.RawParams)
	case protocol.CompletionComplete:
		result, err = server.handleRequestWithCompleteRequest(ctx, request.RawParams)
	case protocol.LoggingSetLevel:
		result, err = server.handleRequestWithSetLogLevel(ctx, request.RawParams)
	default:
		// return protocol.NewJSONRPCErrorResponse(request.ID)
	}

	if err != nil {
		// return &protocol.NewJSONRPCErrorResponse(request.ID, ,err.Error())
	}

	return protocol.NewJSONRPCSuccessResponse(request.ID, result)
}

func (server *Server) receiveNotify(ctx context.Context, sessionID string, notify *protocol.JSONRPCNotification) error {
	return server.handleNotify(ctx, sessionID, notify)
}

func (server *Server) receiveResponse(ctx context.Context, sessionID string, response *protocol.JSONRPCResponse) error {
	// 通过 server.reqID2respChan 将 response 传回发送 request 的协程
	return nil
}
