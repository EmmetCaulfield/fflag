package deque

import (
	"testing"
)

func TestEqual(t *testing.T) {
	testCases := []struct{
		a Deque[int]
		b Deque[int]
	}{
		{Deque[int]{}, Deque[int]{}},
		{Deque[int]{1,2,3}, Deque[int]{1,2,3}},
	}
	for _, test := range testCases {
		if !EqualV(test.a, test.b) {
			t.Errorf("Equality test failed %v != %v", test.a, test.b)
		}
	}
}

func TestNotEqual(t *testing.T) {
	testCases := []struct{
		a Deque[int]
		b Deque[int]
	}{
		{Deque[int]{1,2,3}, Deque[int]{1,2}},
		{Deque[int]{3,2,1}, Deque[int]{1,2,3}},
	}
	for _, test := range testCases {
		if EqualP(&test.a, &test.b) {
			t.Errorf("Inequality test failed %v != %v", test.a, test.b)
		}
	}
}

func TestPushPop(t *testing.T) {
	a := &Deque[int]{1,2}
	b := &Deque[int]{1,2,3}
	v, err := b.Pop()
	if err != nil {
		t.Errorf("Pop() failed: %v %v", a, b)
	}
	if !a.Equal(b) {
		t.Errorf("Equality test failed after Pop() %v != %v", a, b)
	}
	a.Push(v)
	b.Push(3)
	if !a.Equal(b) {
		t.Errorf("Equality test failed after Push() %v != %v", a, b)
	}
}

func TestShiftUnshift(t *testing.T) {
	a := &Deque[int]{2,3}
	b := &Deque[int]{1,2,3}
	v, err := b.Shift()
	if err != nil {
		t.Errorf("Shift() failed: %v %v", a, b)
	}
	if !a.Equal(b) {
		t.Errorf("Equality test failed after Shift() %v != %v", a, b)
	}
	a.Unshift(v)
	b.Unshift(1)
	if !a.Equal(b) {
		t.Errorf("Equality test failed after Unshift() %v != %v", a, b)
	}
}

func TestAppend(t *testing.T) {
	a := &Deque[int]{1,2,3}
	b := &Deque[int]{}
	b.Append(1, 2, 3)
	if !a.Equal(b) {
		t.Errorf("Equality test failed after Append() %v != %v", a, b)
	}
}

func TestPrepend(t *testing.T) {
	a := &Deque[int]{1,2,3}
	b := &Deque[int]{}
	b.Prepend(1, 2, 3)
	if !a.Equal(b) {
		t.Errorf("Equality test failed after Prepend() %v != %v", a, b)
	}
}

func TestPeekFront(t *testing.T) {
	a := &Deque[int]{1,2,3}
	b := &(*a)
	if !a.Equal(b) {
		t.Errorf("Equality test failed: %v != %v", a, b)
	}
	x, err := a.Peek()
	if err != nil || x != 3 {
		t.Errorf("Peek() failed: %v, %v, %v", a, b, err)
	}
	x, err = a.Front()
	if err != nil || x != 1 {
		t.Errorf("Front() failed: %v, %v, %v", a, b, err)
	}
}

func TestError(t *testing.T) {
	a := &Deque[int]{}
	b, err := a.Pop()
	if err == nil {
		t.Errorf("Error test failed after Pop(): %v, %v, %v", a, b, err)
	}
	b, err = a.Shift()
	if err == nil {
		t.Errorf("Error test failed after Shift(): %v, %v, %v", a, b, err)
	}
	b, err = a.Peek()
	if err == nil {
		t.Errorf("Error test failed after Peek(): %v, %v, %v", a, b, err)
	}
	b, err = a.Front()
	if err == nil {
		t.Errorf("Error test failed after Front(): %v, %v, %v", a, b, err)
	}
}
