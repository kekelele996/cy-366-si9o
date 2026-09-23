package util

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// NewRequestID 生成请求追踪 ID。
func NewRequestID() string {
	return uuid.NewString()
}

// WithRequestID 将请求 ID 写入 context。
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID 从 context 读取请求 ID。
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
