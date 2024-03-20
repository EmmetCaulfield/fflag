package fflag

import (
	"reflect"
	"testing"
)

func setup() {
	CommandLine.OnFail.SetContinueBit()
	// CommandLine.OnFail.SetSilentBit()
}

func TestBasicSet(t *testing.T) {
	setup()
	// Bool
	var b bool
	f := NewFlag(b, "foo", "a non-pointer (bad)")
	if f != nil {
		t.Error("unexpected success creating new flag from non-pointer")
	}
	f = NewFlag(&b, "foo", "a boolean flag")
	if f == nil {
		t.Error("failed create boolean flag")
	}
	err := f.Set(true)
	if err != nil || b != true {
		t.Error("failed to set basic boolean from bool constant")
	}
	b = false
	err = f.Set("true")
	if err != nil || b != true {
		t.Error(`failed to set basic boolean from string "true"`)
	}
	b = false
	err = f.Set("1")
	if err != nil || b != true {
		t.Error(`failed to set basic boolean from string "1"`)
	}
	b = false
	err = f.Set(1)
	if err != nil || b != true {
		t.Error(`failed to set basic boolean from int 1`)
	}

	f = NewFlag(&b, "foo", "a boolean flag", WithDefault(true))
	if f == nil {
		t.Error("failed to create boolean flag with default")
	}
	err = f.Set(nil)
	if err != nil || b != false {
		t.Error(`failed to toggle basic boolean with nil`)
	}

	var i8 int8 = 11
	f = NewFlag(&i8, "foo", "an 8-bit int flag", WithDefault(25))
	if f == nil || i8 != 25 {
		t.Errorf("failed to create int8 flag with default (%d != 25)", i8)
	}
	err = f.Set(100)
	if err != nil || i8 != 100 {
		t.Error("failed to set int8 from int 100")
	}
	err = f.Set(128)
	if err == nil {
		t.Error("unexpected success with value out of range")
	}
	err = f.Set(-129)
	if err == nil {
		t.Error("unexpected success with value out of range")
	}
	err = f.Set("50")
	if err != nil || i8 != 50 {
		t.Error(`failed to set int8 with string "50"`)
	}

	var u8 uint8
	f = NewFlag(&u8, "foo", "an 8-bit unsigned int flag")
	err = f.Set(100)
	if err != nil || u8 != 100 {
		t.Error("failed to set basic uint8")
	}
	err = f.Set(256)
	if err == nil {
		t.Error("unexpected success with value out of range")
	}
	err = f.Set(-1)
	if err == nil {
		t.Error("unexpected success with value out of range")
	}
}

func TestVectorSet(t *testing.T) {
	setup()
	b := []bool{true, false, true}
	f := NewFlag(b, "foo", "a non-pointer (bad)")
	if f != nil {
		t.Error("unexpected success creating new flag from non-pointer")
	}
	f = NewFlag(&b, "foo", "a boolean slice")
	if f != nil {
		t.Error("unexpected success creating new flag from non-empty slice")
	}
	a := []bool{}
	f = NewFlag(&a, "foo", "a boolean slice flag")
	if f == nil {
		t.Error("error new flag from empty slice")
	}
	err := f.Set(true)
	if err != nil {
		t.Error("error setting initial value on slice")
	}
	err = f.Set(false)
	if err != nil {
		t.Error("error setting 2nd value on slice")
	}
	err = f.Set(true)
	if err != nil {
		t.Error("error setting 3rd value on slice")
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("value mismatch: expected %v, got %v", b, a)
	}
}
