// Package response 统一 API 响应封装。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
)

// Body 统一响应体。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 成功响应。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: constants.CodeOK, Message: constants.MsgOK, Data: data})
}

// OKMessage 带自定义文案的成功响应。
func OKMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Body{Code: constants.CodeOK, Message: message, Data: data})
}

// Fail 失败响应（自定义 HTTP 状态码）。
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Body{Code: code, Message: message})
}

// FailCode 使用错误码默认文案响应。
func FailCode(c *gin.Context, httpStatus, code int) {
	msg := constants.ErrorMessages[code]
	if msg == "" {
		msg = "未知错误"
	}
	Fail(c, httpStatus, code, msg)
}
