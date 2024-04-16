// Package `fflag` provides POSIX (short) and GNU (long) command-line
// argument parsing with, for the programmer, the functional options
// pattern.
//
// It is somewhat inspired by the `pflag` package in some respects,
// but very significantly different in others. The most significant
// difference is that there is only one `Var()` function (rather than
// six dozen variations): the type of the flag is determined by the
// type of the first argument (rather than the function name), which
// _must_ be a pointer to a basic type, a slice of basic type, or a
// struct implementing the `SetValue` interface (inspired by `pflag`).
//
// The other significant difference is the order of the short flag and
// long flag in the `Var()` argument list, with the short flag coming
// first as a `rune`, most often a single ASCII letter or number. The
// reason for this is that short flags are always listed first in
// manpages and other documentation, so it's actually a bit weird of
// `pflag` to have reversed this de-facto standard order and, in
// practice, I've found it handier to obey the standard order than
// stick to the reversed `flag` argument order.
//
// If there is no short flag, the zero value (0, `'\0'`, or
// `fflag.NoShort`) is used. The usual rules apply to long flags,
// which must consist of letters and numbers, except that the ASCII
// requirement has been relaxed. Any character satisfying
// `unicode.IsLetter()` or `unicode.IsNumber()`, or the ASCII
// hyphen/minus character, '-', are allowed. There is _no attempt
// whatsoever_ at (what `pflag` refers to as) normalization, a very
// dubious utility: just use the long flag you mean to use without
// weird capitalization.
//
// That said, `fflag` meets the onerous GNU expectation that “users
// can abbreviate the option names as long as the abbreviations are
// unique”.
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
// single-short version of the flag as a rune (or `NoShort` if none),
// e.g., 'h'; the long version of the flag (or `NoLong` if none),
// e.g. `--help`; and a very brief description of the flag's
// purpose. For example:
//
//     var value bool
//     fflag.Var(&value, 'h', "help", "prints a help message to stdout")
//
// The first argument to `Var` must be a _pointer_ to one of:
//
//   1) a basic datatype (e.g. `int8`, `float32`, `string`)
//   2) a slice of basic datatype (e.g. `[]int8`, `[]string`)
//   3) something implementing the `SetValue` interface
//
// Non-pointer `value` arguments will cause a `panic()`. As a rule,
// `fflag` will `panic()` in the case of a programmer mistake (during
// setup) and return an `error` otherwise (during argument parsing).
//
// If the value argument implements the `SetValue` interface, `fflag`
// neither modifies the argument itself nor enforces any of its usual
// rules. If you pass something implementing this interface, it's
// assumed that you will take care of everything and don't want
// `fflag` to do anything other than pass along the message “this flag
// appeared with this argument”.
//
// A flag need not have a single-character shortcut. If there is no
// shortcut, a `fflag.NoShort` is given for the shortcut argument:
//
//    fflag.Var(&value, fflag.NoShort, "help", "prints a help message to stdout")
//
// Only letters and numbers are normally allowed as shortcuts. The
// sole exception is '?' due to its widespread use as an alias for
// "help", but this is prohibited by POSIX, so if you want to use
// this, you have to enable it explicitly.
//
// Equally, a flag need not have a long version. If you wanted to have
// `-?` as a short flag with no long version, you would do:
//
//    fflag.PosixRejectQuest = false
//    fflag.Var(&value, '?', fflag.NoLong, "prints a help message to stdout")
//
// There is a special case (and common idiom) where NEITHER a long NOR
// a short form is required: `-NUM` (as in `grep`, `head`, `tail`, and
// several other tools). These special cases are always an alias for
// something else and always refer to a natural number appearing after
// a single hyphen. For example `head`'s `-n/--lines` is best
// represented as:
//
//    uint nlines
//    fflag.Var(&nlines, 'n', 'lines',
//        "print the first NUM lines instead of the first 10",
//        fflag.WithAlias(fflag.NoShort, fflag.NoLong, false),
//        fflag.WithTypeTag("[-]NUM"))
//
// Obviously, this special case can only be used once and attempting
// to use it more than once is a programmer error that results in a
// `panic()`.
//
// The simplest “vanilla” flag is a nullary boolean switch that takes
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
// callback function to be called _every time_ the flag appears on the
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
// be an obvious programmer error, not something that could
// “accidentally” occur at runtime based on user input.
//
// An explicit default can be supplied with `WithDefault()`:
//
//     var hard bool
//     fflag.Var(&hard, 0, "easy", "use easy mode",
//         fflag.WithDefault(true))
//
// In this case, `hard` will default to `true` and become false if
// `--easy` appears on the command line. If repeats are allowed and
// not ignored (`WithRepeat(false)`), the value will toggle between
// `true` and `false`, which is admittedly weird, but if you do stupid
// things, expect to win stupid prizes.
//
// Repeated appearances of a flag are _not_ an error if the value
// argument is (a pointer to) a slice. In this case, successive
// invocations will result in successive values being appended to the
// slice. This makes no sense for a nullary boolean flag.
//
//     values := []int{}
//     NewFlag(&values, 'x', "example", "example flag")
//
// The sole exception to this rule is where a callback function is
// supplied. When a callback is supplied, the callback is responsible
// for _everything_.
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
// The `value` is _not_ set by `fflag` if a callback is supplied (or
// if `value` implements the `SetValue` interface).
//
// For unary (non-boolean) flags, a default can be supplied:
//
//     var file string
//     fflag.Var(&file, 'f', "file", "supply a filename",
//         fflag.WithDefault("/dev/null"))
//
// The value will _always_ be set to the default.
//
// There exist a few utilities with options that differentiate between
// the option appearing (at all) and the option appearing with an
// argument. To support this, you must use `WithOptionalDefault()`
//
//     var outdev := "/dev/stdout"
//     fflag.Var(&outdev, 'o', "out", "supply an output device path",
//         fflag.WithOptionalDefault("/dev/stderr"))
//
// In this case, if `-o` is _not_ supplied at all, `outdev` is not
// changed and keeps the value "/dev/stdout"; if `-o` is supplied
// _without an argument_, `outdev` changes to the default,
// `"/dev/stderr"`, and if an argument is supplied, `outdev` is
// changed to the given argument.
//
// If the value is a (pointer to a) scalar, but the default is a
// slice, the value is constrained to the values in the default, like
// a kind of set or `enum`.
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
// There exist a few utilities with set-constrained options that
// differentiate between the option appearing (at all) and the option
// appearing with an argument. For example, `grep`'s `--color`
// argument (with no option) is different from the option not being
// given at all, but also different from the option being given with
// an argument. This is rare, but (unfortunately) it exists in the
// wild.
//
//     color := "never"
//     fflag.Var(&file, fflag.NoShort, "color", "display matches in color",
//         fflag.WithOptionalDefault([]string{"auto", "never", "always"}),
//         fflag.WithTypeTag("[=WHEN]"))
//
// In this case, if the `color` option is not given, the value is not
// changed and remains at "never". If the argument is given _with no
// argument_, the value is changed to "auto", the first value in the
// default slice, and if the argument is given _with_ an argument, the
// value changes to that argument _provided that_ it is in the default
// slice. If `--color=foo` were given, it would result in an error.
//
package fflag

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/EmmetCaulfield/fflag/pkg/types"
)

