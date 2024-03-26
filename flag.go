// The `fflag` package provides POSIX (short) and GNU (long)
// command-line argument parsing with, for the programmer, the
// functional options pattern.
//
// It is somewhat inspired by the `pflag` package in some respects,
// but very significantly different in others. The most significant
// difference is that there is only one `Var()` function: the type of
// the flag is determined by the type of the first argument (rather
// than the function name), which MUST be a pointer to a basic type, a
// slice of basic type, or a struct implementing the `SetValue`
// interface (inspired by `pflag`).
//
// The other significant difference is the order of the short flag and
// long flag in the `Var()` argument list, with the short flag coming
// first as a `rune`, which must be a single UTF-8 letter, most often
// a single ASCII letter or number. The reason for this is that short
// flags are always listed first in manpages and other documentation,
// so it's actually a bit weird of `pflag` to have reversed this
// de-facto standard order and, in practice, I've found it handier to
// obey the standard order than stick to the `flag` argument order.
//
// If there is no short flag, the zero value (0, `\0', or
// `fflag.NoShort`) is used. The usual rules apply to long flags,
// which must consist of letters and numbers, except that the ASCII
// requirement has been relaxed. Any character satisfying
// unicode.IsLetter() or unicode.IsNumber() or the hyphen '-' are
// allowed. There is no attempt at (what `pflag` refers to as)
// normalization, a very dubious utility: just use the long flag you
// mean to use without weird capitalization.
//
// That said, `fflag` meets the GNU expectation that “users can
// abbreviate the option names as long as the abbreviations are
// unique”. This requires a ludicrous amount of extra effort, but it
// exists as a clearly expressed requirement so we implement it.
//
// `fflag` borrows the `Flag` and `FlagSet` names from `pflag`, adding
// `FlagGroup`. The purpose of a flag group is to enable usage
// information to be generated in a similar format to GNU/POSIX
// utilities like `grep`, with flags grouped in categories. This is an
// additional feature of `fflag` and isn't known to exist elsewhere.
//
// A `Flag` is created and added to the default `FlagGroup` in the
// default `FlagSet` (called `CommandLine` after `pflag`'s equivalent)
// with `Var()`. The minimal call to `Var()` provides: a pointer to a
// variable where the value of the flag is to be stored; the
// single-letter version of the flag as a rune (or 0 if none), e.g.,
// 'h'; the long version of the flag (or "" if none), e.g. `--help`;
// and a very brief description of the flag's purpose. For example:
//
//     fflag.Var(&value, 'h', "help", "prints a help message to stdout")
//
// The first argument to `Var` must be a POINTER to one of:
//
//   1) a basic datatype (e.g. `int8`, `float32`, `string`)
//   2) a slice of basic datatype (e.g. `[]int8`, `[]string`)
//   3) something implementing the `SetValue` interface
//
// Non-pointer arguments are rejected. If the argument implements the
// `SetValue` interface, `fflag` neither modifies the argument itself
// nor enforces any of its usual rules. If you pass something
// implementing this interface, it's assumed that you will take care
// of everything and don't want `fflag` to do anything other than pass
// along the message “this flag appeared with this argument”.
//
// A flag need not have a single-character shortcut. If there is no
// shortcut, a 0 is given for the shortcut argument:
//
//    fflag.Var(&value, 0, "help", "prints a help message to stdout")
//
// Punctuation (or other non-letter, non-number) characters are not
// normally allowed as shortcuts. The sole exception is '?' due to its
// widespread use as an alias for "help", but this is prohibited by
// POSIX, so if you want to use this, you have to enable it
// explicitly.
//
// Equally, a flag need not have a long version. If you wanted to have
// `-?` as a short flag with no long version, you would do:
//
//    fflag.PosixRejectQuest = false
//    fflag.Var(&value, '?', "", "prints a help message to stdout")
//
// There is a special case (and common idiom) where NEITHER a long NOR
// a short form is required: `-NUM` (as in `grep`, `head`, `tail`, and
// several other tools). These special cases are always an alias for
// something else and always refer to a natural number appearing after
// a single hyphen. For example `head`'s `-n/--lines` is best
// represented as:
//
//    int nlines
//    fflag.Var(&nlines, 'n', 'lines',
//        "print the first NUM lines instead of the first 10",
//        fflag.WithAlias(0, "", false), fflag.WithTypeTag("[-]NUM"))
//
// Obviously, this special case can only be used once, but it requires
// no special logic since it is always an error to attempt to create a
// flag (an alias is just a special flag) that shares the short or
// long form of an existing flag.
//
// The simplest ordinary flag is a nullary boolean switch that takes
// no parameter.
//
//     bool value
//     fflag.Var(&value, 'e', "easy", "use easy mode"))
//
// In this case, `value` will default to `false` (the zero value for
// `bool`s) and become `true` if the command-line argument appears in
// either form (long or short). By default, it is an error to repeat a
// scalar flag, but there are 3 options that make an exception:
//
//   * `WithRepeats(ignore bool)`
//   * `AsCounter()`
//   * `WithCallback(callback func(...))`
//
// `WithRepeats()` allows repeat appearances of a flag, `AsCounter()`
// causes the number of occurrences to be counted (if the value
// pointer is a number), and `WithCallback()` causes the given
// callback function to be called EVERY TIME the flag appears on the
// command-line. Much like a `value` (first) argument implementing
// `SetValue`, `fflag` washes its hands of any further involvement,
// and it becomes entirely up to the callback to modify the value
// appropriately, track/ignore repeated appearances, etc.
//
// Several utilities allow `-v/--verbose` to be repeated for
// increasing levels of verbosity.
//
//     int verbosity
//     f := NewFlag(&verbosity, 'v', "verbose", "increase verbosity",
//         AsCounter())
//
// Supplying more than one of these (pairwise redundant or
// contradictory) options would result in a `panic()` since this would
// be an obvious programming error, not something that could
// “accidentally” occur at runtime based on user input.
//
// An explicit default can be supplied with `WithDefault()`:
//
//     var hard bool
//     fflag.Var(&hard, 0, "easy", "use easy mode",
//         fflag.WithDefault(true))
//
// In this case, `hard` will default to `true` and become false if
// `--easy` appears on the command line. If repeats are allowed, the
// value will toggle between `true` and `false`, which is admittedly
// weird, but if you do stupid things, expect stupid results.
//
// Repeated appearances of a flag, while prohibited by default for
// scalar value arguments, are _not_ an error if the value argument is
// (a pointer to) a slice. In this case, successive invocations will
// result in successive values being appended to the slice.
//
//     values := []bool{}
//     NewFlag(&values, 'x', "example", "example flag")
//
// The sole exception to this rule is where a callback function is
// supplied. When a callback is supplied, the callback is responsible
// for EVERYTHING.
//
//     f := NewFlag(&value, 'f', "file", "supply a filename",
//         WithCallback(MyFunc))
//
// The callback function is called with the given pointer, `&value`
// (via the `interface{}` argument), short option, long option,
// command-line argument (as a `string`, if any), and position on the
// command-line. Consider a program `prog`, with the above "file"
// flag, invoked as follows:
//
//     prog -f foo.txt --file bar.txt
//
// Here, `MyFunc` would be called twice as:
//
//     MyFunc(&value, 'f', "file", "foo.txt", 1)
//     MyFunc(&value, 'f', "file", "bar.txt", 3)
//
// The `value` is NOT set by `fflag` if a callback is supplied.
//
// For unary (non-boolean) flags, a default can be supplied:
//
//     var file string
//     fflag.Var(&file, 'f', "file", "supply a filename",
//         WithDefault("/dev/null"))
//
// TODO(emmet): consider what "default" means as bit more.
//
// The value will be set to the default if the argument is NOT
// given. This is EXACTLY equivalent to:
//
//     file := "/dev/null"
//     fflag.Var(&file, 'f', "file", "supply a filename")
//
// Frankly, I'm not sure why anyone would want to use `WithDefault()`
// for this.
//
// However, if the value is a (pointer to a) scalar, but the default
// is a slice, the value is constrained to the values in the default,
// like an enum.
//
// Consider the `--directories` option of GNU `grep`. It can take one
// of 3 values --- `read`, `skip`, and `recurse` --- with the default
// being `read`:
//
//     var string diract
//     f := NewFlag(&diract, 'd', "directories",
//         "if an input file is a directory use ACTION to process it",
//         WithDefault([]string{"read", "skip", "recurse"}),
//         WithTypeTag("ACTION"))
//
// The actual default is the first value in the slice. The remaining
// values in the slice constrain the set of acceptable values. For
// some program, `prog`, with the above flag deefinition, the value of
// `diract` after `fflag.Parse()` would be exactly the same for:
//
//     $ prog
//     $ prog -d read
//     $ prog --directories=read
//
// The following would be fine:
//
//     $ prog -d skip
//     $ prog --directories recurse
//
// But the following would result in a runtime error because `foo` is
// not in the default slice:
//
//     $ prog -d foo
//
package fflag

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/EmmetCaulfield/fflag/pkg/types"
)

