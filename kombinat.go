// Copyright 2024 Dražen Golić. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kombinat

import (
	"fmt"
)

type Generator[T any] interface {
	Current() []T
	CurrentCopy() []T
	Next() bool
	Reset()
	SetDest([]T) error
}

// IntPow calculates n^m as integer
func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}

	result := n

	for i := 2; i <= m; i++ {
		result *= n
	}

	return result
}

// Fac calculates factorial of n (n!)
func Fac(n int) int {
	ret := 1

	for n > 1 {
		ret *= n
		n--
	}

	return ret
}

// Binom calculates the [binomial coefficient].
//
// [binomial coefficient]: https://en.wikipedia.org/wiki/Binomial_coefficient
func Binom(k, n int) int {
	num := 1
	denom := 1

	for k > 0 {
		num *= n
		denom *= k
		n--
		k--
	}

	return num / denom
}

// Creates a low capacity message for generators to panic about it.
func capacityMsg(need, got int) string {
	return fmt.Sprintf("Not enough capacity in the destination slice (need %d, got %d)", need, got)
}
