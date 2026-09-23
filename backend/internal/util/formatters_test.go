package util

import "testing"

func TestStatusText(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "idle", want: "空闲"},
		{in: "using", want: "使用中"},
		{in: "fault", want: "故障"},
		{in: "reserved", want: "已预约"},
		{in: "pending", want: "待确认"},
		{in: "confirmed", want: "已确认"},
		{in: "completed", want: "已完成"},
		{in: "unknown", want: "unknown"},
	}
	for _, tc := range cases {
		if got := StatusText(tc.in); got != tc.want {
			t.Fatalf("StatusText(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRoleText(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "admin", want: "管理员"},
		{in: "staff", want: "店员"},
		{in: "member", want: "会员"},
	}
	for _, tc := range cases {
		if got := RoleText(tc.in); got != tc.want {
			t.Fatalf("RoleText(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestGameTypeText(t *testing.T) {
	if got := GameTypeText("lol"); got != "英雄联盟" {
		t.Fatalf("GameTypeText(lol) = %q", got)
	}
	if got := GameTypeText("csgo"); got != "CSGO" {
		t.Fatalf("GameTypeText(csgo) = %q", got)
	}
}
