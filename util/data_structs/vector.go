package util

import (
	"slices"
)

type Vector[T any] []T

func NewVector[T any](args ...int) Vector[T] {
	var size = 0
	var capacity = 32
	if len(args) > 0 {
		size = args[0]
		if len(args) > 1 {
			capacity = args[1]
		}
	}
	return make(Vector[T], size, capacity)
}

func (v *Vector[T]) Append(x T) *Vector[T] {
	*v = append(*v, x)
	return v
}

func (v *Vector[T]) Extend(other []T) *Vector[T] {
	for _, o := range other {
		*v = append(*v, o)
	}
	return v
}
func (v *Vector[T]) ExtendWith(other ...T) *Vector[T] {
	*v = append(*v, other...)
	return v
}

func (v *Vector[T]) Pop() T {
	var n = len(*v) - 1
	var val = (*v)[n]
	*v = (*v)[:n]
	return val
}

func (v *Vector[T]) Shift() T {
	var val = (*v)[0]
	*v = (*v)[1:]
	return val
}

func (v *Vector[T]) At(i int) T {
	if i < 0 {
		i = len(*v) + i
	}
	return (*v)[i]
}

func (v *Vector[T]) First() T {
	return v.At(0)
}

func (v *Vector[T]) Last() T {
	return v.At(-1)
}

func (v *Vector[T]) IsEmpty() bool {
	return len(*v) == 0
}

func (v *Vector[T]) Sort(cmp func(a, b T) int) *Vector[T] {
	slices.SortFunc(*v, cmp)
	return v
}
