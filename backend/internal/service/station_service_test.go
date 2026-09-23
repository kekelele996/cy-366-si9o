package service

import (
	"log/slog"
	"os"
	"testing"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestStationStatusTransition 机位状态机白盒测试。
func TestStationStatusTransition(t *testing.T) {
	cases := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{name: "idle_to_fault", from: constants.StationIdle, to: constants.StationFault, want: true},
		{name: "idle_to_reserved", from: constants.StationIdle, to: constants.StationReserved, want: true},
		{name: "fault_to_idle", from: constants.StationFault, to: constants.StationIdle, want: true},
		{name: "reserved_to_idle", from: constants.StationReserved, to: constants.StationIdle, want: true},
		{name: "using_to_idle", from: constants.StationUsing, to: constants.StationIdle, want: true},
		{name: "using_to_fault", from: constants.StationUsing, to: constants.StationFault, want: false},
		{name: "fault_to_using", from: constants.StationFault, to: constants.StationUsing, want: false},
		{name: "idle_to_idle", from: constants.StationIdle, to: constants.StationIdle, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := allowedStationTransition(tc.from, tc.to); got != tc.want {
				t.Fatalf("allowedStationTransition(%s->%s) = %v, want %v", tc.from, tc.to, got, tc.want)
			}
		})
	}
}

// TestDefaultGameType 默认游戏类型兜底。
func TestDefaultGameType(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "", want: constants.GameOther},
		{in: constants.GameLOL, want: constants.GameLOL},
		{in: constants.GameKOG, want: constants.GameKOG},
	}
	for _, tc := range cases {
		if got := defaultGameType(tc.in); got != tc.want {
			t.Fatalf("defaultGameType(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

var _ = dto.StationQuery{}
var _ = newTestLogger
