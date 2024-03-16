// The `fflag` package provides GNU-style command-line argument
// parsing with the functional options pattern.
//
// A `Flag` is created and added to the default `FlagSet` with
// `Var()`. The minimal call to `Var()` provides: a pointer to a
// variable where the value of the flag is to be stored; a "label",
// which is the long-form command-line argument that is expected to be
// presented on the command line with two leading dashes,
// e.g. `--help`; and a brief description of the flag's purpose. For
// example:
//
//     fflag.Var(&value, "help", "prints a help message to stdout")
//
// The label must consist solely of Unicode letters, numbers, and
// hyphens. It must not begin with a hyphen.
//
// The first argument to `Var` must be a POINTER to one of:
//
//   1) a basic datatype (e.g. `int8`, `float32`, `string`)
//   2) a slice of basic datatype (e.g. `[]int8`, `[]string`)
//   3) a `struct` implementing the `Set()` interface
//
// A flag may have a single-character shortcut that is expected to be
// presented on the command-line with one leading dash, e.g. `-x`. It
// is introduced with the `WithShortcut(rune r)` option. For example:
//
//     f := fflag.Var(&value, "help", "prints a help message to stdout",
//         fflag.WithShortcut('?'))
//
// Punctuation (or other non-letter, non-number) characters are not
// normally allowed as shortcuts. The sole exception is '?' due to its
// widespread use as an alias for "help".
//
// Note that the above flag definition will not actually cause a help
// message of any kind to be printed: it is generally up to the
// programmer to specify the behavior after command-line parsing is
// complete.
//
// As a convenience, however, the `WithMessage()`, `WithUsage()` and
// `WithExit()` options are provided.
//
//     fflag.Var(&value, "help", "prints a help message to stdout",
//         fflag.WithShortcut('?'),
//         fflag.WithMessage("USAGE: myprog [OPTION]... [FILE]..."),
//         fflag.WithExit(0)
//     )
//
// The difference between `WithMessage()` and `WithUsage()` options is
// that `WithUsage()` prints a flag summary after the given message
// while `WithMessage()` does not. `WithExit()` causes the program to
// exit immediately with the given status code.
//
// The simplest ordinary flag is a nullary boolean switch that takes
// no parameter.
//
//     bool value
//     fflag.Var(&value, "easy", "use easy mode"))
//
// In this case, `value` will default to `false` (the zero value) and
// become `true` if the command-line argument appears in either long
// or shortcut (letter) form (if the `WithShortcut()` option is used)
// or if any aliases are used.
//
//     bool ignoreCase
//     f := NewFlag(&ignoreCase, "ignore-case", "ignore case in patterns",
//         WithShortcut('i'), WithShortcutAlias('y', true))
//
// The second (boolean) argument to `WithShortcutAlias()` says that
// the `-y` alias is obsolete. If marked obsolete, a deprecation
// warning will be printed to `stderr` if it is used and its
// obsolesence will be noted in generated flag summaries.
//
// If a default is supplied, a boolean value will be toggled if the
// flag is given.
//
//     var hard bool
//     fflag.Var(&hard, "easy", "use easy mode", fflag.WithDefault(true))
//
// In this case, `hard` will default to `true` and become false if
// `--easy` appears on the command line.
//
// It is an error to repeat a command-line argument unless the first
// argument to NewFlag is a pointer to a slice:
//
//     values := []bool{}
//     f := NewFlag(&values, "example", WithShortcut('x'))
//
// In this case, successive values of the flag are appended to
// `values` in the order in which they are processed.
//
// The sole exception to this rule is where a callback function is
// supplied:
//
//     var value string
//     f := NewFlag(&value, "file", WithCallback(MyFunc))
//
// In this case, the callback function is called with the given
// pointer (interface), label, parameter, and position on the
// command-line. Thus a program, `prog`, with the above "file" flag,
// could be invoked as follows:
//
//     prog --file foo.txt --file bar.txt
//
// Here, `MyFunc` would be called twice as:
//
//     MyFunc(&value, "file", "foo.txt", 1)
//     MyFunc(&value, "file", "bar.txt", 3)
//
// Notably, `value` is NOT set by `fflag` if a callback is supplied. A
// more usual setup would be:
//
//     var files []string
//     f := NewFlag(&files, "file", WithUsage("files to process"))
//
// After parsing, `files` would have contents equivalent to:
//
//     files := []string{"foo.txt", "bar.txt"}
//
// For unary (non-boolean) flags, a default can be supplied:
//
//     var files []string
//     fflag.Var(&files, "file", WithDefault([]string{"/dev/null"})
//
// The value will be set to the default if the argument is not given.
//
// If the default is not a slice, but the value is (a pointer to)
// slice, the default will be the first element of the slice.
//
// If the value is not a (pointer to a) slice, but the default value
// is a slice, the value is constrained to the values in the default,
// like a kind of enum.
//
// Consider the `--directories` option of GNU `grep`. It can take one
// of 3 values --- `read`, `skip`, and `recurse` --- with the default
// being `read`:
//
//     var string dir
//     f := NewFlag(&dir, "directories", WithDefault([]string{
//         "read", "skip", "recurse"}))
//
// The actual default is the first value in the slice.

