package server

import "context"

type sessionIDKey struct{}

func setSessionIDToCtx(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, sessionIDKey{}, id)
}

func getSessionIDFromCtx(ctx context.Context) (string, bool) {
	id, exist := ctx.Value(sessionIDKey{}).(string)
	return id, exist
}
