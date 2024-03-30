// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fflag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// A flag argument can be:
//
//   * Prohibited
//   * Optional
//   * Required
//
// There are several valid possibilities:
//
//   * `--flag`          --flag()       : prohibited/optional
//   * `--flag=value`    --flag(value)  : optional/required
//   * `--flag value`    --flag(value)  : required/optional
//   * `--flag operand`  --flag()       : prohibited/optional
//   * `-f`              -f()           : prohibited/optional
//   * `-f opd`          -f()           : prohibited/optional
//   * `-f=arg`          -f(=arg)       : optional/required (POSIX)
//   * `-f=arg`          -f(arg)        : optional/required (GNU-ish de-facto, non-POSIX)
//   * `-f arg`          -f(arg)        : optional/required
//   * `-farg`           -f(arg)        : optional/required
//   * `-fgh`            -f() -g() -h() :
//   * `-fgh`            -f() -g(h)     :
//   * `-fgh`            -f(gh)         :
//   * `-fgh=arg`        -f() -g() -h(=arg)
//   * `-fgh=arg`        -f() -g() -h(arg)
//   * `-fgh=arg`        -f() -g(h=arg)
//   * `-fgh=arg`        -f(gh=arg)
//   * `-fgh arg`        -f() -g() -h(arg)
//   * `-fgh arg`        -f(g) -h(arg)
//   * `-fgh opd`        -f() -g() -h()
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
	AMParamBit           = 0b00001000 // The argument has an attached parameter (--flag=param)
	AMHyphenBit          = 0b00010000 // The argument is just hyphens ("-" or "--")
	AMNumberBit          = 0b00100000 // The argument is a number
)

func (am *ArgMask) String() string {
	return fmt.Sprintf("%08b", *am)
}

func (am *ArgMask) SetFlagBit()         { *am = *am | AMFlagBit }
func (am *ArgMask) SetLongBit()         { *am = *am | AMLongBit }
func (am *ArgMask) SetClusterBit()      { *am = *am | AMClusterBit }
func (am *ArgMask) SetParamBit()        { *am = *am | AMParamBit }
func (am *ArgMask) SetHyphenBit()       { *am = *am | AMHyphenBit }
func (am *ArgMask) SetNumberBit()       { *am = *am | AMNumberBit }
func (am *ArgMask) ClrFlagBit()         { *am = *am & ^AMFlagBit }
func (am *ArgMask) ClrLongBit()         { *am = *am & ^AMLongBit }
func (am *ArgMask) ClrClusterBit()      { *am = *am & ^AMClusterBit }
func (am *ArgMask) ClrParamBit()        { *am = *am & ^AMParamBit }
func (am *ArgMask) ClrHyphenBit()       { *am = *am & ^AMHyphenBit }
func (am *ArgMask) ClrNumberBit()       { *am = *am & ^AMNumberBit }
func (am *ArgMask) TstFlagBit() bool    { return *am&AMFlagBit != 0 }
func (am *ArgMask) TstLongBit() bool    { return *am&AMLongBit != 0 }
func (am *ArgMask) TstClusterBit() bool { return *am&AMClusterBit != 0 }
func (am *ArgMask) TstParamBit() bool   { return *am&AMParamBit != 0 }
func (am *ArgMask) TstHyphenBit() bool  { return *am&AMHyphenBit != 0 }
func (am *ArgMask) TstNumberBit() bool  { return *am&AMNumberBit != 0 }

// Tests if an argument mask represents any kind of flag
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

// Tests if an argument mask represents a double-hyphen
func (am *ArgMask) IsDoubleHyphen() bool {
	return am.TstHyphenBit() && am.TstLongBit()
}

// Tests if an argument mask represents a single short flag
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

// Tests if a flag mask represents a flag that is itself a number
func (am *ArgMask) IsNumber() bool {
	return am.TstNumberBit()
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

func parseSingleArg(arg string) (flags string, param string, argType ArgMask) {
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
		flags = flag
		return
	}

	// It could be a single short flag, a cluster of short flags, or a
	// number (the -NUM idiom)
	argType.ClrLongBit()
	r, tail := FirstRune(flag)
	if tail == "" {
		if unicode.IsNumber(r) {
			argType.SetNumberBit()
		}
		flags = string(r)
		return
	}

	// We have a cluster of short flags or a number
	if param == "" {
		// If there's an attached parameter, it can't be the -NUM
		// idiom
		_, err := strconv.ParseUint(flag, 10, 64)
		if err == nil {
			argType.SetNumberBit()
		}
	}
	argType.SetClusterBit()
	flags = flag

	return
}

