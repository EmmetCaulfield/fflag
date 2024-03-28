package fflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutexes(t *testing.T) {
	unmute()
	var c, d bool
	Var(&c, 'c', "cat", "cat flag", InMutex("pet"))
	Var(&d, 'd', "dog", "dog flag", InMutex("pet"))
	CommandLine.Parse([]string{"-c"})
	assert.Equal(t, true, c)
	CommandLine.Parse([]string{"-d"})
	assert.Equal(t, false, d)
}
