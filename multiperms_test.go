// Copyright 2024 Dražen Golić. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package kombinat

import (
	// "fmt"
	"slices"
	"testing"
)

var (
	_mp_data = []struct {
		name  string
		input []string
		reps  []int
		want  [][]string
	}{
		{
			name:  "P(ABCC)",
			input: []string{"A", "B", "C"},
			reps:  []int{1, 1, 2},
			want: [][]string{
				{"C", "C", "B", "A"},
				{"A", "C", "C", "B"},
				{"C", "A", "C", "B"},
				{"C", "C", "A", "B"},
				{"B", "C", "C", "A"},
				{"C", "B", "C", "A"},
				{"A", "C", "B", "C"},
				{"C", "A", "B", "C"},
				{"B", "C", "A", "C"},
				{"A", "B", "C", "C"},
				{"B", "A", "C", "C"},
				{"C", "B", "A", "C"},
			},
		},
		{
			name:  "P(AABB)",
			input: []string{"A", "B"},
			reps:  []int{2, 2},
			want: [][]string{
				{"A", "A", "B", "B"},
				{"A", "B", "A", "B"},
				{"A", "B", "B", "A"},
				{"B", "A", "A", "B"},
				{"B", "A", "B", "A"},
				{"B", "B", "A", "A"},
			},
		},
		{
			name:  "P(A)",
			input: []string{"A"},
			reps:  []int{1},
			want: [][]string{
				{"A"},
			},
		},
		{
			name:  "P(AA)",
			input: []string{"A"},
			reps:  []int{2},
			want: [][]string{
				{"A", "A"},
			},
		},
		{
			name:  "P(AB)",
			input: []string{"A", "B"},
			reps:  []int{1, 1},
			want: [][]string{
				{"B", "A"},
				{"A", "B"},
			},
		},
	}
)

func TestMultiPermutations(t *testing.T) {

	for _, v := range _mp_data {
		v := v

		t.Run(v.name, func(t *testing.T) {
			res, err := MultiPermutations(v.input, v.reps)

			if err != nil {
				t.Errorf("got error: %v", err)
			}

			want := slices.Clone(v.want)

			slices.SortFunc(want, func(e1, e2 []string) int {
				return slices.Compare(e1, e2)
			})

			slices.SortFunc(res, func(e1, e2 []string) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res, want) != 0 {
				t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
			}
		})
	}

	t.Run("errors", func(t *testing.T) {
		_, err := MultiPermutations([]string{}, []int{2})
		if err == nil {
			t.Errorf("didn't err on empty elems")
		}

		_, err = MultiPermutations([]string{"A"}, []int{})
		if err == nil {
			t.Errorf("didn't err on empty reps")
		}

		_, err = MultiPermutations([]string{"A"}, []int{0})
		if err == nil {
			t.Errorf("didn't err on rep = 0")
		}

		_, err = MultiPermutations([]string{"A"}, []int{-1})
		if err == nil {
			t.Errorf("didn't err on rep < 0")
		}
	})

}

func TestMultiPermutationsCount(t *testing.T) {
	if c := MultiPermutationsCount([]int{2, 1, 3, 2}); c != 1680 {
		t.Errorf("got: %v\n, want: 1680", c)
	}

	if MultiPermutationsCount([]int{1, 1, 1}) != PermutationCount(3) {
		t.Errorf("not equivalent")
	}
}

func TestMultiPermutationGenerator(t *testing.T) {
	for _, v := range _mp_data {
		v := v

		t.Run(v.name, func(t *testing.T) {
			gen, err := NewMultiPermutationGenerator(v.input, v.reps)

			if err != nil {
				t.Errorf("got error: %v", err)
			}
			res := make([][]string, 0, len(v.want))

			for gen.Next() {
				res = append(res, gen.CurrentCopy())
			}

			if gen.Next() {
				t.Errorf("didn't return false on end, dest is %v", gen.Current())
			}
			if gen.Next() {
				t.Errorf("didn't return false on end (2), dest is %v", gen.Current())
			}

			want := slices.Clone(v.want)

			slices.SortFunc(want, func(e1, e2 []string) int {
				return slices.Compare(e1, e2)
			})
			slices.SortFunc(res, func(e1, e2 []string) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res, want) != 0 {
				t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
			}

			// reset + dest
			count := MultiPermutationsCount(v.reps)
			res2 := make([][]string, 0, count)
			dest := make([]string, len(v.want[0]))
			err = gen.SetDest(dest)

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

			slices.SortFunc(res2, func(e1, e2 []string) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res2, want) != 0 {
				t.Errorf("Not equal after reset, \ngot: %v, \nwant: %v", res2, want)
			}
		})
	}

	t.Run("errors", func(t *testing.T) {
		_, err := NewMultiPermutationGenerator([]string{}, []int{2})
		if err == nil {
			t.Errorf("didn't err on empty elems")
		}

		_, err = NewMultiPermutationGenerator([]string{"A"}, []int{})
		if err == nil {
			t.Errorf("didn't err on empty reps")
		}

		_, err = NewMultiPermutationGenerator([]string{"A"}, []int{0})
		if err == nil {
			t.Errorf("didn't err on rep = 0")
		}

		_, err = NewMultiPermutationGenerator([]string{"A"}, []int{-1})
		if err == nil {
			t.Errorf("didn't err on rep < 0")
		}
	})
}

func BenchmarkMultiPermutations(b *testing.B) {
	table := []struct {
		name  string
		reps  []int
		input []string
	}{
		{"P(AABCCCDD)", []int{2, 1, 3, 2}, []string{"A", "B", "C", "D"}},
		{"P(AABB)", []int{2, 2}, []string{"A", "B"}},
		{"P(ABCD)", []int{1, 1, 1, 1}, []string{"A", "B", "C", "D"}},
	}

	for _, v := range table {
		v := v
		b.Run(v.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				MultiPermutations(v.input, v.reps)
			}
		})
	}
}

func BenchmarkMultiPermutationGenerator(b *testing.B) {
	gen := new(MultiPermutationGenerator[string])

	table := []struct {
		name  string
		reps  []int
		input []string
	}{
		{"P(AABCCCDD)", []int{2, 1, 3, 2}, []string{"A", "B", "C", "D"}},
		{"P(AABB)", []int{2, 2}, []string{"A", "B"}},
		{"P(ABCD)", []int{1, 1, 1, 1}, []string{"A", "B", "C", "D"}},
	}

	for _, v := range table {
		v := v
		b.Run(v.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gen.Init(v.input, v.reps)
				for gen.Next() {
				}
			}
		})
	}
}
