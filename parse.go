// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fflag

import (
	"fmt"
	"os"
	"strings"
)

// ErrorHandling defines how to handle flag parsing errors.
type ErrorHandling int

const (
	// ContinueOnError will return an err from Parse() if an error is found
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) if an error is found when parsing
	ExitOnError
	// PanicOnError will panic() if an error is found when parsing flags
	PanicOnError
)

// ParseErrorsWhitelist defines the parsing errors that can be ignored
type ParseErrorsWhitelist struct {
	// UnknownFlags will ignore unknown flags errors and continue parsing rest of the flags
	UnknownFlags bool
}

//--unknown (args will be empty)
//--unknown --next-flag ... (args will be --next-flag ...)
//--unknown arg ... (args will be arg ...)
func stripUnknownFlagValue(args []string) []string {
	if len(args) == 0 {
		//--unknown
		return args
	}

	first := args[0]
	if len(first) > 0 && first[0] == '-' {
		//--unknown --next-flag ...
		return args
	}

	//--unknown arg ... (args will be arg ...)
	if len(args) > 1 {
		return args[1:]
	}
	return nil
}

// A flag argument can be:
//
//   * Prohibited
//   * Optional
//   * Required
//
// There are 9 valid possibilities:
//
//   * `--flag`        prohibited/optional
//   * `--flag=value`  optional/required
//   * `--flag value`  optional/required
//   * `-f`            prohibited/optional
//   * `-f=value`      optional/required
//   * `-f value`      optional/required
//   * `-fgh`          -f, -g, -h : prohibited/optional
//   * `-fgh=foo`      -f, -g : prohibited/optional, -h : optional/required
//   * `-fgh foo`      -f, -g : prohibited/optional, -h : optional/required
//
// Each of these can occur:
//
//   * at the end of the command line
//   * followed by another long flag argument
//   * followed by a shortcut flag
//   * followed by a non-flag argument

type ArgMask int8

const (
	AMClear      ArgMask = 0
	AMFlagBit            = 0b00000001 // The argument is a flag (starts with one or two hyphens)
	AMLongBit            = 0b00000010 // The argument is a long flag (starts with two hyphens)
	AMClusterBit         = 0b00000100 // The argument is a cluster of short flags
	AMParamBit           = 0b00001000 // The argument has a parameter (--flag=param)
	AMHyphenBit          = 0b00010000 // The argument is just hyphens ("-" or "--")
)

func (am *ArgMask) String() string {
	return fmt.Sprintf("%08b", *am)
}

func (am *ArgMask) SetFlagBit()         { *am = *am | AMFlagBit }
func (am *ArgMask) SetLongBit()         { *am = *am | AMLongBit }
func (am *ArgMask) SetClusterBit()      { *am = *am | AMClusterBit }
func (am *ArgMask) SetParamBit()        { *am = *am | AMParamBit }
func (am *ArgMask) SetHyphenBit()       { *am = *am | AMHyphenBit }
func (am *ArgMask) ClrFlagBit()         { *am = *am & ^AMFlagBit }
func (am *ArgMask) ClrLongBit()         { *am = *am & ^AMLongBit }
func (am *ArgMask) ClrClusterBit()      { *am = *am & ^AMClusterBit }
func (am *ArgMask) ClrParamBit()        { *am = *am & ^AMParamBit }
func (am *ArgMask) ClrHyphenBit()       { *am = *am & ^AMHyphenBit }
func (am *ArgMask) TstFlagBit() bool    { return *am&AMFlagBit != 0 }
func (am *ArgMask) TstLongBit() bool    { return *am&AMLongBit != 0 }
func (am *ArgMask) TstClusterBit() bool { return *am&AMClusterBit != 0 }
func (am *ArgMask) TstParamBit() bool   { return *am&AMParamBit != 0 }
func (am *ArgMask) TstHyphenBit() bool  { return *am&AMHyphenBit != 0 }

// Tests is a flag mask represents any kind of flag
func (am *ArgMask) IsFlag() bool {
	if am.TstFlagBit() {
		// Sanity check: if it's a flag, it can't be just hyphens, but
		// it can be short or long and with or without a parameter
		if am.TstHyphenBit() {
			panic("flag purports to be just hyphens")
		}
		return true
	}
	return false
}

// Tests if a flag mask represents a single short flag
func (am *ArgMask) IsShortFlag() bool {
	if am.IsFlag() && !am.TstLongBit() {
		// No additional sanity-check to do here
		return !am.TstClusterBit()
	}
	return false
}

// Tests if a flag mask represents a single long flag
func (am *ArgMask) IsLongFlag() bool {
	if am.IsFlag() && am.TstLongBit() {
		// Sanity check: there's no such thing as a cluster of long
		// flags
		if am.TstClusterBit() {
			panic("fatal: purported cluster of long flags")
		}
		return true
	}
	return false
}

// Tests if a flag mask represents a cluster of short flags
func (am *ArgMask) IsCluster() bool {
	if am.IsFlag() && am.TstClusterBit() {
		// Sanity check: if the cluster bit is set, it must be a flag,
		// short, and not just hyphens
		if am.TstLongBit() {
			panic("fatal: purported cluster of long flags")
		}
		return true
	}
	return false
}

