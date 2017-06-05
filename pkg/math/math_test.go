package math

import (
	"reflect"
	"testing"
)

func TestFactor(t *testing.T) {
	var f []int
	type testCase struct {
		i       int
		factors []int
	}
	cases := []testCase{
		testCase{
			i:       -6,
			factors: []int{-1, 2, 3},
		},
		testCase{
			i:       0,
			factors: []int{},
		},
		testCase{
			i:       1,
			factors: []int{},
		},
		testCase{
			i:       2,
			factors: []int{2},
		},
		testCase{
			i:       13,
			factors: []int{13},
		},
		testCase{
			i:       27,
			factors: []int{3, 3, 3},
		},
		testCase{
			i:       30,
			factors: []int{2, 3, 5},
		},
		testCase{
			i:       33,
			factors: []int{3, 11},
		},
		testCase{
			i:       61,
			factors: []int{61},
		},
	}
	for i, c := range cases {
		f = Factor(c.i)
		if len(f) == 0 && len(c.factors) == 0 {
			continue
		}
		if !reflect.DeepEqual(f, c.factors) {
			t.Errorf("test case %d: expected %v, got %v", i, c.factors, f)
		}
	}
}

func TestSquarest(t *testing.T) {
	var a, b int
	type testCase struct {
		i       int
		factors [2]int
	}
	cases := []testCase{
		testCase{
			i:       -6,
			factors: [2]int{0, 0},
		},
		testCase{
			i:       0,
			factors: [2]int{0, 0},
		},
		testCase{
			i:       1,
			factors: [2]int{1, 1},
		},
		testCase{
			i:       2,
			factors: [2]int{1, 2},
		},
		testCase{
			i:       13,
			factors: [2]int{1, 13},
		},
		testCase{
			i:       27,
			factors: [2]int{3, 9},
		},
		testCase{
			i:       30,
			factors: [2]int{5, 6},
		},
		testCase{
			i:       33,
			factors: [2]int{3, 11},
		},
		testCase{
			i:       61,
			factors: [2]int{1, 61},
		},
	}
	for i, c := range cases {
		a, b = Squarest(c.i)
		if a != c.factors[0] || b != c.factors[1] {
			t.Errorf("test case %d: expected %d, %d and got %d, %d", i, c.factors[0], c.factors[1], a, b)
		}
	}
}
