package utils

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](l, r T) T {
	if l > r {
		return l
	}
	return r
}

func Min[T constraints.Ordered](l, r T) T {
	if l < r {
		return l
	}
	return r
}