var DefaultListSeparator string = ","

// POSIX uses '?' for a special purpose in `getopt()`, making it
// unsuitable for use as an option, but some applications use it
// explicitly, often for help, so we allow it, but reject it by
// default.
var PosixRejectQuest bool = true

// POSIX reserves `-W` for vendor options
var PosixRejectW bool = true

// The equals separator, if present, is regarded as part of the
// option-argument under POSIX rules
var PosixEquals bool = true

// A double-hyphen terminates option processing only when it is not an
// option-argument. When set to false, the GNU convention of ALWAYS
// terminating option processing is followed.
var PosixDoubleHyphen bool = true

// POSIX mandates that argument processing is terminated when the
// first non-option argument, which it calls an operand, is
// encountered. When set to false, the GNU convention (allowing a
// mixture of option and non-option) is followed and arguments that
// are not processable as options or option-arguments are appended to
// the operand list as they are encountered.
var PosixOperandStop bool = true

type FlagError struct {
	s string
}

func (fe *FlagError) Error() string {
	return fe.s
}

type CallbackFunction func(f *Flag, arg string, pos int) error

type FlagType uint16

const (
	ClearFlagType     FlagType = 0b0000000000000000
	LongAliasBit      FlagType = 0b0000000000000001
	ShortAliasBit     FlagType = 0b0000000000000010
	ObsoleteBit       FlagType = 0b0000000000000100
	NotImplementedBit FlagType = 0b0000000000001000
	HiddenBit         FlagType = 0b0000000000010000
	ChangedBit        FlagType = 0b0000000000100000
	CounterBit        FlagType = 0b0000000001000000
	RepeatsBit        FlagType = 0b0000000010000000
	IgnoreRepeatsBit  FlagType = 0b0000000100000000
	FileBit           FlagType = 0b0000001000000000
	DefOptionalBit    FlagType = 0b0000010000000000
	SavedFileBit      FlagType = 0b0000100000000000
)

