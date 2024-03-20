package fflag

import(
	"fmt"
	"io"
	"os"
	
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

type FlagSet struct {
	FlagList          []*Flag
	LabelDict         map[string]*Flag
	LetterDict        map[rune]*Flag
	IsSorted          bool
	Output            io.Writer
	IgnoreDoubleDash  bool
	InputArgs        *deque.Deque[string]
	OutputArgs       *deque.Deque[string]
	OnFail            FailOption
	FailExitCode      int
}

type FlagSetOption = func (fs *FlagSet)

func NewFlagSet(opts ...FlagSetOption) *FlagSet {
	fs := &FlagSet {
		FlagList:         []*Flag{},
		LabelDict:        map[string]*Flag{},
		LetterDict:       map[rune]*Flag{},
		IsSorted:         false,
		Output:           os.Stderr,
		IgnoreDoubleDash: false,
		InputArgs:        &deque.Deque[string]{},
		OutputArgs:       &deque.Deque[string]{},
		OnFail:           FailDefault,
		FailExitCode:     2,
	}
	for _, opt := range opts {
		opt(fs)
	}
	return fs
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
	return len(fs.FlagList) > 0
}

func (fs *FlagSet) LookupLabel(label string) *Flag {
	if len(label) == 0 {
		return nil
	}
	if len(label) == 1 {
		return fs.LookupShortcut(rune(label[0]))
	}
	if f, ok := fs.LabelDict[label]; ok {
		return f
	}
	return nil
}

// LookupShortcut returns the Flag structure of the shortcut flag,
// returning nil if none exists.
func (fs *FlagSet) LookupShortcut(r rune) *Flag {
	if r == rune(0) {
		return nil
	}
	if f, ok := fs.LetterDict[r]; ok {
		return f
	}
	return nil
}

func (fs *FlagSet) Lookup(item interface{}) *Flag {
	if label, ok := item.(string); ok {
		return fs.LookupLabel(label)
	}
	if letter, ok := item.(rune); ok {
		return fs.LookupShortcut(letter)
	}
	return nil
}

func Lookup(label string) *Flag {
	return CommandLine.Lookup(label)
}

func (fs *FlagSet) AddFlag(f *Flag) error {
	if f == nil {
		return fmt.Errorf("cannot add nil flag")
	}
	if !IsValidLabel(f.Label) {
		if !IsValidShortcut(f.Letter) {
			return fmt.Errorf("flag has neither a label nor a shortcut letter")
		}
		if !f.IsShortcutAlias() {
			return fmt.Errorf("a label is required except for shortcut aliases")
		}
		f.Label = ""
	}
	if f.Letter != rune(0) && !IsValidShortcut(f.Letter) {
		if !f.IsLabelAlias() {
			return fmt.Errorf("shortcut '%c' for '%s' is invalid", f.Letter, f.Label)
		}
		f.Letter = rune(0)
	}
	fs.Infof("Adding: %s, %c", f.Label, f.Letter)
	if len(f.Label) > 0 {
		if g, ok := fs.LabelDict[f.Label]; ok {
			return fmt.Errorf("shortcut '%c' already used for '%s'", f.Letter, g.Label)		
		}
		fs.LabelDict[f.Label] = f
	}
	if f.Letter != rune(0) {
		if g, ok := fs.LetterDict[f.Letter]; ok {
			return fmt.Errorf("shortcut '%c' already used for '%s'", f.Letter, g.Label)		
		}
		fs.LetterDict[f.Letter] = f
	}

	fs.FlagList = append(fs.FlagList, f)
	fs.IsSorted = false
	return nil
}

func (fs *FlagSet) Var(value interface{}, label string, usage string, opts ...FlagOption) {
	// We need to set the parent flagset early because some of the
	// functions downstream of NewFlag() check that the flag doesn't
	// already exist in the parent flagset
	options := append([]FlagOption{WithParent(fs)}, opts...)
	f := NewFlag(value, label, usage, options...)
	if f == nil {
		fs.Failf("Failed to create new flag %s", label)
	}
	err := fs.AddFlag(f)
	if err != nil {
		fs.Failf("Failed to add new flag %s, %#v: %v", label, f, err)
	}
}

func Var(value interface{}, label string, usage string, opts ...FlagOption) {
	CommandLine.Var(value, label, usage, opts...)
}

func (fs *FlagSet) Dump() {
	fmt.Fprintf(fs.Output, "FlagList: %+v\n", fs.FlagList)
	fmt.Fprintf(fs.Output, "LabelDict: %+v\n", fs.LabelDict)
	fmt.Fprintf(fs.Output, "LetterDict: %+v\n", fs.LetterDict)
	fmt.Fprintf(fs.Output, "IsSorted: %+v\n", fs.IsSorted)
	fmt.Fprintf(fs.Output, "Output: %+v\n", fs.Output)
	fmt.Fprintf(fs.Output, "IgnoreDoubleDash: %+v\n", fs.IgnoreDoubleDash)
	fmt.Fprintf(fs.Output, "InputArgs: %+v\n", fs.InputArgs)
	fmt.Fprintf(fs.Output, "OutputArgs: %+v\n", fs.OutputArgs)
	fmt.Fprintf(fs.Output, "OnFail: %+v\n", fs.OnFail)
	fmt.Fprintf(fs.Output, "FailExitCode: %+v\n", fs.FailExitCode)
}

func (fs *FlagSet) DumpFlags() {
	for _, f := range fs.FlagList {
		char := f.Letter
		if char == rune(0) {
			fmt.Fprintf(fs.Output, "FLAG: --%s = %s\n", f.Label, f.GetValue())
		}
		fmt.Fprintf(fs.Output, "FLAG: -%c/--%s = %s\n", f.Letter, f.Label, f.GetValue())
	}
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
	fmt.Fprintf(fs.Output, "INFO: " + format + "\n", args...)
}


func (fs *FlagSet) FlagStringMaxLen() int {
	maxLen := 0
	for _, f := range fs.FlagList {
		maxLen = max(maxLen, len(f.FlagString()))
	}
	return maxLen
}

func (fs *FlagSet) AlignedFlagDescriptions(pre, mid, post string) []string {
	fstrs := []string{}
	maxl := fs.FlagStringMaxLen()
	for _, f := range fs.FlagList {
		s := fmt.Sprintf("%s%-*s%s%s%s", pre, maxl, f.FlagString(), mid, f.DescString(), post)
		fstrs = append(fstrs, s)
	}
	return fstrs
}