package fflag

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/EmmetCaulfield/fflag/pkg/types"
)

type FlagError struct {
	s string
}

func (fe *FlagError) Error() string {
	return fe.s
}

type CallbackFunction func(value interface{}, label string, arg string, pos int)

type FlagType int8

const (
	ClearFlagType  FlagType = 0b00000000
	LabelAliasBit           = 0b00000001
	LetterAliasBit          = 0b00000010
	ObsoleteBit             = 0b00000100
	HiddenBit               = 0b00001000
	ChangedBit              = 0b00010000
	CounterBit              = 0b00100000
	RepeatableBit           = 0b01000000
)

func (ft *FlagType) TstLabelAliasBit() bool  { return *ft|LabelAliasBit != 0 }
func (ft *FlagType) TstLetterAliasBit() bool { return *ft|LetterAliasBit != 0 }
func (ft *FlagType) TstObsoleteBit() bool    { return *ft|ObsoleteBit != 0 }
func (ft *FlagType) TstHiddenBit() bool      { return *ft|HiddenBit != 0 }
func (ft *FlagType) TstChangedBit() bool     { return *ft|ChangedBit != 0 }
func (ft *FlagType) TstCounterBit() bool     { return *ft|CounterBit != 0 }
func (ft *FlagType) TstRepeatableBit() bool  { return *ft|RepeatableBit != 0 }
func (ft *FlagType) TstAliasBits() bool      { return (*ft|LetterAliasBit)|(*ft|LabelAliasBit) != 0 }

func (ft *FlagType) ClrLabelAliasBit()  { *ft = *ft & ^LabelAliasBit }
func (ft *FlagType) ClrLetterAliasBit() { *ft = *ft & ^LetterAliasBit }
func (ft *FlagType) ClrObsoleteBit()    { *ft = *ft & ^ObsoleteBit }
func (ft *FlagType) ClrHiddenBit()      { *ft = *ft & ^HiddenBit }
func (ft *FlagType) ClrChangedBit()     { *ft = *ft & ^ChangedBit }
func (ft *FlagType) ClrCounterBit()     { *ft = *ft & ^CounterBit }
func (ft *FlagType) ClrRepeatableBit()  { *ft = *ft & ^RepeatableBit }

func (ft *FlagType) SetLabelAliasBit()  { *ft = *ft | LabelAliasBit }
func (ft *FlagType) SetLetterAliasBit() { *ft = *ft | LetterAliasBit }
func (ft *FlagType) SetObsoleteBit()    { *ft = *ft | ObsoleteBit }
func (ft *FlagType) SetHiddenBit()      { *ft = *ft | HiddenBit }
func (ft *FlagType) SetChangedBit()     { *ft = *ft | ChangedBit }
func (ft *FlagType) SetCounterBit()     { *ft = *ft | CounterBit }
func (ft *FlagType) SetRepeatableBit()  { *ft = *ft | RepeatableBit }

