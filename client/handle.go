package client

import (
	"context"
	"encoding/json"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func (client *Client) handleRequestWithPing() (*protocol.PingResult, error) {
	return protocol.NewPingResponse(), nil
}

func (client *Client) handleNotifyWithToolsListChanged(ctx context.Context, rawParams json.RawMessage) error {
	notify := &protocol.ToolListChangedNotification{}
	if err := pkg.JsonUnmarshal(rawParams, notify); err != nil {
		return err
	}
	return client.notifyHandlerWithToolsListChanged(ctx, notify)
}

func (client *Client) handleNotifyWithPromptsListChanged(ctx context.Context, rawParams json.RawMessage) error {
	notify := &protocol.PromptListChangedNotification{}
	if err := pkg.JsonUnmarshal(rawParams, notify); err != nil {
		return err
	}
	return client.notifyHandlerWithPromptListChanged(ctx, notify)
}

func (client *Client) handleNotifyWithResourcesListChanged(ctx context.Context, rawParams json.RawMessage) error {
	notify := &protocol.ResourceListChangedNotification{}
	if err := pkg.JsonUnmarshal(rawParams, notify); err != nil {
		return err
	}
	return client.notifyHandlerWithResourceListChanged(ctx, notify)
}

func (client *Client) handleNotifyWithResourcesUpdated(ctx context.Context, rawParams json.RawMessage) error {
	notify := &protocol.ResourceUpdatedNotification{}
	if err := pkg.JsonUnmarshal(rawParams, notify); err != nil {
		return err
	}
	return client.notifyHandlerWithResourcesUpdated(ctx, notify)
}
