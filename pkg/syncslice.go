package syncslice

import (
	"sync"
	"sync/atomic"
)

// Slice is a concurrent-safe, dynamically-sized slice.
type Slice[T any] struct {
	mu    sync.Mutex
	slice []T
	len   int32
}

// New creates a new concurrent-safe slice.
func New[T any]() *Slice[T] {
	return &Slice[T]{}
}

// Append adds an element to the end of the slice.
func (s *Slice[T]) Append(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice = append(s.slice, value)
	atomic.AddInt32(&s.len, 1)
}

// Get retrieves an element at a given index.
// If the index is out of bounds, it returns the zero value for the type and false.
func (s *Slice[T]) Get(index int) (T, bool) {
	if index < 0 || index >= int(atomic.LoadInt32(&s.len)) {
		var zero T
		return zero, false
	}

	return s.slice[index], true
}

// Set updates the element at a given index.
// If the index is out of bounds, it does nothing and returns false.
func (s *Slice[T]) Set(index int, value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < 0 || index >= int(s.len) {
		return false
	}

	s.slice[index] = value
	return true
}

// Length returns the current length of the slice.
func (s *Slice[T]) Length() int {
	return int(atomic.LoadInt32(&s.len))
}

// Remove removes the element at a given index.
// It returns false if the index is out of bounds.
func (s *Slice[T]) Remove(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < 0 || index >= int(s.len) {
		return false
	}

	s.slice = append(s.slice[:index], s.slice[index+1:]...)
	atomic.AddInt32(&s.len, -1)
	return true
}

// Range calls a function for each element in the slice.
// If the function returns false, it stops the iteration.
func (s *Slice[T]) Range(f func(index int, value T) bool) {
	for i, v := range s.slice {
		if !f(i, v) {
			break
		}
	}
}

// SetSlice replaces the internal slice with the provided slice.
func (s *Slice[T]) SetSlice(newSlice []T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice = make([]T, len(newSlice))
	copy(s.slice, newSlice)
	atomic.StoreInt32(&s.len, int32(len(newSlice)))
}

// GetSlice returns a copy of the internal slice.
func (s *Slice[T]) GetSlice() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	copiedSlice := make([]T, len(s.slice))
	copy(copiedSlice, s.slice)
	return copiedSlice
}
