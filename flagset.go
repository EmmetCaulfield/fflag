package fflag

import(
	"fmt"
	"io"
	"os"
	"strings"
	
	"github.com/EmmetCaulfield/fflag/pkg/deque"
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



// CommandLine is the default FlagSet, named by analogy with `pflag`
var CommandLine *FlagSet = NewFlagSet()

type FlagGroup struct {
	Title string
	FlagList          []*Flag
}

func NewFlagGroup(title string) *FlagGroup {
	return &FlagGroup{
		Title: title,
		FlagList: []*Flag{},
	}
}

type FlagSet struct {
	Groups             []*FlagGroup
	GroupIndex         int
	LongDict           map[string]*Flag
	ShortDict          map[rune]*Flag
	Output             io.Writer
	IgnoreDoubleDash   bool
	HasHyphenNumIdiom  bool
	InputArgs         *deque.Deque[string]
	OutputArgs        *deque.Deque[string]
	OnFail             FailOption
	FailExitCode       int
	OnFileError        FailOption
	FileErrExitCode    int
	Mutex              map[string]*Flag
}

var DefaultFailExitCode int = 2
var DefaultFileErrExitCode int = 2


func (fs *FlagSet) NewFlagGroup(title string) *FlagGroup {
	fg := NewFlagGroup(title)
	fs.Groups = append(fs.Groups, fg)
	fs.GroupIndex = len(fs.Groups) - 1
	return fg
}

func (fs *FlagSet) Group() *FlagGroup {
	// Happy for this to panic
	return fs.Groups[fs.GroupIndex]
}

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

type FlagSetOption = func (fs *FlagSet)

func NewFlagSet(opts ...FlagSetOption) *FlagSet {
	fs := &FlagSet {
		Groups: []*FlagGroup{
			NewFlagGroup("Options"),
		},
		GroupIndex:       0,
		LongDict:         map[string]*Flag{},
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

func WithGroupTitle(title string) FlagSetOption {
	return func(fs *FlagSet) {
		fs.Groups[0].Title = title
	}
}

func WithOutputWriter(w io.Writer) FlagSetOption {
	return func(fs *FlagSet) {
		fs.Output = w
	}
}

func WithPanicOnFail() FlagSetOption {
	return func(fs *FlagSet) {
		fs.OnFail.SetPanicBit()
	}
}

func WithContinueOnFail() FlagSetOption {
	return func(fs *FlagSet) {
		fs.OnFail.SetContinueBit()
	}
}

func WithSilentFail() FlagSetOption {
	return func(fs *FlagSet) {
		fs.OnFail.SetSilentBit()
	}
}

func IgnoringDoubleDash() FlagSetOption {
	return func(fs *FlagSet) {
		fs.IgnoreDoubleDash = true
	}
}

// HasFlags returns a bool to indicate if the FlagSet has any flags
// defined.
func (fs *FlagSet) HasFlags() bool {
	n := 0
	for _, g := range fs.Groups {
		n += len(g.FlagList)
	}
	return n > 0
}

func (fs *FlagSet) LookupLong(long string) *Flag {
	if len(long) == 0 {
		if fs.HasHyphenNumIdiom {
			return fs.LongDict[""]
		}
		return nil
	}
	if len(long) == 1 {
		return fs.LookupShort(rune(long[0]))
	}
	if f, ok := fs.LongDict[long]; ok {
		return f
	}
	return nil
}

// LookupShort returns the Flag structure of the shortcut flag,
// returning nil if none exists.
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

func (fs *FlagSet) Lookup(item interface{}) *Flag {
	if long, ok := item.(string); ok {
		return fs.LookupLong(long)
	}
	if short, ok := item.(rune); ok {
		return fs.LookupShort(short)
	}
	return nil
}

func Lookup(long string) *Flag {
	return CommandLine.Lookup(long)
}

func (fs *FlagSet) AddFlag(f *Flag) error {
	if f == nil {
		return fmt.Errorf("cannot add nil flag")
	}
	if !IsValidPair(f.Short, f.Long) {
		return fmt.Errorf("flag '%s' has invalid short/long flags", f)
	}
	if g, ok := fs.LongDict[f.Long]; f.Long != NoLong && ok {
		return fmt.Errorf("long flag '%s' already used for '%s'", f.Long, g)
	}
	if f.Long != NoLong || f.Short == NoShort {
		fs.LongDict[f.Long] = f
	} 

	if g, ok := fs.ShortDict[f.Short]; f.Short != NoShort && ok {
		return fmt.Errorf("shortcut '%c' already used for '%s'", f.Short, g.Long)
	}
	if f.Short != NoShort || f.Long == NoLong {
		fs.ShortDict[f.Short] = f
	}

	fs.Group().FlagList = append(fs.Group().FlagList, f)
	return nil
}

func (fs *FlagSet) Var(value interface{}, short rune, long string, usage string, opts ...FlagOption) {
	// We need to set the parent flagset early because some of the
	// functions downstream of NewFlag() check that the flag doesn't
	// already exist in the parent flagset
	options := append([]FlagOption{WithParent(fs)}, opts...)
	f := NewFlag(value, short, long, usage, options...)
	if f == nil {
		fs.Failf("Failed to create new flag %s", long)
	}
	err := fs.AddFlag(f)
	if err != nil {
		fs.Failf("Failed to add new flag %s, %#v: %v", long, f, err)
	}
}

func Var(value interface{}, short rune, long string, usage string, opts ...FlagOption) {
	CommandLine.Var(value, short, long, usage, opts...)
}

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

func Equ(short rune, long string, equiv string, value string) {
	CommandLine.Equ(short, long, equiv, value)
}

func (fs *FlagSet) Dump() {
	fmt.Fprintf(fs.Output, "Groups: %+v\n", fs.Groups)
	fmt.Fprintf(fs.Output, "GroupIndex: %+v\n", fs.GroupIndex)
	fmt.Fprintf(fs.Output, "LongDict: %+v\n", fs.LongDict)
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


func (fs *FlagSet) FlagStringMaxLen() int {
	maxLen := 0
	for _, g := range fs.Groups {
		for _, f := range g.FlagList {
			maxLen = max(maxLen, len(f.FlagString()))
		}
	}
	return maxLen
}

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
