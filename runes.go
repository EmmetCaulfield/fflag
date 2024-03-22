package fflag

import (
	"unicode/utf8"
)

// Negative runes in the range -16 to -1 (U+FFF0 to U+FFFF) are
// Unicode "Specials", notably U+FFFD (-3), the "Unicode replacement
// character", so we avoid this range for error indication, but
// otherwise use negative rune values to indicate an error in cases
// where it's more convenient than a separate `error` return.
const (
	ErrRuneEmptyStr   rune = -17
	ErrRuneIdSepBad        = -18
	ErrRuneShortBad        = -19
	ErrRuneIdPartsBad      = -20
)

func FirstRune(s string) (rune, string) {
	if len(s) == 0 {
		return ErrRuneEmptyStr, ""
	}
	for _, char := range s {
		return char, s[utf8.RuneLen(char):]
	}
	// This is impossible
	panic("unreachable code")
}