type Flag struct {
	Value         interface{}
	Label         string
	Letter        rune
	Type          FlagType
	ValueTypeTag  string
	Default       interface{}
	AliasFor      *Flag
	FileFlag      *Flag
	Usage         string
	Callback      CallbackFunction
	ParentFlagSet *FlagSet
}

func (f *Flag) Set(value interface{}) error {
	// Prefer the SetValue interface if present:
	if setter, ok := f.Value.(types.SetValue); ok {
		if str, ok := value.(string); ok {
			return setter.Set(str)
		}
		f.Failf("Cannot pass non-string to SetValue.Set(string) in flag.Set() for flag '%s'", f.Label)
		return &FlagError{"failed to pass non-string to SetValue.Set()"}
	}

	if value == nil {
		var boolp *bool
		var ok, def bool
		if boolp, ok = f.Value.(*bool); !ok {
			f.Failf("flag.Set(nil) called for non-boolean flag '%s' of type %T", f.Label, f.Value)
			return &FlagError{"cannot set nil value for non-bool"}
		}
		// If a default was given, use it, otherwise the zero
		// value (`false`) returned by the type assertion is the
		// default we want in the absence of a stipulated default
		def, _ = f.GetDefault().(bool)
		*boolp = !def
		return nil
	}

	if str, ok := value.(string); ok {
		err := types.FromStr(f.Value, str)
		if err != nil {
			f.Failf("failed to convert '%s' to %T: %v", str, f.Value, err)
		}
		return err
	}

	f.Failf("type %T not handled in flag.Set() for flag '%s'", value, f.Label)
	return &FlagError{"type not handled in flag.Set()"}
}

func (f *Flag) GetDefaultLen() int {
	return types.SliceLen(f.Default)
}

func (f *Flag) GetDefault() interface{} {
	if f.GetDefaultLen() > 0 {
		return types.ItemAt(f.Default, 0)
	}
	return f.Default
}

func (f *Flag) GetDefaultDescription() string {
	buf := &bytes.Buffer{}
	if f.Default == nil {
		return ""
	}
	enum, ok := f.Default.([]interface{})
	if ok {
		if len(enum) > 0 {
			buf.WriteString(types.StrConv(enum[0]) + "(default)")
		}
		if len(enum) > 1 {
			buf.WriteString("|" + types.StrConv(enum[1:len(enum)-1], types.WithSep("|")))
		}
	} else {
		buf.WriteString(types.StrConv(f.Default) + "(default)")
	}
	return buf.String()
}

func (f *Flag) GetTypeTag() string {
	if len(f.ValueTypeTag) > 0 {
		return f.ValueTypeTag
	}
	if f.GetDefaultLen() > 1 {
		return "ENUM"
	}
	if types.IsInt(f.Value) {
		return "INT"
	}
	if types.IsUint(f.Value) {
		return "NUM"
	}
	if types.IsFloat(f.Value) {
		return "FLT"
	}
	if types.IsString(f.Value) {
		return "STR"
	}
	return ""
}

func (f *Flag) FlagString() string {
	buf := &bytes.Buffer{}
	if f.Letter == rune(0) {
		buf.WriteString("    ")
	} else {
		buf.WriteRune('-')
		buf.WriteRune(f.Letter)
		if len(f.Label) > 1 {
			buf.WriteString(", ")
		}
	}
	if len(f.Label) > 1 {
		fmt.Fprintf(buf, "--%s", f.Label)
	}
	tag := f.GetTypeTag()
	if len(tag) > 0 {
		fmt.Fprintf(buf, "=%s", tag)
	}
	return buf.String()
}

type FlagOption = func(f *Flag)

func WithShortcut(letter rune) FlagOption {
	return func(f *Flag) {
		f.Letter = letter
	}
}

