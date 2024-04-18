package fflag

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/EmmetCaulfield/fflag/pkg/deque"
	"github.com/EmmetCaulfield/fflag/pkg/trie"
)

// What to do on error. The default is the zero value: (not silent,
// don't continue, don't panic), i.e. produce a message and exit.
type FailOption int8

const(
	FailDefault  FailOption = 0b00000000
	FailSilent              = 0b00000001
	FailContinue            = 0b00000010
	FailPanic               = 0b00000100
)

func (fb *FailOption) TstSilentBit() bool   { return *fb&FailSilent != 0 }
func (fb *FailOption) TstContinueBit() bool { return *fb&FailContinue != 0 }
func (fb *FailOption) TstPanicBit() bool    { return *fb&FailPanic != 0 }
func (fb *FailOption) SetSilentBit()   { *fb |= FailSilent }
func (fb *FailOption) SetContinueBit() { *fb |= FailContinue }
func (fb *FailOption) SetPanicBit()    { *fb |= FailPanic }
func (fb *FailOption) ClrSilentBit()   { *fb &= ^FailSilent }
func (fb *FailOption) ClrContinueBit() { *fb &= ^FailContinue }
func (fb *FailOption) ClrPanicBit()    { *fb &= ^FailPanic }



// `CommandLine` is the default `FlagSet`, named by analogy with
// `pflag`'s variable with the same purpose.
var CommandLine *FlagSet = NewFlagSet()

// The `FlagGroup` provides a way of grouping help/usage text about
// flags into groups with a title.
type FlagGroup struct {
	Title string
	FlagList          []*Flag
}

// Creates a new `FlagGroup` given a title for the group.
func NewFlagGroup(title string) *FlagGroup {
	return &FlagGroup{
		Title: title,
		FlagList: []*Flag{},
	}
}

// A `FlagSet` is the top-level object containing all the flags and
// mechanisms for looking up flags by short/long option string.
//
// In the most commonly expected use case, there will be one `FlagSet`
// for a program, the default `CommandLine`.
type FlagSet struct {
	Groups             []*FlagGroup
	GroupIndex         int
	LongTrie          *trie.TrieNode[Flag]
	ShortDict          map[rune]*Flag
	Output             io.Writer
	IgnoreDoubleDash   bool
	HasHyphenNumIdiom  bool
	HasNumberShorts    bool
	InputArgs         *deque.Deque[string]
	OutputArgs        *deque.Deque[string]
	OnFail             FailOption
	FailExitCode       int
	OnFileError        FailOption
	FileErrExitCode    int
	Mutex              map[string]*Flag
}

// DefaultFailExitCode is the exit code that will be used when
// argument processing fails and `OnFail` is not `Continue`
var DefaultFailExitCode int = 2

// Some programs (e.g. grep) use a different exit code for file errors
// than for other errors.
var DefaultFileErrExitCode int = 2

// Function `NewFlagGroup()` creates a new titled flag group within a
// flagset and makes it the default `FlagGroup` to which subsequent
// flags will be added.
func (fs *FlagSet) NewFlagGroup(title string) *FlagGroup {
	fg := NewFlagGroup(title)
	fs.Groups = append(fs.Groups, fg)
	fs.GroupIndex = len(fs.Groups) - 1
	return fg
}

// Function `Group()` returns a pointer to the current default
// `FlagGroup`.
func (fs *FlagSet) Group() *FlagGroup {
	// Happy for this to panic
	return fs.Groups[fs.GroupIndex]
}

// Function `Group()` creates a new titled flag group in the default
// `FlagSet`. If the `FlagSet` is empty, it just renames the default
// flag group.
func Group(title string) {
	fs := CommandLine
	// If there's only one group and no flags yet, just rename the
	// group
	if len(fs.Groups) == 1 && !fs.HasFlags() {
		fs.Groups[0].Title = title
		return
	}
	_ = CommandLine.NewFlagGroup(title)
}

// Functional option type for `FlagSet` options.
type FlagSetOption = func (fs *FlagSet)

