package util

type Pair[T, U any] struct {
	A T
	B U
}

type KeyValPair[T, U any] struct {
	Key T
	Val U
}
