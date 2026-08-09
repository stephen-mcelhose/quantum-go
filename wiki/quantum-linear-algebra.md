---
type: concept
title: Quantum Linear Algebra
description: The math package — flat complex128 matrix representation, Kronecker (tensor) products, conjugate transpose, and why these operations are the mathematical substrate of all quantum gates.
resource: math/matrix.go
tags: [linear-algebra, matrix, kronecker, tensor-product, complex128, math]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Linear Algebra

The `math` package is the foundation of quantum-go. It knows nothing about qubits — it's pure linear algebra over `complex128`. Every gate matrix, every state vector transformation, and every verification comparison ultimately reduces to the operations defined here.

## The Matrix Type

```go
type Matrix struct {
    Rows, Cols int
    Data       []complex128  // row-major flat slice
}
```

**Row-major flat array:** element at (row, col) lives at `Data[row*Cols + col]`. This layout is cache-friendly for row-wise traversal (which `Mul` does) and avoids the pointer indirection of `[][]complex128`. See [[simulator-optimizations]] for why this matters.

**Complex128:** Go's `complex128` is two `float64` values (real + imaginary). This gives 64-bit precision per component — sufficient for quantum gate matrices where exact values like 1/√2 appear frequently. The package predefines:

```go
var (
    Zero = complex(0, 0)
    One  = complex(1, 0)
    I    = complex(0, 1)    // imaginary unit
    HV   = 1.0 / math.Sqrt(2.0)   // 1/√2
    HC   = complex(HV, 0)
    HCN  = complex(-HV, 0)
)
```

## Core Operations

### `Mul(a, b Matrix) Matrix` — Matrix Multiplication
Standard O(n³) matrix multiplication via three nested loops. Used when the simulator falls back to full matrix-vector multiply (see [[gate-application]]). For 2×2 gate matrices this is trivial; for n-qubit full system matrices this is O(4ⁿ) — which is why the bitwise optimizations exist.

### `Tensor(a, b Matrix) Matrix` — Kronecker Product
The Kronecker product is how you combine quantum systems. If A is m×n and B is p×q:
- Result is (mp)×(nq)
- Element (i, j) of result = A[i/p, j/q] × B[i%p, j%q]

```go
// Zero elements are skipped — key optimization for sparse gate matrices
if va == 0 { continue }
if vb == 0 { continue }
res.Set(i*b.Rows+k, j*b.Cols+l, va*vb)
```

**Physical meaning:** To apply a 1-qubit gate H to qubit 0 of a 2-qubit system, you need the 4×4 matrix H ⊗ I (or I ⊗ H for qubit 1). `Tensor(H, I)` constructs this. The bitwise optimizations in [[gate-application]] avoid this expansion entirely for most gates.

**Why skip zeros?** Single-qubit gate matrices are typically sparse (H has no zeros, but X/Z have 50% zeros; identity is maximally sparse). Skipping zeros in the inner loop avoids 0×something multiplications, which matters when computing large Kronecker products.

### `ConjugateTranspose(m Matrix) Matrix` — Hermitian Adjoint (†)
Transposes the matrix and complex-conjugates every element:

```go
res.Set(j, i, complex(real(val), -imag(val)))  // (a+bi)† = (a-bi)
```

**Physical meaning:** For a unitary gate matrix U, U† = U⁻¹. This is how inverse gates are computed — when `IsInverse()` is true, the gate returns `ConjugateTranspose(GetMatrix())`. The T†, S†, V† gates are defined this way.

### `IdentityMatrix(dim int) Matrix`
Returns the dim×dim identity matrix (1s on diagonal). Used for padding single-qubit gates to act on multi-qubit systems, and as the no-op gate.

## Relationship to State Vectors

State vectors in quantum-go are `[]complex128` — a column vector of 2ⁿ complex amplitudes. The `Mul` function handles matrix-vector multiplication when a gate is represented as a Matrix and applied via the fallback path in [[gate-application]]. In the optimized path, the matrix entries are read directly (`matrix.Get(0,0)`, `matrix.Get(0,1)`, etc.) without constructing intermediate representations.

## Thermodynamics Extensions

The `math` package also contains density matrix operations used by [[quantum-thermodynamics]]:
- `ToDensityMatrix(v []complex128) Matrix` — ρ = |ψ⟩⟨ψ| (outer product)
- `PartialTrace(rho Matrix, n, traceOut int) Matrix` — trace over a subsystem
- `ExpectationValue(rho, H Matrix) complex128` — ⟨H⟩ = Tr(ρH)
- `VonNeumannEntropy(rho Matrix) float64` — S = -Tr(ρ log ρ)

These are not used by the core gate simulation but enable thermodynamic analysis of quantum states.

## Key Points

- Row-major flat `[]complex128` is the single memory layout for all matrices — gates, state vectors (as 1D), and Kronecker intermediates.
- `Tensor` is the mathematical operation that combines subsystems; the bitwise tricks in [[gate-application]] avoid it at runtime.
- `ConjugateTranspose` (†) is how every inverse gate is implemented — unitary gates are self-inverse under †.
- `HV = 1/√2` is the most important constant in the codebase — it appears in every Hadamard operation.
- Zero-skipping in `Tensor` is a meaningful optimization for sparse gate matrices.

## Sources

- `math/matrix.go`