var DefaultListSeparator string = ","
var PosixRejectQuest bool = true
var PosixRejectW bool = true

type FlagError struct {
	s string
}

func (fe *FlagError) Error() string {
	return fe.s
}

type CallbackFunction func(value interface{}, letter rune, label string, arg string, pos int) error

type FlagType uint16

const (
	ClearFlagType     FlagType = 0b0000000000000000
	LabelAliasBit     FlagType = 0b0000000000000001
	LetterAliasBit    FlagType = 0b0000000000000010
	ObsoleteBit       FlagType = 0b0000000000000100
	NotImplementedBit FlagType = 0b0000000000001000
	HiddenBit         FlagType = 0b0000000000010000
	ChangedBit        FlagType = 0b0000000000100000
	CounterBit        FlagType = 0b0000000001000000
	RepeatsBit        FlagType = 0b0000000010000000
	IgnoreRepeatsBit  FlagType = 0b0000000100000000
)

func (ft *FlagType) TstLabelAliasBit() bool     { return *ft&LabelAliasBit != 0 }
func (ft *FlagType) TstLetterAliasBit() bool    { return *ft&LetterAliasBit != 0 }
func (ft *FlagType) TstObsoleteBit() bool       { return *ft&ObsoleteBit != 0 }
func (ft *FlagType) TstNotImplementedBit() bool { return *ft&NotImplementedBit != 0 }
func (ft *FlagType) TstHiddenBit() bool         { return *ft&HiddenBit != 0 }
func (ft *FlagType) TstChangedBit() bool        { return *ft&ChangedBit != 0 }
func (ft *FlagType) TstCounterBit() bool        { return *ft&CounterBit != 0 }
func (ft *FlagType) TstRepeatsBit() bool        { return *ft&RepeatsBit != 0 }
func (ft *FlagType) TstIgnoreRepeatsBit() bool  { return *ft&IgnoreRepeatsBit != 0 }
func (ft *FlagType) TstAliasBits() bool         { return (*ft&LetterAliasBit)|(*ft&LabelAliasBit) != 0 }

