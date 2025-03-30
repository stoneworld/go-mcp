package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"go-mcp/pkg"
	"go-mcp/protocol"
)

func (server *Server) ping(ctx context.Context) error {
	sessionID, exist := getSessionIDFromCtx(ctx)
	if !exist {
		return pkg.ErrLackSession
	}

	response, err := server.callClient(ctx, sessionID, protocol.Ping, protocol.NewPingRequest())
	if err != nil {
		return err
	}

	var result protocol.PingResult
	if err := pkg.JsonUnmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}

func (server *Server) SendNotification4ToolListChanges(ctx context.Context) error {
	var errList error
	server.sessionID2session.Range(func(sessionID string, value interface{}) bool {
		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationToolsListChanged, protocol.NewToolListChangedNotification()); err != nil {
			errList = errors.Join(fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return errList
}

func (server *Server) SendNotification4PromptListChanges(ctx context.Context) error {
	var errList error
	server.sessionID2session.Range(func(sessionID string, value interface{}) bool {
		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationPromptsListChanged, protocol.NewPromptListChangedNotification()); err != nil {
			errList = errors.Join(fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return errList
}

func (server *Server) SendNotification4ResourceListChanges(ctx context.Context) error {
	var errList error
	server.sessionID2session.Range(func(sessionID string, value interface{}) bool {
		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationResourcesListChanged, protocol.NewResourceListChangedNotification()); err != nil {
			errList = errors.Join(fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return errList
}

func (server *Server) SendNotification4ResourcesUpdated(ctx context.Context, notify *protocol.ResourceUpdatedNotification) error {
	var errList error
	server.sessionID2session.Range(func(sessionID string, value interface{}) bool {
		s, ok := value.(*session)
		if !ok {
			return true
		}

		if _, ok := s.subscribedResources.Get(notify.URI); !ok {
			return true
		}

		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationResourcesUpdated, notify); err != nil {
			errList = errors.Join(fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return errList
}

// 负责request和response的拼接
func (server *Server) callClient(ctx context.Context, sessionID string, method protocol.Method, params protocol.ServerRequest) (json.RawMessage, error) {
	value, ok := server.sessionID2session.Load(sessionID)
	if !ok {
		return nil, pkg.ErrLackSession
	}
	session := value.(*session)

	requestID := strconv.FormatInt(session.requestID.Add(1), 10)

	if err := server.sendMsgWithRequest(ctx, sessionID, requestID, method, params); err != nil {
		return nil, err
	}

	respChan := make(chan *protocol.JSONRPCResponse)
	session.reqID2respChan.Set(requestID, respChan)

	select {
	case <-ctx.Done():
		session.reqID2respChan.Remove(requestID)
		return nil, ctx.Err()
	case response := <-respChan:
		if err := response.Error; err != nil {
			return nil, pkg.NewResponseError(err.Code, err.Message, err.Data)
		}
		return response.RawResult, nil
	}
}