func (ft *FlagType) TstLongAliasBit() bool      { return *ft&LongAliasBit != 0 }
func (ft *FlagType) TstShortAliasBit() bool     { return *ft&ShortAliasBit != 0 }
func (ft *FlagType) TstObsoleteBit() bool       { return *ft&ObsoleteBit != 0 }
func (ft *FlagType) TstNotImplementedBit() bool { return *ft&NotImplementedBit != 0 }
func (ft *FlagType) TstHiddenBit() bool         { return *ft&HiddenBit != 0 }
func (ft *FlagType) TstChangedBit() bool        { return *ft&ChangedBit != 0 }
func (ft *FlagType) TstCounterBit() bool        { return *ft&CounterBit != 0 }
func (ft *FlagType) TstRepeatsBit() bool        { return *ft&RepeatsBit != 0 }
func (ft *FlagType) TstIgnoreRepeatsBit() bool  { return *ft&IgnoreRepeatsBit != 0 }
func (ft *FlagType) TstFileBit() bool           { return *ft&FileBit != 0 }
func (ft *FlagType) TstDefOptionalBit() bool    { return *ft&DefOptionalBit != 0 }
func (ft *FlagType) TstSavedFileBit() bool      { return *ft&SavedFileBit != 0 }
func (ft *FlagType) TstAliasBits() bool         { return (*ft&ShortAliasBit)|(*ft&LongAliasBit) != 0 }

func (ft *FlagType) ClrLongAliasBit()      { *ft = *ft & ^LongAliasBit }
func (ft *FlagType) ClrShortAliasBit()     { *ft = *ft & ^ShortAliasBit }
func (ft *FlagType) ClrObsoleteBit()       { *ft = *ft & ^ObsoleteBit }
func (ft *FlagType) ClrNotImplementedBit() { *ft = *ft & ^NotImplementedBit }
func (ft *FlagType) ClrHiddenBit()         { *ft = *ft & ^HiddenBit }
func (ft *FlagType) ClrChangedBit()        { *ft = *ft & ^ChangedBit }
func (ft *FlagType) ClrCounterBit()        { *ft = *ft & ^CounterBit }
func (ft *FlagType) ClrRepeatsBit()        { *ft = *ft & ^RepeatsBit }
func (ft *FlagType) ClrIgnoreRepeatsBit()  { *ft = *ft & ^IgnoreRepeatsBit }
func (ft *FlagType) ClrFileBit()           { *ft = *ft & ^FileBit }
func (ft *FlagType) ClrDefOptionalBit()    { *ft = *ft & ^DefOptionalBit }
func (ft *FlagType) ClrSavedFileBit()      { *ft = *ft & ^SavedFileBit }

func (ft *FlagType) SetLongAliasBit()      { *ft = *ft | LongAliasBit }
func (ft *FlagType) SetShortAliasBit()     { *ft = *ft | ShortAliasBit }
func (ft *FlagType) SetObsoleteBit()       { *ft = *ft | ObsoleteBit }
func (ft *FlagType) SetNotImplementedBit() { *ft = *ft | NotImplementedBit }
func (ft *FlagType) SetHiddenBit()         { *ft = *ft | HiddenBit }
func (ft *FlagType) SetChangedBit()        { *ft = *ft | ChangedBit }
func (ft *FlagType) SetCounterBit()        { *ft = *ft | CounterBit }
func (ft *FlagType) SetRepeatsBit()        { *ft = *ft | RepeatsBit }
func (ft *FlagType) SetIgnoreRepeatsBit()  { *ft = *ft | IgnoreRepeatsBit }
func (ft *FlagType) SetFileBit()           { *ft = *ft | FileBit }
func (ft *FlagType) SetDefOptionalBit()    { *ft = *ft | DefOptionalBit }
func (ft *FlagType) SetSavedFileBit()      { *ft = *ft | SavedFileBit }

