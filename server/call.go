package server

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/pkg"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func (server *Server) sendNotification4ToolListChanges(ctx context.Context) error {
	if server.capabilities.Tools == nil || !server.capabilities.Tools.ListChanged {
		return pkg.ErrServerNotSupport
	}

	var errList []error
	server.sessionID2session.Range(func(sessionID string, value *session) bool {
		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationToolsListChanged, protocol.NewToolListChangedNotification()); err != nil {
			errList = append(errList, fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return pkg.JoinErrors(errList)
}

func (server *Server) sendNotification4PromptListChanges(ctx context.Context) error {
	if server.capabilities.Prompts == nil || !server.capabilities.Prompts.ListChanged {
		return pkg.ErrServerNotSupport
	}

	var errList []error
	server.sessionID2session.Range(func(sessionID string, value *session) bool {
		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationPromptsListChanged, protocol.NewPromptListChangedNotification()); err != nil {
			errList = append(errList, fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return pkg.JoinErrors(errList)
}

func (server *Server) sendNotification4ResourceListChanges(ctx context.Context) error {
	if server.capabilities.Resources == nil || !server.capabilities.Resources.ListChanged {
		return pkg.ErrServerNotSupport
	}

	var errList []error
	server.sessionID2session.Range(func(sessionID string, value *session) bool {
		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationResourcesListChanged, protocol.NewResourceListChangedNotification()); err != nil {
			errList = append(errList, fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return pkg.JoinErrors(errList)
}

func (server *Server) SendNotification4ResourcesUpdated(ctx context.Context, notify *protocol.ResourceUpdatedNotification) error {
	if server.capabilities.Resources == nil || !server.capabilities.Resources.Subscribe {
		return pkg.ErrServerNotSupport
	}

	var errList []error
	server.sessionID2session.Range(func(sessionID string, s *session) bool {
		if _, ok := s.subscribedResources.Get(notify.URI); !ok {
			return true
		}

		if err := server.sendMsgWithNotification(ctx, sessionID, protocol.NotificationResourcesUpdated, notify); err != nil {
			errList = append(errList, fmt.Errorf("sessionID=%s, err: %w", sessionID, err))
		}
		return true
	})
	return pkg.JoinErrors(errList)
}

// Responsible for request and response assembly
// func (server *Server) callClient(ctx context.Context, sessionID string, method protocol.Method, params protocol.ServerRequest) (json.RawMessage, error) {
// 	session, ok := server.sessionID2session.Load(sessionID)
// 	if !ok {
// 		return nil, pkg.ErrLackSession
// 	}
//
// 	requestID := strconv.FormatInt(session.requestID.Add(1), 10)
//
// 	if err := server.sendMsgWithRequest(ctx, sessionID, requestID, method, params); err != nil {
// 		return nil, err
// 	}
//
// 	respChan := make(chan *protocol.JSONRPCResponse)
// 	session.reqID2respChan.Set(requestID, respChan)
//  defer session.reqID2respChan.Remove(requestID)
//
// 	select {
// 	case <-ctx.Done():
// 		return nil, ctx.Err()
// 	case response := <-respChan:
// 		if err := response.Error; err != nil {
// 			return nil, pkg.NewResponseError(err.Code, err.Message, err.Data)
// 		}
// 		return response.RawResult, nil
// 	}
// }
