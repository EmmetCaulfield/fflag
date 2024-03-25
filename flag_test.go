package fflag

import (
	"fmt"
	"reflect"
	"testing"
)

func setup() {
	CommandLine.OnFail.SetContinueBit()
	// CommandLine.OnFail.SetSilentBit()
}

func runesToAscii(src string) string {
	dst := ""
	for _, r := range src {
		if r > 31 && r < 127 {
			dst = dst + string(r)
		} else {
			dst = dst + fmt.Sprintf("{%d}", r)
		}
	}
	return dst
}

func TestIdRoundtrip(t *testing.T) {
	table := []struct {
		letter rune
		label  string
		id     string
		expect bool
	}{
		{'e', "example", "e/example", true},
		{'e', "", "e/", true},
		{0, "example", "\u0000/example", true},
		{0, "", "\u0000/", true},
		{ErrRuneEmptyStr, "", "", true},
	}
	for i, test := range table {
		id := ID(test.letter, test.label)
		result := (id == test.id)
		if result != test.expect {
			t.Errorf("test ID() roundtrip %d failed, expected '{%d}/%s', got '%s'", i, test.letter, runesToAscii(test.label), runesToAscii(id))
		}

		unletter, unlabel := UnID(test.id)
		result = (unletter == test.letter && unlabel == test.label)
		if result != test.expect {
			t.Errorf("test UnID() roundtrip %d failed, expected '%s', got '{%d}/%s'", i, runesToAscii(test.id), unletter, unlabel)
		}
	}
}

func TestBasicSet(t *testing.T) {
	setup()
	// Bool
	var b bool
	f := NewFlag(b, 0, "foo", "a non-pointer (bad)")
	if f != nil {
		t.Error("unexpected success creating new flag from non-pointer")
	}
	f = NewFlag(&b, 0, "foo", "a boolean flag")
	if f == nil {
		t.Error("failed create boolean flag")
	}
	err := f.Set(true, 0)
	if err != nil || b != true {
		t.Error("failed to set basic boolean from bool constant")
	}
	err = f.Set(true, 0)
	if err == nil {
		t.Error("unexpected success on repeat Set() of flag not marked repeatable")
	}
	f.Type.SetRepeatsBit()
	b = false
	err = f.Set("true", 0)
	if err != nil || b != true {
		t.Error(`failed to set basic boolean from string "true"`)
	}
	b = false
	err = f.Set("1", 0)
	if err != nil || b != true {
		t.Error(`failed to set basic boolean from string "1"`)
	}
	b = false
	err = f.Set(1, 0)
	if err != nil || b != true {
		t.Error(`failed to set basic boolean from int 1`)
	}
	if f.Count != 5 {
		t.Errorf(`wrong repeat count; expected 5, got %d`, f.Count)
	}

	f = NewFlag(&b, 0, "foo", "a boolean flag", WithDefault(true))
	if f == nil {
		t.Error("failed to create boolean flag with default")
	}
	err = f.Set(nil, 0)
	if err != nil || b != false {
		t.Error(`failed to toggle basic boolean with nil`)
	}
	if f.Count != 1 {
		t.Errorf(`wrong repeat count; expected 1, got %d`, f.Count)
	}

	var i8 int8 = 11
	f = NewFlag(&i8, 0, "foo", "an 8-bit int flag", WithDefault(int8(25)), WithRepeats(false))
	if f == nil || i8 != 25 {
		t.Errorf("failed to create int8 flag with default (%d != 25)", i8)
	}
	err = f.Set(100, 0)
	if err != nil || i8 != 100 {
		t.Error("failed to set int8 from int 100")
	}
	err = f.Set(128, 0)
	if err == nil {
		t.Error("unexpected success with value out of range")
	}
	err = f.Set(-129, 0)
	if err == nil {
		t.Error("unexpected success with value out of range")
	}
	err = f.Set("50", 0)
	if err != nil || i8 != 50 {
		t.Error(`failed to set int8 with string "50"`)
	}
	if f.Count != 4 {
		t.Errorf(`wrong repeat count; expected 4, got %d`, f.Count)
	}

	def := []int8{5, 10, 15, 20}
	f = NewFlag(&i8, 0, "foo", "an 8-bit int flag", WithDefault(def), WithRepeats(false))
	if f == nil || i8 != 5 {
		t.Errorf("failed to create int8 flag with slice default (%d != 5)", i8)
	}
	err = f.Set(11, 0)
	if err == nil || i8 == 11 {
		t.Errorf("unexpected success setting int8 flag with slice default (%d !in %v)", i8, def)
	}
	err = f.Set("15", 0)
	if err != nil || i8 != 15 {
		t.Errorf("failed to create int8 flag with value from default slice (%d != 15)", i8)
	}
	err = f.Set(300, 0)
	if err == nil {
		t.Errorf("unexpected success setting out-of-range value %v: %v", i8, err)
	}
	if i8 != 15 {
		t.Errorf("unexpected value change %d<%T> != %d", i8, i8, 15)
	}
	t.Logf("%d", int8(0b00101100))

	var u8 uint8
	f = NewFlag(&u8, 0, "foo", "an 8-bit unsigned int flag", WithRepeats(true))
	err = f.Set(100, 0)
	if err != nil || u8 != 100 {
		t.Error("failed to set basic uint8")
	}
	err = f.Set(50, 0)
	if err != nil || u8 != 100 {
		t.Error("repeat set not ignored")
	}
	if f.Count != 2 {
		t.Errorf(`wrong repeat count; expected 3, got %d`, f.Count)
	}

	var u16 uint16
	f = NewFlag(&u16, 0, "foo", "a 16-bit counter", AsCounter())
	err = f.Set(100, 0)
	if err != nil || u16 != 1 {
		t.Errorf("failed to set basic counter; expected 1, got %d", u16)
	}
	// The value argument of a counter set is ignored
	f.Set(nil, 0)
	f.Set("something", 0)
	f.Set(-100000000, 0)
	f.Set(3.14159, 0)
	f.Set(nil, 0)

	if f.Count != 6 {
		t.Errorf(`wrong repeat count; expected 6, got %d`, f.Count)
	}
}

