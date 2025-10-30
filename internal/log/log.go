package log

import (
	"log/slog"
	"reflect"
)

// ErrAttr returns an slog.Attr for error with `error` key
func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

// SafeValueAttr returns an slog.Attr that safely dereferences v if it's non-nil.
func SafeValueAttr[T any](k string, v *T) slog.Attr {
	if v == nil {
		return slog.Any(k, v)
	}
	return slog.Any(k, *v)
}

// SafeValuesAttr returns an slog.Attr that safely dereferences any pointer elements.
func SafeValuesAttr(k string, vs []any) slog.Attr {
	o := make([]any, len(vs))
	for i, v := range vs {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Pointer && !rv.IsNil() {
			o[i] = rv.Elem().Interface()
		} else {
			o[i] = v
		}
	}
	return slog.Any(k, o)
}
