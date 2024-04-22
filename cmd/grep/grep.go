package main

import (
	"fmt"
	"strings"

	"github.com/EmmetCaulfield/fflag"
)

// Usage: grep [OPTION]... PATTERNS [FILE]...
// Search for PATTERNS in each FILE.
// Example: grep -i 'hello world' menu.h main.c
// PATTERNS can contain multiple patterns separated by newlines.

type OptStruct struct {
	ExtendedRegexp       bool
	FixedStrings         bool
	BasicRegexp          bool
	PerlRegexp           bool
	Regexp               []string
	File                 string
	IgnoreCase           bool
	NoIgnoreCase         bool
	WordRegexp           bool
	LineRegexp           bool
	NullData             bool
	NoMessages           bool
	InvertMatch          bool
	Version              bool
	Help                 bool
	MaxCount             int
	ByteOffset           bool
	LineNumber           bool
	LineBuffered         bool
	WithFilename         bool
	NoFilename           bool
	Label                string
	OnlyMatching         bool
	Quiet                bool
	BinaryFiles          string
	Text                 bool
	Directories          string
	Devices              string
	Recursive            bool
	DereferenceRecursive bool
	Include              string   // glob
	Exclude              []string // glob
	// ExcludeFrom          string
	ExcludeDir        string // glob
	FilesWithoutMatch bool
	FilesWithMatches  bool
	Count             uint
	InitialTab        bool
	Null              bool
	BeforeContext     uint
	AfterContext      uint
	Context           uint
	GroupSeparator    string
	NoGroupSeparator  bool
	Color             string
	Binary            bool
	// -NUM
}

func (o *OptStruct) Dump() {
	fields := strings.Split(fmt.Sprintf("%+v", *o), " ")
	fields[0] = fields[0][1:]
	end := len(fields) - 1
	fields[end] = fields[end][:len(fields[end])-1]
	maxColonPosition := 0
	for _, field := range fields {
		maxColonPosition = max(maxColonPosition, strings.Index(field, ":"))
	}
	maxColonPosition++
	fmt.Println("{")
	for _, field := range fields {
		kvp := strings.Split(field, ":")
		if len(kvp) == 1 {
			fmt.Printf("\t%*s   %s\n", maxColonPosition, "", kvp[0])
		} else {
			fmt.Printf("\t%-*s: %s\n", maxColonPosition, kvp[0], kvp[1])
		}
	}
	fmt.Println("}")
}

func ValidateRegex(f *fflag.Flag, arg string, pos int) error {
	fmt.Printf("validating: '%s', '%s', %d\n", f, arg, pos)
	f.SetOnly(arg, pos)
	return nil
}

