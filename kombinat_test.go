package kombinat

import (
	"cmp"
	"slices"
	"testing"
)

// Utility function for comparing slices of slices
func compareSliceOfSlices[T cmp.Ordered](s1, s2 [][]T) int {
	return slices.CompareFunc(s1, s2, func(e1, e2 []T) int {
		return slices.Compare(e1, e2)
	})
}

func TestIntPow(t *testing.T) {
	if p := IntPow(3, 2); p != 9 {
		t.Errorf("IntPow error, want: 9, got: %v", p)
	}
	if p := IntPow(1, 2); p != 1 {
		t.Errorf("IntPow error, want: 1, got: %v", p)
	}
	if p := IntPow(10, 4); p != 10000 {
		t.Errorf("IntPow error, want: 10000, got: %v", p)
	}
}

func TestFactorial(t *testing.T) {
	if a := Fac(-2); a != 1 {
		t.Errorf("Factorial error, want: 1, got: %v", a)
	}
	if a := Fac(0); a != 1 {
		t.Errorf("Factorial error, want: 1, got: %v", a)
	}
	if a := Fac(1); a != 1 {
		t.Errorf("Factorial error, want: 1, got: %v", a)
	}
	if a := Fac(2); a != 2 {
		t.Errorf("Factorial error, want: 2, got: %v", a)
	}
	if a := Fac(7); a != 5040 {
		t.Errorf("Factorial error, want: 5040, got: %v", a)
	}
}

func TestBinom(t *testing.T) {
	if a := Binom(1, 4); a != 4 {
		t.Errorf("Binom error, want: 4, got: %v", a)
	}
	if a := Binom(2, 4); a != 6 {
		t.Errorf("Binom error, want: 6, got: %v", a)
	}
	if a := Binom(3, 4); a != 4 {
		t.Errorf("Binom error, want: 4, got: %v", a)
	}
	if a := Binom(4, 4); a != 1 {
		t.Errorf("Binom error, want: 1, got: %v", a)
	}
}
