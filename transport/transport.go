package transport

import (
	"context"
)

/*
* Transport 是对底层传输层的抽象。
* GO-MCP 需要能够在 server/client 间传递 JSON-RPC 消息。
 */

// Message 定义基础消息接口
type Message []byte

type ClientTransport interface {
	// Start 启动传输连接
	Start() error

	// Send 发送消息
	Send(ctx context.Context, msg Message) error
	// SetReceiver 设置对对端消息的处理器
	SetReceiver(receiver ClientReceiver)

	// Close 关闭传输连接
	Close() error
}

type ClientReceiver interface {
	Receive(ctx context.Context, msg []byte)
}

type ServerTransport interface {
	// Run 开始监听, 这是同步的, 在 Shutdown 之前, 不能返回
	Run() error

	// Send 发送消息
	Send(ctx context.Context, sessionID string, msg Message) error
	// SetReceiver 设置对对端消息的处理器
	SetReceiver(ServerReceiver)

	// Shutdown 关闭监听, 内部实现时需要使用 ctx 控制超时
	Shutdown(ctx context.Context) error
}

type ServerReceiver interface {
	Receive(ctx context.Context, sessionID string, msg []byte)
}
