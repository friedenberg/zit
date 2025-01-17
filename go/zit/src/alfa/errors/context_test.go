package errors

import (
	ConTeXT "context"
	"testing"
)

func TestContextCancelled(t *testing.T) {
	ctx := MakeContext(ConTeXT.Background())

	var must1, must2, after1 bool

	if err := ctx.Run(
		func(ctx Context) {
			didPanic := false

			defer func() {
				t.Log("defer1")

				if r := recover(); r != nil {
					t.Log("recover")
					didPanic = true

					if r != errContextCancelled {
						t.Errorf("expected recover to be %q", errContextCancelled)
					}
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

			ctx.After(func() error {
				after1 = true
				return nil
			})

			ctx.Cancel()
			ctx.ContinueOrPanicOnDone()

			t.Errorf("expected to not get here")

			if !didPanic {
				t.Errorf("expected to panic")
			}
		},
	); err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	if !must1 || !must2 || !after1 {
		t.Errorf("expected all must and after functions to execute")
	}
}
