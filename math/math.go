package math

import "math"

// Factor returns a slice consisting of the prime factors of a given integer.
func Factor(n int) []int {
	var factors []int
	if n == 0 {
		return factors
	}

	if n < 0 {
		factors = append(factors, -1)
		n *= -1
	}

	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}

	for i := 3; i <= n; i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}

	return factors
}

// Squarest finds the two closest numbers that multiply
// to n.
func Squarest(n int) (int, int) {
	if n <= 0 {
		return 0, 0
	}

	root := math.Sqrt(float64(n))
	max := int(root)
	// If we get a perfect square, then quit early.
	if max*max == n {
		return max, max
	}

	var candidate int
	candidates := []int{1}
	factors := Factor(n)
	for _, f := range factors {
		var newCandidates []int
		for _, c := range candidates {
			candidate = f * c
			if candidate > max {
				continue
			}
			newCandidates = append(newCandidates, candidate)
		}
		candidates = append(candidates, newCandidates...)
	}
	c := candidates[len(candidates)-1]
	return c, n / c
}

// NiceSquarest finds the two closest numbers that multiply
// to n and differ by less than one order of magnitude.
func NiceSquarest(n int) (int, int) {
	if n <= 0 {
		return 0, 0
	}
	var a, b int
	for a, b = Squarest(n); b/a > 10 || a < 20; n++ {
		a, b = Squarest(n)
	}
	return a, b
}
