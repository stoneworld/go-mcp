package server

import (
	"context"
	"encoding/json"
	"errors"

	"go-mcp/protocol"
)

func (server *Server) handleRequestWithInitialize(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
}

func (server *Server) handleRequestWithPing(ctx context.Context, rawParams json.RawMessage) (protocol.ServerResponse, error) {
	return nil, nil
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
	return nil, nil
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

func (server *Server) handleNotify(ctx context.Context, sessionID string, notify *protocol.JSONRPCNotification) error {
	if notify.Method == "" {
		return errors.New(`notify method can't is ""`)
	}

	if notify.Method == protocol.NotificationInitialized {
		// TODO: 官方文档约定“服务器在接收到 initialized 通知前，不应发送除 ping 和日志记录之外的其他请求。” 如果这样相当于server层要感知 SessionID。
		close(server.sessionID2session[sessionID].readyChan)
		return nil
	}

	// TODO: 使用server里定义一个 notifyMethod2handler 对通知进行处理
	handler, ok := server.notifyMethod2handler[notify.Method]
	if !ok {
		// 打印 warn/info 日志
		// 此处也可以向上抛error，在上层识别error统一打日志
		return nil
	}
	handler(notify.RawParams)
	return nil
}
