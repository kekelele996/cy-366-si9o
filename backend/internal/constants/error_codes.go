package constants

// 全局错误码，集中维护。0 表示成功，非 0 表示业务错误。
const (
	CodeOK           = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeConflict     = 40900
	CodeInternal     = 50000
	CodeValidation   = 42200
	CodeUserExists   = 40001
	CodeUserNotFound = 40002
	CodeWrongPass    = 40003
	CodeStationBusy  = 40004
	CodeStationFault = 40005
	CodeInsufficient = 40006
	CodeReservation  = 40007
	CodeSessionOpen  = 40008
	CodeTournament   = 40009
)

// 错误码与默认文案映射。
var ErrorMessages = map[int]string{
	CodeOK:           "ok",
	CodeBadRequest:   "请求参数错误",
	CodeUnauthorized: "未登录或登录已过期",
	CodeForbidden:    "没有操作权限",
	CodeNotFound:     "资源不存在",
	CodeConflict:     "资源状态冲突",
	CodeInternal:     "服务器内部错误",
	CodeValidation:   "参数校验失败",
	CodeUserExists:   "用户名已存在",
	CodeUserNotFound: "用户不存在",
	CodeWrongPass:    "用户名或密码错误",
	CodeStationBusy:  "机位正在使用中，无法执行该操作",
	CodeStationFault: "机位处于故障状态，无法执行该操作",
	CodeInsufficient: "会员余额不足",
	CodeReservation:  "预约状态不允许该操作",
	CodeSessionOpen:  "该机位已有进行中的上机记录",
	CodeTournament:   "赛事状态不允许该操作",
}
