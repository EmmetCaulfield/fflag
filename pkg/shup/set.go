// `shup` implements a trie-based algorithm for determining the
// shortest unique prefix of each string in a set
package shup

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	s := Set[T]{}
	s.Add(items...)
	return s
}

func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

func (s Set[T]) Has(items ...T) bool {
	for _, item := range items {
		if _, ok := s[item]; !ok {
			return false
		}
	}
	return true
}

func (s Set[T]) Del(items ...T) {
	for _, item := range items {
		if _, ok := s[item]; ok {
			delete(s, item)
		}
	}
}

func (s Set[T]) Slice() []T {
	keys := make([]T, len(s))
	i := 0
	for key := range s {
		keys[i] = key
		i++
	}
	return keys
}
