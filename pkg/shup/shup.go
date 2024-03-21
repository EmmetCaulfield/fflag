// `shup` implements a trie-based algorithm for determining the
// shortest unique prefix of each string in a set
package shup

func ShortestUniquePrefixMap(ss []string) map[string]string {
	trie := NewTrie(0, ss...)
	trie.Populate()
	supm := make(map[string]string, len(trie.Strs))
	for _, s := range trie.Strs.Slice() {
		pfx, _ := trie.Descend(s)
		if _, ok := supm[s]; ok {
			panic("duplicate in set/map")
		}
		supm[s] = pfx
	}
	return supm
}