type Flag struct {
	Value         interface{}
	Long          string
	Short         rune
	Type          FlagType
	Count         int
	ValueTypeTag  string
	Default       interface{}
	AliasFor      *Flag
	Usage         string
	Callback      CallbackFunction
	ListSeparator string
	Mutexes       map[string]struct{}
	parentFlagSet *FlagSet
	savedCallback CallbackFunction
}

const IdSep string = "/"

// A non-numeric, non-alphabetic ASCII character (other than '?') used
// as a placeholder meaning "there is no short option" in a variety of
// contexts.
const NoShort rune = rune(0)
const NoLong string = ""

// Only allow letters, numbers, and the question-mark as shortcut
// letters
func IsValidShort(r rune) bool {
	if PosixRejectQuest && r == '?' {
		log.Panicf("cannot use '-?' as a short option if `fflag.PosixRejectQuest` is `true`")
	}
	if PosixRejectW && r == 'W' {
		log.Panicf("cannot use '-W' as a short option if `fflag.PosixRejectW` is `true`")
	}
	return r == '?' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

// Only allow letters, numbers, and hyphens in labels
func IsValidLong(s string) bool {
	// A long must be longer than one byte:
	if len(s) < 2 {
		return false
	}
	// A long can't begin with a hyphen
	if s[0] == '-' {
		return false
	}
	// Longs must otherwise consist entirely of letters, numbers, and
	// hyphens
	for _, r := range s {
		if r == '-' || unicode.IsLetter(r) || unicode.IsNumber(r) {
			continue
		}
		return false
	}
	return true
}

func IsValidPair(short rune, long string) bool {
	if short == NoShort && long == NoLong {
		// We use this for the -NUM special case used by a few utilities
		// (e.g. head, tail), which has NEITHER a normal valid shortcut
		// nor a normal valid long flag
		return true
	}
	goodShort := IsValidShort(short)
	if long == NoLong && goodShort {
		return true
	}
	goodLong := IsValidLong(long)
	if short == NoShort && goodLong {
		return true
	}
	if goodShort && goodLong {
		return true
	}
	return false
}

func ID(short rune, long string) string {
	if IsValidPair(short, long) {
		return string(short) + IdSep + long
	}
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
	if parts[1] == NoLong && emptyOrNoShort(parts[0]) {
		return NoShort, NoLong
	}

	short, tail := FirstRune(parts[0])
	if short < 0 || tail != "" {
		return ErrRuneShortBad, NoLong
	}
	if IsValidPair(short, parts[1]) {
		return short, parts[1]
	}
	return ErrRuneIdPartsBad, ""
}

func (f *Flag) ParentFlagSet() *FlagSet {
	if f.parentFlagSet == nil {
		return CommandLine
	}
	return f.parentFlagSet
}

func (f *Flag) MutexCollides() *Flag {
	if f.Mutexes == nil || len(f.Mutexes) == 0 {
		return nil
	}
	fs := f.ParentFlagSet()
	for name, _ := range f.Mutexes {
		if flag, ok := fs.Mutex[name]; ok {
			if flag != nil {
				return flag
			}
			fs.Mutex[name] = f
		}
	}
	return nil
}

// Function `EnterCallback()` saves and alters some internal state of
// a flag that would otherwise cause recursion in `fflag.Set()` if and
// when it is called by a callback (`fflag.Set()` calls the
// callback). It is intended to be called inside a callback before any
// call to `Set()`. It does not need to be called if the callback does
// not call `Set()` on the flag pointer.
func (f *Flag) EnterCallback() {
	if f.Callback == nil {
		panic("cannot enter callback if there is no callback")
	}
	f.savedCallback = f.Callback
	f.Callback = nil
	if f.Type.TstFileBit() {
		f.Type.SetSavedFileBit()
	} else {
		f.Type.ClrSavedFileBit()
	}
	f.Type.ClrFileBit()
}

// Function `ExitCallback()` restores the state of a flag that was
// modified by `EnterCallback()` to prevent recursion in
// `fflag.Set()`. It is intended to be called inside a callback after
// the last call to `Set()`. It does not need to be called if the
// callback does not call `Set()` on the flag pointer and
// `EnterCallback()` has not been called.
func (f *Flag) ExitCallback() {
	if f.Callback != nil {
		panic("cannot exit callback without prior call to EnterCallback()")
	}
	f.Callback = f.savedCallback
	if f.Type.TstSavedFileBit() {
		f.Type.SetFileBit()
	} else {
		f.Type.ClrFileBit()
	}
}

func (f *Flag) Test(value interface{}, argPos int) error {
	return f.testOrSet(value, argPos, false)
}

func (f *Flag) Set(value interface{}, argPos int) error {
	return f.testOrSet(value, argPos, true)
}

// TestOrSet() sets `f.Value` to `value` if `doSet` is `true`,
// otherwise it silently tests, insofar as possible, whether the set
// would succeed or not.
func (f *Flag) testOrSet(value interface{}, argPos int, doSet bool) error {
	prev := f.MutexCollides()
	if prev != nil {
		f.Failf("flag '%s' conflicts with previously given flag '%s'", f, prev)
		return &FlagError{"mutex collision in Flag.Set()"}
	}
	// Prefer the SetValue interface if present:
	if setter, ok := f.Value.(types.SetValue); ok {
		if str, ok := value.(string); ok {
			f.Count++
			if doSet {
				return setter.Set(str)
			}
			return nil
		}
		if doSet {
			f.Failf("Cannot pass non-string to SetValue.Set(string) in flag.Set() for flag '%s'", f)
		}
		return &FlagError{"cannot pass non-string to SetValue.Set()"}
	}

	if f.AliasFor != nil {
		f = f.AliasFor
	}
	if f.AliasFor != nil {
		log.Panic("double alias in Flag.Set(...)")
	}

	if doSet {
		f.Count++
	}
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
		err := f.testOrSetOnly(str, argPos, doSet)
		if err != nil && doSet {
			f.Failf("failed to set counter '%s' from %d: %v", f, f.Count, err)
		}
		return err
	}

	if f.IsFileReader() {
		filename, ok := value.(string)
		if !ok {
			if doSet {
				f.Failf("file-reader flag %s expects a filename (string) argument", f)
			}
			return &FlagError{"argument to flag-reader not a string"}
		}
		file, err := os.Open(filename)
		if err != nil {
			if doSet {
				f.Failf("failed to open file '%s' for flag '%s': %v", filename, f, err)
			}
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			if f.HasCallback() {
				if doSet {
					err := f.Callback(f, line, lineNo)
					if err != nil {
						f.Failf("callback failed for '%s' in '%s': %v", line, f, err)
						return err
					}
				}
				continue
			}
			err := f.testOrSetOnly(line, lineNo, doSet)
			if err != nil {
				if doSet {
					f.Failf("failed to set '%s' from line %d in '%s': %v", f, lineNo, filename, err)
				}
				return err
			}
			if err = scanner.Err(); err != nil {
				if doSet {
					f.Failf("error scanning '%s': %v", filename, err)
				}
				return err
			}
		}
		return nil
	}

	if f.HasCallback() {
		v, _ := value.(string)
		if doSet {
			return f.Callback(f, v, argPos)
		}
		return nil
	}

	if f.Count > 1 && !f.IsRepeatable() {
		if doSet {
			f.Failf("flag '%s' is not repeatable", f.String())
		}
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
			if doSet {
				*boolp = !def
			}
			return nil
		}
		value = f.GetDefault()
		if value == nil {
			if doSet {
				f.Failf("flag.Set(nil) called for flag '%s' with no default", f)
			}
			return &FlagError{"cannot set nil value for non-bool with no default"}
		}
	} else if !f.InDefaults(value) {
		// TODO(emmet): consider supporting constrained defaults in a
		// command-line list. It's reasonable to expect that each item
		// in a list optarg would be checked against the list in
		// f.Default (if any), but in reality, a non-scalar optarg
		// (e.g. `-x foo,bar,baz`) will fail here.
		if doSet {
			f.Failf("value %v not found in defaults %v for '%s'", value, f.Default, f)
		}
		return &FlagError{"value constrained by defaults"}
	}
	return f.testOrSetOnly(value, argPos, doSet)
}

