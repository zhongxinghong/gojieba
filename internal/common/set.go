package common

import "strings"

// StringSet represents a collection of strings like
// std::unordered_set<std::string> in C++. It implemented by hash table.
type StringSet map[string]struct{}

// NewStringSet create an instance of StringSet.
func NewStringSet() StringSet {
	return make(StringSet)
}

func (s StringSet) String() string {
	vlist := make([]string, 0, s.Len())
	for v := range s {
		vlist = append(vlist, v)
	}
	return "[" + strings.Join(vlist, " ") + "]"
}

// Len returns the size of this string set.
func (s StringSet) Len() int {
	return len(s)
}

// Add add a string to this string set.
func (s StringSet) Add(v string) bool {
	if s.Has(v) {
		return false
	}
	s[v] = struct{}{}
	return true
}

// Del delete a string from this string set.
func (s StringSet) Del(v string) bool {
	if !s.Has(v) {
		return false
	}
	delete(s, v)
	return true
}

// Has tells you whether this string set contains the given string.
func (s StringSet) Has(v string) bool {
	_, ok := s[v]
	return ok
}