func setup() *OptStruct {
	opt := &OptStruct{}
	fflag.Group("Pattern selection and interpretation")
	fflag.Var(&opt.ExtendedRegexp, 'E', "extended-regexp", "PATTERNS are extended regular expressions",
		fflag.InMutex("pat-type"))
	fflag.Var(&opt.FixedStrings, 'F', "fixed-strings", "PATTERNS are strings",
		fflag.InMutex("pat-type"))
	fflag.Var(&opt.BasicRegexp, 'G', "basic-regexp", "PATTERNS are basic regular expressions",
		fflag.InMutex("pat-type"))
	fflag.Var(&opt.PerlRegexp, 'P', "perl-regexp", "PATTERNS are Perl regular expressions",
		fflag.InMutex("pat-type"))
	fflag.Var(&opt.Regexp, 'e', "regexp", "use PATTERNS for matching", fflag.WithTypeTag("PATTERNS"),
		fflag.WithCallback(ValidateRegex))
	fflag.Var(&opt.Regexp, 'f', "file", "take PATTERNS from FILE", fflag.WithTypeTag("FILE"),
		fflag.ReadFile(), fflag.WithCallback(ValidateRegex))
	fflag.Var(&opt.IgnoreCase, 'i', "ignore-case", "ignore case distinctions in patterns and data",
		fflag.InMutex("case"))
	fflag.Var(&opt.NoIgnoreCase, fflag.NoShort, "no-ignore-case", "do not ignore case distinctions (default)",
		fflag.InMutex("case"))
	fflag.Var(&opt.WordRegexp, 'w', "word-regexp", "match only whole words",
		fflag.InMutex("word/line"))
	fflag.Var(&opt.LineRegexp, 'x', "line-regexp", "match only whole lines",
		fflag.InMutex("word/line"))
	fflag.Var(&opt.NullData, 'z', "null-data", "a data line ends in 0 byte, not newline")

	fflag.Group("Miscellaneous")
	fflag.Var(&opt.NoMessages, 's', "no-messages", "suppress error messages")
	fflag.Var(&opt.InvertMatch, 'v', "invert-match", "select non-matching lines")
	fflag.Var(&opt.Version, 'V', "version", "display version information and exit")
	fflag.Var(&opt.Help, fflag.NoShort, "help", "display this help text and exit")

	fflag.Group("Output control")
	fflag.Var(&opt.MaxCount, 'm', "max-count", "stop after NUM selected lines", fflag.WithTypeTag("NUM"))
	fflag.Var(&opt.ByteOffset, 'b', "byte-offset", "print the byte offset with output lines")
	fflag.Var(&opt.LineNumber, 'n', "line-number", "print line number with output lines")
	fflag.Var(&opt.LineBuffered, fflag.NoShort, "line-buffered", "flush output on every line")
	fflag.Var(&opt.WithFilename, 'H', "with-filename", "print file name with output lines",
		fflag.InMutex("with/no-filename"))
	fflag.Var(&opt.NoFilename, 'h', "no-filename", "suppress the file name prefix on output",
		fflag.InMutex("with/no-filename"))
	fflag.Var(&opt.Label, fflag.NoShort, "label", "use LABEL as the standard input file name prefix", fflag.WithTypeTag("LABEL"))
	fflag.Var(&opt.OnlyMatching, 'o', "only-matching", "show only nonempty parts of lines that match")
	fflag.Var(&opt.Quiet, 'q', "quiet", "suppress all normal output",
		fflag.WithAlias(fflag.NoShort, "silent", false))
	fflag.Var(&opt.BinaryFiles, fflag.NoShort, "binary-files", "assume that binary files are TYPE;", fflag.WithTypeTag("TYPE"),
		fflag.WithDefault([]string{"binary", "text", "without-match"}))
	fflag.Equ('a', "text", "binary-files", "text")
	fflag.Equ('I', fflag.NoLong, "binary-files", "without-match")
	fflag.Var(&opt.Directories, 'd', "directories", "how to handle directories;", fflag.WithTypeTag("ACTION"),
		fflag.WithDefault([]string{"read", "recurse", "skip"}))
	fflag.Var(&opt.Devices, 'D', "devices", "how to handle devices, FIFOs and sockets;", fflag.WithTypeTag("ACTION"),
		fflag.WithDefault([]string{"read", "skip"}))
	fflag.Equ('r', "recursive", "directories", "recurse")
	fflag.Var(&opt.DereferenceRecursive, 'R', "dereference-recursive", "likewise, but follow all symlinks")
	fflag.Var(&opt.Include, fflag.NoShort, "include", "search only files that match GLOB (a file pattern)", fflag.WithTypeTag("GLOB"))
	fflag.Var(&opt.Exclude, fflag.NoShort, "exclude", "skip files that match GLOB", fflag.WithTypeTag("GLOB"))
	fflag.Var(&opt.Exclude, fflag.NoShort, "exclude-from", "skip files that match any file pattern from FILE",
		fflag.WithTypeTag("FILE"), fflag.ReadFile())
	fflag.Var(&opt.ExcludeDir, fflag.NoShort, "exclude-dir", "skip directories that match GLOB", fflag.WithTypeTag("GLOB"))
	fflag.Var(&opt.FilesWithoutMatch, 'L', "files-without-match", "print only names of FILEs with no selected lines",
		fflag.InMutex("with(out)-match"))
	fflag.Var(&opt.FilesWithMatches, 'l', "files-with-matches", "print only names of FILEs with selected lines",
		fflag.InMutex("with(out)-match"))
	fflag.Var(&opt.Count, 'c', "count", "print only a count of selected lines per FILE", fflag.AsCounter())
	fflag.Var(&opt.InitialTab, 'T', "initial-tab", "make tabs line up (if needed)")
	fflag.Var(&opt.Null, 'Z', "null", "print 0 byte after FILE name")

	fflag.Group("Context control")
	fflag.Var(&opt.BeforeContext, 'B', "before-context", "print NUM lines of leading context", fflag.WithTypeTag("NUM"))
	fflag.Var(&opt.AfterContext, 'A', "after-context", "print NUM lines of trailing context", fflag.WithTypeTag("NUM"))
	fflag.Var(&opt.Context, 'C', "context", "print NUM lines of output context", fflag.WithTypeTag("NUM"),
		fflag.WithAlias(fflag.NoShort, fflag.NoLong, false))
	fflag.Var(&opt.GroupSeparator, fflag.NoShort, "group-separator", "print SEP on line between matches with context",
		fflag.WithTypeTag("SEP"), fflag.InMutex("group-sep"))
	fflag.Var(&opt.NoGroupSeparator, fflag.NoShort, "no-group-separator", "do not print separator for matches with context",
		fflag.InMutex("group-sep"))
	fflag.Var(&opt.Color, fflag.NoShort, "color", "use markers to highlight the matching strings;", fflag.WithTypeTag("WHEN"),
		fflag.WithOptionalDefault([]string{"always", "never", "auto"}), fflag.WithAlias(fflag.NoShort, "colour", false))
	fflag.Var(&opt.Binary, 'U', "binary", "do not strip CR characters at EOL (MSDOS/Windows)")

	// When FILE is '-', read standard input.  With no FILE, read '.' if
	// recursive, '-' otherwise.  With fewer than two FILEs, assume -h.
	// Exit status is 0 if any line is selected, 1 otherwise;
	// if any error occurs and is not given, the exit status is 2.
	return opt
}

func main() {
	opt := setup()
	fflag.Parse()
	opt.Dump()
	fflag.CommandLine.DumpUsage()
}
