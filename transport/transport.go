package transport

import (
	"context"
)

// Message 定义基础消息接口
type Message []byte

// Transport 是对底层传输层的抽象。
// GO-MCP 需要能够在 server/client 间传递 JSON-RPC 消息。
type ClientTransport interface {
	// Start 启动传输连接
	Start() error

	// Send 发送消息
	Send(ctx context.Context, msg Message) error
	// SetReceiver 设置对对端消息的处理器
	SetReceiver(receiver clientReceiver)

	// Close 关闭传输连接
	Close() error
}

type clientReceiver interface {
	Receive(ctx context.Context, msg []byte)
}

type ServerTransport interface {
	// Start 启动传输连接
	Start() error

	// Send 发送消息
	Send(ctx context.Context, sessionID string, msg Message) error
	// SetReceiver 设置对对端消息的处理器
	SetReceiver(serverReceiver)

	// Close 关闭传输连接
	Close() error
}

type serverReceiver interface {
	Receive(ctx context.Context, sessionID string, msg []byte)
}

// // Config 定义传输层配置
// type Config struct {
// 	// BufferSize 定义缓冲区大小
// 	BufferSize int
// 	// HeartbeatInterval 心跳间隔时间(秒)
// 	HeartbeatInterval int
// 	// ReadTimeout 读取超时时间(秒)
// 	ReadTimeout int
// 	// WriteTimeout 写入超时时间(秒)
// 	WriteTimeout int
// }
//
// // DefaultConfig 返回默认配置
// func DefaultConfig() *Config {
// 	return &Config{
// 		BufferSize:        4096,
// 		HeartbeatInterval: 30,
// 		ReadTimeout:       60,
// 		WriteTimeout:      60,
// 	}
// }
//
// // TransportError 定义传输层错误
// type TransportError struct {
// 	Code    int
// 	Message string
// }
//
// // Error 实现error接口
// func (e *TransportError) Error() string {
// 	return e.Message
// }
//
// // Common transport errors
// var (
// 	ErrConnectionClosed = &TransportError{Code: 1001, Message: "connection closed"}
// 	ErrReadTimeout      = &TransportError{Code: 1002, Message: "read timeout"}
// 	ErrWriteTimeout     = &TransportError{Code: 1003, Message: "write timeout"}
// 	ErrInvalidMessage   = &TransportError{Code: 1004, Message: "invalid message"}
// )