func (am *ArgMask) HasParam() bool {
	if am.TstParamBit() {
		// If the param bit is set, it must be a flag and not just
		// hyphens
		if !am.TstFlagBit() {
			panic("fatal: purported non-flag with parameter")
		}
		if am.TstHyphenBit() {
			panic("fatal: purported hyphens with parameter")
		}
		return true
	}
	return false
}

func (am *ArgMask) IsNonFlag() bool {
	if !am.IsFlag() {
		// If the flag bit is clear, it may be just hyphens and short
		// (single hyphen) or long (double hyphen), but can't have a
		// parameter or be a cluster.
		if am.TstClusterBit() {
			panic("fatal: purported non-flag cluster")
		}
		if am.TstParamBit() {
			panic("fatal: purported non-flag with parameter")
		}
		return true
	}
	return false
}

func parseSingleArg(arg string) (flags []string, param string, argType ArgMask) {
	// minimum length flag is 2 (e.g. "-x")
	if len(arg) < 2 {
		// argType.ClrLongBit()
		if arg == "-" {
			argType.SetHyphenBit()
			// argType.ClrFlagBit()
		}
		// flags = []string{}
		param = arg
		return
	}

	var flag string
	if arg[0:2] == "--" {
		// A long flag
		argType.SetLongBit()
		flag = arg[2:len(arg)]
	} else if arg[0] == '-' {
		// A short flag
		argType.ClrLongBit()
		flag = arg[1:len(arg)]
	} else {
		// Not a flag, must be a param
		// argType.ClrFlagBit()
		// flags = []string{}
		param = arg
		return
	}

	// If we get here, we have a flag or cluster of flags (stripped of
	// leading hyphens) in `flag`, potentially with a parameter
	if len(flag) == 0 {
		// There's nothing after the flag start indicator
		argType.SetHyphenBit()
		if argType.TstLongBit() {
			// flags = []string{}
			param = "--"
			return
		}
		panic("impossible condition")
	}

	argType.SetFlagBit()
	parts := strings.SplitN(flag, "=", 2)
	if len(parts) == 2 {
		argType.SetParamBit()
		flag = parts[0]
		param = parts[1]
	} else {
		argType.ClrParamBit()
	}

	if argType.IsLongFlag() {
		// We have a single long flag
		flags = []string{flag}
		return
	}

	// It could be a single short flag or a cluster of short flags
	argType.ClrLongBit()
	if len(flag) == 1 {
		// We have a single short flag
		flags = []string{flag}
		return
	}

	// We have a cluster of short flags
	argType.SetClusterBit()
	flags = strings.Split(flag, "")

	return
}

func (fs *FlagSet) parse() error {
	var err error
	var i int = 0
NEXTARG:
	for arg, err := fs.InputArgs.Shift(); err == nil; arg, err = fs.InputArgs.Shift() {
		i++
		flags, param, argType := parseSingleArg(arg)
		if !argType.IsFlag() {
			fs.OutputArgs.Push(param)
			continue NEXTARG
		}
		var flag *Flag = nil
		if argType.IsCluster() {
			for _, s := range flags[:len(flags)-1] {
				flag = fs.Lookup(s)
				if flag == nil {
					fs.Failf("short flag '%s' in cluster '%s' is not defined", s, arg)
					continue
				}
				err = flag.Set(nil)
				if err != nil {
					fs.Failf("failed to set short flag '%s' in cluster '%s'", s, arg)
					continue
				}
				continue NEXTARG
			}
			flag = fs.Lookup(flags[len(flags)-1])
			if flag == nil {
				fs.Failf("shortcut '%s' not defined in cluster '%s'", flags[len(flags)-1], arg)
				continue
			}
		} else {
			flag = fs.Lookup(flags[0])
			if flag == nil {
				fs.Failf("flag '%s' not defined", flags[0])
				continue
			}
		}
		if argType.HasParam() {
			err = flag.Set(param)
			if err != nil {
				fs.Failf("failed to set flag `%s` with '%s'", flag.String(), param)
			}
			continue NEXTARG
		}
		// Peek at the next argument
		next, err := fs.InputArgs.Front()
		if err != nil {
			// End of InputArgs
			err = flag.Set(nil)
			if err != nil {
				fs.Failf("failed to set flag `%s` at EOL with no parameter", flag.String())
			}
			continue NEXTARG
		}
		// Have next arg, might be a parameter
		flags, param, nextArgType := parseSingleArg(next)
		if !nextArgType.IsFlag() {
			// Not a flag, try it as a parameter
			err = flag.Set(param)
			if err != nil {
				fs.Failf("failed to set flag `%s` with '%s'", flag.String(), param)
			}
			continue NEXTARG
		}
		// Next arg is a flag, current flag has no parameter
		err = flag.Set(nil)
		if err != nil {
			fs.Failf("failed to set flag `%s` with no parameter", flag.String())
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (fs *FlagSet) Parse(arguments []string) error {
	fs.InputArgs.Init(arguments...)
	return fs.parse()
}

func Parse() {
	CommandLine.Parse(os.Args[1:])
}
