package kombinat

import (
	"fmt"
	"slices"
)

func MultiPermutationsCount(reps []int) int {
	n := 0
	for _, k := range reps {
		n += k
	}

	n = Fac(n)

	for _, k := range reps {
		n /= Fac(k)
	}

	return n
}

// This function generates permutations of a multiset by following
// the algorithm described below. It accepts slices
//   - elems: set of elements to make permutations from
//   - reps: number of repetitions per element, must be 1 or higher.
//
// Number of items in both slices must match, as reps correspond to the elems.
//
// [Algorithm 1]
// Visits the permutations of multiset E. The permutations are stored
// in a singly-linked list pointed to by head pointer h. Each node in the linked
// list has a value field v and a next field n. The init(E) call creates a
// singly-linked list storing the elements of E in non-increasing order with h, i,
// and j pointing to its first, second-last, and last nodes, respectively. The
// null pointer is given by φ. Note: If E is empty, then init(E) should exit.
// Also, if E contains only one element, then init(E) does not need to provide a
// value for i.
//
//	[h, i, j] ← init(E)
//	visit(h)
//	while j.n ≠ φ orj.v < h.v do
//		if j.n ≠ φ and i.v ≥ j.n.v then
//	    	s←j
//		else
//	    	s←i
//		end if
//		t ← s.n
//		s.n ← t.n
//		t.n ← h
//		if t.v < h.v then
//	    	i←t
//		end if
//		j←i.n
//		h←t
//		visit(h)
//	end while
//
// Found in "Loopless Generation of Multiset Permutations using a Constant Number
// of Variables by Prefix Shifts."  Aaron Williams, 2009
//
// For implementations in other languages, check out [multipermute].
//
// [Algorithm 1]: https://dl.acm.org/doi/10.5555/1496770.1496877
// [multipermute]: https://github.com/ekg/multipermute
func MultiPermutations[T any](elems []T, reps []int) ([][]T, error) {
	lelems := len(elems)
	lreps := len(reps)

	if lelems == 0 || lreps == 0 {
		return nil, fmt.Errorf("empty input slice(s)")
	}

	if lelems != lreps {
		return nil, fmt.Errorf("input lengths do not match")
	}

	if lelems == 1 && lreps == 1 && reps[0] == 1 {
		return [][]T{{elems[0]}}, nil
	}

	n := 0

	for i := 0; i < len(reps); i++ {
		if reps[i] <= 0 {
			return nil, fmt.Errorf("value of a rep must be >= 1")
		}

		n += reps[i]
	}

	// init

	res := make([][]T, 0, n)
	init := make([]int, 0, n)

	for i := 0; i < len(reps); i++ {
		for j := 0; j < reps[i]; j++ {
			init = append(init, i)
		}
	}

	h := &listElement[int]{init[0], nil}

	for _, el := range init[1:] {
		h = &listElement[int]{el, h}
	}

	i := h.nth(len(init) - 2)
	j := h.nth(len(init) - 1)
	var s, t *listElement[int]

	temp := make([]T, n)
	dump(h, elems, n, temp)
	res = append(res, slices.Clone(temp))

	// loop

	for j.next != nil || j.value < h.value {

		if j.next != nil && i.value >= j.next.value {
			s = j
		} else if i.next != nil { // <- fix for [A,B]
			s = i
		} else {
			s = h
		}

		t = s.next
		s.next = t.next
		t.next = h

		if t.value < h.value {
			i = t
		}

		j = i.next
		h = t

		dump(h, elems, n, temp)
		res = append(res, slices.Clone(temp))
	}

	return res, nil
}

// MultiPermutationGenerator implements a [Generator] interface for generating multiset
// permutations by an algorithm described in [MultiPermutations].
type MultiPermutationGenerator[T any] struct {
	elems, dest             []T
	reps                    []int
	n                       int
	i, j, s, t, h           *listElement[int]
	first, returned, single bool
}

