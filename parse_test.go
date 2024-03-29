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
	assert.Equal(t, expected, s)
}

