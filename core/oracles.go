package core

import (
	"strconv"
	smath "github.com/stephen-mcelhose/quantum-go/math"
)

// NewConstantOracle creates an oracle for a constant function f(x) = value.
// n is the number of input qubits.
func NewConstantOracle(n int, value int) *Oracle {
	dim := 1 << (n + 1)
	m := smath.NewMatrix(dim, dim)
	for x := 0; x < (1 << n); x++ {
		for y := 0; y < 2; y++ {
			// |y>|x> -> |y ^ value>|x>
			input := (y << n) | x
			output := ((y ^ value) << n) | x
			m.Set(output, input, 1)
		}
	}
	return NewOracle(0, m)
}

// NewBalancedOracle creates an oracle for a balanced function f(x) = x_0.
func NewBalancedOracle(n int) *Oracle {
	dim := 1 << (n + 1)
	m := smath.NewMatrix(dim, dim)
	for x := 0; x < (1 << n); x++ {
		fx := x & 1 // Use first bit as function output
		for y := 0; y < 2; y++ {
			input := (y << n) | x
			output := ((y ^ fx) << n) | x
			m.Set(output, input, 1)
		}
	}
	return NewOracle(0, m)
}

// NewInnerProductOracle creates an oracle for f(x) = s . x (mod 2).
// Used in Bernstein-Vazirani.
func NewInnerProductOracle(s string) *Oracle {
	n := len(s)
	dim := 1 << (n + 1)
	m := smath.NewMatrix(dim, dim)
	
	sBits := make([]int, n)
	for i, r := range s {
		if r == '1' {
			sBits[n-1-i] = 1 // LSB of string is at index 0
		}
	}

	for x := 0; x < (1 << n); x++ {
		fx := 0
		for i := 0; i < n; i++ {
			if (x>>i)&1 == 1 && sBits[i] == 1 {
				fx ^= 1
			}
		}
		for y := 0; y < 2; y++ {
			input := (y << n) | x
			output := ((y ^ fx) << n) | x
			m.Set(output, input, 1)
		}
	}
	return NewOracle(0, m)
}

// NewSimonOracle creates an oracle for Simon's algorithm with hidden string s.
// f(x) = f(y) iff y = x ^ s.
func NewSimonOracle(s string) *Oracle {
	n := len(s)
	sVal, _ := strconv.ParseInt(s, 2, 64)
	
	// Simon's oracle maps |x>|0> -> |x>|f(x)>
	// A simple f(x) satisfying the condition is:
	// If x < x ^ s, f(x) = x
	// If x > x ^ s, f(x) = x ^ s
	// Input register |x> (bits 0..n-1), output register |y> (bits n..2n-1)
	dim := 1 << (2 * n)
	m := smath.NewMatrix(dim, dim)
	
	for x := 0; x < (1 << n); x++ {
		fx := x
		if int64(x ^ int(sVal)) < int64(x) {
			fx = x ^ int(sVal)
		}
		for y := 0; y < (1 << n); y++ {
			// |y>|x> -> |y ^ fx>|x>
			input := (y << n) | x
			output := ((y ^ fx) << n) | x
			m.Set(output, input, 1)
		}
	}
	return NewOracle(0, m)
}