func (f *Flag) TestOnly(value interface{}, argPos int) error {
	return f.testOrSetOnly(value, argPos, false)
}

func (f *Flag) SetOnly(value interface{}, argPos int) error {
	return f.testOrSetOnly(value, argPos, true)
}

// Function testOrSetOnly() sets `f.Value` to `value` if `doSet` is
// `true`, otherwise it silently tests, insofar as possible, whether
// the set would succeed or not.
func (f *Flag) testOrSetOnly(value interface{}, argPos int, doSet bool) error {
	// Convert the value to a string if it's not already one
	var ok bool
	var str string
	if str, ok = value.(string); !ok {
		str = types.StrConv(value, types.WithSep(f.ListSeparator))
		if str == "" {
			if doSet {
				f.Failf("failed to convert '%v' to a nonempty string in '%s'", value, f)
			}
			return &FlagError{"cannot convert value to string"}
		}
	}

	// Set the value from the string version
	err := types.FromStr(f.Value, str, doSet, types.WithSep(f.ListSeparator))
	if err != nil {
		if doSet {
			f.Failf("failed to convert '%s' to %T: %v", str, f.Value, err)
		}
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
		log.Panic("double alias in Flag.GetDefault()")
	}
	if types.IsSlice(f.Default) {
		if types.SliceLen(f.Default) > 0 {
			return types.ItemAt(f.Default, 0)
		}
		// Default is an empty slice, which should be impossible here
		log.Panic("f.Default is an empty slice in Flag.GetDefault()")
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
		log.Panic("double alias in Flag.InDefaults()")
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
	if f.Short != NoShort && len(f.Long) > 1 {
		return "-" + string(f.Short) + ", --" + f.Long
	}
	if len(f.Long) > 1 {
		return "--" + f.Long
	}
	if f.Short != NoShort {
		return "-" + string(f.Short)
	}
	return ""
}

func (f *Flag) FormatShort() string {
	if f.Short == NoShort {
		if f.Long == NoLong {
			return "-NUM"
		}
		return ""
	}

	tag := f.GetTypeTag()
	if len(tag) == 0 || f.IsAlias() {
		return "-" + string(f.Short)
	}

	if f.Type.TstDefOptionalBit() {
		return "-" + string(f.Short) + "[" + tag + "]"
	}
	return "-" + string(f.Short) + " " + tag
}

func (f *Flag) FormatLong() string {
	if f.Long == NoLong {
		return ""
	}

	tag := f.GetTypeTag()
	if len(tag) == 0 || f.IsAlias() {
		return "--" + f.Long
	}

	if f.Type.TstDefOptionalBit() {
		return "--" + f.Long + "[=" + tag + "]"
	}
	return "--" + f.Long + "=" + tag
}

// Returns f.String() wrapped in extra stuff for help/usage output
func (f *Flag) FlagString() string {
	short := f.FormatShort()
	if short == "" {
		short = "    "
	}
	if f.Long == NoLong {
		return short
	}
	if f.Short == NoShort {
		return short + f.FormatLong()
	}
	return short + ", " + f.FormatLong()
}

func (f *Flag) DescString() string {
	if f.Type.TstAliasBits() && f.AliasFor != nil {
		buf := &bytes.Buffer{}
		if f.Type.TstObsoleteBit() {
			buf.WriteString("obsolete ")
		}
		buf.WriteString("synonym for ")
		buf.WriteString(f.AliasFor.String())
		if f.Value != nil {
			if value, ok := f.Value.(string); ok {
				buf.WriteString("=" + value)
			}
		}
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
	if f.Short == NoShort {
		return f.Long
	}
	// TODO(emmet): think about special case of no long and no short
	// used for -NUM
	return string(f.Short) + f.Long
}

type FlagOption = func(f *Flag) error

type AliasOption = func(f *Flag) error

func WithParent(fs *FlagSet) FlagOption {
	return func(f *Flag) error {
		if f.parentFlagSet != nil && f.parentFlagSet != fs {
			log.Panicf("attempt to change parent flagset in fflag.WithParent() for %s", f)
		}
		f.parentFlagSet = fs
		return nil
	}
}

func withValue(value string) AliasOption {
	return func(f *Flag) error {
		f.Value = value
		return nil
	}
}

func WithListSeparator(sep rune) FlagOption {
	return func(f *Flag) error {
		if f.IsHyphenNum() {
			log.Panicf("hyphen-num idiom cannot have a list sepearator")
		}
		if !types.IsSlice(f.Value) {
			log.Panicf("cannot set separator for non-list value %s", f)
		}
		f.ListSeparator = string(sep)
		return nil
	}
}

func WithAlias(short rune, long string, obsolete bool) FlagOption {
	return func(f *Flag) error {
		var flag *Flag = nil
		fs := f.ParentFlagSet()
		if long != NoLong {
			flag = fs.LookupLong(long)
			if flag != nil {
				log.Panicf("long flag '%s' already exists for alias '%s'", flag, long)
			}
		}
		if short != NoShort {
			flag = fs.LookupShort(short)
			if flag != nil {
				log.Panicf("short flag '%s' already exists for alias '%c'", flag, short)
			}
		}
		flag = f.NewAlias(short, long)
		if flag == nil {
			log.Panicf("error creating alias -%c/--%s for `%s`", short, long, f)
		}
		if obsolete {
			flag.Type.SetObsoleteBit()
		} else {
			flag.Type.ClrObsoleteBit()
		}
		err := fs.AddFlag(flag)
		if err != nil {
			log.Panicf("Error adding alias '%s': %v", flag, err)
		}
		return nil
	}
}

func (f *Flag) setupDefault(def interface{}, optional bool) error {
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
			log.Panicf("non-string default for '%s' where value implements SetValue interface", f)
		}
		if !types.SameBaseType(valType, defType) {
			log.Panicf("type mismatch: default type <%T> for value type <%T> in '%s'", def, f.Value, f)
		}
	}
	if optional {
		f.Type.SetDefOptionalBit()
	}
	// Really set the default value only if `optional` is `false`
	// (i.e. we want `doSet` argument to testOrSetOnly() to be true)
	f.Default = def
	def = f.GetDefault()
	err := f.testOrSetOnly(def, 0, !optional)
	if err != nil {
		log.Panicf("failed to set/test value to default (%v) for '%s': %v", f.Default, f, err)
	}
	return nil
}

