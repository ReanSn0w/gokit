package tool_test

import (
	"context"
	"testing"
	"time"

	"git.papkovda.ru/library/gokit/pkg/tool"
)

func TestTokenRefresher_New(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	token := tool.NewTokenRefresher(t, "test").MustStart(ctx, time.Second, func(m tool.TokenRefresherGetSet) error {
		v, ok := m.Get(tool.MainToken)
		if !ok {
			v = 1
		}

		m.Set(tool.MainToken, v.(int)+1)
		return nil
	})

	time.Sleep(time.Second * 5)

	cancel()

	val, ok := token.Main()
	if !ok {
		t.Errorf("token is empty. val = %v", val)
		t.FailNow()
	}

	intVal, ok := val.(int)
	if !ok {
		t.Error("token val not int")
		t.FailNow()
	}

	if intVal < 4 {
		t.Errorf("val is %v, %v < 4", intVal, intVal)
	}
}
