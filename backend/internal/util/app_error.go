package util

import "fmt"

// AppError 业务错误，携带错误码与上下文。
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现 error 接口。
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s err=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is/As 链式判断。
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建业务错误。
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WrapAppError 包裹底层错误并携带业务信息。
func WrapAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
