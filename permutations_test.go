// Copyright 2024 Dražen Golić. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kombinat

import (
	"fmt"
	"slices"
	"testing"
)

var (
	_perm_items = []int{1, 2, 3, 4}

	_perm_table = map[int][][]int{
		0: {},
		1: {{1}},
		4: {
			{1, 2, 3, 4},
			{1, 2, 4, 3},
			{1, 3, 2, 4},
			{1, 3, 4, 2},
			{1, 4, 2, 3},
			{1, 4, 3, 2},
			{2, 1, 3, 4},
			{2, 1, 4, 3},
			{2, 3, 1, 4},
			{2, 3, 4, 1},
			{2, 4, 1, 3},
			{2, 4, 3, 1},
			{3, 1, 2, 4},
			{3, 1, 4, 2},
			{3, 2, 1, 4},
			{3, 2, 4, 1},
			{3, 4, 1, 2},
			{3, 4, 2, 1},
			{4, 1, 2, 3},
			{4, 1, 3, 2},
			{4, 2, 1, 3},
			{4, 2, 3, 1},
			{4, 3, 1, 2},
			{4, 3, 2, 1},
		},
	}
)

func TestPermutations(t *testing.T) {
	for i, want := range _perm_table {
		i := i
		want := want
		input := _perm_items[0:i]
		count := PermutationCount(len(input))

		t.Run(fmt.Sprintf("p(%d)=%d", i, count), func(t *testing.T) {
			res, _ := Permutations(input)
			slices.SortFunc(res, func(e1, e2 []int) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res, want) != 0 {
				t.Errorf("Not equal, \ngot: %v, \nwant: %v", res, want)
			}
		})
	}
}

func TestPermutationGenerator(t *testing.T) {
	for i, want := range _perm_table {
		i := i
		want := want
		input := _perm_items[0:i]
		count := PermutationCount(len(input))

		t.Run(fmt.Sprintf("p(%d)=%d", i, count), func(t *testing.T) {
			gen := new(PermutationGenerator[int])
			gen.Init(input)
			res := make([][]int, 0, count)

			for i := 0; i < count; i++ {
				gen.Next()
				res = append(res, gen.CurrentCopy())
			}

			if gen.Next() {
				t.Errorf("Didn't return false on end, dest is %v", gen.Current())
			}
			if gen.Next() {
				t.Errorf("Didn't return false on end (2), dest is %v", gen.Current())
			}

			slices.SortFunc(res, func(e1, e2 []int) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res, want) != 0 {
				t.Errorf("Not equal, \ngot: %v, \nwant: %v", res, want)
			}

			// reset + dest
			res2 := make([][]int, 0, count)
			dest := make([]int, i)
			err := gen.SetDest(dest)

			if err != nil {
				t.Errorf("%v", err)
			}

			gen.Reset()

			for i := 0; i < count; i++ {
				gen.Next()
				res2 = append(res2, slices.Clone(dest))
			}

			if gen.Next() {
				t.Errorf("Didn't return false on end after reset, dest is %v", dest)
			}
			if gen.Next() {
				t.Errorf("Didn't return false on end after reset (2), dest is %v", dest)
			}

			slices.SortFunc(res2, func(e1, e2 []int) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res2, want) != 0 {
				t.Errorf("Not equal after reset, \ngot: %v, \nwant: %v", res2, want)
			}
		})
	}
}

func BenchmarkPermutations(b *testing.B) {
	items := []int{1, 2, 3, 4, 5, 6}

	for k := 2; k <= 6; k++ {
		k := k
		b.Run(fmt.Sprintf("p(%d)=%d", k, PermutationCount(k)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Permutations(items[0:k])
			}
		})
	}
}

func BenchmarkPermutationGenerator(b *testing.B) {
	items := []int{1, 2, 3, 4, 5, 6}

	for k := 2; k <= 6; k++ {
		k := k

		b.Run(fmt.Sprintf("p(%d)=%d", k, PermutationCount(k)), func(b *testing.B) {
			gen := new(PermutationGenerator[int])

			for i := 0; i < b.N; i++ {
				gen.Init(items[0:k])
				for gen.Next() {
					gen.Current()
				}
			}
		})
	}
}
