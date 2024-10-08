package store

type try[T any] struct {
	value T
	err   error
}
