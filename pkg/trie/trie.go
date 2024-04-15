package trie

import(
	"fmt"
	"unicode/utf8"
)

type TrieNode[T any] struct {
	Item *T
	Tail string
	Nodes map[rune]*TrieNode[T]
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

func (t *TrieNode[T]) Get(key string) (*T, error) {
	if len(key) == 0 && len(t.Nodes) == 0 {
		// We've exhausted the search key and there are no sub-nodes
		// to look at:
		return t.Item, nil
	}
	if len(t.Tail) >= len(key) && t.Tail[:len(key)] == key {
		// Found an unambiguous substring match:
		return t.Item, nil
	}
	r, tail := firstRune(key)
	if r == utf8.RuneError {
		// We either got a bona-fide utf8 rune error or passed the
		// empty string to firstRune(), neither of which we're
		// prepared to handle
		panic("unexpected string error")
	}
	if node, ok := t.Nodes[r]; ok {
		return node.Get(tail)
	}
	// Search key not found
	return nil, nil
}
	

func (t *TrieNode[T]) Add(key string, item *T) error {
	if len(key) == 0 {
		// We've exhausted the key
		if t.Item == nil {
			t.Item = item
			return nil
		}
		return fmt.Errorf("duplicate key in trie")
	}
	if len(t.Nodes) == 0 && len(t.Tail) == 0 && t.Item == nil {
		t.Tail = key
		t.Item = item
		return nil
	}
	// If there's a tail here, we need to move it down:
	if len(t.Tail) > 0 {
		tR, tTail := firstRune(t.Tail)
		if tR == utf8.RuneError {
			panic("unexpected string error")
		}
		t.Nodes[tR] = &TrieNode[T]{
			Item: t.Item,
			Tail: tTail,
			Nodes: map[rune]*TrieNode[T]{},
		}
		t.Item = nil
		t.Tail = ""
	}

	r, tail := firstRune(key)
	if r == utf8.RuneError {
		panic("unexpected string error")
	}
	if node, ok := t.Nodes[r]; ok {
		return node.Add(tail, item)
	}

	t.Nodes[r] = &TrieNode[T]{
		Item: item,
		Tail: tail,
		Nodes: map[rune]*TrieNode[T]{},
	}
	return nil
}

func NewTrie[T any]() *TrieNode[T] {
	return &TrieNode[T]{
		Item: nil,
		Tail: "",
		Nodes: map[rune]*TrieNode[T]{},
	}
}

func (t *TrieNode[T]) Dump(pfx string) {
	if t == nil {
		fmt.Printf("%s<nil>\n", pfx)
		return
	}
	fmt.Printf("%s(%v, '%s')\n", pfx, t.Item, t.Tail)
	for r, t := range t.Nodes {
		fmt.Printf(pfx+"- %c: ", r)
		t.Dump(pfx+"    ")
	}
}
