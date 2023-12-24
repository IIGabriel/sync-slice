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

// GetUnsafe retrieves an element at a given index without checking for out-of-bounds.
// This method is unsafe and can Panic if the index is out of bounds.
func (s *Slice[T]) GetUnsafe(index int) T {
	return s.slice[index]
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

// SetUnsafe updates the element at a given index without checking for out-of-bounds.
// This method is unsafe and can Panic if the index is out of bounds.
func (s *Slice[T]) SetUnsafe(index int, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice[index] = value
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

// RemoveUnsafe removes the element at a given index without checking for out-of-bounds.
// This method is unsafe and can Panic if the index is out of bounds.
func (s *Slice[T]) RemoveUnsafe(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice = append(s.slice[:index], s.slice[index+1:]...)
	atomic.AddInt32(&s.len, -1)
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

// ManipulateAtIndex applies a manipulation function to the element at the specified index.
// It returns true if the index is within bounds and the operation is successful.
func (s *Slice[T]) ManipulateAtIndex(index int, manipulateFunc func(*T)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < 0 || index >= int(s.len) {
		return false
	}
	manipulateFunc(&s.slice[index])

	return true
}

// ManipulateAtIndexUnsafe applies a manipulation function to the element at the specified index without bounds checking.
// This method is unsafe and can panic if the index is out of bounds.
func (s *Slice[T]) ManipulateAtIndexUnsafe(index int, manipulateFunc func(*T)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	manipulateFunc(&s.slice[index])
}
