package shup

import (
	"testing"
)

func equal[T comparable](a, b []T) bool {
	for i, j := range a {
		if j != b[i] {
			return false
		}
	}
	return true
}

func TestBasics(t *testing.T) {
	s := NewSet[int]([]int{1, 1, 2, 2, 2, 3, 3, 3, 3}...)
	if len(s) != 3 {
		t.Errorf("wrong set length: expected 3, got %d", len(s))
	}

	for i := 1; i <= 3; i++ {
		if !s.Has(i) {
			t.Errorf("expected %d in set", i)
		}
	}

	for i := 4; i <= 6; i++ {
		if s.Has(i) {
			t.Errorf("unexpected %d in set", i)
		}
	}

	s.Del(2)
	if len(s) != 2 || s.Has(2) {
		t.Error("delete failed")
	}

	c := s.Slice()
	if len(c) != 2 {
		t.Errorf("wrong slice length: expected 2, got %d", len(c))
	}
	if !(equal(c, []int{1,3}) || equal(c, []int{3,1})) {
		t.Errorf("wrong elements in slice, expected []int{1,3}, got %#v", c)
	}
}
