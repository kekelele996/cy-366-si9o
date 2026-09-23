package handler

import (
	"github.com/esportsbar/backend/internal/constants"
)

// httpStatusFor 错误码到 HTTP 状态码映射（handler 层再次集中）。
func httpStatusFor(code int) int {
	switch code {
	case constants.CodeUnauthorized:
		return 401
	case constants.CodeForbidden:
		return 403
	case constants.CodeNotFound:
		return 404
	case constants.CodeValidation, constants.CodeBadRequest:
		return 400
	case constants.CodeConflict, constants.CodeUserExists, constants.CodeStationBusy,
		constants.CodeStationFault, constants.CodeInsufficient, constants.CodeReservation,
		constants.CodeSessionOpen, constants.CodeTournament:
		return 409
	default:
		return 500
	}
}
