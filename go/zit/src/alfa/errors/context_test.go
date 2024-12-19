package errors

import (
	"context"
	"testing"
)

func TestContextCancelled(t *testing.T) {
	ctx := MakeContext(context.Background())

	didPanic := false

	var must1, must2 bool

	defer func() {
		t.Log("defer1")

		if r := recover(); r != nil {
			t.Log("recover")
			didPanic = true

			if r != ErrContextCancelled {
				t.Errorf("expected recover to be %q", ErrContextCancelled)
			}
		}

		if !must1 || !must2 {
			t.Errorf("expected both must functions to execute")
		}
	}()

	defer ctx.Must(func() error {
		t.Log("must1")
		must1 = true
		return nil
	})

	defer ctx.Must(func() error {
		t.Log("must2")
		must2 = true
		return nil
	})

	ctx.Cancel(nil)
	ctx.Heartbeat()

	t.Errorf("expected to not get here")

	if !didPanic {
		t.Errorf("expected to panic")
	}
}