// DisambiguateCluster decides how an apparent/possible cluster is to
// be disambiguated.
//
// The cluster `-fgh` has 3 possible legal interpretations under POSIX
// rules.
//
//   * f() g() h()
//   * f() g(h)
//   * f(gh)
//
// We work from left-to-right, giving precedence to interpretation as
// a flag. Once a non-flag is encountered, the rest of the string is
// assumed to be an option-argument to the last flag.
func (fs *FlagSet) disambiguateCluster(flags string, param string, argType ArgMask, pos int) *Flag {
	// Process clusters by POSIX rules where the last flag in
	// the cluster can have an option-argument.
	var curr *Flag
	for i, s := range flags {
		prev := curr
		curr = fs.Lookup(s)
		if curr == nil {
			// Could be a number:
			if argType.IsNumber() && param == "" {
				curr = fs.Lookup(NoShort)
				if curr != nil {
					err := curr.Set(flags, pos)
					if err != nil {
						fs.Failf("failed to set '%s' with '%s' (-NUM idiom): %v", curr, flags, err)
					}
					return nil
				}
			}
			// Non-flag: this and whatever follows must be an attached
			// option-argument to the previous flag
			optarg := flags[i:]
			if param != "" {
				optarg += "=" + param
			}
			err := prev.Set(optarg, pos)
			if err != nil {
				// We may return (or not) after Fail depending on OnFail setting
				fs.Failf("failed to set '%s' with '%s': %v", prev, optarg, err)
			}
			return nil
		}
		if prev != nil {
			err := prev.Set(nil, pos)
			if err != nil {
				fs.Failf("failed to set '%s' with nil: %v", prev, err)
			}
		}
	}
	// We now have the last flag in `curr` that hasn't been acted on:
	// return it in case there's an unattached option-argument (aka
	// parameter) in the next argument
	return curr
}

// Function StopParsing moves all remaining input arguments to the
// output slice, optionally discarding the first element of the input
func (fs *FlagSet) stopParsing(shift bool) {
	if shift {
		_, _ = fs.InputArgs.Shift()
	}
	fs.OutputArgs.Append([]string(*fs.InputArgs)...)
	fs.InputArgs.Clear()
}

func (fs *FlagSet) parse() error {
	var err error
	var i int = 0

	for arg, err := fs.InputArgs.Shift(); err == nil; arg, err = fs.InputArgs.Shift() {
		i++
		flags, param, argType := parseSingleArg(arg)
		if !argType.IsFlag() {
			fs.OutputArgs.Push(param)
			if PosixOperandStop {
				fs.stopParsing(false)
				return nil
			}
			continue
		}
		if argType.IsDoubleHyphen() {
			// arg can't be an option-argument at this point, so we
			// terminate processing under either POSIX or GNU rules
			fs.stopParsing(false)
			return nil
		}
		var flag *Flag = nil
		if argType.IsCluster() {
			// It's parsed as a cluster, but that doesn't mean it
			// is. It could be a flag with an attached argument.
			flag = fs.disambiguateCluster(flags, param, argType, i)
			if flag == nil {
				// Fully handled in fs.disambiguateCluster()
				continue
			}
		} else {
			flag = fs.Lookup(flags)
			if flag == nil {
				if !argType.IsNumber() {
					fs.Failf("flag '%s' not defined (NaN)", flags)
					continue
				}
				flag = fs.Lookup(NoShort)
				if flag == nil {
					fs.Failf("flag '-NUM' not defined for '%s'", flags)
					continue
				}
				err = flag.Set(flags, i)
				if err != nil {
					fs.Failf("failed to set -NUM flag with '%s': %v", flags, err)
				}
				continue
			}
		}
		if argType.HasParam() {
			// This must've been attached with an '=', so if it's a
			// short flag, the '=' is part of the argument under POSIX
			// rules and it we don't have to check if the flag takes
			// an argument: if it fails, there's a mistake on the
			// command-line.
			if (argType.IsShortFlag() || argType.IsCluster()) && PosixEquals {
				err = flag.Set("="+param, i)
			} else {
				err = flag.Set(param, i)
			}
			if err != nil {
				fs.Failf("failed to set flag `%s` with '%s': %v", flag.String(), param, err)
			}
			continue
		}
		// Peek at the next argument to see if it's a parameter (aka
		// option-argument)
		next, err := fs.InputArgs.Front()
		if err != nil {
			// End of InputArgs
			err = flag.Set(nil, i)
			if err != nil {
				fs.Failf("failed to set flag `%s` at EOL with no parameter: %v", flag.String(), err)
			}
			// At EOL
			return nil
		}
		// Have next arg, might be a parameter
		flags, param, nextArgType := parseSingleArg(next)
		if !nextArgType.IsFlag() {
			if nextArgType.IsDoubleHyphen() {
				// Under GNU (not POSIX) rules, we terminate if the
				// double-hyphen appears anywhere:
				if !PosixDoubleHyphen {
					fs.stopParsing(true)
					return nil
				}
				// See if the flag will accept "--" as an argument:
				err = flag.Test("--", i)
				if err != nil {
					// flag wouldn't eat it, so not an option-argument
					fs.stopParsing(true)
					return nil
				}
				err = flag.Set("--", i)
				if err != nil {
					fs.Failf("failed to set flag `%s` with `--` after Test(): %v", flag, err)
				}
				// It worked as a parameter/optarg, so consume it
				_, _ = fs.InputArgs.Shift()
				i++
				continue
			}
			// Not a flag, try it as a parameter
			if !flag.IsBool() {
				err = flag.Set(param, i)
				if err == nil {
					// It worked as a parameter, so consume it
					_, _ = fs.InputArgs.Shift()
					i++
				}
				continue
			}
		}
		// Next arg is a flag, current flag has no parameter
		err = flag.Set(nil, i)
		if err != nil {
			fs.Failf("failed to set flag `%s` with no parameter", flag.String())
		}
	}
	return err
}

func (fs *FlagSet) Parse(arguments []string) error {
	fs.InputArgs.Init(arguments...)
	return fs.parse()
}

func Parse() {
	CommandLine.Parse(os.Args[1:])
}
