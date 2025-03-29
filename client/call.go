package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

// 请求

// 1. 请求构造
// 2. 发送请求 client.callServer(ctx)
// 3. 响应解析

func (client *Client) initialization(ctx context.Context, request protocol.InitializeRequest) (*protocol.InitializeResult, error) {
	response, err := client.callServer(ctx, protocol.Initialize, request)
	if err != nil {
		return nil, err
	}
	var result protocol.InitializeResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	//TODO add meta
	notify := &protocol.InitializedNotification{
		Meta: map[string]interface{}{},
	}
	if err := client.sendNotification4Initialized(ctx, notify); nil != err {
		return nil, fmt.Errorf("failed to send InitializedNotification: %w", err)
	}

	client.capabilities = result.Capabilities
	client.initialized = true
	return &result, nil
}

func (client *Client) Ping(ctx context.Context, request protocol.PingRequest) (*protocol.PingResult, error) {
	response, err := client.callServer(ctx, protocol.Ping, request)
	if err != nil {
		return nil, err
	}

	var result protocol.PingResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) ListPrompts(ctx context.Context, request protocol.ListPromptsRequest) (*protocol.ListPromptsResult, error) {
	response, err := client.callServer(ctx, protocol.PromptsList, request)
	if err != nil {
		return nil, err
	}

	var result protocol.ListPromptsResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) GetPrompt(ctx context.Context, request protocol.GetPromptRequest) (*protocol.GetPromptResult, error) {
	response, err := client.callServer(ctx, protocol.PromptsGet, request)
	if err != nil {
		return nil, err
	}

	var result protocol.GetPromptResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

func (client *Client) ListResources(ctx context.Context, request protocol.ListResourcesRequest) (*protocol.ListResourcesResult, error) {
	response, err := client.callServer(ctx, protocol.ResourcesList, request)
	if err != nil {
		return nil, err
	}

	var result protocol.ListResourcesResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, err
}

func (client *Client) ReadResource(ctx context.Context, request protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
	response, err := client.callServer(ctx, protocol.ResourcesRead, request)
	if err != nil {
		return nil, err
	}

	var result protocol.ReadResourceResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) ListResourceTemplates(ctx context.Context, request protocol.ListResourceTemplatesRequest) (*protocol.ListResourceTemplatesResult, error) {
	response, err := client.callServer(ctx, protocol.ResourceListTemplates, request)
	if err != nil {
		return nil, err
	}

	var result protocol.ListResourceTemplatesResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) SubscribeResourceChange(ctx context.Context, request protocol.SubscribeRequest) (*protocol.SubscribeResult, error) {
	response, err := client.callServer(ctx, protocol.ResourcesSubscribe, request)
	if err != nil {
		return nil, err
	}

	var result protocol.SubscribeResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) UnSubscribeResourceChange(ctx context.Context, request protocol.UnsubscribeRequest) (*protocol.UnsubscribeResult, error) {
	response, err := client.callServer(ctx, protocol.ResourcesUnsubscribe, request)
	if err != nil {
		return nil, err
	}

	var result protocol.UnsubscribeResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) ListTools(ctx context.Context, request protocol.ListToolsRequest) (*protocol.ListToolsResult, error) {
	response, err := client.callServer(ctx, protocol.ToolsList, request)
	if err != nil {
		return nil, err
	}

	var result protocol.ListToolsResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) CallTool(ctx context.Context, request protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	response, err := client.callServer(ctx, protocol.ToolsCall, request)
	if err != nil {
		return nil, err
	}

	var result protocol.CallToolResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) CompleteRequest(ctx context.Context, request protocol.CompleteRequest) (*protocol.CompleteResult, error) {
	response, err := client.callServer(ctx, protocol.CompletionComplete, request)
	if err != nil {
		return nil, err
	}

	var result protocol.CompleteResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

func (client *Client) SetLogLevel(ctx context.Context, request protocol.SetLoggingLevelResult) (*protocol.SetLoggingLevelResult, error) {
	response, err := client.callServer(ctx, protocol.LoggingSetLevel, request)
	if err != nil {
		return nil, err
	}
	var result protocol.SetLoggingLevelResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}

// 通知
// 1. 构造通知结构体
// 2. 发送通知 client.sendMsgWithNotification(ctx)

func (client *Client) sendNotification4Initialized(ctx context.Context, notify *protocol.InitializedNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationInitialized, notify)
}

func (client *Client) SendNotification4Cancelled(ctx context.Context, notify *protocol.CancelledNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationCancelled, notify)
}

func (client *Client) SendNotification4Progress(ctx context.Context, notify *protocol.ProgressNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationProgress, notify)
}

func (client *Client) SendNotification4RootListChanges(ctx context.Context, notify *protocol.RootsListChangedNotification) error {
	return client.sendMsgWithNotification(ctx, protocol.NotificationRootsListChanged, notify)
}

func (client *Client) callAndParse(ctx context.Context, method protocol.Method, request protocol.ClientRequest, result protocol.ServerResponse) error {
	rawResult, err := client.callServer(ctx, method, request)
	if err != nil {
		return err
	}

	return pkg.JsonUnmarshal(rawResult, &result)
}

// 负责request和response的拼接
func (client *Client) callServer(ctx context.Context, method protocol.Method, params protocol.ClientRequest) (json.RawMessage, error) {
	if !client.initialized && method != protocol.Initialize {
		return nil, fmt.Errorf("client not initialized")
	}

	requestID := strconv.FormatInt(client.requestID.Add(1), 10)
	// 发送请求
	if err := client.sendMsgWithRequest(ctx, requestID, method, params); err != nil {
		return nil, fmt.Errorf("callServer: %w", err)
	}

	respChan := make(chan *protocol.JSONRPCResponse)

	client.reqID2respChan.Set(requestID, respChan)

	select {
	case <-ctx.Done():
		client.reqID2respChan.Remove(requestID)
		return nil, ctx.Err()
	case response := <-respChan:
		if err := response.Error; err != nil {
			return nil, pkg.NewResponseError(err.Code, err.Message, err.Data)
		}
		return response.RawResult, nil
	}
}
