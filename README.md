# `fflag`

## Introduction

`fflag` is a command-line option processing library, similar in effect
to `flag` or `pflag`, but it's easier to use.

Rather than using dozens of different functions, `fflag` detects the
type of flags from the pointer passed to store the value.

`fflag` uses the functional options pattern. For example, if you were
writing grep and wanted to implement the `-q`/`--quiet` switch, you
would use:

    var quiet bool
    fflag.NewFlag(&quiet, "quiet",
        WithShortOption("q"),
        WithAlias("silent", false),
        WithDefault(false),
        WithArgument(false),
        WithMessage("do not write anything to stdout")
    )

Repeating the `quiet` switch is considered an error _unless_ the
`WithCallback(...)` option provides a function to call when the switch
appears on the command-line. The function has the following signature:

```go
type Callback func(varp interface{}, flag string, arg string, pos int)
}
```

The function will be called with the interface (containing a pointer
to the variable to be set) supplied as the first argument to
`NewFlag`, the flag as encountered on the command-line, the argument
to the flag as encountered on the command-line, and the position of
the flag in the command-line argument list.

If a callback is supplied, `fflag` will _not_ set any variable: it is
up to the callback to do that via the supplied pointer.

Repeatable flags are detected based on the type of the first
argument. For example, `fflag` can tell that `grep`'s `-e/--regext`
is repeatable because the supplied pointer is a pointer to a slice:
there is no need for separate `...Slice...` constructors:

    patterns := []string{}
    fflag.NewFlag(&patterns, "pattern",
        WithShortFlag("e"),
        WithArgument(true),
        WithMessage("do not write anything to stdout")
        WithFileFlag("--file", WithShortOption("f"))
    )

The `WithFileFlag` option recognizes a pattern for command-line
arguments where a repeatable flag has a one-per-line file
equivalent. This is deemed common enough to be worth supporting. It is
an error to specify `WithArgument(false)` or `WithFileFlag(...)`
unless the first argument points to a slice.
