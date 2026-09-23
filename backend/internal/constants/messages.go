package constants

// 接口返回文案、日志文案、错误提示文案统一维护。
const (
	MsgOK            = "ok"
	MsgCreateSuccess = "创建成功"
	MsgUpdateSuccess = "更新成功"
	MsgDeleteSuccess = "删除成功"
	MsgLoginSuccess  = "登录成功"
	MsgLogoutSuccess = "退出登录成功"
	MsgRegisterOK    = "注册成功"

	MsgRechargeOK   = "充值成功"
	MsgBuyPackageOK = "时长包购买成功"
	MsgReserveOK    = "预约成功"
	MsgWaitlistOK   = "该时段已被预约，已为你加入候补，请留意顺位变化"
	MsgCheckInOK    = "开机成功"
	MsgRenewOK      = "续费成功"
	MsgCheckoutOK   = "下机成功"
	MsgDrawOK       = "抽签分组完成"
)

// 日志文案模板（非格式化部分）。
const (
	LogUserLogin     = "用户登录"
	LogUserRegister  = "用户注册"
	LogStationChange = "机位状态变更"
	LogRecharge      = "会员充值"
	LogReservation   = "机位预约"
	LogSessionStart  = "会员上机"
	LogSessionEnd    = "会员下机"
	LogTournament    = "赛事操作"
	LogAuditAction   = "审计动作"
)
