package sidecar

func asSliceNoZeroes[T comparable](values ...T) []T {
	var out []T
	var zero T

	for _, s := range values {
		if s != zero {
			out = append(out, s)
		}
	}

	return out
}
