package kombinat

import (
	"fmt"
	"slices"
	"testing"
)

var (
	_combi_items = []int{1, 2, 3, 4}

	_combi_table = map[int][][]int{
		1: {
			{1},
			{2},
			{3},
			{4},
		},
		2: {
			{1, 2},
			{1, 3},
			{1, 4},
			{2, 3},
			{2, 4},
			{3, 4},
		},
		3: {
			{1, 2, 3},
			{1, 2, 4},
			{1, 3, 4},
			{2, 3, 4},
		},
		4: {
			{1, 2, 3, 4},
		},
	}
)

func TestCombinationCount(t *testing.T) {
	if n := CombinationCount(1, 4); n != 4 {
		t.Errorf("Want 4, got %v", n)
	}
	if n := CombinationCount(2, 4); n != 6 {
		t.Errorf("Want 6, got %v", n)
	}
	if n := CombinationCount(4, 4); n != 1 {
		t.Errorf("Want 1, got %v", n)
	}
}

func TestCombinations(t *testing.T) {
	for k, want := range _combi_table {
		k := k
		want := want

		t.Run(fmt.Sprintf("c(%d,4)=%d", k, CombinationCount(k, 4)), func(t *testing.T) {
			res, _ := Combinations(k, _combi_items)
			slices.SortFunc(res, func(e1, e2 []int) int {
				return slices.Compare(e1, e2)
			})

			if compareSliceOfSlices(res, want) != 0 {
				t.Errorf("Not equal, \ngot: %v, \nwant: %v", res, want)
			}
		})
	}
}

func TestCombinationGenerator(t *testing.T) {
	for k, want := range _combi_table {
		k := k
		want := want
		count := CombinationCount(k, 4)

		t.Run(fmt.Sprintf("c(%d,4)=%d", k, count), func(t *testing.T) {
			res := make([][]int, 0, count)
			gen, _ := NewCombinationGenerator(k, _combi_items)

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
			dest := make([]int, k)
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

func BenchmarkCombinations(b *testing.B) {
	items := []int{1, 2, 3, 4, 5, 6}

	for k := 2; k <= 6; k++ {
		k := k
		b.Run(fmt.Sprintf("c(%d,6)=%d", k, CombinationCount(k, 6)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Combinations(k, items)
			}
		})
	}
}

func BenchmarkCombinationGenerator(b *testing.B) {
	items := []int{1, 2, 3, 4, 5, 6}

	for k := 2; k <= 6; k++ {
		k := k

		b.Run(fmt.Sprintf("c(%d,6)=%d", k, CombinationCount(k, 6)), func(b *testing.B) {
			gen := new(CombinationGenerator[int])

			for i := 0; i < b.N; i++ {
				gen.Init(k, items)
				for gen.Next() {
					gen.Current()
				}
			}
		})
	}
}
