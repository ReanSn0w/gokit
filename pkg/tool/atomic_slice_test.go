package tool_test

import (
	"sync"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/tool"
)

func TestAtomicSlice_Push(t *testing.T) {
	as := tool.NewAtomicSlice[int]()

	// Test pushing to end
	as.Push(-1, 1)
	as.Push(-1, 2)
	as.Push(-1, 3)
	if as.Sprint() != "[1 2 3]" {
		t.Errorf("Expected [1 2 3], got %s", as.Sprint())
	}

	// Test pushing to start
	as.Push(0, 0)
	if as.Sprint() != "[0 1 2 3]" {
		t.Errorf("Expected [0 1 2 3], got %s", as.Sprint())
	}

	// Test pushing to middle
	as.Push(2, 5)
	if as.Sprint() != "[0 1 5 2 3]" {
		t.Errorf("Expected [0 1 5 2 3], got %s", as.Sprint())
	}

	// Test pushing with out of range index
	as.Push(100, 6)
	if as.Sprint() != "[0 1 5 2 3 6]" {
		t.Errorf("Expected [0 1 5 2 3 6], got %s", as.Sprint())
	}

	// Test pushing with last element index
	as.Push(as.Len(), 7)
	if as.Sprint() != "[0 1 5 2 3 6 7]" {
		t.Errorf("Expected [0 1 5 2 3 6 7], got %s", as.Sprint())
	}
}

func TestAtomicSlice_Pop(t *testing.T) {
	as := tool.NewAtomicSlice[int]()
	as.Push(-1, 1)
	as.Push(-1, 2)
	as.Push(-1, 3)

	// Test popping from end
	v := as.Pop(-1)
	if v != 3 || as.Sprint() != "[1 2]" {
		t.Errorf("Expected 3 and [1 2], got %d and %s", v, as.Sprint())
	}

	// Test popping from start
	v = as.Pop(0)
	if v != 1 || as.Sprint() != "[2]" {
		t.Errorf("Expected 1 and [2], got %d and %s", v, as.Sprint())
	}

	// Test popping last element
	v = as.Pop(-1)
	if v != 2 || as.Sprint() != "[]" {
		t.Errorf("Expected 2 and [], got %d and %s", v, as.Sprint())
	}

	// Test popping from empty slice
	v = as.Pop(-1)
	if v != 0 || as.Sprint() != "[]" {
		t.Errorf("Expected 0 and [], got %d and %s", v, as.Sprint())
	}
}

func TestAtomicSlice_Len(t *testing.T) {
	as := tool.NewAtomicSlice[int]()
	if as.Len() != 0 {
		t.Errorf("Expected length 0, got %d", as.Len())
	}

	as.Push(-1, 1)
	as.Push(-1, 2)
	if as.Len() != 2 {
		t.Errorf("Expected length 2, got %d", as.Len())
	}

	as.Pop(-1)
	if as.Len() != 1 {
		t.Errorf("Expected length 1, got %d", as.Len())
	}
}

func TestAtomicSlice_Sort(t *testing.T) {
	as := tool.NewAtomicSlice[int]()
	as.Push(-1, 3)
	as.Push(-1, 1)
	as.Push(-1, 4)
	as.Push(-1, 2)

	as.Sort(func(i, j int) bool { return i < j })
	if as.Sprint() != "[1 2 3 4]" {
		t.Errorf("Expected [1 2 3 4], got %s", as.Sprint())
	}
}

func TestAtomicSlice_Concurrency(t *testing.T) {
	as := tool.NewAtomicSlice[int]()
	var wg sync.WaitGroup
	n := 1000

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			as.Push(-1, i)
		}(i)
	}
	wg.Wait()

	if as.Len() != n {
		t.Errorf("Expected length %d, got %d", n, as.Len())
	}
}