func (ft *FlagType) ClrLabelAliasBit()     { *ft = *ft & ^LabelAliasBit }
func (ft *FlagType) ClrLetterAliasBit()    { *ft = *ft & ^LetterAliasBit }
func (ft *FlagType) ClrObsoleteBit()       { *ft = *ft & ^ObsoleteBit }
func (ft *FlagType) ClrNotImplementedBit() { *ft = *ft & ^NotImplementedBit }
func (ft *FlagType) ClrHiddenBit()         { *ft = *ft & ^HiddenBit }
func (ft *FlagType) ClrChangedBit()        { *ft = *ft & ^ChangedBit }
func (ft *FlagType) ClrCounterBit()        { *ft = *ft & ^CounterBit }
func (ft *FlagType) ClrRepeatsBit()        { *ft = *ft & ^RepeatsBit }
func (ft *FlagType) ClrIgnoreRepeatsBit()  { *ft = *ft & ^IgnoreRepeatsBit }

func (ft *FlagType) SetLabelAliasBit()     { *ft = *ft | LabelAliasBit }
func (ft *FlagType) SetLetterAliasBit()    { *ft = *ft | LetterAliasBit }
func (ft *FlagType) SetObsoleteBit()       { *ft = *ft | ObsoleteBit }
func (ft *FlagType) SetNotImplementedBit() { *ft = *ft | NotImplementedBit }
func (ft *FlagType) SetHiddenBit()         { *ft = *ft | HiddenBit }
func (ft *FlagType) SetChangedBit()        { *ft = *ft | ChangedBit }
func (ft *FlagType) SetCounterBit()        { *ft = *ft | CounterBit }
func (ft *FlagType) SetRepeatsBit()        { *ft = *ft | RepeatsBit }
func (ft *FlagType) SetIgnoreRepeatsBit()  { *ft = *ft | IgnoreRepeatsBit }

