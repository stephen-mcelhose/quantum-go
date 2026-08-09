// Package math contains fuzz tests for matrix operations.
// These tests verify correctness across a wide range of matrix dimensions.
package math

import (
	"testing"
)

// FuzzMatrixMul uses fuzzing to test matrix multiplication with various dimensions.
// It verifies that the operation produces correct output dimensions and doesn't panic.
func FuzzMatrixMul(f *testing.F) {
	// Add seed corpus
	f.Add(2, 2, 2, 2)
	f.Fuzz(func(t *testing.T, r1, c1, r2, c2 int) {
		// Restrict dimensions to reasonable values for fuzzing
		if r1 < 1 || r1 > 10 || c1 < 1 || c1 > 10 || r2 < 1 || r2 > 10 || c2 < 1 || c2 > 10 {
			return
		}
		// Mul requires a.Cols == b.Rows
		c1 = r2

		a := NewMatrix(r1, c1)
		b := NewMatrix(r2, c2)

		// Fill with some values
		for i := range a.Data {
			a.Data[i] = complex(float64(i), 0)
		}
		for i := range b.Data {
			b.Data[i] = complex(float64(i), 1)
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Mul panicked with %v for dims (%d,%d) and (%d,%d)", r, r1, c1, r2, c2)
			}
		}()

			res := Mul(a, b)
			if res.Rows != r1 || res.Cols != c2 {

			t.Errorf("result dimensions mismatch: expected (%d,%d), got (%d,%d)", r1, c2, res.Rows, res.Cols)
		}
	})
}

// FuzzMatrixTensor uses fuzzing to test the tensor (Kronecker) product operation
// with various matrix dimensions.
func FuzzMatrixTensor(f *testing.F) {
	f.Add(2, 2, 2, 2)
	f.Fuzz(func(t *testing.T, r1, c1, r2, c2 int) {
		if r1 < 1 || r1 > 5 || c1 < 1 || c1 > 5 || r2 < 1 || r2 > 5 || c2 < 1 || c2 > 5 {
			return
		}

		a := NewMatrix(r1, c1)
		b := NewMatrix(r2, c2)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Tensor panicked with %v", r)
			}
		}()

		res := Tensor(a, b)
		if res.Rows != r1*r2 || res.Cols != c1*c2 {
			t.Errorf("result dimensions mismatch")
		}
	})
}
