package syncslice_test

import (
	syncslice "github.com/IIGabriel/sync-slice/pkg"
	"reflect"
	"sync"
	"testing"
)

// TestAppend tests the Append method for both success and failure scenarios.
func TestAppend(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	if s.Length() != 2 {
		t.Errorf("Expected length 2, got %d", s.Length())
	}
}

// TestGet tests the Get method.
func TestGet(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	if val, ok := s.Get(1); !ok || val != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", val)
	}

	// Test out of bounds
	if _, ok := s.Get(3); ok {
		t.Error("Expected false for out of bounds index")
	}
}

// TestSet tests the Set method.
func TestSet(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	if !s.Set(1, 3) {
		t.Error("Failed to set value at index 1")
	}

	if val, _ := s.Get(1); val != 3 {
		t.Errorf("Expected value 3 at index 1, got %d", val)
	}

	// Test out of bounds
	if s.Set(3, 4) {
		t.Error("Expected false for out of bounds index")
	}
}

// TestLength tests the Length method.
func TestLength(t *testing.T) {
	s := syncslice.New[int]()

	if s.Length() != 0 {
		t.Errorf("Expected length 0, got %d", s.Length())
	}

	s.Append(1)
	s.Append(2)

	if s.Length() != 2 {
		t.Errorf("Expected length 2, got %d", s.Length())
	}
}

// TestRemove tests the Remove method.
func TestRemove(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	if !s.Remove(0) {
		t.Error("Failed to remove value at index 0")
	}

	if s.Length() != 1 {
		t.Errorf("Expected length 1 after removal, got %d", s.Length())
	}

	// Test out of bounds
	if s.Remove(2) {
		t.Error("Expected false for out of bounds index")
	}
}

// TestRange tests the Range method for both complete and early termination scenarios.
func TestRange(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)
	s.Append(3)

	// Test complete iteration
	var results []int
	s.Range(func(index int, value int) bool {
		results = append(results, value)
		return true
	})

	expectedComplete := []int{1, 2, 3}
	if !reflect.DeepEqual(results, expectedComplete) {
		t.Errorf("Expected complete iteration results %v, got %v", expectedComplete, results)
	}

	// Test early termination
	results = []int{}
	s.Range(func(index int, value int) bool {
		results = append(results, value)
		return index < 1 // Terminate early at the second element
	})

	expectedEarly := []int{1, 2}
	if !reflect.DeepEqual(results, expectedEarly) {
		t.Errorf("Expected early termination results %v, got %v", expectedEarly, results)
	}
}

// TestConcurrency tests the concurrency safety of the Slice.
func TestConcurrency(t *testing.T) {
	s := syncslice.New[int]()
	var wg sync.WaitGroup

	// Perform concurrent appends
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			s.Append(val)
		}(i)
	}

	wg.Wait()

	if s.Length() != 1000 {
		t.Errorf("Expected length 1000, got %d", s.Length())
	}
}

func TestSetSlice(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	newSlice := []int{3, 4, 5}
	s.SetSlice(newSlice)

	if s.Length() != 3 {
		t.Errorf("Expected length 3 after SetSlice, got %d", s.Length())
	}

	if val, _ := s.Get(0); val != 3 {
		t.Errorf("Expected first element to be 3 after SetSlice, got %d", val)
	}
}

func TestGetSlice(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	sliceCopy := s.GetSlice()

	if len(sliceCopy) != 2 {
		t.Errorf("Expected copied slice length 2, got %d", len(sliceCopy))
	}

	if sliceCopy[0] != 1 || sliceCopy[1] != 2 {
		t.Errorf("Expected copied slice to be [1, 2], got %v", sliceCopy)
	}

	// Modifying the original slice to check if the copied slice remains unchanged
	s.Append(3)
	if len(sliceCopy) != 2 {
		t.Errorf("Expected copied slice length to remain 2 after modifying original, got %d", len(sliceCopy))
	}
}

// TestSetUnsafe tests the SetUnsafe method.
func TestSetUnsafe(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for out of bounds index")
		}
	}()
	s.SetUnsafe(2, 3) // should panic
}

// TestRemoveUnsafe tests the RemoveUnsafe method.
func TestRemoveUnsafe(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Did not expect panic for valid index")
		}
	}()
	s.RemoveUnsafe(1)

	if s.Length() != 1 {
		t.Errorf("Expected length 1 after removal, got %d", s.Length())
	}
}

// TestManipulateAtIndex tests the ManipulateAtIndex method.
func TestManipulateAtIndex(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	// Teste bem-sucedido de manipulação
	success := s.ManipulateAtIndex(1, func(val *int) {
		*val = 5
	})

	if !success {
		t.Errorf("Manipulation should have succeeded")
	}

	if val, _ := s.Get(1); val != 5 {
		t.Errorf("Expected value 5 at index 1, got %d", val)
	}

	if s.ManipulateAtIndex(10, func(val *int) {
		*val = 5
	}) {
		t.Errorf("Manipulation should have failed for out of bounds index")
	}
}

// TestGetUnsafe tests the GetUnsafe method.
func TestGetUnsafe(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	// Teste bem-sucedido de obtenção
	if val := s.GetUnsafe(1); val != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", val)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for out of bounds index")
		}
	}()
	_ = s.GetUnsafe(10) // should panic
}

// TestManipulateAtIndexUnsafe tests the ManipulateAtIndexUnsafe method.
func TestManipulateAtIndexUnsafe(t *testing.T) {
	s := syncslice.New[int]()
	s.Append(1)
	s.Append(2)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Did not expect panic for valid index")
		}
	}()
	s.ManipulateAtIndexUnsafe(1, func(val *int) {
		*val = 10
	})

	if val, _ := s.Get(1); val != 10 {
		t.Errorf("Expected value 10 at index 1, got %d", val)
	}
}
