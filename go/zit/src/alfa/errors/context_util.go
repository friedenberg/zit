package errors

import (
	"time"
)

func RunChildContextWithPrintTicker(
	parentContext Context,
	runFunc func(Context),
	printFunc func(time.Time),
	duration time.Duration,
) (err error) {
	context := MakeContext(parentContext)

	if err = context.Run(
		func(ctx Context) {
			ticker := time.NewTicker(duration)
			ctx.After(MakeNilFunc(ticker.Stop))

			go func() {
				for {
					select {
					case <-ctx.Done():
						return

					case t := <-ticker.C:
						printFunc(t)
					}
				}
			}()

			runFunc(ctx)
		},
	); err != nil {
		err = Wrap(err)
		return
	}

	return
}
