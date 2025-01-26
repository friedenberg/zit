package quiter

type ElementOrError[E any] struct {
	Element E
	Error   error
}