func WithDefault(def interface{}) FlagOption {
	return func(f *Flag) error {
		return f.setupDefault(def, false)
	}
}

func WithOptionalDefault(def interface{}) FlagOption {
	return func(f *Flag) error {
		if f.IsBool() {
			log.Panicf("optional default makes no sense for boolean flag '%s'", f)
		}
		return f.setupDefault(def, true)
	}
}

func WithRepeats(ignore bool) FlagOption {
	return func(f *Flag) error {
		if f.IsHyphenNum() {
			log.Panicf("hyphen-num idiom cannot be repeated")
		}
		if f.HasCallback() {
			log.Panicf("WithRepeats() is redundant if WithCallback() is used (%s)", f)
		}
		if f.IsCounter() {
			log.Panicf("WithRepeats() is redundant if AsCounter() is used (%s)", f)
		}
		if !f.IsScalar() {
			log.Panicf("WithRepeats() is redundant if the value is not a scalar (%s)", f)
		}

		f.Type.SetRepeatsBit()
		if ignore {
			f.Type.SetIgnoreRepeatsBit()
		} else {
			// Shouldn't be necessary, but...
			f.Type.ClrIgnoreRepeatsBit()
		}
		return nil
	}
}

func AsCounter() FlagOption {
	return func(f *Flag) error {
		if f.IsHyphenNum() {
			log.Panicf("hyphen-num idiom cannot be a counter")
		}
		if f.HasCallback() {
			log.Panicf("cannot use flag '%s' with callback as counter", f)
		}
		if !f.IsScalar() {
			log.Panicf("cannot use non-scalar (type %T) as counter for '%s'", f.Value, f)
		}
		if !f.IsNumber() {
			log.Panicf("counter variable must be a number, not type %T, for '%s'", f.Value, f)
		}
		if f.IsRepeatable() {
			log.Panicf("WithRepeats() is irrelevant for counter '%s'", f)
		}
		f.Type.SetCounterBit()
		return nil
	}
}

