package shup

import(
	"unicode/utf8"
)

type TrieNode struct {
	Char rune
	Strs Set[string]
	Subs map[rune]*TrieNode
}

func NewTrie(char rune, items ...string) *TrieNode {
	// Dedupe the strings:
	ns := NewSet[string](items...)
	// If the empty string is in the set, get rid of it
	ns.Del("")
	// If the remaining set is empty, don't create a new trie:
	if len(ns) == 0 {
		return nil
	}
	return &TrieNode{
		Char: char,
		Strs: ns,
		Subs: map[rune]*TrieNode{},
	}
}

// Returns the first rune and the rest of the string
func firstRune(s string) (rune, string) {
	if len(s) == 0 {
		return utf8.RuneError, ""
	}
	for _, char := range s {
		return char, s[utf8.RuneLen(char):]
	}
	// Impossible
	panic("unreachable code")
}

func (tn *TrieNode) Populate() {
	// We stop when there is only one string (containing the remnant)
	// in `Strs`. It shouldn't be possible to have nothing in `Strs`
	// since, in that case, the node shouldn't have been created at
	// all.
	if len(tn.Strs) < 2 {
		return
	}
	for _, str := range tn.Strs.Slice() {
		if len(str) == 0 {
			// Don't propagate the empty string
			continue
		}
		first, rest := firstRune(str)
		if first == utf8.RuneError {
			// We either got a bona-fide utf8 rune error or passed the
			// empty string to firstRune(), neither of which we're
			// prepared to handle
			panic("unexpected string error")
		}
		if trie, ok := tn.Subs[first]; ok {
			// If len(rest) == 0, do we need to propagate to another
			// level of trie? --- not if the trie already exists
			if len(rest) == 0 {
				continue
			}
			// trie could still be `nil` (from below)
			if trie == nil {
				trie = NewTrie(first, rest)
				if trie == nil {
					// This isn't possible
					panic("nil tree from non-empty `rest` string")
				}
				tn.Subs[first] = trie
				continue
			}
			trie.Strs.Add(rest)
		} else {
			// If len(rest) == 0, put in a nil Trie pointer:
			if len(rest) == 0 {
				tn.Subs[first] = nil
				continue
			}
			tn.Subs[first] = NewTrie(first, rest)
		}
	}

	for _, trie := range tn.Subs {
		if trie == nil {
			continue
		}
		trie.Populate()
	}
}

func (tn *TrieNode) Descend(s string) (string, string) {
	var remnant string
	prefix := ""
	for _, char := range s {
		prefix += string(char)
		if len(tn.Strs) == 1 {
			_, remnant = firstRune(tn.Strs.Slice()[0])
			break
		}
		var ok bool
		if tn, ok = tn.Subs[char]; ok {
			continue
		}
		break
	}
	return prefix, remnant 
}
