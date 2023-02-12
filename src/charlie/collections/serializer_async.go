package collections

type AsyncSerializer[T any] struct {
	chError <-chan error
	chE     chan<- T
	chDone  <-chan struct{}
}

func MakeAsyncSerializer[T any](wf WriterFunc[T]) AsyncSerializer[T] {
	chError := make(chan error)
	chE := make(chan T)
	chDone := make(chan struct{})

	go func(chError chan<- error, chE <-chan T, chDone chan<- struct{}) {
		defer func() {
			chDone <- struct{}{}
		}()

		for e := range chE {
			// TODO what how error
			err := wf(e)
			if err != nil {
				chError <- err
			}
		}
	}(chError, chE, chDone)

	return AsyncSerializer[T]{
		chError: chError,
		chE:     chE,
		chDone:  chDone,
	}
}

func (s AsyncSerializer[T]) Do(e T) (err error) {
	select {
	case err = <-s.chError:
	case s.chE <- e:
	}

	return
}

func (s AsyncSerializer[T]) Wait() (err error) {
	select {
	case err = <-s.chError:
	}

	close(s.chE)
	<-s.chDone

	return
}