// We distinguish "not implemented" from "obsolete" or "deprecated"
// for those cases where end users might reasonably expect a
// particular option to be implemented, but it hasn't been for some
// reason other than deprecation/obsolesence.
func NotImplemented() FlagOption {
	return func(f *Flag) error {
		f.Type.SetNotImplementedBit()
		return nil
	}
}

func Deprecated() FlagOption {
	return func(f *Flag) error {
		f.Type.SetObsoleteBit()
		return nil
	}
}

func Obsolete() FlagOption {
	return Deprecated()
}

func WithTypeTag(tag string) FlagOption {
	return func(f *Flag) error {
		f.ValueTypeTag = tag
		return nil
	}
}

func WithCallback(callback CallbackFunction) FlagOption {
	return func(f *Flag) error {
		if f.IsCounter() {
			log.Panicf("callback supplied for counter '%s'", f)
		}
		f.Callback = callback
		return nil
	}
}

// A file-reading flag can't be a counter, have a callback, or be an
// alias:
func ReadFile() FlagOption {
	return func(f *Flag) error {
		if f.IsHyphenNum() {
			log.Panicf("hyphen-num idiom cannot be a file reader")
		}
		if f.IsCounter() {
			log.Panicf("counter flag '%s' cannot be a file reader", f)
		}
		if f.HasCallback() {
			log.Panicf("flag '%s' with callback cannot be a file reader", f)
		}
		if f.IsAlias() {
			log.Panicf("alias flag '%s' cannot be a file reader", f)
		}
		if !types.IsSlice(f.Value) {
			log.Panicf("value of file reader flag '%s' must point at a slice", f)
		}
		f.Type.SetFileBit()
		return nil
	}
}

