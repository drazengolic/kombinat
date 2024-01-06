package kombinat

import (
	"fmt"
	"slices"
	"testing"
)

var (
	_var_data = []struct {
		input []int
		want  [][]int
		k     int
	}{
		{
			input: []int{1, 2, 3},
			want:  [][]int{{1}, {2}, {3}},
			k:     1,
		},
		{
			input: []int{1, 2, 3},
			want: [][]int{
				{1, 1},
				{1, 2},
				{1, 3},
				{2, 1},
				{2, 2},
				{2, 3},
				{3, 1},
				{3, 2},
				{3, 3},
			},
			k: 2,
		},
		{
			input: []int{1, 2},
			want: [][]int{
				{1, 1, 1},
				{1, 1, 2},
				{1, 2, 1},
				{1, 2, 2},
				{2, 1, 1},
				{2, 1, 2},
				{2, 2, 1},
				{2, 2, 2},
			},
			k: 3,
		},
		{
			input: []int{1, 2, 3, 4},
			want: [][]int{
				{1, 1, 1},
				{1, 1, 2},
				{1, 1, 3},
				{1, 1, 4},
				{1, 2, 1},
				{1, 2, 2},
				{1, 2, 3},
				{1, 2, 4},
				{1, 3, 1},
				{1, 3, 2},
				{1, 3, 3},
				{1, 3, 4},
				{1, 4, 1},
				{1, 4, 2},
				{1, 4, 3},
				{1, 4, 4},

				{2, 1, 1},
				{2, 1, 2},
				{2, 1, 3},
				{2, 1, 4},
				{2, 2, 1},
				{2, 2, 2},
				{2, 2, 3},
				{2, 2, 4},
				{2, 3, 1},
				{2, 3, 2},
				{2, 3, 3},
				{2, 3, 4},
				{2, 4, 1},
				{2, 4, 2},
				{2, 4, 3},
				{2, 4, 4},

				{3, 1, 1},
				{3, 1, 2},
				{3, 1, 3},
				{3, 1, 4},
				{3, 2, 1},
				{3, 2, 2},
				{3, 2, 3},
				{3, 2, 4},
				{3, 3, 1},
				{3, 3, 2},
				{3, 3, 3},
				{3, 3, 4},
				{3, 4, 1},
				{3, 4, 2},
				{3, 4, 3},
				{3, 4, 4},

				{4, 1, 1},
				{4, 1, 2},
				{4, 1, 3},
				{4, 1, 4},
				{4, 2, 1},
				{4, 2, 2},
				{4, 2, 3},
				{4, 2, 4},
				{4, 3, 1},
				{4, 3, 2},
				{4, 3, 3},
				{4, 3, 4},
				{4, 4, 1},
				{4, 4, 2},
				{4, 4, 3},
				{4, 4, 4},
			},
			k: 3,
		},
	}
)

func TestVariations(t *testing.T) {
	for _, d := range _var_data {
		d := d

		t.Run(fmt.Sprintf("v(%d,%d)", d.k, len(d.input)), func(t *testing.T) {
			res, err := Variations(d.k, d.input)

			if err != nil {
				t.Errorf("Error'd with: %v", err)
			}

			if compareSliceOfSlices(res, d.want) != 0 {
				t.Errorf("Not equal, \ngot: %v, \nwant: %v", res, d.want)
			}
		})
	}
}

func TestVariationGenerator(t *testing.T) {
	for _, vd := range _var_data {
		vd := vd

		t.Run(fmt.Sprintf("v(%d,%d)", vd.k, len(vd.input)), func(t *testing.T) {
			gen, err := NewVariationGenerator(vd.k, vd.input)

			if err != nil {
				t.Errorf("Error'd with: %v", err)
			}

			for i, w := range vd.want {
				if gen.Next(); slices.Compare(w, gen.Current()) != 0 {
					t.Errorf("Not equal at %v, \ngot: %v, \nwant: %v", i, gen.Current(), w)
				}
			}
			if gen.Next() {
				t.Errorf("Didn't return false on end, dest is %v", gen.Current())
			}
			if gen.Next() {
				t.Errorf("Didn't return false on end (2), dest is %v", gen.Current())
			}

			dest := make([]int, vd.k)
			err = gen.SetDest(dest)

			if err != nil {
				t.Errorf("%v", err)
			}

			gen.Reset()

			for i, w := range vd.want {
				if gen.Next(); slices.Compare(w, dest) != 0 {
					t.Errorf("Not equal at %v after reset, \ngot: %v, \nwant: %v", i, dest, w)
				}
			}
			if gen.Next() {
				t.Errorf("Didn't return false on end after reset, dest is %v", dest)
			}
			if gen.Next() {
				t.Errorf("Didn't return false on end after reset (2), dest is %v", dest)
			}
		})
	}
}

func BenchmarkVariations(b *testing.B) {
	items := []int{1, 2, 3, 4}

	for k := 2; k <= 6; k++ {
		k := k

		b.Run(fmt.Sprintf("v(%d,4)=%d", k, VariationCount(k, 4)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Variations(k, items)
			}
		})
	}
}

func BenchmarkVariationGenerator(b *testing.B) {
	items := []int{1, 2, 3, 4}

	for k := 2; k <= 6; k++ {
		k := k
		gen := new(VariationGenerator[int])

		b.Run(fmt.Sprintf("v(%d,4)=%d", k, VariationCount(k, 4)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gen.Init(k, items)
				for gen.Next() {
					gen.Current()
				}
			}
		})
	}
}
