package util

import (
	"fmt"
	"time"
)

// 状态文本、类型文本、日期格式化工具，全项目多处耦合引用。

// StatusText 将状态枚举转为中文展示文本。
func StatusText(status string) string {
	switch status {
	case "idle":
		return "空闲"
	case "using":
		return "使用中"
	case "fault":
		return "故障"
	case "reserved":
		return "已预约"
	case "pending":
		return "待确认"
	case "confirmed":
		return "已确认"
	case "checked_in":
		return "已开机"
	case "completed":
		return "已完成"
	case "cancelled":
		return "已取消"
	case "active":
		return "进行中"
	case "draft":
		return "草稿"
	case "open":
		return "报名中"
	case "ready":
		return "已分组"
	case "finished":
		return "已结束"
	case "paid":
		return "已支付"
	case "joined":
		return "已通过"
	case "rejected":
		return "已拒绝"
	case "playing":
		return "进行中"
	}
	return status
}

// RoleText 将角色枚举转为中文。
func RoleText(role string) string {
	switch role {
	case "admin":
		return "管理员"
	case "staff":
		return "店员"
	case "member":
		return "会员"
	}
	return role
}

// GameTypeText 将游戏类型转为中文。
func GameTypeText(gameType string) string {
	switch gameType {
	case "lol":
		return "英雄联盟"
	case "csgo":
		return "CSGO"
	case "kog":
		return "王者荣耀"
	case "other":
		return "其他"
	}
	return gameType
}

// FormatTime 格式化时间。
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDuration 将分钟数格式化为 x小时y分钟。
func FormatDuration(minutes int) string {
	if minutes < 0 {
		minutes = 0
	}
	return fmt.Sprintf("%d小时%d分钟", minutes/60, minutes%60)
}
