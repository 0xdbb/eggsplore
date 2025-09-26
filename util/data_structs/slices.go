package util

import (
	"golang.org/x/exp/constraints"
)

func At[T any](v []T, i int) T {
	if i < 0 {
		i = len(v) + i
	}
	return v[i]
}

func First[T any](v []T) T {
	return At(v, 0)
}

func Last[T any](v []T) T {
	return At(v, -1)
}

func IsEmpty[T any](v []T) bool {
	return len(v) == 0
}

func AllIsEqualTo[T constraints.Ordered](slice []T, n T) bool {
	for _, v := range slice {
		if v != n {
			return false
		}
	}
	return true
}
