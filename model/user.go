package model

import (
	"context"
	"net/http"
)

const UserSession = "userSession"

type vkIDKey struct{}

func VKIDFromContext(ctx context.Context) (int, bool) {
	vkID, found := ctx.Value(vkIDKey{}).(int)
	return vkID, found
}

func VKIDFromRequest(r *http.Request) (int, bool) {
	return VKIDFromContext(r.Context())
}

func SetVKIDOnContext(ctx context.Context, vkID int) context.Context {
	return context.WithValue(ctx, vkIDKey{}, vkID)
}

func SetVKIDOnRequest(r *http.Request, vkID int) *http.Request {
	return r.WithContext(SetVKIDOnContext(r.Context(), vkID))
}
