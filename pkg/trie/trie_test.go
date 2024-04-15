package trie

import (
	"strings"
	"testing"
)

func TestTrieBasics(t *testing.T) {
	trie := NewTrie[string]()
	contents := []string{"foo", "bar", "baz", "bazaar", "fop", "quux"}

	for _, s := range contents {
		v := strings.ToUpper(s)
		err := trie.Add(s, &v)
		if err != nil {
			t.Errorf("error: %v", err)
		}
	}

	// Exact keys (should succeed):
	for _, s := range contents {
		n, err := trie.Get(s)
		if err != nil {
			t.Errorf("error retrieving node for '%s': %v", s, err)
		}
		if n == nil {
			t.Errorf("failed to retrieve node for '%s'", s)
		}
		if *n != strings.ToUpper(s) {
			t.Errorf("got wrong value: expected '%s', got '%s'", strings.ToUpper(s), *n)
		}
	}

	// Ambiguous keys (should fail):
	for _, s := range []string{"f", "fo", "b", "ba"} {
		n, err := trie.Get(s)
		if err != nil {
			t.Errorf("error retrieving node for '%s': %v", s, err)
		}
		if n != nil {
			t.Errorf("retrieved node %+v for ambiguous key '%s'", *n, s)
		}
	}

	// Short unique keys (should succeed)
	for _, s := range []string{"baza", "bazaa"} {
		n, err := trie.Get(s)
		if err != nil {
			t.Errorf("error retrieving node for '%s': %v", s, err)
		}
		if n == nil {
			t.Errorf("failed to retrieve node for '%s'", s)
		}
		if *n != "BAZAAR" {
			t.Errorf("got wrong value: expected 'BAZAAR', got '%s'", *n)
		}
	}

	for _, s := range []string{"q", "qu", "quu"} {
		n, err := trie.Get(s)
		if err != nil {
			t.Errorf("error retrieving node for '%s': %v", s, err)
		}
		if n == nil {
			t.Errorf("failed to retrieve node for '%s'", s)
		}
		if *n != "QUUX" {
			t.Errorf("got wrong value: expected 'QUUX', got '%s'", *n)
		}
	}

	// Duplicate keys (should fail):
	for _, s := range []string{"foo", "fop", "bar", "baz", "quux"} {
		value := "DUPE-" + s
		err := trie.Add(s, &value)
		if err == nil {
			t.Errorf("unexpected success adding duplicate key '%s'", s)
		}
	}
}
