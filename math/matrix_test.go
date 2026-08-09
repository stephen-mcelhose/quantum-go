// Package math contains tests for complex matrix operations.
// It verifies correctness of matrix multiplication, tensor products,
// conjugate transpose, and other linear algebra operations.
package math

import (
	"math/cmplx"
	"testing"
)

// TestIdentityMatrix verifies that IdentityMatrix creates a proper identity matrix
// with 1's on the diagonal and 0's elsewhere.
func TestIdentityMatrix(t *testing.T) {
	dim := 3
	m := IdentityMatrix(dim)
	if m.Rows != dim || m.Cols != dim {
		t.Errorf("expected %dx%d, got %dx%d", dim, dim, m.Rows, m.Cols)
	}
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			expected := Zero
			if i == j {
				expected = One
			}
			if m.Get(i, j) != expected {
				t.Errorf("at (%d, %d) expected %v, got %v", i, j, expected, m.Get(i, j))
			}
		}
	}
}

// TestMul verifies matrix multiplication with a known example.
func TestMul(t *testing.T) {
	// A = [1 2; 3 4]
	// B = [5 6; 7 8]
	// AB = [19 22; 43 50]
	a := NewMatrix(2, 2)
	a.Data = []complex128{1, 2, 3, 4}
	b := NewMatrix(2, 2)
	b.Data = []complex128{5, 6, 7, 8}

	res := Mul(a, b)
	expected := []complex128{19, 22, 43, 50}
	for i, v := range res.Data {
		if v != expected[i] {
			t.Errorf("at index %d expected %v, got %v", i, expected[i], v)
		}
	}
}

// TestTensor verifies the Kronecker tensor product with a known example.
func TestTensor(t *testing.T) {
	// A = [1 2; 3 4]
	// B = [0 1; 1 0]
	// A x B = [0 1 0 2; 1 0 2 0; 0 3 0 4; 3 0 4 0]
	a := NewMatrix(2, 2)
	a.Data = []complex128{1, 2, 3, 4}
	b := NewMatrix(2, 2)
	b.Data = []complex128{0, 1, 1, 0}

	res := Tensor(a, b)
	expected := []complex128{
		0, 1, 0, 2,
		1, 0, 2, 0,
		0, 3, 0, 4,
		3, 0, 4, 0,
	}
	if len(res.Data) != 16 {
		t.Fatalf("expected length 16, got %d", len(res.Data))
	}
	for i, v := range res.Data {
		if v != expected[i] {
			t.Errorf("at index %d expected %v, got %v", i, expected[i], v)
		}
	}
}

// TestConjugateTranspose verifies the Hermitian conjugate (adjoint) operation.
func TestConjugateTranspose(t *testing.T) {
	// A = [1+i 2-i]
	// A^H = [1-i; 2+i]
	a := NewMatrix(1, 2)
	a.Set(0, 0, complex(1, 1))
	a.Set(0, 1, complex(2, -1))

	res := ConjugateTranspose(a)
	if res.Rows != 2 || res.Cols != 1 {
		t.Fatalf("expected 2x1, got %dx%d", res.Rows, res.Cols)
	}
	if res.Get(0, 0) != complex(1, -1) {
		t.Errorf("expected 1-i, got %v", res.Get(0, 0))
	}
	if res.Get(1, 0) != complex(2, 1) {
		t.Errorf("expected 2+i, got %v", res.Get(1, 0))
	}
}

// TestHVConstants verifies the mathematical constants used in quantum gates.
func TestHVConstants(t *testing.T) {
	// Verify HC and HCN
	expectedHV := 1.0 / cmplx.Sqrt(2.0)
	if real(HC) != real(expectedHV) {
		t.Errorf("HC real part mismatch: expected %v, got %v", real(expectedHV), real(HC))
	}
	if real(HCN) != -real(expectedHV) {
		t.Errorf("HCN real part mismatch: expected %v, got %v", -real(expectedHV), real(HCN))
	}
}
