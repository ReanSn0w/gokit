package tool_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ReanSn0w/gokit/pkg/tool"
)

func TestNewLoop(t *testing.T) {
	task := func() {}
	loop := tool.NewLoop(task)

	if loop == nil {
		t.Error("NewLoop returned nil")
	}
}

func TestLoop_Once(t *testing.T) {
	counter := 0
	task := func() {
		counter++
	}

	loop := tool.NewLoop(task)
	loop.Once()

	if counter != 1 {
		t.Errorf("Expected counter to be 1, got %d", counter)
	}
}

func TestLoop_Run(t *testing.T) {
	counter := 0
	mutex := sync.Mutex{}
	task := func() {
		mutex.Lock()
		defer mutex.Unlock()
		counter++
	}

	loop := tool.NewLoop(task)
	loop.Run(50 * time.Millisecond)

	// Wait for a bit to allow multiple executions
	time.Sleep(200 * time.Millisecond)
	loop.Stop()

	mutex.Lock()
	defer mutex.Unlock()
	if counter < 3 {
		t.Errorf("Expected counter to be at least 3, got %d", counter)
	}
}

func TestLoop_Stop(t *testing.T) {
	counter := 0
	mutex := sync.Mutex{}
	task := func() {
		mutex.Lock()
		defer mutex.Unlock()
		counter++
	}

	loop := tool.NewLoop(task)
	loop.Run(50 * time.Millisecond)

	// Wait for a bit to allow some executions
	time.Sleep(100 * time.Millisecond)
	loop.Stop()

	// Store the current counter value
	mutex.Lock()
	currentCounter := counter
	mutex.Unlock()

	// Wait a bit more to ensure no more executions occur
	time.Sleep(100 * time.Millisecond)

	mutex.Lock()
	defer mutex.Unlock()
	if counter != currentCounter {
		t.Errorf("Expected counter to remain at %d, but it changed to %d", currentCounter, counter)
	}
}
