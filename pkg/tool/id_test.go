package tool_test

import (
	"testing"
	"time"

	"git.papkovda.ru/library/gokit/pkg/tool"
)

func TestNewID(t *testing.T) {
	id := tool.NewID()
	if len(id) != 24 {
		t.Errorf("Expected ID length to be 24, got %d", len(id))
	}
}

func TestNewIDFromTimestamp(t *testing.T) {
	timestamp := time.Now()
	id := tool.NewIDFromTimestamp(timestamp)
	if len(id) != 24 {
		t.Errorf("Expected ID length to be 24, got %d", len(id))
	}
}

func TestNewIDUniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := tool.NewID()
		if ids[id] {
			t.Errorf("Generated duplicate ID: %s", id)
		}
		ids[id] = true
	}
}

func TestNewIDFromTimestampOrder(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Second)

	id1 := tool.NewIDFromTimestamp(t1)
	id2 := tool.NewIDFromTimestamp(t2)

	if id1 >= id2 {
		t.Errorf("Expected ID1 (%s) to be less than ID2 (%s)", id1, id2)
	}
}
