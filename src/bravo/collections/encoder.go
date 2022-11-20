package collections

type EncoderLike[T any] interface {
	Encode(*T) (int64, error)
}

func EncoderToWriter[T any](el EncoderLike[T]) WriterFunc[*T] {
  return func(e *T) (err error) {
    _, err = el.Encode(e)
    return
  }
}
