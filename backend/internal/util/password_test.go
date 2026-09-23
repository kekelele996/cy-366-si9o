package util

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		expected bool
	}{
		{name: "valid", raw: "admin123456", expected: true},
		{name: "wrong", raw: "wrongpass", expected: false},
		{name: "empty", raw: "", expected: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := HashPassword("admin123456")
			if err != nil {
				t.Fatalf("HashPassword error: %v", err)
			}
			if got := CheckPassword(hash, tc.raw); got != tc.expected {
				t.Fatalf("CheckPassword(%q) = %v, want %v", tc.raw, got, tc.expected)
			}
		})
	}
}
