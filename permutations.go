// Copyright 2024 Dražen Golić. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kombinat

import (
	"fmt"
	"slices"
)

func PermutationCount(n int) int {
	if n <= 0 {
		return 0
	}

	return Fac(n)
}

// Permutations without repetition by using non-recursive [Heap's algorithm].
// Returns an error if elems is empty or nil.
//
// [Heap's algorithm]: https://en.wikipedia.org/wiki/Heap's_algorithm
func Permutations[T any](elems []T) ([][]T, error) {
	n := len(elems)

	if n == 0 {
		return nil, fmt.Errorf("input slice is nil or empty")
	}

	count := PermutationCount(n)
	res := make([][]T, 0, count)

	c := make([]int, n)
	a := slices.Clone(elems)
	res = append(res, slices.Clone(elems))

	i := 1
	for i < n {
		if c[i] < i {
			if i%2 == 0 {
				s := a[0]
				a[0] = a[i]
				a[i] = s
			} else {
				s := a[c[i]]
				a[c[i]] = a[i]
				a[i] = s
			}
			res = append(res, slices.Clone(a))
			c[i]++
			i = 1
		} else {
			c[i] = 0
			i++
		}
	}

	return res, nil
}

// PermutationGenerator implements a [Generator] interface for generating
// permutations by an algorithm described in [Permutations].
type PermutationGenerator[T any] struct {
	n, i     int
	c        []int
	a, elems []T
}

// Init initializes a generator of permutations.
// Returns an error if input slice is nil or empty.
func (gen *PermutationGenerator[T]) Init(elems []T) error {
	if len(elems) == 0 {
		return fmt.Errorf("input slice is nil or empty")
	}

	gen.n = len(elems)
	gen.elems = elems
	gen.c = make([]int, gen.n)
	gen.a = slices.Clone(elems)
	gen.i = 0

	return nil
}

// Reset resets the generator to the beginning of the sequence.
func (gen *PermutationGenerator[T]) Reset() {
	s := gen.a
	gen.Init(gen.elems)
	gen.SetDest(s)
}

// Current returns the internal slice that holds the current permutation.
// If you need to modify the returned slice, use [PermutationGenerator.CurrentCopy] instead.
func (gen *PermutationGenerator[T]) Current() []T {
	return gen.a
}

// CurrentCopy returns a copy of the internal slice that holds the current permutation.
// If you don't need to modify the returned slice, use [PermutationGenerator.Current] to avoid allocation.
func (gen *PermutationGenerator[T]) CurrentCopy() []T {
	return slices.Clone(gen.a)
}

// SetDest sets a destination slice that will receive the results.
// Returns an error if there's not enough capacity in the slice.
//
// After the destination slice is set, subsequent calls to [PermutationGenerator.Current]
// will return the provided slice.
func (gen *PermutationGenerator[T]) SetDest(dest []T) error {
	if got := cap(dest); got < gen.n {
		return fmt.Errorf(capacityMsg(gen.n, got))
	}

	copy(dest, gen.a)
	gen.a = dest

	return nil
}

// Next produces a new permutation in the generator. If it returns false,
// there are no more permutations available.
func (gen *PermutationGenerator[T]) Next() bool {

	if gen.n == 0 {
		return false
	}

	if gen.i == 0 {
		gen.i = 1
		return true
	}

	for gen.i < gen.n && gen.c[gen.i] >= gen.i {
		gen.c[gen.i] = 0
		gen.i++
	}

	if gen.i >= gen.n {
		return false
	}

	if gen.i%2 == 0 {
		s := gen.a[0]
		gen.a[0] = gen.a[gen.i]
		gen.a[gen.i] = s
	} else {
		s := gen.a[gen.c[gen.i]]
		gen.a[gen.c[gen.i]] = gen.a[gen.i]
		gen.a[gen.i] = s
	}

	gen.c[gen.i]++

	if gen.i < gen.n {
		gen.i = 1
	}

	return true
}

// NewPermutationGenerator creates and initializes a new PermutationGenerator.
// Arguments and returned errors are the same ones from the [PermutationGenerator.Init] method.
func NewPermutationGenerator[T any](elems []T) (*PermutationGenerator[T], error) {
	gen := new(PermutationGenerator[T])
	err := gen.Init(elems)

	if err != nil {
		return nil, err
	}

	return gen, nil
}