// Creates a new flagset, applying the supplied functional options.
func NewFlagSet(opts ...FlagSetOption) *FlagSet {
	fs := &FlagSet {
		Groups: []*FlagGroup{
			NewFlagGroup("Options"),
		},
		GroupIndex:       0,
		LongTrie:         trie.NewTrie[Flag](),
		ShortDict:        map[rune]*Flag{},
		Output:           os.Stderr,
		IgnoreDoubleDash: false,
		InputArgs:        &deque.Deque[string]{},
		OutputArgs:       &deque.Deque[string]{},
		OnFail:           FailDefault,
		FailExitCode:     DefaultFailExitCode,
		OnFileError:      FailDefault,
		FileErrExitCode:  DefaultFileErrExitCode,
		Mutex:            map[string]*Flag{},
	}
	for _, opt := range opts {
		opt(fs)
	}
	return fs
}

// Option `WithGroupTitle()` sets the title of the default (first)
// flag group.
func WithGroupTitle(title string) FlagSetOption {
	return func(fs *FlagSet) {
		fs.Groups[0].Title = title
	}
}

// Option `WithOutputWriter()` sets the output writer for error
// messages used if `OnFail` is not `Silent`.
func WithOutputWriter(w io.Writer) FlagSetOption {
	return func(fs *FlagSet) {
		fs.Output = w
	}
}

// Option `WithPanicOnFail()` causes argument processing to panic on
// any failure.
func WithPanicOnFail() FlagSetOption {
	return func(fs *FlagSet) {
		fs.OnFail.SetPanicBit()
	}
}

// Option `WithContinueOnFail()` causes argument processing to
// continue on failure, likely causing unpredictable results.
func WithContinueOnFail() FlagSetOption {
	return func(fs *FlagSet) {
		fs.OnFail.SetContinueBit()
	}
}

// Option `WithSilentFail()` suppresses printing error messages due to
// argument processing failure.
func WithSilentFail() FlagSetOption {
	return func(fs *FlagSet) {
		fs.OnFail.SetSilentBit()
	}
}

// Function `HasFlags()` returns `true` if the `FlagSet` has any flags
// defined and `false` if the `FlagSet` is empty.
func (fs *FlagSet) HasFlags() bool {
	n := 0
	for _, g := range fs.Groups {
		n += len(g.FlagList)
	}
	return n > 0
}

// Function `LookupLong()` looks up a long flag in the `FlagSet`. If
// the string is a unique prefix match for a long option, a pointer to
// the corresponding `Flag` is returned, otherwise `nil` is returned.
//
// However, we treat single-rune longs as shorts preferentially,
// otherwise we can have the situation where -x and --x are different,
// which could happen if "x" was the shortest unique prefix of a long,
// but 'x' was also  defined as a short for a different flag.
func (fs *FlagSet) LookupLong(long string) *Flag {
	r, tail := FirstRune(long)
	if len(tail) == 0 {
		f := fs.LookupShort(r)
		if f != nil {
			return f
		}
	}

	f, err := fs.LongTrie.Get(long)
	if err != nil {
		return nil
	}
	return f
}

// Function `LookupShort()` returns a pointer to the `Flag`
// corresponding to the given rune or nil if none exists.
func (fs *FlagSet) LookupShort(r rune) *Flag {
	if r == NoShort {
		if fs.HasHyphenNumIdiom {
			return fs.ShortDict[NoShort]
		}
		return nil
	}
	if f, ok := fs.ShortDict[r]; ok {
		return f
	}
	return nil
}

// Function `Lookup()` takes either a string or a rune and looks it up
// in the `FlagSet` as a long or as a short option as appropriate,
// returning a pointer to the corresponding `Flag` if it exists or
// `nil` if it doesn't.
func (fs *FlagSet) Lookup(item interface{}) *Flag {
	if long, ok := item.(string); ok {
		return fs.LookupLong(long)
	}
	if short, ok := item.(rune); ok {
		return fs.LookupShort(short)
	}
	return nil
}

// Function `Lookup()` takes either a string or a rune and looks it up
// as a long or as a short option as appropriate in the default `FlagSet`.
func Lookup(item interface{}) *Flag {
	return CommandLine.Lookup(item)
}

