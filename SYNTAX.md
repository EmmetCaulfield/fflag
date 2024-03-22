# Command-line Argument Syntax and Grammar

The grammar of POSIX-style and GNU-style options are summarized online:

  * POSIX: <https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html>
  * POSIX `getopt()`: <https://pubs.opengroup.org/onlinepubs/9699919799/functions/getopt.html>
  * GNU: <https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html>
  * GNU `getopt()`: <https://www.gnu.org/software/libc/manual/html_node/Getopt.html>

## POSIX Options

Importantly, POSIX only deals with “short” options consisting of a
single ASCII alphanumeric character introduced by a single hyphen.
The specification gives the following prototype:

    utility_name[-a][-b][-c option_argument] [-d|-e][-f[option_argument]][operand...]

The option `-W` is reserved.

### Space Between an Option and its Argument

It requires that

> [...] the option be a separate argument from its option-argument and
> that option-arguments not be optional...

i.e, option-arguments should be separate arguments and mandatory, except

> [for] a mandatory option-argument [...], a conforming implementation
> shall also permit applications to specify the option and
> option-argument in the same argument string without intervening
> <blank> characters

So, they _don't_ have to be separate arguments.

> [for] an optional option-argument [...], a conforming application
> shall place any option-argument for that option directly adjacent to
> the option in the same argument string, without intervening <blank>
> characters. If the utility receives an argument containing only the
> option, it shall behave as specified in its description for an
> omitted option-argument

So, they _can't_ be separate arguments.

In short, the POSIX (short) options can have:

  * no option argument (option argument prohibited)
  * a mandatory option-argument
  * an optional option-argument (not recommended)

The space(s) between an option and its argument (if it has one) is
optional if the argument is mandatory (recommended, supported by
`getopt()`), and _prohibited_ if the argument is optional
(discommended, not supported by `getopt()`).

In the prototype from the specification, then, the argument to `-c` is
mandatory (because a space is shown), but the option to `-f` could be
optional (because no space is shown).

    utility_name [-a][-b][-c option_argument] [-d|-e][-f[option_argument]][operand...]

Consider a utility, `bar`, with the following options:

  * `o` - optional argument

In the command

    bar -ofoo

`foo` _must_ be interpreted as the (optional) option argument to `-o`,
while in

    bar -o foo

`foo` must be interpreted as an operand. 

If, however, the argument to `-o` was _mandatory_, then the two (`bar
-ofoo` and `bar -o foo`) _must_ be treated the same, `foo` is the
argument to `-o` and there are _no_ operands.


### Clustering

> One or more options without option-arguments, followed by at most
> one option that takes an option-argument, should be accepted when
> grouped behind one '-' delimiter.

In other words, single-letter options (the only kind POSIX addresses)
that don't have arguments can be clustered, but the last option in the
cluster _may_ have an argument.


## GNU Options

“Long” options, i.e. more than a single character introduced by two
hyphens, are strictly a GNU thing, but it also deals with “short”
(POSIX-like) options.

### Space Between an Option and its Argument

It says:

> An option and its argument may or may not appear as separate
> tokens. (In other words, the whitespace separating them is
> optional.) Thus, `-o foo` and `-ofoo` are equivalent.

Note that this conflicts with the POSIX requirement when the
option-argument is optional, since the space between and option and
its argument it _prohibited_ when the argument is optional, which
means that `-o foo` and `-ofoo` are definitely _not_ equivalent under
POSIX (as explained above) unless there are no optional
option-arguments.

This suggests that, under the GNU scheme, if a single-letter option
takes an argument at all, that argument is mandatory.

In the discussion of long options, it says:

> To specify an argument for a long option, write `--name=value`. This
> syntax enables a long option to accept an argument that is itself
> optional.

This obliquely suggests that short options aren’t enabled “to accept
an argument that is itself optional”, which (arguably) supports the
idea that arguments for short GNU options are either prohibited or
required, but never optional.

### Clustering

> Multiple options may follow a hyphen delimiter in a single token if
> the options do not take arguments. Thus, `-abc` is equivalent to `-a
> -b -c`.

Strictly, `-c` here can't take an option, while POSIX explicitly says
that the last option in the cluster _may_ take an option.

### Long Options

#### Long Option Abbreviation

The GNU conventions document contains the bombshell:

> Users can abbreviate the option names as long as the abbreviations are unique.





