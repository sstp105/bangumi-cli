package libs

import "fmt"

// Set is a generic set type that can hold elements of any type.
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates and returns a new Set for a given type T
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{elements: make(map[T]struct{})}
}

// Add adds an element to the set
func (s *Set[T]) Add(value T) {
	s.elements[value] = struct{}{}
}

// Contains checks if an element exists in the set
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.elements[value]
	return exists
}

// Remove removes an element from the set
func (s *Set[T]) Remove(value T) {
	delete(s.elements, value)
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Print displays all elements in the set
func (s *Set[T]) Print() {
	for item := range s.elements {
		fmt.Println(item)
	}
}
