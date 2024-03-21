package shup

import (
	"testing"
)

func equalMaps[T comparable](a, b map[T]T) bool {
	for i, j := range a {
		if j != b[i] {
			return false
		}
	}
	return true
}

func TestTrieBasics(t *testing.T) {
	trie := NewTrie(0, "foo", "bar", "baz", "bazaar", "fop", "quux", "foo", "foo")
	if len(trie.Strs) != 6 {
		t.Errorf("wrong set length: expected 6, got %d", len(trie.Strs))
	}

	trie.Populate()
	if trie.Subs == nil {
		t.Error("nil tree in root after population")
	}
	if len(trie.Subs) != 3 {
		t.Errorf("wrong number of subtries, expected 3, got %d", len(trie.Subs))
	}
	for _, item := range []rune{'f', 'b', 'q'} {
		if _, ok := trie.Subs[item]; !ok {
			t.Errorf("subtree '%c' not found", item)
		}
	}

	t.Log("Start descent")
	pfx, rem := trie.Descend("bazaar")
	if pfx != "baza" || rem != "ar" {
		t.Errorf("descent failed: expected 'baza','ar'; got '%s','%s'", pfx, rem)
	} else {
		t.Logf("descent succeeded: expected 'baza','ar'; got '%s','%s'", pfx, rem)
	}		
	pfx, rem = trie.Descend("foo")
	if pfx != "foo" || rem != "" {
		t.Errorf("descent failed: expected 'foo',''; got '%s','%s'", pfx, rem)
	}
}