// Function `AddFlag()` adds a flag to the default `FlagGroup` in a
// `FlagSet`.
func (fs *FlagSet) AddFlag(f *Flag) error {
	if f == nil {
		return fmt.Errorf("cannot add nil flag")
	}
	if !IsValidPair(f.Short, f.Long) {
		return fmt.Errorf("flag '%s' has invalid short/long flags", f)
	}
	if f.Long != NoLong {
		err := fs.LongTrie.Add(f.Long, f)
		if err != nil {
			return fmt.Errorf("error adding long flag '%s': %w", f.Long, err)
		}
	} 

	if g, ok := fs.ShortDict[f.Short]; f.Short != NoShort && ok {
		return fmt.Errorf("shortcut '%c' already used for '%s'", f.Short, g.Long)
	}
	if unicode.IsNumber(f.Short) {
		if fs.HasHyphenNumIdiom {
			return fmt.Errorf("cannot have '-%c' with -NUM idiom defined", f.Short)
		}
		fs.HasNumberShorts = true
	}
	if f.Short != NoShort || f.Long == NoLong {
		fs.ShortDict[f.Short] = f
	}
	if f.Short == NoShort && f.Long == NoLong {
		// -NUM idiom
		if fs.HasHyphenNumIdiom {
			return fmt.Errorf("cannot use -NUM idiom twice")
		}
		if fs.HasNumberShorts {
			return fmt.Errorf("cannot use -NUM idiom if digits are used as flags")
		}
		fs.HasHyphenNumIdiom = true
	}
	fs.Group().FlagList = append(fs.Group().FlagList, f)
	return nil
}


// Function `Var()` is the principal method of creating new flags in a
// `FlagSet`. Its usage is discussed at length in the introductory
// documentation.
func (fs *FlagSet) Var(value interface{}, short rune, long string, usage string, opts ...FlagOption) {
	// We need to set the parent flagset early because some of the
	// functions downstream of NewFlag() check that the flag doesn't
	// already exist in the parent flagset
	options := append([]FlagOption{WithParent(fs)}, opts...)
	f := NewFlag(value, short, long, usage, options...)
	if f == nil {
		log.Panicf("failed to create new flag -%c/--%s", short, long)
	}
	err := fs.AddFlag(f)
	if err != nil {
		log.Panicf("failed to add new flag '%s': %v", f, err)
	}
}

// Function `Var()` creates a new `Flag` in the current `FlagGroup` in
// the default `FlagSet`.
func Var(value interface{}, short rune, long string, usage string, opts ...FlagOption) {
	CommandLine.Var(value, short, long, usage, opts...)
}

// Function `Equ()` creates an equivalent to an extant flag,
// identified by a long option (`equiv`), with the given argument
// `value`. For example, `grep`'s `-I` is equivalent to
// `--binary-files=without-match`. I've never seen these kinds of
// shortcut equivalents defined in terms of short options, so the
// target flag is identified by the long option only.
func (fs *FlagSet) Equ(short rune, long string, equiv string, value string) {
	var f *Flag = nil
	f = fs.LookupLong(equiv)
	if f == nil {
		panic("flag not found in equivalent lookup")
	}
	a := f.NewAlias(short, long, withValue(value))
	err := fs.AddFlag(a)
	if err != nil {
		f.Failf("Error adding alias: %v", err)
	}
}

// Function `Equ()` creates an equivalent to an existing flag with a
// value (optarg) in the default `FlagSet`.
func Equ(short rune, long string, equiv string, value string) {
	CommandLine.Equ(short, long, equiv, value)
}

