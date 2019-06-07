package core

import (
	"context"
	"time"
)

const ContextKeyNow = "core:now"

func WithNow(ctx context.Context) context.Context { return context.WithValue(ctx, ContextKeyNow, time.Now()) }

func Now(ctx context.Context) (time.Time, context.Context) {
	t, ok := ctx.Value(ContextKeyNow).(time.Time)
	if !ok {
		t = time.Now()
		return t, context.WithValue(ctx, ContextKeyNow, t)
	}
	return t, ctx
}
