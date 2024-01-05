package kombinat

import (
	// "fmt"
	"slices"
	"testing"
)

func TestMultiPermutations(t *testing.T) {
	t.Run("P(1233)", func(t *testing.T) {
		res, err := MultiPermutations([]int{1, 2, 3}, []int{1, 1, 2})

		if err != nil {
			t.Errorf("got error: %v", err)
		}

		want := [][]int{
			{3, 3, 2, 1},
			{1, 3, 3, 2},
			{3, 1, 3, 2},
			{3, 3, 1, 2},
			{2, 3, 3, 1},
			{3, 2, 3, 1},
			{1, 3, 2, 3},
			{3, 1, 2, 3},
			{2, 3, 1, 3},
			{1, 2, 3, 3},
			{2, 1, 3, 3},
			{3, 2, 1, 3},
		}

		slices.SortFunc(want, func(e1, e2 []int) int {
			return slices.Compare(e1, e2)
		})
		slices.SortFunc(res, func(e1, e2 []int) int {
			return slices.Compare(e1, e2)
		})

		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

	t.Run("P(AABB)", func(t *testing.T) {
		res, err := MultiPermutations([]string{"A", "B"}, []int{2, 2})

		if err != nil {
			t.Errorf("got error: %v", err)
		}

		want := [][]string{
			{"A", "A", "B", "B"},
			{"A", "B", "A", "B"},
			{"A", "B", "B", "A"},
			{"B", "A", "A", "B"},
			{"B", "A", "B", "A"},
			{"B", "B", "A", "A"},
		}

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

	t.Run("P(A)", func(t *testing.T) {
		res, err := MultiPermutations([]string{"A"}, []int{1})

		if err != nil {
			t.Errorf("got error: %v", err)
		}
		want := [][]string{{"A"}}
		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

	t.Run("P(AA)", func(t *testing.T) {
		res, err := MultiPermutations([]string{"A"}, []int{2})
		if err != nil {
			t.Errorf("got error: %v", err)
		}
		want := [][]string{{"A", "A"}}
		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

	t.Run("P(AB)", func(t *testing.T) {
		res, err := MultiPermutations([]string{"A", "B"}, []int{1, 1})
		if err != nil {
			t.Errorf("got error: %v", err)
		}
		want := [][]string{{"B", "A"}, {"A", "B"}}
		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

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
	t.Run("P(1233)", func(t *testing.T) {
		gen, err := NewMultiPermutationGenerator([]int{1, 2, 3}, []int{1, 1, 2})

		if err != nil {
			t.Errorf("got error: %v", err)
		}

		want := [][]int{
			{3, 3, 2, 1},
			{1, 3, 3, 2},
			{3, 1, 3, 2},
			{3, 3, 1, 2},
			{2, 3, 3, 1},
			{3, 2, 3, 1},
			{1, 3, 2, 3},
			{3, 1, 2, 3},
			{2, 3, 1, 3},
			{1, 2, 3, 3},
			{2, 1, 3, 3},
			{3, 2, 1, 3},
		}

		res := make([][]int, 0, 12)

		for gen.Next() {
			res = append(res, gen.CurrentCopy())
		}

		if gen.Next() {
			t.Errorf("didn't return false on end, dest is %v", gen.Current())
		}
		if gen.Next() {
			t.Errorf("didn't return false on end (2), dest is %v", gen.Current())
		}

		slices.SortFunc(want, func(e1, e2 []int) int {
			return slices.Compare(e1, e2)
		})
		slices.SortFunc(res, func(e1, e2 []int) int {
			return slices.Compare(e1, e2)
		})

		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

	t.Run("P(A)", func(t *testing.T) {
		gen, err := NewMultiPermutationGenerator([]string{"A"}, []int{1})

		if err != nil {
			t.Errorf("got error: %v", err)
		}
		res := make([][]string, 0, 1)

		for gen.Next() {
			res = append(res, gen.CurrentCopy())
		}

		if gen.Next() {
			t.Errorf("didn't return false on end, dest is %v", gen.Current())
		}
		if gen.Next() {
			t.Errorf("didn't return false on end (2), dest is %v", gen.Current())
		}

		want := [][]string{{"A"}}
		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

	t.Run("P(AA)", func(t *testing.T) {
		gen, err := NewMultiPermutationGenerator([]string{"A"}, []int{2})

		if err != nil {
			t.Errorf("got error: %v", err)
		}

		want := [][]string{{"A", "A"}}
		res := make([][]string, 0, 1)

		for gen.Next() {
			res = append(res, gen.CurrentCopy())
		}
		if gen.Next() {
			t.Errorf("didn't return false on end, dest is %v", gen.Current())
		}
		if gen.Next() {
			t.Errorf("didn't return false on end (2), dest is %v", gen.Current())
		}
		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

	t.Run("P(AB)", func(t *testing.T) {
		gen, err := NewMultiPermutationGenerator([]string{"A", "B"}, []int{1, 1})

		if err != nil {
			t.Errorf("got error: %v", err)
		}

		want := [][]string{{"B", "A"}, {"A", "B"}}
		res := make([][]string, 0, 1)

		for gen.Next() {
			res = append(res, gen.CurrentCopy())
		}
		if gen.Next() {
			t.Errorf("didn't return false on end, dest is %v", gen.Current())
		}
		if gen.Next() {
			t.Errorf("didn't return false on end (2), dest is %v", gen.Current())
		}
		if compareSliceOfSlices(res, want) != 0 {
			t.Errorf("not equal, \ngot: %v, \nwant: %v", res, want)
		}
	})

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
	elems := []string{"A", "B", "C", "D"}
	reps := []int{2, 1, 3, 2}

	for i := 0; i < b.N; i++ {
		MultiPermutations(elems, reps)
	}
}

func BenchmarkMultiPermutationGenerator(b *testing.B) {
	elems := []string{"A", "B", "C", "D"}
	reps := []int{2, 1, 3, 2}
	gen := new(MultiPermutationGenerator[string])

	for i := 0; i < b.N; i++ {
		gen.Init(elems, reps)
		for gen.Next() {
			gen.Current()
		}
	}
}