func TestVectorSet(t *testing.T) {
	setup()
	b := []bool{true, false, true}
	f := NewFlag(b, 0, "foo", "a non-pointer (bad)")
	if f != nil {
		t.Error("unexpected success creating new flag from non-pointer")
	}
	f = NewFlag(&b, 0, "foo", "a boolean slice")
	if f != nil {
		t.Error("unexpected success creating new flag from non-empty slice")
	}
	a := []bool{}
	f = NewFlag(&a, 0, "foo", "a boolean slice flag")
	if f == nil {
		t.Error("error new flag from empty slice")
	}
	err := f.Set(true, 0)
	if err != nil {
		t.Error("error setting initial value on slice")
	}
	err = f.Set(false, 0)
	if err != nil {
		t.Error("error setting 2nd value on slice")
	}
	err = f.Set(true, 0)
	if err != nil {
		t.Error("error setting 3rd value on slice")
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("value mismatch: expected %v, got %v", b, a)
	}

	a = a[:0]
	err = f.Set([]string{"true", "false", "true"}, 0)
	if err != nil {
		t.Error("error setting bool slice from string slice")
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("value mismatch: expected %v, got %v", b, a)
	}

	i8a := []int8{3, 1, 4, 1, 5, 9, 2}
	i8v := []int8{}
	f = NewFlag(&i8v, 'l', "list", "a list of digits", WithListSeparator('|'))
	if f == nil {
		t.Error("error creating int8 slice flag")
	}
	f.Set("3|1|4|1|5|9|2", 0)
	if !reflect.DeepEqual(i8a, i8v) {
		t.Errorf("value mismatch: expected %v, got %v", i8a, i8v)
	}
}
