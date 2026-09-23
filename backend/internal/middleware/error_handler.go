package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/util"
)

// ErrorHandler 统一错误处理中间件：捕获 panic 与业务错误。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered", "err", rec, "path", c.FullPath())
				util.AbortJSON(c, http.StatusInternalServerError, constants.CodeInternal, "服务器内部错误")
			}
		}()
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *util.AppError
			if errors.As(err, &appErr) {
				util.AbortJSON(c, httpStatus(appErr.Code), appErr.Code, appErr.Message)
				return
			}
			logger.Error("request error", "err", err, "path", c.FullPath())
			util.AbortJSON(c, http.StatusInternalServerError, constants.CodeInternal, "服务器内部错误")
		}
	}
}

// httpStatus 错误码到 HTTP 状态码映射。
func httpStatus(code int) int {
	switch code {
	case constants.CodeUnauthorized:
		return http.StatusUnauthorized
	case constants.CodeForbidden:
		return http.StatusForbidden
	case constants.CodeNotFound:
		return http.StatusNotFound
	case constants.CodeValidation, constants.CodeBadRequest:
		return http.StatusBadRequest
	case constants.CodeConflict, constants.CodeUserExists, constants.CodeStationBusy,
		constants.CodeStationFault, constants.CodeInsufficient, constants.CodeReservation,
		constants.CodeSessionOpen, constants.CodeTournament:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
