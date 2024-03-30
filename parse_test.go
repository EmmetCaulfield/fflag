package fflag

import (
	"testing"

	"github.com/EmmetCaulfield/fflag/pkg/deque"
	"github.com/stretchr/testify/assert"
)

func mute() {
	CommandLine.OnFail.SetSilentBit()
}

func unmute() {
	CommandLine.OnFail.ClrSilentBit()
}

func TestPosixRejects(u *testing.T) {
	t := assert.TestingT(u)
	unmute()
	b := false
	assert.Panics(t, func() { Var(&b, '?', "", "should panic") })
	PosixRejectQuest = false
	Var(&b, '?', "", "should not panic now")
	CommandLine.Parse([]string{"-?"})
	assert.Equal(t, true, b)

	b = false
	assert.Panics(t, func() { Var(&b, 'W', "", "should panic") })
	PosixRejectW = false
	Var(&b, 'W', "", "should not panic now")
	CommandLine.Parse([]string{"-W"})
	assert.Equal(t, true, b)
}

func TestPosixDoubleHyphen(u *testing.T) {
	t := assert.TestingT(u)
	var a, b, c bool
	var s string
	fs := NewFlagSet()
	fs.Var(&a, 'a', "ant", "six legs")
	fs.Var(&b, 'b', "bat", "two legs, two wings")
	fs.Var(&c, 'c', "cow", "four legs")
	fs.Var(&s, 's', "snake", "no legs")

	args := []string{"-a", "-s", "--", "-c", "operand"}

	// Under GNU rules, "--" stops argument processing immediately,
	// leaving "-c" and "operand":
	expected := &deque.Deque[string]{}
	expected.Init("-c", "operand")
	PosixDoubleHyphen = false
	fs.Parse(args)
	assert.Equal(t, expected, fs.OutputArgs, "GNU rule")

	// Under POSIX rules, "--" is the argument to "-s" and processing
	// does not stop, leaving just "operand" in the output.
	_, _ = expected.Shift()
	PosixDoubleHyphen = true
	fs.Reset()
	fs.Parse(args)
	assert.Equal(t, expected, fs.OutputArgs, "POSIX rule")
}

func TestPosixEquals(u *testing.T) {
	t := assert.TestingT(u)
	var s string
	fs := NewFlagSet()
	fs.Var(&s, 's', "snake", "no legs")

	args := []string{"-s=python"}

	// Under de-facto GNU-like rules, the "=" is not part of the optarg
	expected := "python"
	PosixEquals = false
	fs.Parse(args)
	assert.Equal(t, expected, s, "GNU rule")

	fs.Reset()
	// Under POSIX rules, the "=" is not special and is part of the
	// optarg
	expected = "=python"
	PosixEquals = true
	fs.Parse(args)
	assert.Equal(t, expected, s, "GNU rule")
}
	
func TestPosixOperandStop(u *testing.T) {
	t := assert.TestingT(u)
	var a, b, c bool
	var s string
	fs := NewFlagSet()

	fs.Var(&a, 'a', "ant", "six legs")
	fs.Var(&b, 'b', "bat", "two legs, two wings")
	fs.Var(&c, 'c', "cow", "four legs")
	fs.Var(&s, 's', "snake", "no legs")

	args := []string{"-b", "-c", "operand", "-a", "-s", "python"}

	// Under GNU rules, 'a' is true, 's' is "python", and "operand" is in the output
	expected := &deque.Deque[string]{}
	expected.Init("operand")
	PosixOperandStop = false
	fs.Parse(args)
	assert.Equal(t, true, a)
	assert.Equal(t, true, b)
	assert.Equal(t, true, c)
	assert.Equal(t, "python", s)
	assert.Equal(t, expected, fs.OutputArgs, "GNU operand stop")

	fs.Reset(); a=false; b=false; c=false; s=""
	// Under POSIX rules, 'a' is false, 's' is "", and "operand", "-a", "-s", "python" is in the output
	expected.Append("-a", "-s", "python")
	PosixOperandStop = true
	fs.Parse(args)
	assert.Equal(t, false, a)
	assert.Equal(t, true, b)
	assert.Equal(t, true, c)
	assert.Equal(t, "", s)
	assert.Equal(t, expected, fs.OutputArgs, "POSIX operand stop")
}

func TestCluster(u *testing.T) {
	t := assert.TestingT(u)
	var a, b, c bool
	var s string
	fs := NewFlagSet()
	fs.Var(&a, 'a', "ant", "six legs")
	fs.Var(&b, 'b', "bat", "two legs, two wings")
	fs.Var(&c, 'c', "cow", "four legs")
	fs.Var(&s, 's', "snake", "no legs")

	args := []string{"-abs", "python"}
	expected := args[1]
	fs.Parse(args)
	assert.Equal(t, true, a)
	assert.Equal(t, true, b)
	assert.Equal(t, false, c)
	assert.Equal(t, expected, s, "with detached argument")

	fs.Reset(); a=false; b=false
	args = []string{"-abspython"}
	expected = "python"
	fs.Parse(args)
	assert.Equal(t, true, a)
	assert.Equal(t, true, b)
	assert.Equal(t, false, c)
	assert.Equal(t, expected, s, "with attached argument")

	fs.Reset(); a=false; b=false
	args = []string{"-abs=python"}
	PosixEquals = true
	expected = "=python"
	fs.Parse(args)
	assert.Equal(t, true, a)
	assert.Equal(t, true, b)
	assert.Equal(t, false, c)
	assert.Equal(t, expected, s, "with equals attached argument, POSIX mode")

	fs.Reset(); a=false; b=false
	args = []string{"-abs=python"}
	PosixEquals = false
	expected = "python"
	fs.Parse(args)
	assert.Equal(t, true, a)
	assert.Equal(t, true, b)
	assert.Equal(t, false, c)
	assert.Equal(t, expected, s, "with equals attached argument, GNU mode")
}

func TestCluster2(u *testing.T) {
	t := assert.TestingT(u)
	var a, b, c bool
	var s string
	fs := NewFlagSet()
	fs.Var(&a, 'a', "ant", "six legs")
	fs.Var(&b, 'b', "bat", "two legs, two wings")
	fs.Var(&c, 'c', "cow", "four legs")
	fs.Var(&s, 's', "snake", "no legs")

	args := []string{"-ac", "operand"}
	expected := &deque.Deque[string]{}
	expected.Init("operand")
	fs.Parse(args)

	assert.Equal(t, true, a)
	assert.Equal(t, false, b)
	assert.Equal(t, true, c)
	assert.Equal(t, "", s)
	assert.Equal(t, expected, fs.OutputArgs, "GNU rule")
}

func TestHyphenNumIdiom(u *testing.T) {
	t := assert.TestingT(u)
	var n uint
	fs := NewFlagSet()
	fs.Var(&n, NoShort, NoLong, "a natural number")

	// Looks like a flag
	args := []string{"-7"}
	fs.Parse(args)
	assert.Equal(t, uint(7), n)

	// Looks like a cluster
	fs.Reset()
	args = []string{"-371"}
	fs.Parse(args)
	assert.Equal(t, uint(371), n)
}

// How should `-42a5` be interpreted? <-42, a(5)>? -4(2a5)? -(0x42a5)
// Should -NUM handle non-decimal NUM (e.g. octal, hex)?
// Rule: if you use the -NUM idiom, you can't define numeric short flags?