func InMutex(name string) FlagOption {
	return func(f *Flag) error {
		fs := f.ParentFlagSet()
		fs.Mutex[name] = nil
		if _, ok := f.Mutexes[name]; ok {
			log.Panicf("flag '%s' added to mutex '%s' more than once", f, name)
		} else {
			f.Mutexes[name] = struct{}{}
		}
		return nil
	}
}

func NewFlag(value interface{}, short rune, long string, usage string, opts ...FlagOption) *Flag {
	// We potentially use the type identifier several times
	valType := types.Type(value)

	// Require pointers as storage targets:
	if !valType.TstPointerBit() {
		return nil
	}
	// Don't allow non-empty slices as storage targets:
	if valType.TstSliceBit() && types.SliceLen(value) != 0 {
		return nil
	}
	// We don't know what to do with things that are neither basic
	// types nor implement the SetValue interface:
	if valType.TstOtherBit() {
		log.Panicf("value type <%T> is not supported", value)
	}
	if !IsValidPair(short, long) {
		log.Panicf("flag pair '-%c/--%s' is not valid or not permitted", short, long)
	}
	if short == NoShort && long == NoLong {
		// Special -NUM idiom
		if valType.TstSliceBit() || !valType.TstUintBit() {
			log.Panicf("a scalar unsigned integer is required for the -NUM idiom")
		}
	}
	f := &Flag{
		Value:         value,
		Long:          long,
		Short:         short,
		Usage:         usage,
		Count:         0,
		ListSeparator: DefaultListSeparator,
		Mutexes:       map[string]struct{}{},
	}
	if valType.TstSliceBit() {
		f.Type.SetRepeatsBit()
	}
	for i, opt := range opts {
		err := opt(f)
		if err != nil {
			log.Panicf("error setting option %d for flag '%s'", i, f)
		}
	}
	return f
}

func (f *Flag) NewAlias(short rune, long string, opts ...FlagOption) *Flag {
	// alias has same type as target except that the appropriate alias
	// bits are set
	flagType := f.Type
	if !IsValidPair(short, long) {
		log.Panicf("short/long pair -%c/--%s not valid in NewAlias()", short, long)
	}
	if short == 0 {
		short = NoShort
	}
	if short == NoShort || IsValidShort(short) {
		flagType.SetShortAliasBit()
	}
	if !flagType.TstAliasBits() {
		return nil
	}

	a := &Flag{
		Value:         nil, // stored in `AliasFor` target
		Long:          long,
		Short:         short,
		AliasFor:      f,
		Type:          flagType,
		Count:         -1, // count in `AliasFor` target
		parentFlagSet: f.parentFlagSet,
	}

	for i, opt := range opts {
		err := opt(a)
		if err != nil {
			log.Panicf("error setting option %d for alias '%s'", i, f)
		}
	}

	return a
}

func (f *Flag) IsLongAlias() bool {
	return f.Type.TstLongAliasBit()
}
func (f *Flag) IsShortAlias() bool {
	return f.Type.TstShortAliasBit()
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
func (f *Flag) IsFileReader() bool {
	return f.Type.TstFileBit()
}
func (f *Flag) IgnoreRepeats() bool {
	return f.Type.TstIgnoreRepeatsBit()
}
func (f *Flag) IsScalar() bool {
	return !types.IsSlice(f.Value)
}
func (f *Flag) IsBool() bool {
	return types.IsBool(f.Value)
}
func (f *Flag) IsNumber() bool {
	return types.IsNum(f.Value)
}
func (f *Flag) IsHyphenNum() bool {
	return f.Short == NoShort && f.Long == NoLong
}
func (f *Flag) HasCallback() bool {
	return f.Callback != nil
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
