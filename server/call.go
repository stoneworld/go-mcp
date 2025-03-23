package server

import (
	"context"
	"encoding/json"

	"go-mcp/protocol"
)

// 请求

// 1. 请求构造
// 2. 发送请求 server.callClient(ctx)
// 3. 响应解析

func (server *Server) ping(ctx context.Context) error {
	// server.callClient(ctx)
	return nil
}

func (server *Server) listRoots(ctx context.Context) error {
	// server.callClient(ctx)
	return nil
}

func (server *Server) CreateMessagesSample(ctx context.Context) error {
	// 可以从 ctx 里取得 session id, sessionID, exist := getSessionIDFromCtx(ctx)
	// server.callClient(ctx)
	return nil
}

// 通知
// 1. 构造通知结构体
// 2. Cancelled、Progress、LoggingMessage类型的通知，可以从 ctx 里取得 session id, sessionID, exist := getSessionIDFromCtx(ctx)
// 3. 发送通知 server.sendMsgWithNotification(ctx)

func (server *Server) SendNotification4Cancelled(ctx context.Context) error {
	// sessionID, exist := getSessionIDFromCtx(ctx)
	// if !exist {
	// 	return pkg.errors
	// }
	// server.sendMsgWithNotification(ctx,sessionID)
	return nil
}

func (server *Server) SendNotification4Progress(ctx context.Context) error {
	return nil
}

func (server *Server) SendNotification4ToolListChanges(ctx context.Context) error {
	return nil
}

func (server *Server) SendNotification4PromptListChanges(ctx context.Context) error {
	return nil
}

func (server *Server) SendNotification4ResourceListChanges(ctx context.Context) error {
	return nil
}

func (server *Server) SendNotification4ResourcesUpdated(ctx context.Context) error {
	return nil
}

func (server *Server) SendNotification4LoggingMessage(ctx context.Context) error {
	return nil
}

// 负责request和response的拼接
func (server *Server) callClient(ctx context.Context, sessionID string, method protocol.Method, params protocol.ServerRequest) (json.RawMessage, error) {
	requestID := server.requestID.Add(1)
	// 发送请求
	if err := server.sendMsgWithRequest(ctx, sessionID, requestID, method, params); err != nil {
		return nil, err
	}

	// TODO：
	// 通过chan阻塞等待response
	// <- server.sessionID2session[sessionID].reqID2respChan
	// 使用ctx进行超时控制
	return nil, nil
}
