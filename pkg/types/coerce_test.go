package types

import (
	"testing"
)

func TestCoerce(t *testing.T) {
	foo, err := CoerceScalar(int8(0), 100)
	if err != nil {
		t.Errorf("unexpected error coercing int(100) to int(8): %v", err)
	}
	if _, ok := foo.(int8); !ok {
		t.Errorf("unexpected return type %T coercing int(100) to int(8)", foo)
	}
	if foo != int8(100) {
		t.Errorf("unexpected mismatch; expected int8(100), got %d<%T>", foo, foo)
	}

	foo, err = CoerceScalar(int8(0), 500)
	if err == nil {
		t.Errorf("expected error missing coercing int(500) to int(8): %v", err)
	}
	if _, ok := foo.(int8); !ok {
		t.Errorf("unexpected return type %T coercing int(100) to int(8)", foo)
	}
	if foo != int8(-12) {
		t.Errorf("unexpected mismatch; expected int8(100), got %d<%T>", foo, foo)
	}
}
