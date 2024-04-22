# `fflag`

Package `fflag` provides POSIX (short) and GNU (long) command-line
argument parsing with, for the programmer, the functional options
pattern.

## Overview

It is somewhat inspired by the `pflag` package in some respects,
but very significantly different in others. The most significant
difference is that there is only one `Var()` function (rather than
six dozen variations): the type of the flag is determined by the
type of the first argument (rather than the function name), which
_must_ be a pointer to a basic type, a slice of basic type, or a
struct implementing the `SetValue` interface (inspired by `pflag`).

The other significant difference is the order of the short flag and
long flag in the `Var()` argument list, with the short flag coming
first as a `rune`, most often a single ASCII letter or number. The
reason for this is that short flags are always listed first in
manpages and other documentation, so it's actually a bit weird of
`pflag` to have reversed this de-facto standard order and, in
practice, I've found it handier to obey the standard order than
stick to the reversed `flag` argument order.

If there is no short flag, the zero value (0, `'\0'`, or
`fflag.NoShort`) is used. The usual rules apply to long flags,
which must consist of letters and numbers, except that the ASCII
requirement has been relaxed. Any character satisfying
`unicode.IsLetter()` or `unicode.IsNumber()`, or the ASCII
hyphen/minus character, '-', are allowed. There is _no attempt
whatsoever_ at (what `pflag` refers to as) normalization, a very
dubious utility: just use the long flag you mean to use without
weird capitalization.

That said, `fflag` meets the onerous GNU expectation that “users
can abbreviate the option names as long as the abbreviations are
unique”. This gives rise to the issue that a long flag `--xyzzy`
could be unique as `--x`, but a short flag `-x` could be defined
for something completely different. `fflag` resolves this by giving
priority to the short flag interpretation if a long flag is one
character long, so that `-x` and `--x` are interpreted the same
even if `--xyzzy` is different and `--x` would be unique: you would
have to type `--xy` to get `--xyzzy`.

`fflag` borrows the `Flag` and `FlagSet` names from `pflag`, adding
`FlagGroup`. The purpose of a flag group is to enable usage
information to be generated in a similar format to GNU/POSIX
utilities like `grep`, with flags grouped in categories. This is an
additional feature of `fflag` and isn't known to exist elsewhere.

A `Flag` is created and added to the default `FlagGroup` in the
default `FlagSet` (called `CommandLine` after `pflag`'s equivalent)
with `Var()`. The minimal call to `Var()` provides: a pointer to a
variable where the value of the flag is to be stored; the
single-short version of the flag as a rune (or `NoShort` if none),
e.g., 'h'; the long version of the flag (or `NoLong` if none),
e.g. `--help`; and a very brief description of the flag's
purpose. For example:

    var value bool
    fflag.Var(&value, 'h', "help", "prints a help message to stdout")

The first argument to `Var` must be a _pointer_ to one of:

  1) a basic datatype (e.g. `int8`, `float32`, `string`)
  2) a slice of basic datatype (e.g. `[]int8`, `[]string`)
  3) something implementing the `SetValue` interface

Non-pointer `value` arguments will cause a `panic()`. As a rule,
`fflag` will `panic()` in the case of a programmer mistake (during
setup) and return an `error` otherwise (during argument parsing).

If the value argument implements the `SetValue` interface, `fflag`
neither modifies the argument itself nor enforces any of its usual
rules. If you pass something implementing this interface, it's
assumed that you will take care of everything and don't want
`fflag` to do anything other than pass along the message “this flag
appeared with this argument”.

A flag need not have a single-character shortcut. If there is no
shortcut, a `fflag.NoShort` is given for the shortcut argument:

    fflag.Var(&value, fflag.NoShort, "help", "prints a help message")

Only letters and numbers are normally allowed as shortcuts. The
sole exception is '?' due to its widespread use as an alias for
"help", but this is prohibited by POSIX, so if you want to use
this, you have to enable it explicitly.

Equally, a flag need not have a long version. If you wanted to have
`-?` as a short flag with no long version, you would do:

    fflag.PosixRejectQuest = false
    fflag.Var(&value, '?', fflag.NoLong, "prints a help message to stdout")

There is a special case (and common idiom) where NEITHER a long NOR
a short form is required: `-NUM` (as in `grep`, `head`, `tail`, and
several other tools). These special cases are always an alias for
something else and always refer to a natural number appearing after
a single hyphen. For example `head`'s `-n/--lines` is best
represented as:

    uint nlines
    fflag.Var(&nlines, 'n', 'lines',
        "print the first NUM lines instead of the first 10",
        fflag.WithAlias(fflag.NoShort, fflag.NoLong, false))

Obviously, this special case can only be used once and attempting
to use it more than once is a programmer error that results in a
`panic()`.

The simplest “vanilla” flag is a nullary boolean switch that takes
no parameter.

    bool value
    fflag.Var(&value, 'e', "easy", "use easy mode"))

In this case, `value` will default to `false` (the zero value for
`bool`s) and become `true` if the command-line argument appears in
either form (long or short). By default, it is an error to repeat a
scalar flag, but there are 3 options that make an exception:

  * `WithRepeats(ignore bool)`
  * `AsCounter()`
  * `WithCallback(callback func(...))`

