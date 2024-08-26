package tool_test

import (
	"testing"

	"github.com/ReanSn0w/gokit/pkg/tool"
)

func Test_Stack(t *testing.T) {
	cases := []struct {
		mode   tool.StackMode
		input  []int
		output []int
	}{
		{
			mode:   tool.StackModeFIFO,
			input:  []int{1, 2, 3, 4, 5},
			output: []int{1, 2, 3, 4, 5},
		},
		{
			mode:   tool.StackModeFILO,
			input:  []int{1, 2, 3, 4, 5},
			output: []int{5, 4, 3, 2, 1},
		},
	}

	for i, c := range cases {
		stack := tool.NewStack[int](c.mode)

		t.Logf("Test: %v, Mode: %v, Values: %v\n", i, func() string {
			switch c.mode {
			case tool.StackModeFIFO:
				return "fifo"
			default:
				return "filo"
			}
		}(), len(c.input))

		for _, val := range c.input {
			stack.Push(val)
		}

		t.Logf("stack sprint: %v\n", stack.Sprint())

		if len(c.output) != stack.Len() {
			t.Errorf("case %v faled. have: %v want: %v", i, stack.Len(), len(c.output))
		}

		for j := 0; j < len(c.output); j++ {
			value := stack.Pop()
			if value != c.output[j] {
				t.Errorf("case %v failed. element %v not eq. have: %v want: %v", i, j, value, c.output[j])
			}
		}
	}
}
