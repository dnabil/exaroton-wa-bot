package helper

func If[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}

func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

func Ptr[T any](v T) *T {
	return &v
}

func Deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
