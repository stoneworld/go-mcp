package server

// 对 request、notification 路由到对应 handler

// MessageType 定义消息类型
type MessageType int

const (
	// TypeRequest 请求消息类型
	TypeRequest MessageType = iota + 1
	// TypeResponse 响应消息类型
	TypeResponse
	// TypeNotification 通知消息类型
	TypeNotification
)
