package testutil

// Ptr returns the address of the given value
func Ptr[T any](v T) *T {
	return &v
}
