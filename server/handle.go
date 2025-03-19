package server

import (
	"context"
	"encoding/json"
	"errors"

	"go-mcp/protocol"
)

func (server *Server) handleRequestWithInitialize(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithPing(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListPrompts(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithGetPrompt(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListResources(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithReadResource(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithReadResourceTemplates(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithSubscribeResourceChange(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithUnSubscribeResourceChange(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithListTools(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithCallTool(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithCompleteRequest(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleRequestWithSetLogLevel(ctx context.Context, rawParams json.RawMessage) (protocol.Result, error) {
	return nil, nil
}

func (server *Server) handleNotify(ctx context.Context, notify *protocol.JSONRPCNotification) error {
	if notify.Method == "" {
		return errors.New("notify method can't is \"\"")
	}

	// TODO: 使用server里定义一个 notifyMethod2handler 对通知进行处理
	handler, ok := server.notifyMethod2handler[notify.Method]
	if !ok {
		// 打印 warn/info 日志
		// 此处也可以向上抛error，在上层识别error统一打日志
		return nil
	}
	handler(notify.Params)
	return nil
}