func (fs *FlagSet) Dump() {
	fmt.Fprintf(fs.Output, "Groups: %+v\n", fs.Groups)
	fmt.Fprintf(fs.Output, "GroupIndex: %+v\n", fs.GroupIndex)
	fmt.Fprintf(fs.Output, "LongTrie: %+v\n", fs.LongTrie)
	fmt.Fprintf(fs.Output, "ShortDict: %+v\n", fs.ShortDict)
	fmt.Fprintf(fs.Output, "Output: %+v\n", fs.Output)
	fmt.Fprintf(fs.Output, "IgnoreDoubleDash: %+v\n", fs.IgnoreDoubleDash)
	fmt.Fprintf(fs.Output, "InputArgs: %+v\n", fs.InputArgs)
	fmt.Fprintf(fs.Output, "OutputArgs: %+v\n", fs.OutputArgs)
	fmt.Fprintf(fs.Output, "OnFail: %+v\n", fs.OnFail)
	fmt.Fprintf(fs.Output, "FailExitCode: %+v\n", fs.FailExitCode)
}

func (fs *FlagSet) DumpFlags() {
	for _, g := range fs.Groups {
		fmt.Fprintf(fs.Output, "Group: %s\n", g.Title)
		for _, f := range g.FlagList {
			fmt.Fprintf(fs.Output, "\tFLAG: %s = %s\n", f, f.GetValue())
		}
	}
}

func (fs *FlagSet) DumpUsage() {
	fmt.Println(strings.Join(fs.AlignedFlagDescriptions("  ", "  ", ""), "\n"))
}

func (fs *FlagSet) Failf(format string, args ...interface{}) {
	if !fs.OnFail.TstSilentBit() {
		fmt.Fprintf(fs.Output, "ERROR: " + format + "\n", args...)
	}
	if fs.OnFail.TstContinueBit() {
		return
	}
	if fs.OnFail.TstPanicBit() {
		panic(fmt.Sprintf(format, args...))
	}
	os.Exit(fs.FailExitCode)
}

func (fs *FlagSet) Infof(format string, args ...interface{}) {
	if !fs.OnFail.TstSilentBit() {
		fmt.Fprintf(fs.Output, "INFO: " + format + "\n", args...)
	}
}

func (fs *FlagSet) Warnf(format string, args ...interface{}) {
	if !fs.OnFail.TstSilentBit() {
		fmt.Fprintf(fs.Output, "WARNING: " + format + "\n", args...)
	}
}

// Function `FlagStringMaxLen()` determines and returns the maximum
// length of any FlagString() in a `FlagSet` (without regard to
// `FlagGroup` membership). The flag string is a formatted
// representation of the long and/or short options for a flag used in
// help/usage output.
func (fs *FlagSet) FlagStringMaxLen() int {
	maxLen := 0
	for _, g := range fs.Groups {
		for _, f := range g.FlagList {
			maxLen = max(maxLen, len(f.FlagString()))
		}
	}
	return maxLen
}

// Function `AlignedFlagDescriptions()` returns a slice of
// similarly-formatted string descriptions of the `Flag`s in a
// `FlagSet`, separated by `FlagGroup` titles.
func (fs *FlagSet) AlignedFlagDescriptions(pre, mid, post string) []string {
	fstrs := []string{}
	maxl := fs.FlagStringMaxLen()
	for _, g := range fs.Groups {
		fstrs = append(fstrs, "\n" + g.Title + "\n")
		for _, f := range g.FlagList {
			s := fmt.Sprintf("%s%-*s%s%s%s", pre, maxl, f.FlagString(), mid, f.DescString(), post)
			fstrs = append(fstrs, s)
		}
	}
	return fstrs
}

// Function `Reset()` clears the input & output args, mutexes, and
// counts while keeping the flag setup. It mostly exists for testing.
func (fs *FlagSet) Reset() {
	// fmt.Fprintf(os.Stderr, "FlagSet has %d groups and %d mutexes\n", len(fs.Groups), len(fs.Mutex))
	fs.InputArgs.Clear()
	fs.OutputArgs.Clear()
	for name, _ := range fs.Mutex {
		fs.Mutex[name] = nil
	}
	for _, g := range fs.Groups {
		// fmt.Fprintf(os.Stderr, "Visiting group '%s' with %d flags\n", g.Title, len(g.FlagList))
		for _, f := range g.FlagList {
			// fmt.Fprintf(os.Stderr, "Clearing flag '%s'\n", f)
			f.Count = 0
		}
	}
}
