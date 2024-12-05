package repository

type result[T any] struct {
	ok  T
	err error
}

type Order int
