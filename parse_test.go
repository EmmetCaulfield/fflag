package fflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mute() {
	CommandLine.OnFail.SetSilentBit()
}

func unmute() {
	CommandLine.OnFail.ClrSilentBit()
}

func TestPosixisms(t *testing.T) {
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
