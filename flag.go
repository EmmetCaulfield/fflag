// The `fflag` package provides GNU/POSIX style command-line argument
// parsing with the functional options pattern.
//
// It is somewhat inspired by the `pflag` package in some respects,
// but significantly different in others. The most significant
// difference is that there is only one `Var()` function: the type of
// the flag is determined by the type of the first argument, which
// MUST be a pointer to a basic type, a slice of basic type, or a
// struct implementing the `SetValue` interface (inspired by
// `pflag`).
//
// The next significant difference is the order of the short flag and
// long flag in the `Var()` argument list, with the short flag coming
// first as a `rune`, which must be a single UTF-8 letter, most often
// a single ASCII letter or number. If there is no short flag, the
// zero value (0, or `\0') is used. The usual rules apply to long
// flags, which must consist of letters and numbers, except that the
// ASCII requirement has been relaxed. Any character satisfying
// unicode.IsLetter() or unicode.IsNumber() or the hyphen '-' are
// allowed. There is no attempt at normalization, a dubious utility: just use the long
// flag you mean to use.
//
// `fflag` borrows the `Flag` and `FlagSet` names from `pflag`, adding
// `FlagGroup`. The purpose of a flag group is to enable usage
// information to be generated in a similar format to GNU/POSIX
// utilities like `grep`.
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
//   3) a `struct` implementing the `Set()` interface
//
// A flag need not have a single-character shortcut. If there is no
// shortcut, a 0 is given for the shortcut argument:
//
//    fflag.Var(&value, 0, "help", "prints a help message to stdout")
//
// Punctuation (or other non-letter, non-number) characters are not
// normally allowed as shortcuts. The sole exception is '?' due to its
// widespread use as an alias for "help".
//
// Equally, a flag need not have a long version either:
//
//    fflag.Var(&value, '?', "", "prints a help message to stdout")
//
// Indeed, there is a special case (and common idiom) where NEITHER a
// long nor a short form is required: `-NUM` (as in `grep`, `head`,
// `tail`, and several other tools). These special cases are always an
// alias for something else and always refer to an integer appearing
// after a single hyphen. For example `head`'s `-n/--lines` is best
// represented as:
//
//    int nlines
//    fflag.Var(&nlines, 'n', 'lines', "print the first NUM lines instead of the first 10",
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
// causes the number of occurrences to be counted, and
// `WithCallback()` causes the given callback function to be called
// every time the flag appears on the command-line.
//
// Several utilities allow `-v/--verbose` to be repeated for
// increasing levels of verbosity.
//
//     int verbosity
//     f := NewFlag(&verbosity, 'v', "verbose", "increase verbosity", AsCounter())
//
// Note that it would be an error to supply BOTH the `AsCounter()` and
// `WithCallback()` options for the same flag because `AsCounter()`
// must modify the `value` while `WithCallback()` leaves modification
// of the value to the callback function for flexibility.
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
// While it is an error to repeat a command-line argument whose
// `value` argument is a pointer to a scalar (but see `WithRepeats()`
// and `AsCounter()` options), if the value argument is a pointer to a
// slice, successive invocations will result in successive values
// being appended to the slice.
//
//     values := []bool{}
//     NewFlag(&values, 'x', "example", "example flag")
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
	"strings"
	"unicode"

	"github.com/EmmetCaulfield/fflag/pkg/types"
)

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
	parentFlagSet *FlagSet
}

const IdSep string = "/"

// A non-numeric, non-alphabetic ASCII character other than '?' used
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

	// Last-ditch attempt: round-trip the value
	str := types.StrConv(value)
	err := types.FromStr(f.Value, str)
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
	if f.GetDefaultLen() > 0 {
		return types.ItemAt(f.Default, 0)
	}
	return f.Default
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
		if f.parentFlagSet != nil {
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
		f.Default = def
		err := types.FromStr(f.Value, types.StrConv(def))
		if err != nil {
			f.Failf("failed to set default value for '%s'", f.Label)
		}
	}
}

func WithRepeats(ignore bool) FlagOption {
	return func(f *Flag) {
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

	if !(IsValidLabel(label) || IsValidShortcut(letter)) {
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
		Letter: letter,
		Usage:  usage,
		Count:  0,
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
