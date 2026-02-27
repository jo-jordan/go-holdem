package internal

func Map[T any, R any](in []T, fn func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}

func Reduce[T any, R any](in []T, acc R, fn func(R, T) R) R {
	for _, v := range in {
		acc = fn(acc, v)
	}
	return acc
}