#### Long Option Arguments

All the [GNU `libc`
document](https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html)
says about long option arguments is:

> To specify an argument for a long option, write `--name=value`. This
> syntax enables a long option to accept an argument that is itself
> optional.

This suggests that this is the only way to supply an argument for a
long option. This is, of course, not true. It obliquely hints that a
short option may not be enabled to “accept an argument that is itself
optional”, which is also untrue: GNU `getopt()` itself supports
optional arguments (“this is a GNU extension”)


### De Facto Options

#### Arguments to Long Options

The GNU documents don't say _anything_ about long options separated
from their arguments by spaces (i.e. option-arguments given as
separate command-line arguments). If you lived in a vacuum, you might
assume that the syntax for supplying any argument to a long option is
just what is given (with the `=`), but every tool and utility accepts
arguments to long options separated by a space (i.e. as a separate
command-line argument).

#### The Single Hyphen

The POSIX document says that the single-hyphen should be interpreted
to mean “standard input” for consistency with other tools &
utilities. This is, indeed, a longstanding convention, but it has
nothing to do with argument processing _per se_. It should be
perfectly acceptable for any option processing library to treat it as
a non-option argument (i.e. as an option-argument or operand).

However, many tools and utilities accept `-NUM` command-line syntax,
usually to indicate some (positive integer) number of lines to be
shown from a file:

  * `grep -3 ${USER} /etc/passwd`
  * `tail -100 /var/log/messages`
  * `head -50 /var/log/messages`

We could regard this as “the empty short flag”, but it represents a
special case that requires special handling to support although one
could, of course, just treat it as an operand (non-option argument)
and leave its handling up to the client code.

This is further complicated by the fact that some tools _do_ use
numbers as options, e.g. `xargs -0`.

#### Equals-sign as Argument Separator

Many tools allow the equals sign, `=` as an argument separator for
POSIX-style _short_ options, not just GNU-style long options.

This introduces an ambiguity.


#### Mixed Options

It is often the case that options are “long in general but some of
them have shortcuts” or “short in general but some of them have long
versions”.

This is not the case in practice.

Utilities often have long options with no short version _and_ short
options with no long version and, as an occasional special case, an
“option” with neither. GNU `grep` falls into this category with
`--help` (no short option), `-I` no long option, and `-NUM` syntax.

Thus, neither the short option nor the long option can be considered
“primary” or “always defined” with the other being “extra”.

### Terminology

To avoid language like “the option has an optional option-argument”
and resolve differences in terminology between the GNU documents and
the POSIX one, we define some terminology.

We reserve the term _argument_ for what are, more specifically,
“program arguments” or “command-line arguments” that will appear as
separate elements of `os.Argv` (in the _Go_ context).

An argument can be:

  * an _operand_ (left over after argument processing, not any kind of option or optarg)
  * an _optarg_ (an option-argument), which is either
    * a _detached_ optarg (an argument by itself)
    * an _attached_ optarg, part of the same argument as an option,
      appearing to the right of it in the argument
  * an _option_, which is one of:
    * a _short_ (short option), introduced by a single hyphen, e.g. `-a`, `-b`, ...
      * potentially with an _attached optarg_, e.g. `-cfoo`;
    * a _cluster_ (of shorts), introduced by a single hyphen, e.g. `-abc`
      * potentially with an _attached optarg_, e.g. `-abcfoo`; or
    * a _long_ (long option), introduced by a double-hyphen, e.g. `--opt`
      * potentially with an _attached optarg_, e.g. `--opt=foo`
  * a _special token_ including
    * the _double-hyphen_, `--`
    * the _single-hyphen_, `-`

Notably, the single-hyphen isn’t really a special token _per se_, but
its correct handling requires some thought.


### Divergence

#### Double-hyphen

The double-hyphen “terminates all options” or “forces in all cases the
end of option scanning” according to GNU, but POSIX qualifies it as
“the first [double-hyphen] that is not an [optarg]”.

#### Cluster vs. Optarg

Suppose two options, `-f` and `-o` are defined. How is `-foo` to be
interpreted? Is it equivalent to `-f oo` or `-f -o -o`?

#### Equals-sign Separator

If an option, `-o`, takes an argument, then `-o=foo` (_de facto_ often
works) means that the option-argument is `=foo` under POSIX rules, but
`foo` under this _de facto_ rule.

