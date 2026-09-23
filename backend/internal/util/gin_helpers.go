package util

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/pkg/response"
)

// AbortJSON 输出统一错误响应。
func AbortJSON(c *gin.Context, httpStatus, code int, message string) {
	response.Fail(c, httpStatus, code, message)
}