// Init initializes a generator of multiset permutations.
//
// It accepts slices
//   - elems: set of elements to make permutations from
//   - reps: number of repetitions per element, must be 1 or higher.
//
// Number of items in both slices must match, as reps correspond to the elems.
func (gen *MultiPermutationGenerator[T]) Init(elems []T, reps []int) error {
	lelems := len(elems)
	lreps := len(reps)

	if lelems == 0 || lreps == 0 {
		return fmt.Errorf("empty input slice(s)")
	}

	if lelems != lreps {
		return fmt.Errorf("input lengths do not match")
	}

	gen.single = lelems == 1 && lreps == 1

	n := 0

	for i := 0; i < len(reps); i++ {
		if reps[i] <= 0 {
			return fmt.Errorf("value of a rep must be >= 1")
		}

		n += reps[i]
	}

	init := make([]int, 0, n)

	if len(gen.dest) != n {
		gen.dest = make([]T, n)
	}

	for i := 0; i < len(reps); i++ {
		for j := 0; j < reps[i]; j++ {
			init = append(init, i)
		}
	}

	gen.h = &listElement[int]{init[0], nil}

	for _, el := range init[1:] {
		gen.h = &listElement[int]{el, gen.h}
	}

	gen.i = gen.h.nth(len(init) - 2)
	gen.j = gen.h.nth(len(init) - 1)

	gen.first = true
	gen.returned = false
	gen.n = n
	gen.elems = elems
	gen.reps = reps

	return nil
}

// Reset resets the generator to the beginning of the sequence.
func (gen *MultiPermutationGenerator[T]) Reset() {
	gen.Init(gen.elems, gen.reps)
}

// Current returns the internal slice that holds the current permutation.
// If you need to modify the returned slice, use [CurrentCopy] instead.
func (gen *MultiPermutationGenerator[T]) Current() []T {
	return gen.dest
}

// CurrentCopy returns a copy of the internal slice that holds the current permutation.
// If you don't need to modify the returned slice, use [Current] to avoid allocation.
func (gen *MultiPermutationGenerator[T]) CurrentCopy() []T {
	return slices.Clone(gen.dest)
}

// SetDest sets a destination slice that will receive the results.
// Returns an error if there's not enough capacity in the slice.
//
// After the destination slice is set, subsequent calls to [Current]
// will return the provided slice.
func (gen *MultiPermutationGenerator[T]) SetDest(dest []T) error {
	if got := cap(dest); got < gen.n {
		return fmt.Errorf(capacityMsg(gen.n, got))
	}

	copy(dest, gen.dest)
	gen.dest = dest

	return nil
}

// Next produces a new permutation in the generator. If it returns false,
// there are no more permutations available.
func (gen *MultiPermutationGenerator[T]) Next() bool {
	// single element case
	if gen.single {
		if !gen.returned {
			for i := 0; i < gen.reps[0]; i++ {
				gen.dest[i] = gen.elems[0]
			}
			gen.returned = true
			return true
		}

		return false
	}

	if gen.first {
		dump(gen.h, gen.elems, gen.n, gen.dest)
		gen.first = false
		return true
	}

	if gen.j.next != nil || gen.j.value < gen.h.value {

		if gen.j.next != nil && gen.i.value >= gen.j.next.value {
			gen.s = gen.j
		} else if gen.i.next != nil { // <- fix for [A,B]
			gen.s = gen.i
		} else {
			gen.s = gen.h
		}

		gen.t = gen.s.next
		gen.s.next = gen.t.next
		gen.t.next = gen.h

		if gen.t.value < gen.h.value {
			gen.i = gen.t
		}

		gen.j = gen.i.next
		gen.h = gen.t

		dump(gen.h, gen.elems, gen.n, gen.dest)

		return true
	}

	return false
}

// NewMultiPermutationGenerator creates and initializes a new MultiPermutationGenerator.
// Arguments and returned errors are the same ones from the [Init] method.
func NewMultiPermutationGenerator[T any](elems []T, reps []int) (*MultiPermutationGenerator[T], error) {
	gen := new(MultiPermutationGenerator[T])
	err := gen.Init(elems, reps)

	if err != nil {
		return nil, err
	}

	return gen, nil
}

// simple linked list element
type listElement[T any] struct {
	value T
	next  *listElement[T]
}

// nth element after
func (el *listElement[T]) nth(n int) *listElement[T] {
	if el == nil {
		return nil
	}

	v := el.next

	for i := 1; i < n; i++ {
		v = v.next

		if v == nil {
			break
		}
	}

	return v
}

// dump list into a dest slice
func dump[T any](e *listElement[int], items []T, max int, dest []T) {
	if e == nil {
		return
	}

	dest[0] = items[e.value]

	v := e.next

	for i := 1; i < max && v != nil; i++ {
		dest[i] = items[v.value]
		v = v.next
	}
}
