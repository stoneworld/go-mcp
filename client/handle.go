package client

import (
	"context"
	"encoding/json"
	"fmt"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

func (client *Client) handleRequestWithPing() (*protocol.PingResult, error) {
	return protocol.NewPingResponse(), nil
}

func (client *Client) handleRequestWithListRoots(ctx context.Context, rawParams json.RawMessage) (*protocol.ListRootsResult, error) {
	request := &protocol.ListRootsRequest{}
	if err := parse(rawParams, request); err != nil {
		return nil, err
	}
	// TODO: 需要处理request.Cursor的翻页操作
	return &protocol.ListRootsResult{
		Roots: client.roots,
	}, nil
}

func (client *Client) handleRequestWithCreateMessagesSampling(ctx context.Context, rawParams json.RawMessage) (*protocol.CreateMessageResult, error) {
	request := &protocol.CreateMessageRequest{}
	if err := parse(rawParams, request); err != nil {
		return nil, err
	}
	return client.createMessagesSampleHandler(ctx, request)
}

func (client *Client) handleNotifyWithCancelled(ctx context.Context, rawParams json.RawMessage) error {
	param := &protocol.CancelledNotification{}
	if err := parse(rawParams, param); err != nil {
		return err
	}
	return client.cancelledNotifyHandler(ctx, param)
}

func parse(rawParams json.RawMessage, request interface{}) error {
	if err := pkg.JsonUnmarshal(rawParams, &request); err != nil {
		return fmt.Errorf("JsonUnmarshal: rawParams=%s, err=%w", rawParams, err)
	}
	return nil
}
