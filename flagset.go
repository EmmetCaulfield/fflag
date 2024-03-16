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
	if f, ok := fs.LabelDict[label]; ok {
		return f
	}
	return nil
}

// LookupShortcut returns the Flag structure of the shortcut flag,
// returning nil if none exists.
func (fs *FlagSet) LookupShortcut(r rune) *Flag {
	if f, ok := fs.LetterDict[r]; ok {
		return f
	}
	return nil
}

func (fs *FlagSet) Lookup(item interface{}) *Flag {
	if label, ok := item.(string); ok {
		if len(label) == 0 {
			return nil
		}
		if len(label) == 1 {
			return fs.LookupShortcut(rune(label[0]))
		}
		return fs.LookupLabel(label)
	}
	if letter, ok := item.(rune); ok {
		if letter == rune(0) {
			return nil
		}
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
	if !IsValidShortcut(f.Letter) {
		if !f.IsLabelAlias() {
			return fmt.Errorf("shortcut '%c' for '%s' is invalid", f.Letter, f.Label)
		}
		f.Letter = rune(0)
	}

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
	f.ParentFlagSet = fs
	return nil
}

func (fs *FlagSet) Var(value interface{}, label string, usage string, opts ...FlagOption) {
	f := NewFlag(value, label, usage, opts...)
	_ = fs.AddFlag(f)
}

func Var(value interface{}, label string, usage string, opts ...FlagOption) {
	CommandLine.Var(value, label, usage, opts...)
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