type Flag struct {
	Value         interface{}
	Label         string
	Letter        rune
	Type          FlagType
	Count         int
	ValueTypeTag  string
	Default       interface{}
	AliasFor      *Flag
	FileFlag      *Flag
	Usage         string
	Callback      CallbackFunction
	ListSeparator string
	parentFlagSet *FlagSet
}

const IdSep string = "/"

// A non-numeric, non-alphabetic ASCII character (other than '?') used
// as a placeholder meaning "there is no short option" in a variety of
// contexts.
const NoShort rune = rune(0)

func ID(letter rune, label string) string {
	// We use this for the -NUM special case used by a few utilities
	// (e.g. head, tail), which has NEITHER a normal valid shortcut
	// nor a normal valid long flag
	if letter == NoShort && len(label) == 0 {
		return string(NoShort) + IdSep
	}
	// Otherwise, we require either the shortcut or the long flag to
	// be valid
	if IsValidShortcut(letter) || IsValidLabel(label) {
		return string(letter) + IdSep + label
	}
	// An empty ID string is always an error
	return ""
}

func emptyOrNoShort(s string) bool {
	if len(s) == 0 {
		return true
	}
	// The NoShort indicator is always a single byte, so we don't have
	// to extract the first rune, we can just treat it like a byte.
	if len(s) == 1 && rune(s[0]) == NoShort {
		return true
	}
	return false
}

func UnID(id string) (rune, string) {
	if id == "" {
		return ErrRuneEmptyStr, ""
	}
	parts := strings.Split(id, IdSep)
	if len(parts) != 2 {
		return ErrRuneIdSepBad, ""
	}
	if parts[1] == "" && emptyOrNoShort(parts[0]) {
		return NoShort, ""
	}

	letter, tail := FirstRune(parts[0])
	if letter < 0 || tail != "" {
		return ErrRuneShortBad, ""
	}

	if IsValidShortcut(letter) || IsValidLabel(parts[1]) {
		return letter, parts[1]
	}
	return ErrRuneIdPartsBad, ""
}

func (f *Flag) ParentFlagSet() *FlagSet {
	if f.parentFlagSet == nil {
		return CommandLine
	}
	return f.parentFlagSet
}

func (f *Flag) Set(value interface{}, argPos int) error {
	// Prefer the SetValue interface if present:
	if setter, ok := f.Value.(types.SetValue); ok {
		if str, ok := value.(string); ok {
			f.Count++
			return setter.Set(str)
		}
		f.Failf("Cannot pass non-string to SetValue.Set(string) in flag.Set() for flag '%s'", f.Label)
		return &FlagError{"failed to pass non-string to SetValue.Set()"}
	}

	if f.AliasFor != nil {
		f = f.AliasFor
	}
	if f.AliasFor != nil {
		panic("Double aliases are not permitted")
	}

	if f.HasCallback() {
		v, _ := value.(string)
		return f.Callback(f.Value, f.Letter, f.Label, v, argPos)
	}

	f.Count++
	if f.IsCounter() {
		// TODO(emmet): think about this. It might be useful to be
		// able to stick the count into a string and it would work
		// fine, but OTOH, do we really want to be putting counts into
		// strings just because we can? At the other extreme, should
		// we be insisting that count must be an int "just because"?
		//
		// if !types.IsNum(f.Value) {
		//     panic("non-numeric value cannot be a counter")
		// }
		str := types.StrConv(f.Count)
		err := types.FromStr(f.Value, str)
		if err != nil {
			f.Failf("failed to set counter '%s' from %d", f.String(), f.Count)
		}
		return err
	}

	if f.Count > 1 && !f.IsRepeatable() {
		f.Failf("flag '%s' is not repeatable", f.String())
		return &FlagError{"flag not repeatable"}
	}

	if f.Count > 1 && f.IgnoreRepeats() {
		return nil
	}

	if value == nil {
		if boolp, ok := f.Value.(*bool); ok {
			// If a default was given, use it, otherwise the zero
			// value (`false`) returned by the type assertion is the
			// default we want in the absence of a stipulated default
			def, _ := f.GetDefault().(bool)
			*boolp = !def
			return nil
		}
		value = f.GetDefault()
		if value == nil {
			f.Failf("flag.Set(nil) called for flag '%s' with no default", f)
			return &FlagError{"cannot set nil value for non-bool with no default"}
		}
	} else if !f.InDefaults(value) {
		// TODO(emmet): consider supporting constrained defaults in a
		// command-line list. It's reasonable to expect that each item
		// in a list optarg would be checked against the list in
		// f.Default (if any), but in reality, a non-scalar optarg
		// (e.g. `-x foo,bar,baz`) will fail here.
		f.Failf("value %v not found in defaults %v for '%s'", value, f.Default, f)
		return &FlagError{"value constrained by defaults"}
	}

	// TODO(emmet): look at doing this this other than by
	// round-tripping via a string; OTOH, the value will usually be a
	// string anyway.

	// Convert the value to a string if it's not already one
	var ok bool
	var str string
	if str, ok = value.(string); !ok {
		str = types.StrConv(value, types.WithSep(f.ListSeparator))
		if str == "" {
			f.Failf("failed to convert '%v' to a nonempty string in '%s'", value, f)
			return &FlagError{"cannot convert value to string"}
		}
	}

	// Set the value from the string version
	err := types.FromStr(f.Value, str, types.WithSep(f.ListSeparator))
	if err != nil {
		f.Failf("failed to convert '%s' to %T: %v", str, f.Value, err)
		return err
	}
	return nil
}