func WithAlias(label string, letter rune, obsolete bool) FlagOption {
	return func(f *Flag) {
		flag := f.ParentFlagSet.LookupLabel(label)
		if flag != nil {
			f.Failf("long flag already exists for alias '%s'", label)
			panic("alias cannot be created")
		}
		flag = f.ParentFlagSet.LookupShortcut(letter)
		if flag != nil {
			f.Failf("short flag already exists for alias '%c'", letter)
			panic("alias cannot be created")
		}
		flag = NewAlias(label, letter, f)
		f.ParentFlagSet.AddFlag(f)
	}
}

func WithDefault(def interface{}) FlagOption {
	return func(f *Flag) {
		if types.IsPointerTo(f.Value, def) {
			f.Default = def
		} else {
			f.Default = nil
		}
	}
}

func AsCounter() FlagOption {
	return func(f *Flag) {
		if types.IsNum(f.Value) {
			f.Type.SetCounterBit()
		} else {
			f.Failf("cannot use non-numeric flag '%s' as counter", f.Label)
		}
	}
}

func WithTypeTag(tag string) FlagOption {
	return func(f *Flag) {
		f.ValueTypeTag = tag
	}
}

func WithCallback(callback CallbackFunction) FlagOption {
	return func(f *Flag) {
		f.Callback = callback
	}
}

func NewFlag(value interface{}, label string, usage string, opts ...FlagOption) *Flag {
	if !types.IsPointer(value) {
		return nil
	}
	if !IsValidLabel(label) {
		return nil
	}
	if types.IsOtherT(value) {
		if _, ok := value.(types.SetValue); !ok {
			return nil
		}
	}
	f := &Flag{
		Value:  value,
		Label:  label,
		Letter: rune(0),
		Usage:  usage,
	}
	if types.IsSlice(value) {
		f.Type.SetRepeatableBit()
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func NewAlias(label string, letter rune, target *Flag) *Flag {
	flagType := ClearFlagType
	if len(label) > 1 && IsValidLabel(label) {
		flagType.SetLabelAliasBit()
	}
	if letter != rune(0) && IsValidShortcut(letter) {
		flagType.SetLetterAliasBit()
	}
	if !flagType.TstAliasBits() {
		return nil
	}

	return &Flag{
		Value:    nil, // stored in target
		Label:    label,
		Letter:   letter,
		AliasFor: target,
		Type:     flagType,
	}
}

func (f *Flag) IsLabelAlias() bool {
	return f.Type.TstLabelAliasBit()
}
func (f *Flag) IsShortcutAlias() bool {
	return f.Type.TstLetterAliasBit()
}
func (f *Flag) IsAlias() bool {
	return f.Type.TstAliasBits()
}
func (f *Flag) IsHidden() bool {
	return f.Type.TstHiddenBit()
}
func (f *Flag) IsChanged() bool {
	return f.Type.TstChangedBit()
}
func (f *Flag) IsCounter() bool {
	return f.Type.TstCounterBit()
}
func (f *Flag) IsRepeatable() bool {
	return f.Type.TstRepeatableBit()
}

// Only allow letters and numbers as shortcut letters
func IsValidShortcut(r rune) bool {
	return r == '?' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

// Only allow letters, numbers, and underscore in labels
func IsValidLabel(label string) bool {
	// A label must be longer than one letter:
	if len(label) < 2 {
		return false
	}
	// A label can't begin with a hyphen
	if label[0] == '-' {
		return false
	}
	// Labels must otherwise consist entirely of letters, numbers, and
	// hyphens
	for _, r := range label {
		if r == '-' || unicode.IsLetter(r) || unicode.IsNumber(r) {
			continue
		}
		return false
	}
	return true
}

func (f *Flag) Failf(format string, args ...interface{}) {
	if f.ParentFlagSet == nil {
		panic("parent flagset not set")
	}
	f.ParentFlagSet.Failf(format, args...)
}
