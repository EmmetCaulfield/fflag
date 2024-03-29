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

func TestPosixisms(u *testing.T) {
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

func TestDoubleHyphen(u *testing.T) {
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
	u.Logf("Calling fs.Parse()")
	fs.Parse(args)
	assert.Equal(t, expected, fs.OutputArgs, "GNU rule")

	// Under POSIX rules, "--" is the argument to "-s" and processing
	// does not stop, leaving just "operand" in the output.
	u.Logf("Calling fs.Reset()")
	fs.Reset()
	_, _ = expected.Shift()
	PosixDoubleHyphen = true
	fs.Parse(args)
	assert.Equal(t, expected, fs.OutputArgs, "POSIX rule")
}
