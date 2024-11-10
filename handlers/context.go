package handlers

import (
	"context"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	adminKey  contextKey = "admin"
)

func setUserContext(ctx context.Context, userID string, isAdmin bool) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	return context.WithValue(ctx, adminKey, isAdmin)
}

func getUserFromContext(ctx context.Context) (string, bool) {
	userID, _ := ctx.Value(userIDKey).(string)
	isAdmin, _ := ctx.Value(adminKey).(bool)
	return userID, isAdmin
}