func (f *Flag) GetValue() string {
	if f.AliasFor != nil {
		f = f.AliasFor
	}
	return types.StrConv(f.Value)
}

func (f *Flag) GetDefaultLen() int {
	if f.AliasFor != nil {
		f = f.AliasFor
	}
	return types.SliceLen(f.Default)
}

func (f *Flag) GetDefault() interface{} {
	if f.AliasFor != nil {
		f = f.AliasFor
	}
	if f.AliasFor != nil {
		panic("double aliases are not allowed")
	}
	if types.IsSlice(f.Default) {
		if types.SliceLen(f.Default) > 0 {
			return types.ItemAt(f.Default, 0)
		}
		// Default is an empty slice. Odd.
		panic("default cannot be an empty slice")
	}
	return f.Default
}

// InDefaults() returns true if the argument is in the f.Default slice
// or if f.Default is not a slice, otherwise it returns false.
func (f *Flag) InDefaults(ix interface{}) bool {
	if f.AliasFor != nil {
		f = f.AliasFor
	}
	if f.AliasFor != nil {
		panic("double aliases are not allowed")
	}
	if !types.IsSlice(f.Default) {
		return true
	}
	for i := 0; i < types.SliceLen(f.Default); i++ {
		d := types.ItemAt(f.Default, i)
		v, err := types.CoerceScalar(d, ix)
		if err != nil {
			// TODO(emmet): think this through
			f.Failf("error coercing %T (arg) to %T (defaults): %v", ix, d, err)
			return false
		}
		// fmt.Fprintf(os.Stderr, "%+v<%T> ?= %+v<%T> (%t) %+v<%T>\n", d, d, v, v, d == v, ix, ix)
		if d == v {
			return true
		}
	}
	return false
}

