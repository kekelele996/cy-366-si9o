package service

import (
	"testing"
	"time"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/model"
)

func mkTime(h, m int) time.Time {
	return time.Date(2026, 9, 23, h, m, 0, 0, time.UTC)
}

func mkReservation(id uint, startH, startM, endH, endM int, status string, createdAt time.Time) model.Reservation {
	return model.Reservation{
		ID:        id,
		StationID: 1,
		StartTime: mkTime(startH, startM),
		EndTime:   mkTime(endH, endM),
		Status:    status,
		CreatedAt: createdAt,
	}
}

// TestOverlaps 时段重叠判定。
func TestOverlaps(t *testing.T) {
	cases := []struct {
		name           string
		s1, e1, s2, e2 [2]int
		want           bool
	}{
		{name: "overlap", s1: [2]int{10, 0}, e1: [2]int{12, 0}, s2: [2]int{11, 0}, e2: [2]int{13, 0}, want: true},
		{name: "adjacent_not_overlap", s1: [2]int{10, 0}, e1: [2]int{12, 0}, s2: [2]int{12, 0}, e2: [2]int{13, 0}, want: false},
		{name: "disjoint", s1: [2]int{10, 0}, e1: [2]int{11, 0}, s2: [2]int{12, 0}, e2: [2]int{13, 0}, want: false},
		{name: "contained", s1: [2]int{10, 0}, e1: [2]int{14, 0}, s2: [2]int{11, 0}, e2: [2]int{12, 0}, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := overlaps(
				mkTime(tc.s1[0], tc.s1[1]), mkTime(tc.e1[0], tc.e1[1]),
				mkTime(tc.s2[0], tc.s2[1]), mkTime(tc.e2[0], tc.e2[1]),
			)
			if got != tc.want {
				t.Fatalf("overlaps() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestWaitlistPositionInQueue 候补顺位按提交顺序计算，且只统计时段重叠者。
func TestWaitlistPositionInQueue(t *testing.T) {
	base := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	// 同机位候补队列：W1 10-12（最早）、W2 10-12（次之）、W3 15-17（不同时段，最早）。
	w1 := mkReservation(1, 10, 0, 12, 0, constants.ReservationWaitlisted, base)
	w2 := mkReservation(2, 10, 0, 12, 0, constants.ReservationWaitlisted, base.Add(time.Minute))
	w3 := mkReservation(3, 15, 0, 17, 0, constants.ReservationWaitlisted, base.Add(2*time.Minute))
	queue := []model.Reservation{w1, w2, w3}

	if pos := waitlistPositionInQueue(&w1, queue); pos != 1 {
		t.Fatalf("w1 position = %d, want 1", pos)
	}
	if pos := waitlistPositionInQueue(&w2, queue); pos != 2 {
		t.Fatalf("w2 position = %d, want 2", pos)
	}
	if pos := waitlistPositionInQueue(&w3, queue); pos != 1 {
		t.Fatalf("w3 position = %d, want 1（不同时段不排队）", pos)
	}
}

// TestPickPromotable 取消兑现规则：最早提交且不再冲突者；冲突者跳过不淘汰；只兑现一人。
func TestPickPromotable(t *testing.T) {
	base := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	cancelled := mkReservation(100, 10, 0, 12, 0, constants.ReservationCancelled, base.Add(-time.Hour))

	t.Run("earliest_non_conflicting_wins", func(t *testing.T) {
		w1 := mkReservation(1, 10, 0, 12, 0, constants.ReservationWaitlisted, base)
		w2 := mkReservation(2, 10, 30, 11, 30, constants.ReservationWaitlisted, base.Add(time.Minute))
		w3 := mkReservation(3, 11, 0, 13, 0, constants.ReservationWaitlisted, base.Add(2*time.Minute))
		queue := []model.Reservation{w1, w2, w3}
		// w1 仍与另一笔有效预约冲突（如同时段的 checked_in 或另一确认单），w2 不冲突。
		conflictIDs := map[uint]bool{1: true, 3: true}
		picked := pickPromotable(queue, &cancelled, func(c *model.Reservation) (bool, error) {
			return conflictIDs[c.ID], nil
		})
		if picked == nil || picked.ID != 2 {
			got := uint(0)
			if picked != nil {
				got = picked.ID
			}
			t.Fatalf("picked = %v, want w2(id=2)", got)
		}
	})

	t.Run("only_one_picked", func(t *testing.T) {
		w1 := mkReservation(1, 10, 0, 12, 0, constants.ReservationWaitlisted, base)
		w2 := mkReservation(2, 10, 0, 12, 0, constants.ReservationWaitlisted, base.Add(time.Minute))
		picked := pickPromotable([]model.Reservation{w1, w2}, &cancelled, func(*model.Reservation) (bool, error) {
			return false, nil
		})
		if picked == nil || picked.ID != 1 {
			t.Fatalf("picked should be earliest w1(id=1)")
		}
	})

	t.Run("no_overlap_no_pick", func(t *testing.T) {
		w1 := mkReservation(1, 15, 0, 17, 0, constants.ReservationWaitlisted, base)
		picked := pickPromotable([]model.Reservation{w1}, &cancelled, func(*model.Reservation) (bool, error) {
			return false, nil
		})
		if picked != nil {
			t.Fatalf("non-overlapping waitlist must not be picked, got id=%d", picked.ID)
		}
	})

	t.Run("all_still_conflict_no_pick", func(t *testing.T) {
		w1 := mkReservation(1, 10, 0, 12, 0, constants.ReservationWaitlisted, base)
		w2 := mkReservation(2, 10, 30, 11, 30, constants.ReservationWaitlisted, base.Add(time.Minute))
		picked := pickPromotable([]model.Reservation{w1, w2}, &cancelled, func(*model.Reservation) (bool, error) {
			return true, nil
		})
		if picked != nil {
			t.Fatalf("all conflicting, want nil, got id=%d", picked.ID)
		}
	})
}

// TestCancellable 取消状态白名单。
func TestCancellable(t *testing.T) {
	cases := map[string]bool{
		constants.ReservationPending:    true,
		constants.ReservationWaitlisted: true,
		constants.ReservationConfirmed:  true,
		constants.ReservationCheckedIn:  true,
		constants.ReservationCompleted:  false,
		constants.ReservationCancelled:  false,
	}
	for status, want := range cases {
		if got := cancellable(status); got != want {
			t.Fatalf("cancellable(%s) = %v, want %v", status, got, want)
		}
	}
}

var _ = constants.ReservationWaitlisted
