package kombinat

import (
	"fmt"
	"slices"
)

func VariationCount(k, n int) int {
	return IntPow(n, k)
}

// Variations with repetitions.
// Returns an error when k <= 0 or when elems is nil or empty.
func Variations[T any](k int, elems []T) ([][]T, error) {
	n := len(elems)

	if k <= 0 {
		return nil, fmt.Errorf("k must be bigger than 0")
	}

	if n == 0 {
		return nil, fmt.Errorf("input slice is nil or empty")
	}

	count := VariationCount(k, n)

	res := make([][]T, 0, count)
	row := 0
	pows := make([]int, k)

	for i := 0; i < k; i++ {
		pows[i] = IntPow(n, k-i-1)
	}

	for {
		if row == count {
			break
		}

		ret := make([]T, k)

		for col := 0; col < k; col++ {
			i := (row / pows[col]) % n
			ret[col] = elems[i]
		}

		res = append(res, ret)

		row++
	}

	return res, nil
}

// VariationGenerator implements a [Generator] interface
// for generating variations.
type VariationGenerator[T any] struct {
	n, k, count, row int
	elems, ret       []T
	pows             []int
}

// Init initializes a generator for variations of size k out of elements
// in the elems slice. Returns an error if input slice is nil or empty,
// or if k < 1.
func (gen *VariationGenerator[T]) Init(k int, elems []T) error {
	if k <= 0 {
		return fmt.Errorf("k must be >= 1")
	}

	if len(elems) == 0 {
		return fmt.Errorf("input slice is nil or empty")
	}

	gen.n = len(elems)

	if k != gen.k {
		gen.ret = make([]T, k)
		gen.pows = make([]int, k)
		gen.count = IntPow(gen.n, k)

		for i := 0; i < k; i++ {
			gen.pows[i] = IntPow(gen.n, k-i-1)
		}
	}

	gen.row = 0
	gen.k = k
	gen.elems = elems

	return nil
}

// Reset resets the generator to the beginning of the sequence.
func (gen *VariationGenerator[T]) Reset() {
	gen.Init(gen.k, gen.elems)
}

// Current returns the internal slice that holds the current variation.
// If you need to modify the returned slice, use [CurrentCopy] instead.
func (gen *VariationGenerator[T]) Current() []T {
	return gen.ret
}

// CurrentCopy returns a copy of the internal slice that holds the current variation.
// If you don't need to modify the returned slice, use [Current] to avoid allocation.
func (gen *VariationGenerator[T]) CurrentCopy() []T {
	return slices.Clone(gen.ret)
}

// SetDest sets a destination slice that will receive the results.
// Returns an error if there's not enough capacity in the slice.
//
// After the destination slice is set, subsequent calls to [Current]
// will return the provided slice.
func (gen *VariationGenerator[T]) SetDest(dest []T) error {
	if got := cap(dest); got < gen.n {
		return fmt.Errorf(capacityMsg(gen.k, got))
	}

	copy(dest, gen.ret)
	gen.ret = dest

	return nil
}

// Next produces a new variation in the generator. If it returns false,
// there are no more variations available.
func (gen *VariationGenerator[T]) Next() bool {
	if gen.k <= 0 || gen.row == gen.count {
		return false
	}

	for col := 0; col < gen.k; col++ {
		i := (gen.row / gen.pows[col]) % gen.n
		gen.ret[col] = gen.elems[i]
	}

	gen.row++

	return true
}

// NewVariationGenerator creates and initializes a new VariationGenerator.
// Arguments and returned errors are the same ones from the [Init] method.
func NewVariationGenerator[T any](k int, elems []T) (*VariationGenerator[T], error) {
	gen := new(VariationGenerator[T])
	err := gen.Init(k, elems)

	if err != nil {
		return nil, err
	}

	return gen, nil
}
