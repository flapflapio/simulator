package app

import (
	"testing"

	"github.com/flapflapio/simulator/internal/simtest"
)

func TestNonNilStr(t *testing.T) {
	for _, tc := range []struct {
		name   string
		s1, s2 *string
		expect func(s1, s2 *string) *string
	}{
		{
			name:   "both non-nil",
			s1:     simtest.StringPointer("s1"),
			s2:     simtest.StringPointer("s2"),
			expect: func(s1, s2 *string) *string { return s2 },
		},
		{
			name:   "s1 non-nil s2 nil",
			s1:     simtest.StringPointer("s1"),
			expect: func(s1, s2 *string) *string { return s1 },
		},
		{
			name:   "s1 nil s2 non-nil",
			s2:     simtest.StringPointer("s2"),
			expect: func(s1, s2 *string) *string { return s2 },
		},
		{
			name:   "both nil",
			expect: func(s1, s2 *string) *string { return s1 },
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := takeNonNilStr(tc.s1, tc.s2)
			if expected := tc.expect(tc.s1, tc.s2); actual != expected {
				exp := "s1"
				notExp := "s2"
				if expected == tc.s2 {
					exp = "s2"
					exp = "s1"
				}
				t.Fatalf("Expected %v to be returned but got %v", exp, notExp)
			}
		})
	}
}
