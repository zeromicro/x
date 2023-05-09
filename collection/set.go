package collection

import "github.com/zeromicro/go-zero/core/lang"

// Set is not thread-safe, for concurrent use, make sure to use it with synchronization.
type Set[T comparable] struct {
	data map[T]lang.PlaceholderType
}

// NewSet returns a managed Set, can only put the values with the same type.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]lang.PlaceholderType),
	}
}

// Add adds i into s.
func (s *Set[T]) Add(i ...T) {
	for _, each := range i {
		s.add(each)
	}
}

// Contains checks if i is in s.
func (s *Set[T]) Contains(i T) bool {
	if len(s.data) == 0 {
		return false
	}

	_, ok := s.data[i]
	return ok
}

// Keys returns the keys in s.
func (s *Set[T]) Keys() []T {
	var keys []T

	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

// Remove removes i from s.
func (s *Set[T]) Remove(i T) {
	delete(s.data, i)
}

// Count returns the number of items in s.
func (s *Set[T]) Count() int {
	return len(s.data)
}

func (s *Set[T]) add(i T) {
	s.data[i] = lang.Placeholder
}
