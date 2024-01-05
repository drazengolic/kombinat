package kombinat

import (
	"fmt"
	"slices"
)

func CombinationCount(m, n int) int {
	return Binom(m, n)
}

// Combinations by the Chase's Twiddle algorithm (1970), based on [implementation in C]
// by Matthew Belmonte in 1996.
//
// Returns error if the input slice is empty or nil, or if m is less than 1 or bigger than len(elems).
//
// [implementation in C]: https://web.archive.org/web/20221024045742/http://www.netlib.no/netlib/toms/382
func Combinations[T any](m int, elems []T) ([][]T, error) {
	switch {
	case m <= 0:
		return nil, fmt.Errorf("m must be >= 1")
	case len(elems) == 0:
		return nil, fmt.Errorf("empty input slice")
	case m > len(elems):
		return nil, fmt.Errorf("m is too large")
	}

	n := len(elems)

	count := CombinationCount(m, len(elems))
	res := make([][]T, 0, count)

	// init twiddle
	i := 0
	p := make([]int, n+2)

	// current combination, init to elems[n-m:n]
	c := slices.Clone(elems[n-m : n])

	p[0] = n + 1
	for i = 1; i != n-m+1; i++ {
		p[i] = 0
	}

	for i != n+1 {
		p[i] = i + m - n
		i++
	}

	p[n+1] = -2

	if m == 0 {
		p[1] = 1
	}

	res = append(res, slices.Clone(c))

	// twiddle loop
	var x, z, k int

	for {
		j := 1

		for p[j] <= 0 {
			j++
		}

		if p[j-1] == 0 {
			for i := j - 1; i != 1; i-- {
				p[i] = -1
			}
			p[j] = 0
			x = 0
			z = 0
			p[1] = 1
		} else {
			if j > 1 {
				p[j-1] = 0
			}

			j++
			for p[j] > 0 {
				j++
			}

			k = j - 1
			i = j

			for p[i] == 0 {
				p[i] = -1
				i++
			}

			if p[i] == -1 {
				p[i] = p[k]
				z = p[k] - 1
				x = i - 1
				p[k] = -1
			} else {
				if i == p[0] {
					break
				} else {
					p[j] = p[i]
					z = p[i] - 1
					p[i] = 0
					x = j - 1
				}
			}
		}

		c[z] = elems[x]
		res = append(res, slices.Clone(c))
	}

	return res, nil
}

// CombinationGenerator implements a [Generator] interface for generating combinations of
// size m out of elements in elems slice on every invocation of the [Next] method.
// For algorithm details see [Combinations].
type CombinationGenerator[T any] struct {
	n, i, x, z, k, m int
	p                []int
	c, elems         []T
	first, done      bool
}

// Init initializes a generator of combinations of size m out of elements in the elems slice.
func (gen *CombinationGenerator[T]) Init(m int, elems []T) error {

	switch {
	case m <= 0:
		return fmt.Errorf("m must be >= 1")
	case len(elems) == 0:
		return fmt.Errorf("empty input slice")
	case m > len(elems):
		return fmt.Errorf("m is too large")
	}

	gen.m = m
	gen.elems = elems

	// init twiddle
	gen.n, gen.i = len(gen.elems), 0
	gen.p = make([]int, gen.n+2)

	// current combination, init to elems[n-m:n]
	gen.c = slices.Clone(gen.elems[gen.n-gen.m : gen.n])

	gen.p[0] = gen.n + 1
	for gen.i = 1; gen.i != gen.n-gen.m+1; gen.i++ {
		gen.p[gen.i] = 0
	}

	for gen.i != gen.n+1 {
		gen.p[gen.i] = gen.i + gen.m - gen.n
		gen.i++
	}

	gen.p[gen.n+1] = -2

	if gen.m == 0 {
		gen.p[1] = 1
	}

	gen.first = true
	gen.done = false
	gen.x, gen.z, gen.k = 0, 0, 0

	return nil
}

// Reset resets the generator to the beginning of the sequence.
func (gen *CombinationGenerator[T]) Reset() {
	gen.Init(gen.m, gen.elems)
}

// Next produces a new combination in the generator. If it returns false,
// there are no more combinations available.
func (gen *CombinationGenerator[T]) Next() bool {
	if gen.done {
		return false
	}

	if gen.first {
		gen.first = false
		return true
	}

	j := 1

	for gen.p[j] <= 0 {
		j++
	}

	if gen.p[j-1] == 0 {
		for gen.i = j - 1; gen.i != 1; gen.i-- {
			gen.p[gen.i] = -1
		}
		gen.p[j] = 0
		gen.x = 0
		gen.z = 0
		gen.p[1] = 1
	} else {
		if j > 1 {
			gen.p[j-1] = 0
		}

		j++
		for gen.p[j] > 0 {
			j++
		}

		gen.k = j - 1
		gen.i = j

		for gen.p[gen.i] == 0 {
			gen.p[gen.i] = -1
			gen.i++
		}

		if gen.p[gen.i] == -1 {
			gen.p[gen.i] = gen.p[gen.k]
			gen.z = gen.p[gen.k] - 1
			gen.x = gen.i - 1
			gen.p[gen.k] = -1
		} else {
			if gen.i == gen.p[0] {
				gen.done = true
				return false
			} else {
				gen.p[j] = gen.p[gen.i]
				gen.z = gen.p[gen.i] - 1
				gen.p[gen.i] = 0
				gen.x = j - 1
			}
		}
	}

	gen.c[gen.z] = gen.elems[gen.x]
	return true
}

// Current returns the internal slice that holds the current combination.
// If you need to modify the returned slice, use [CurrentCopy] instead.
func (gen *CombinationGenerator[T]) Current() []T {
	return gen.c
}

// CurrentCopy returns a copy of the internal slice that holds the current combination.
// If you don't need to modify the returned slice, use [Current] to avoid allocation.
func (gen *CombinationGenerator[T]) CurrentCopy() []T {
	return slices.Clone(gen.c)
}

// SetDest sets a destination slice that will receive the results.
// Returns an error if there's not enough capacity in the slice.
//
// After the destination slice is set, subsequent calls to [Current]
// will return the provided slice.
func (gen *CombinationGenerator[T]) SetDest(dest []T) error {
	if got := cap(dest); got < gen.m {
		return fmt.Errorf(capacityMsg(gen.m, got))
	}

	copy(dest, gen.c)
	gen.c = dest

	return nil
}

// NewCombinationGenerator creates and initializes a new CombinationGenerator.
// Arguments and returned errors are the same ones from the [Init] method.
func NewCombinationGenerator[T any](m int, elems []T) (*CombinationGenerator[T], error) {
	gen := new(CombinationGenerator[T])
	err := gen.Init(m, elems)

	if err != nil {
		return nil, err
	}

	return gen, nil
}