func (f *Flag) GetDefaultDescription() string {
	// TODO(emmet): handle aliases
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
	if f.AliasFor != nil {
		f = f.AliasFor
	}
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

// Returns short and long version of flags in the format of GNU
// utilities i.e., where both long and short versions are defined,
// the short version first, followed by a comma, followed by long
// version, e.g. "-x, --example", otherwise just the version that's
// defined, e.g. "-x" or "--example"
func (f *Flag) String() string {
	if f.Letter != NoShort && len(f.Label) > 1 {
		return "-" + string(f.Letter) + ", --" + f.Label
	}
	if len(f.Label) > 1 {
		return "--" + f.Label
	}
	if f.Letter != NoShort {
		return "-" + string(f.Letter)
	}
	return ""
}

// Returns f.String() wrapped in extra stuff for help/usage output
func (f *Flag) FlagString() string {
	buf := &bytes.Buffer{}
	if f.Letter == NoShort {
		buf.WriteString(`    `)
	}
	buf.WriteString(f.String())

	tag := f.GetTypeTag()
	if len(tag) > 0 {
		buf.WriteRune('=')
		buf.WriteString(tag)
	}
	return buf.String()
}

func (f *Flag) DescString() string {
	if f.Type.TstAliasBits() && f.AliasFor != nil {
		buf := &bytes.Buffer{}
		if f.Type.TstObsoleteBit() {
			buf.WriteString("obsolete ")
		}
		buf.WriteString("synonym for ")
		buf.WriteString(f.AliasFor.String())
		return buf.String()
	}
	if f.Type.TstNotImplementedBit() {
		return "not implemented"
	}
	// TODO(emmet): handle non-aliases
	return f.Usage
}

// Provides a sort key for sorting flags in the conventional order
// based on the short and long versions of the flag.
//
// GNU manpages present flags (within a group) in lexicographic order,
// ignoring the distinction between long and short flags. That is, a
// flag `--bat` with no short will appear after `-a, --ant` and before
// `-c, --cat`. Case is ignored in the sort, but uppercase shorts are
// presented before lowercase shorts.
func (f *Flag) SortKey() string {
	if f.Letter == NoShort {
		return f.Label
	}
	// TODO(emmet): think about special case of no long and no short
	// used for -NUM
	return string(f.Letter) + f.Label
}

type FlagOption = func(f *Flag)

type AliasOption = func(f *Flag)

func WithParent(fs *FlagSet) FlagOption {
	return func(f *Flag) {
		if f.parentFlagSet != nil && f.parentFlagSet != fs {
			panic("parent flagset already set")
		}
		f.parentFlagSet = fs
	}
}

func WithValue(value string) AliasOption {
	return func(f *Flag) {
		f.Value = value
	}
}

func WithListSeparator(sep rune) FlagOption {
	return func(f *Flag) {
		if !types.IsSlice(f.Value) {
			panic("cannot set separator for non-list value")
		}
		f.ListSeparator = string(sep)
	}
}

func WithAlias(letter rune, label string, obsolete bool) FlagOption {
	return func(f *Flag) {
		var flag *Flag = nil
		flag = f.ParentFlagSet().LookupLabel(label)
		if flag != nil {
			f.Failf("long flag already exists for alias '%s'", label)
			panic("alias cannot be created")
		}
		flag = f.ParentFlagSet().LookupShortcut(letter)
		if flag != nil {
			f.Failf("short flag already exists for alias '%c'", letter)
			panic("alias cannot be created")
		}
		flag = f.NewAlias(letter, label)
		if flag == nil {
			f.Failf("error creating alias -%c/--%s for `%s`", letter, label, f)
			panic("alias cannot be created")
		}
		if obsolete {
			flag.Type.SetObsoleteBit()
		} else {
			flag.Type.ClrObsoleteBit()
		}
		err := f.ParentFlagSet().AddFlag(flag)
		if err != nil {
			f.Failf("Error adding alias: %v", err)
		}
	}
}

func WithDefault(def interface{}) FlagOption {
	return func(f *Flag) {
		defType := types.Type(def)
		// Always allow the default to be a string or a slice of
		// strings since the value-to-set will come from the
		// command-line and be a string anyway
		if !defType.TstStringBit() {
			valType := types.Type(f.Value)
			if valType.TstSetterBit() {
				// We shouldn't be here if f.Value implements the
				// SetValue interface, because that always takes a
				// string so the default should be a string.
				panic("non-string default for object implementing SetValue interface")
			}
			if !types.SameBaseType(valType, defType) {
				fmt.Fprintf(os.Stderr, "%016b != %016b for '%#v'", valType, defType, f)
				panic("non-string default has different base type")
			}
		}
		// Set the default value
		f.Default = def
		def = f.GetDefault()
		err := types.FromStr(f.Value, types.StrConv(def))
		if err != nil {
			panic("failed to set default")
		}
	}
}

func WithRepeats(ignore bool) FlagOption {
	return func(f *Flag) {
		if f.HasCallback() {
			f.Warnf("WithRepeats() is redundant if WithCallback() is used (%s)", f)
			return
		}
		if f.IsCounter() {
			f.Warnf("WithRepeats() is redundant if AsCounter() is used (%s)", f)
			return
		}
		if !f.IsScalar() {
			f.Warnf("WithRepeats() is redundant if the value is not a scalar (%s)", f)
			return
		}

		f.Type.SetRepeatsBit()
		if ignore {
			f.Type.SetIgnoreRepeatsBit()
		} else {
			// Shouldn't be necessary, but...
			f.Type.ClrIgnoreRepeatsBit()
		}
	}
}

func AsCounter() FlagOption {
	return func(f *Flag) {
		if f.HasCallback() {
			panic("cannot use flag with callback as counter")
		}
		if !f.IsScalar() {
			panic("cannot use non-scalar (slice/object) as counter")
		}
		if !f.IsNumber() {
			panic("counter variable must be a number")
		}
		if f.IsRepeatable() {
			f.Warnf("WithRepeats() is irrelevant if AsCounter() is used (%s)", f)
		}
		f.Type.SetCounterBit()
	}
}

// We distinguish "not implemented" from "obsolete" or "deprecated"
// for those cases where end users might reasonably expect a
// particular option to be implemented, but it hasn't been for some
// reason other than deprecation/obsolesence.
func NotImplemented() FlagOption {
	return func(f *Flag) {
		f.Type.SetNotImplementedBit()
	}
}

func Deprecated() FlagOption {
	return func(f *Flag) {
		f.Type.SetObsoleteBit()
	}
}

func Obsolete() FlagOption {
	return func(f *Flag) {
		f.Type.SetObsoleteBit()
	}
}

func WithTypeTag(tag string) FlagOption {
	return func(f *Flag) {
		f.ValueTypeTag = tag
	}
}

func WithCallback(callback CallbackFunction) FlagOption {
	return func(f *Flag) {
		if f.IsCounter() {
			panic("callback supplied for counter")
		}
		f.Callback = callback
	}
}

func NewFlag(value interface{}, letter rune, label string, usage string, opts ...FlagOption) *Flag {
	// Require pointers as storage targets:
	if !types.IsPointer(value) {
		return nil
	}
	// Don't allow non-empty slices as storage targets:
	if types.IsSlice(value) && types.SliceLen(value) != 0 {
		return nil
	}
	// Require `PosixRejectQuest` to be set to `false` before we
	// accept `-?` (prohibited) as a short option
	if PosixRejectQuest && letter == '?' {
		panic("POSIX disallows '-?' as a single-letter option")
	}
	// Require `PosixRejectW` to be set to `false` before we accept
	// `-W` (reserved) as a short option
	if PosixRejectW && letter == 'W' {
		panic("POSIX reserves '-W' as a single-letter option")
	}

	if !(IsValidLabel(label) || IsValidShortcut(letter)) {
		return nil
	}
	if types.IsOtherT(value) {
		return nil
	}
	f := &Flag{
		Value:         value,
		Label:         label,
		Letter:        letter,
		Usage:         usage,
		Count:         0,
		ListSeparator: DefaultListSeparator,
	}
	if types.IsSlice(value) {
		f.Type.SetRepeatsBit()
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *Flag) NewAlias(letter rune, label string, opts ...FlagOption) *Flag {
	// alias has same type as target except that the appropriate alias
	// bits are set
	flagType := f.Type
	if IsValidLabel(label) {
		flagType.SetLabelAliasBit()
	}
	if letter == 0 {
		letter = NoShort
	}
	if letter == NoShort || IsValidShortcut(letter) {
		flagType.SetLetterAliasBit()
	}
	if !flagType.TstAliasBits() {
		return nil
	}

	a := &Flag{
		Value:         nil, // stored in `AliasFor` target
		Label:         label,
		Letter:        letter,
		AliasFor:      f,
		Type:          flagType,
		Count:         -1, // count in `AliasFor` target
		parentFlagSet: f.parentFlagSet,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
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
	return f.Type.TstRepeatsBit()
}
func (f *Flag) IgnoreRepeats() bool {
	return f.Type.TstIgnoreRepeatsBit()
}
func (f *Flag) IsScalar() bool {
	return !types.IsSlice(f.Value)
}
func (f *Flag) IsNumber() bool {
	return types.IsNum(f.Value)
}
func (f *Flag) HasCallback() bool {
	return f.Callback != nil
}

// Only allow letters, numbers, and the question-mark as shortcut
// letters
func IsValidShortcut(r rune) bool {
	if PosixRejectQuest && r == '?' {
		return false
	}
	if PosixRejectW && r == 'W' {
		return false
	}
	return r == '?' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

// Only allow letters, numbers, and underscore in labels
func IsValidLabel(label string) bool {
	// A label must be longer than one byte:
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
	f.ParentFlagSet().Failf(format, args...)
}

func (f *Flag) Infof(format string, args ...interface{}) {
	f.ParentFlagSet().Infof(format, args...)
}

func (f *Flag) Warnf(format string, args ...interface{}) {
	f.ParentFlagSet().Warnf(format, args...)
}
