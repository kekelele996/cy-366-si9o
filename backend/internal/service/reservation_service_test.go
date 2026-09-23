package service

import (
	"testing"

	"github.com/esportsbar/backend/internal/constants"
)

// TestReservationCancelable 可取消状态白盒测试。
func TestReservationCancelable(t *testing.T) {
	cases := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "pending", status: constants.ReservationPending, want: true},
		{name: "confirmed", status: constants.ReservationConfirmed, want: true},
		{name: "waitlisted", status: constants.ReservationWaitlisted, want: true},
		{name: "checked_in", status: constants.ReservationCheckedIn, want: true},
		{name: "cancelled_no_second_cancel", status: constants.ReservationCancelled, want: false},
		{name: "completed", status: constants.ReservationCompleted, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := cancelable(tc.status); got != tc.want {
				t.Fatalf("cancelable(%s) = %v, want %v", tc.status, got, tc.want)
			}
		})
	}
}
