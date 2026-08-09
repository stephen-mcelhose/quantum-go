// Package math provides linear algebra operations for quantum computing.
// It includes complex matrix operations, tensor products, and common matrix transformations
// required for quantum gate operations. Matrices are stored in row-major order using
// flat 1D slices for better performance.
package math

import (
	"fmt"
	"math"
	"strings"
)

// Matrix represents a complex-valued matrix.
// Data is stored in row-major order as a flat 1D slice for memory efficiency and cache locality.
// Element at (row, col) is stored at index row*Cols + col.
type Matrix struct {
	// Rows is the number of rows in the matrix.
	Rows, Cols int
	// Data contains all matrix elements in row-major order.
	Data []complex128
}

// NewMatrix creates a new Matrix with the given dimensions.
// All elements are initialized to zero.
func NewMatrix(rows, cols int) Matrix {
	return Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([]complex128, rows*cols),
	}
}

// Get returns the complex value at position (row, col).
// No bounds checking is performed for efficiency.
func (m Matrix) Get(row, col int) complex128 {
	return m.Data[row*m.Cols+col]
}

// Set assigns a complex value to position (row, col).
// No bounds checking is performed for efficiency.
func (m *Matrix) Set(row, col int, val complex128) {
	m.Data[row*m.Cols+col] = val
}

// Common complex constants used in quantum computing.
var (
	// Zero is the complex number 0+0i.
	Zero = complex(0, 0)
	// One is the complex number 1+0i.
	One = complex(1, 0)
	// I is the imaginary unit 0+1i.
	I = complex(0, 1)
	// HV is 1/√2, the value used in the Hadamard gate.
	HV = 1.0 / math.Sqrt(2.0)
	// HC is the complex number (1/√2)+0i.
	HC = complex(HV, 0)
	// HCN is the complex number -(1/√2)+0i.
	HCN = complex(-HV, 0)
)

// IdentityMatrix returns an identity matrix of dimension dim x dim.
// The identity matrix has 1's on the main diagonal and 0's elsewhere.
func IdentityMatrix(dim int) Matrix {
	m := NewMatrix(dim, dim)
	for i := 0; i < dim; i++ {
		m.Set(i, i, One)
	}
	return m
}

// Mul computes the matrix multiplication C = A × B.
// The number of columns in A must equal the number of rows in B.
// Panics if dimensions are incompatible.
//
// The resulting matrix has dimensions (A.Rows × B.Cols).
func Mul(a, b Matrix) Matrix {
	if a.Cols != b.Rows {
		panic(fmt.Sprintf("dimension mismatch: a.Cols (%d) != b.Rows (%d)", a.Cols, b.Rows))
	}
	res := NewMatrix(a.Rows, b.Cols)
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			var sum complex128
			for k := 0; k < a.Cols; k++ {
				sum += a.Get(i, k) * b.Get(k, j)
			}
			res.Set(i, j, sum)
		}
	}
	return res
}

// Tensor computes the Kronecker (tensor) product of matrices a and b.
// The Kronecker product is fundamental in quantum computing for combining
// quantum states and gates operating on different qubits.
//
// If A is m×n and B is p×q, the result is mp×nq.
// Element (i,j) of the result is A[i/p,j/q] * B[i%p,j%q].
//
// Zero elements are skipped for efficiency.
func Tensor(a, b Matrix) Matrix {
	res := NewMatrix(a.Rows*b.Rows, a.Cols*b.Cols)
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			va := a.Get(i, j)
			if va == 0 {
				continue
			}
			for k := 0; k < b.Rows; k++ {
				for l := 0; l < b.Cols; l++ {
					vb := b.Get(k, l)
					if vb == 0 {
						continue
					}
					res.Set(i*b.Rows+k, j*b.Cols+l, va*vb)
				}
			}
		}
	}
	return res
}

// ConjugateTranspose returns the Hermitian conjugate (adjoint) of the matrix.
// This operation transposes the matrix and takes the complex conjugate of each element.
// For a quantum gate matrix, the conjugate transpose represents the inverse gate.
func ConjugateTranspose(m Matrix) Matrix {
	res := NewMatrix(m.Cols, m.Rows)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			val := m.Get(i, j)
			res.Set(j, i, complex(real(val), -imag(val)))
		}
	}
	return res
}

// String returns a human-readable string representation of the matrix.
// Each row is displayed on a separate line with complex values formatted as (real, imag).
func (m Matrix) String() string {
	var sb strings.Builder
	for i := 0; i < m.Rows; i++ {
		sb.WriteString("[")
		for j := 0; j < m.Cols; j++ {
			if j > 0 {
				sb.WriteString(", ")
			}
			v := m.Get(i, j)
			sb.WriteString(fmt.Sprintf("(%.4f, %.4f)", real(v), imag(v)))
		}
		sb.WriteString("]\n")
	}
	return sb.String()
}