`WithRepeats()` allows repeat appearances of a flag, `AsCounter()`
causes the number of occurrences to be counted (if the value
pointer is a number), and `WithCallback()` causes the given
callback function to be called _every time_ the flag appears on the
command-line. Much like a `value` (first) argument implementing
`SetValue`, `fflag` washes its hands of any further involvement,
and it becomes entirely up to the callback to modify the value
appropriately, track/ignore repeated appearances, etc.

Several utilities allow `-v/--verbose` to be repeated for
increasing levels of verbosity.

    int verbosity
    f := NewFlag(&verbosity, 'v', "verbose", "increase verbosity",
        AsCounter())

Supplying more than one of these (pairwise redundant or
contradictory) options would result in a `panic()` since this would
be an obvious programmer error, not something that could
“accidentally” occur at runtime based on user input.

An explicit default can be supplied with `WithDefault()`:

    var hard bool
    fflag.Var(&hard, fflag.NoShort, "easy", "use easy mode",
        fflag.WithDefault(true))

In this case, `hard` will default to `true` and become false if
`--easy` appears on the command line. If repeats are allowed and
not ignored (`WithRepeat(false)`), the value will toggle between
`true` and `false`, which is admittedly weird, but if you do stupid
things, expect to win stupid prizes.

Repeated appearances of a flag are _not_ an error if the value
argument is (a pointer to) a slice. In this case, successive
invocations will result in successive values being appended to the
slice. This makes no sense for a nullary boolean flag.

    values := []int{}
    NewFlag(&values, 'x', "example", "example flag")

The sole exception to this rule is where a callback function is
supplied. When a callback is supplied, the callback is responsible
for _everything_.

    f := NewFlag(&value, 'f', "file", "supply a filename",
        WithCallback(MyFunc))

The callback function is called with the `Flag` pointer, string
argument and flag position (in os.Argv). The underlying variable,
short option, long option, etc. can be retrieved via the `Flag`
pointer. Consider a program `prog`, with the above "file" flag,
invoked as follows:

    prog -f foo.txt --file bar.txt

Here, `MyFunc` would be called twice as:

    MyFunc(f, "foo.txt", 1)
    MyFunc(f, "bar.txt", 3)

The `value` is _not_ set by `fflag` if a callback is supplied (or
if `value` implements the `SetValue` interface). You cannot call
f.Set() in the callback as this would lead to infinite recursion:
it is f.Set() that calls the callback. If you want to set the value
"normally" inside a callback, call `f.SetOnly()` instead, which
just sets the value and bypasses the usual flag type logic.

For unary (non-boolean) flags, a default can be supplied:

    var file string
    fflag.Var(&file, 'f', "file", "supply a filename",
        fflag.WithDefault("/dev/null"))

In this case, the value is set to the default immediately.

There exist a few utilities with options that differentiate between
the option appearing (at all) and the option appearing with an
argument. To support this, you must use `WithOptionalDefault()`

    var outdev := "/dev/stdout"
    fflag.Var(&outdev, 'o', "out", "supply an output device path",
        fflag.WithOptionalDefault("/dev/stderr"))

In this case, the value is not immediately set to the default. If
`-o` is _not_ supplied at all, `outdev` is not changed and keeps
the value "/dev/stdout"; if `-o` is supplied _without an argument_,
`outdev` changes to the default, `"/dev/stderr"`, and if an
argument is supplied, `outdev` is changed to the given
argument. This (rather strange) behavior is required to support
`grep`'s `--color` option, for example (more below).

If the default is a slice, valid values are constrained to the
values in the default, like a kind of set or `enum`.

Consider the `--directories` option of GNU `grep`. It can take one
of 3 values --- `read`, `skip`, and `recurse` --- with the default
being `read`:

    var string diract
    f := NewFlag(&diract, 'd', "directories",
        "if an input file is a directory use ACTION to process it",
        WithDefault([]string{"read", "skip", "recurse"}),
        WithTypeTag("ACTION"))

The actual default is the first value in the slice. The remaining
values in the slice constrain the set of acceptable values. For
some program, `prog`, with the above flag definition, the value of
`diract` after `fflag.Parse()` would be exactly the same for:

    $ prog
    $ prog -d read
    $ prog --directories=read

The following would be fine:

    $ prog -d skip
    $ prog --directories recurse

But the following would result in a runtime error because `foo` is
not one of "read", "skip", or "recurse", which are in the default
slice:

    $ prog -d foo

There exist a few utilities with set-constrained options that
differentiate between the option appearing (at all) and the option
appearing with an argument. For example, `grep`'s `--color`
argument (with no option) is different from the option not being
given at all, but also different from the option being given with
an argument. This is rare, but (unfortunately) it exists in the
wild.

    color := "never"
    fflag.Var(&color, fflag.NoShort, "color", "display matches in color",
        fflag.WithOptionalDefault([]string{"auto", "never", "always"}),
        fflag.WithTypeTag("[=WHEN]"))

In this case, if the `color` option is not given, the value is not
changed and remains at "never". If the argument is given _with no
argument_, the value is changed to "auto", the first value in the
default slice, and if the argument is given _with_ an argument, the
value changes to that argument _provided that_ it is in the default
slice. If `--color=foo` were given, it would result in an error.

## Package Options



