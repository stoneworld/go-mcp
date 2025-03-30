package client

import (
	"context"
	"errors"
	"fmt"

	"go-mcp/pkg"
	"go-mcp/protocol"

	"github.com/tidwall/gjson"
)

func (client *Client) Receive(ctx context.Context, msg []byte) error {
	defer pkg.Recover()

	if !gjson.GetBytes(msg, "id").Exists() {
		notify := &protocol.JSONRPCNotification{}
		if err := pkg.JsonUnmarshal(msg, &notify); err != nil {
			return err
		}
		go func() {
			defer pkg.Recover()

			if err := client.receiveNotify(ctx, notify); err != nil {
				client.logger.Errorf("receive notify:%+v error: %s", notify, err.Error())
				return
			}
		}()
		return nil
	}

	// 判断 request和response
	if !gjson.GetBytes(msg, "method").Exists() {
		resp := &protocol.JSONRPCResponse{}
		if err := pkg.JsonUnmarshal(msg, &resp); err != nil {
			return err
		}
		go func() {
			defer pkg.Recover()

			if err := client.receiveResponse(resp); err != nil {
				client.logger.Errorf("receive response:%+v error: %s", resp, err.Error())
				return
			}
		}()
		return nil
	}

	req := &protocol.JSONRPCRequest{}
	if err := pkg.JsonUnmarshal(msg, &req); err != nil {
		return err
	}
	if !req.IsValid() {
		return pkg.ErrRequestInvalid
	}
	go func() {
		defer pkg.Recover()

		if err := client.receiveRequest(ctx, req); err != nil {
			client.logger.Errorf("receive request:%+v error: %s", req, err.Error())
			return
		}
	}()
	return nil
}

func (client *Client) receiveRequest(ctx context.Context, request *protocol.JSONRPCRequest) error {
	var (
		result protocol.ClientResponse
		err    error
	)

	switch request.Method {
	case protocol.Ping:
		result, err = client.handleRequestWithPing()
	// case protocol.RootsList:
	// 	result, err = client.handleRequestWithListRoots(ctx, request.RawParams)
	// case protocol.SamplingCreateMessage:
	// 	result, err = client.handleRequestWithCreateMessagesSampling(ctx, request.RawParams)
	default:
		err = fmt.Errorf("%w: method=%s", pkg.ErrMethodNotSupport, request.Method)
	}

	if err != nil {
		if errors.Is(err, pkg.ErrMethodNotSupport) {
			return client.sendMsgWithError(ctx, request.ID, protocol.METHOD_NOT_FOUND, err.Error())
		} else if errors.Is(err, pkg.ErrRequestInvalid) {
			return client.sendMsgWithError(ctx, request.ID, protocol.INVALID_REQUEST, err.Error())
		} else if errors.Is(err, pkg.ErrJsonUnmarshal) {
			return client.sendMsgWithError(ctx, request.ID, protocol.PARSE_ERROR, err.Error())
		}
		return client.sendMsgWithError(ctx, request.ID, protocol.INTERNAL_ERROR, err.Error())
	}
	return client.sendMsgWithResponse(ctx, request.ID, result)
}

func (client *Client) receiveNotify(ctx context.Context, notify *protocol.JSONRPCNotification) error {
	switch notify.Method {
	// case protocol.NotificationCancelled:
	// 	return client.handleNotifyWithCancelled(ctx, notify.RawParams)
	// case protocol.NotificationProgress:
	// 	return nil
	// case protocol.NotificationLogMessage:
	//	return nil
	case protocol.NotificationToolsListChanged:
		return client.handleNotifyWithToolsListChanged(ctx, notify.RawParams)
	case protocol.NotificationPromptsListChanged:
		return client.handleNotifyWithPromptsListChanged(ctx, notify.RawParams)
	case protocol.NotificationResourcesListChanged:
		return client.handleNotifyWithResourcesListChanged(ctx, notify.RawParams)
	case protocol.NotificationResourcesUpdated:
		return client.handleNotifyWithResourcesUpdated(ctx, notify.RawParams)
	default:
		return fmt.Errorf("%w: method=%s", pkg.ErrMethodNotSupport, notify.Method)
	}
}

func (client *Client) receiveResponse(response *protocol.JSONRPCResponse) error {
	respChan, ok := client.reqID2respChan.Get(fmt.Sprint(response.ID))
	if !ok {
		return fmt.Errorf("%w: requestID=%+v", pkg.ErrLackResponseChan, response.ID)
	}

	select {
	case respChan <- response:
	default:
		return fmt.Errorf("%w: response=%+v", pkg.ErrDuplicateResponseReceived, response)
	}
	return nil
}
