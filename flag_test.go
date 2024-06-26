package fflag

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	CommandLine.OnFail.SetContinueBit()
	CommandLine.OnFail.SetSilentBit()
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
		short  rune
		long   string
		id     string
		expect bool
	}{
		{'e', "example", "e" + IdSep + "example", true},
		{'e', "", "e" + IdSep, true},
		{NoShort, "example", string(NoShort) + IdSep + "example", true},
		{NoShort, "", string(NoShort) + IdSep, true},
		{ErrRuneEmptyStr, "", "", true},
	}
	for i, test := range table {
		id := ID(test.short, test.long)
		result := (id == test.id)
		if result != test.expect {
			t.Errorf("test ID() roundtrip %d failed, expected '{%d}/%s', got '%s'", i, test.short, runesToAscii(test.long), runesToAscii(id))
		}

		unletter, unlabel := UnID(test.id)
		result = (unletter == test.short && unlabel == test.long)
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
		t.Error("failed to create a boolean flag")
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

func TestWithDefault(t *testing.T) {
	b := false
	// Default should be the same type as the value or a string
	assert.Panics(t, func() { Var(&b, 'b', "", "should panic", WithDefault(100)) })
	assert.Panics(t, func() { Var(&b, 'b', "", "should panic", WithDefault("foo")) })
	Var(&b, 'b', "", "should NOT panic", WithDefault("true"))
	assert.Equal(t, true, b)
}

func TestHyphenNumIdiomVar(t *testing.T) {
	b := false
	var s string
	var u uint
	assert.Panics(t, func() { Var(&b, NoShort, NoLong, "non-number, should panic") })
	Var(&u, NoShort, NoLong, "unsigned, should not panic")
	assert.Panics(t, func() { Var(&u, NoShort, NoLong, "not twice") })
	assert.Panics(t, func() { Var(&s, '7', "seven", "no numeric shorts with -NUM defined") })

	fs := NewFlagSet()
	fs.Var(&s, '7', "seven", "numeric shorts allowed")
	assert.Panics(t, func() { fs.Var(&u, NoShort, NoLong, "now ambiguous") })
}

func TestLongGet(t *testing.T) {
	fs := NewFlagSet()
	var a, b, c bool
	fs.Var(&a, NoShort, "ant", "six legs")
	fs.Var(&b, 'b', "bat", "two legs, two wings")
	fs.Var(&c, 'k', "cat", "four legs")

	f := fs.Lookup('a') // Should fail
	if f != nil && f.Value == &a {
		t.Error("unexpected success looking up rune('a')")
	}
	f = fs.Lookup("a") // Should succeed
	if f == nil || f.Value != &a {
		t.Error("error looking up string(\"a\")")
	}

	f = fs.Lookup('b') // Should succeed
	if f == nil || f.Value != &b {
		t.Error("error looking up rune('b')")
	}
	f = fs.Lookup("b") // Should succeed
	if f == nil || f.Value != &b {
		t.Error("error looking up string(\"b\")")
	}
	f = fs.Lookup("ba") // Should succeed
	if f == nil || f.Value != &b {
		t.Error("error looking up string(\"ba\")")
	}
	f = fs.Lookup("bat") // Should succeed
	if f == nil || f.Value != &b {
		t.Error("error looking up string(\"bat\")")
	}
	f = fs.Lookup("batx") // Should fail
	if f != nil {
		t.Error("unexpected success looking up string(\"batx\")")
	}

	f = fs.Lookup('c') // Should fail
	if f != nil {
		t.Error("unexpected success looking up rune('c')")
	}
	f = fs.Lookup('k') // Should succeed
	if f == nil || f.Value != &c {
		t.Error("failed looking up rune('k')")
	}
	f = fs.Lookup("k") // Should succeed
	if f == nil || f.Value != &c {
		t.Error("failed looking up string(\"k\")")
	}
	f = fs.Lookup("c") // Should succeed
	if f == nil || f.Value != &c {
		t.Error("error looking up string(\"ca\")")
	}
	// Priority to shorts in one-rune longs
	var d bool
	fs.Var(&d, 'c', "cow", "four legs, moos")
	f = fs.Lookup("c") // Should succeed, but find `d`
	if f == nil || f.Value != &d {
		t.Error("error looking up string(\"c\")")
	}
}
