/*
Package math provides complex matrix operations for quantum computing.

This package implements linear algebra operations on complex-valued matrices,
which are fundamental to quantum gate operations and state transformations.

# Matrix Representation

Matrices are stored in row-major order using a flat 1D slice for performance
and memory efficiency. For a matrix M with dimensions (rows × cols), element
M[i,j] is stored at index i*cols + j in the Data slice.

# Usage

	import "github.com/stephen-mcelhose/quantum-go/math"

	// Create a 2x2 Hadamard matrix
	h := math.NewMatrix(2, 2)
	hv := math.HV  // 1/√2
	h.Data = []complex128{hv, hv, hv, -hv}

	// Matrix multiplication
	result := math.Mul(h, h)  // H² = I

	// Tensor product (Kronecker)
	combined := math.Tensor(h, h)  // H ⊗ H

	// Conjugate transpose (adjoint)
	hermitian := math.ConjugateTranspose(h)

# Common Operations

NewMatrix(rows, cols): Creates a zero-initialized matrix.

Get(row, col): Retrieves element at position (row, col).

Set(row, col, val): Sets element at position (row, col).

IdentityMatrix(dim): Creates a dim×dim identity matrix.

Mul(a, b): Matrix multiplication A×B. Requires a.Cols == b.Rows.

Tensor(a, b): Kronecker (tensor) product A⊗B. Used to combine single-qubit
gates into multi-qubit gates.

ConjugateTranspose(m): Returns the Hermitian conjugate (adjoint) M†.
For quantum gates, this represents the inverse operation.

# Quantum Gate Construction

Single-qubit gates are 2×2 matrices. To create a multi-qubit gate that applies
a single-qubit operation to one qubit while leaving others unchanged, use the
tensor product with identity matrices:

	// Apply Hadamard to qubit 1 in a 3-qubit system
	gate := math.Tensor(math.IdentityMatrix(2),
		math.Tensor(h, math.IdentityMatrix(2)))

However, for simulation efficiency, the local package applies gates directly
to state vectors without constructing these large matrices.

# Constants

The package provides common quantum computing constants:

	Zero = 0+0i
	One = 1+0i
	I = 0+1i (imaginary unit)
	HV = 1/√2 ≈ 0.7071067811865476
	HC = HV+0i (complex form of 1/√2)
	HCN = -HV+0i (complex form of -1/√2)

# Performance

Matrix operations are not optimized for large dimensions. The flat array
representation provides good cache locality for small to medium matrices.
For quantum simulation, direct state vector manipulation (as in the local
package) is preferred over explicit matrix operations when possible.

	Zero-valued elements are skipped in tensor product computation to improve
	performance for sparse matrices.

	# Quantum Thermodynamics

	The package provides tools for analyzing quantum systems from a thermodynamic
	perspective, including density matrix formalism and observable measurement.

	# Density Matrices

	A pure state vector |ψ⟩ can be converted to a density matrix ρ = |ψ⟩⟨ψ|:

		rho := math.ToDensityMatrix(stateVector)

	# Partial Trace

	To analyze a subsystem of a larger quantum system, use PartialTrace to
	integrate out the degrees of freedom of the rest of the system:

		// Trace out qubit 0 in a 2-qubit system to get the reduced density matrix
		// of qubit 1.
		rho_sub := math.PartialTrace(rho_full, 2, []int{0})

	# Observables and Expectation Values

	The expectation value of an observable (represented by a Hermitian matrix H)
	is given by ⟨H⟩ = tr(ρH):

		energy := math.ExpectationValue(rho, hamiltonian)

	# Entropy

	The von Neumann entropy S(ρ) = -tr(ρ ln ρ) measures the degree of mixedness
	or entanglement in a quantum state:

		s := math.VonNeumannEntropy(rho)

	# Matrix Exponential

	Continuous-time evolution is supported via the matrix exponential e^M:

		unitary := math.Exp(m)

	# Example: Hadamard Gate


	// Define Hadamard matrix
	h := math.NewMatrix(2, 2)
	h.Data = []complex128{
		math.HC, math.HC,
		math.HC, math.HCN,
	}

	// Verify H² = I
	h2 := math.Mul(h, h)
	// h2.Data should equal [1, 0, 0, 1] (identity)

# Mathematical Properties

The package ensures that:

1. Matrix multiplication produces correct dimensions: (m×n) × (n×p) → (m×p)
2. Tensor products produce correct dimensions: (m×n) ⊗ (p×q) → (mp×nq)
3. Conjugate transpose swaps dimensions: (m×n)† → (n×m)
4. Operations preserve complex number precision
*/
package math
