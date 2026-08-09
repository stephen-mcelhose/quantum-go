---
type: concept
title: Arithmetic Gate Implementations
description: How core/arithmetic.go implements Fourier (QFT), Add, AddInteger, AddIntegerModulus, and MulModulus as nested BlockGate trees — the engine behind Shor's algorithm.
resource: core/arithmetic.go
tags: [arithmetic, qft, fourier, add, addinteger, mulmodulus, shor, block-gate]
timestamp: 2026-08-09T03:26:15Z
---

# Arithmetic Gate Implementations

`core/arithmetic.go` implements the quantum arithmetic stack as nested [[composite-gates|BlockGate]] trees. All operations are structured circuits — not standalone unitary matrices — enabling `ApplyOptimize` to recurse rather than materializing exponentially large matrices.

## The Hierarchy

```
MulModulus
└── ControlledBlockGate (conditioned on input bits)
    └── AddIntegerModulus
        ├── AddInteger (QFT + phase rotations + IQFT)
        ├── ControlledBlockGate (conditioned on carry)
        └── AddInteger (inverse)
```

## Fourier — QFT / IQFT

```go
func NewFourier(dim, idx int) *Fourier
```

Builds the QFT as a Block with:
1. For each qubit `i` from `dim−1` down to 0:
   - `H(i)` — Hadamard puts qubit into superposition
   - `CR(i+1−j, i, 2π/2ʲ)` for j=2..i+1 — controlled rotations from lower qubits
2. `Swap(i, dim−1−i)` — bit-reversal permutation at the end

**Inverse QFT:**
```go
invQFT := NewFourier(m, 0)
invQFT.SetInverse(true)
```
The `SetInverse(true)` flag propagates through `BlockGate.ApplyOptimize` as `g.Block.ApplyOptimize(v, g.Inverse)` — causing the block to execute its steps in **reverse order** and each gate inverted.

QASM: `NewFourier` has no QASM equivalent (it's not a standard gate), but its components (H, CR) are QASM-parseable.

## Add — Quantum-Quantum Addition

```go
func NewAdd(x0, x1, y0, y1 int) *Add
// |x⟩|y⟩ → |x⟩|x+y mod 2^m⟩
```

Three stages in the Block:
1. `NewFourier(m, 0)` — QFT on the x register
2. For each pair of qubits: `NewCr(i, cr0, 2π/2^{1+j})` — phase rotations controlled by y-register bits
3. `NewFourier(m, 0).SetInverse(true)` — IQFT to return to computational basis

The `cr0` index walks the y register; when `cr0 < m+n`, the y-register qubit exists and the CR gate fires. See [[quantum-arithmetic]] for the conceptual explanation.

**AffectedQubits** are set to absolute indices after block construction:
```go
for i := range a.AffectedQubits {
    a.AffectedQubits[i] = x0 + i
}
```

## AddInteger — Classical Integer into Quantum Register

```go
func NewAddInteger(x0, x1, num int) *AddInteger
// |x⟩ → |x + num mod 2^m⟩  (num is classical)
```

Stages:
1. `NewFourier(m, 0)` — QFT on register
2. For each qubit `i`: accumulate phase rotations from classical bits of `num` using `math.Mul` on combined phase matrices → `NewSingleQubitMatrixGate(i, mat)` — one combined gate per qubit
3. `NewFourier(m, 0).SetInverse(true)` — IQFT

The combined matrix approach (`math.Mul(mat, gateMat)`) merges all phase contributions to qubit `i` into a single `SingleQubitMatrixGate`, reducing the step count.

## AddIntegerModulus — Classical Integer mod N

```go
func NewAddIntegerModulus(x0, x1, a, N int) *AddIntegerModulus
// |x⟩|0⟩ → |(x + a) mod N⟩|0⟩   (uses ancilla carry qubit at index x1+1)
```

Seven stages in the Block (1 extra qubit for carry):
1. `AddInteger(a)` — tentative addition
2. `AddInteger(N).Inverse()` — subtract N to check for overflow
3. `Cnot(n−1, carry)` — copy sign bit to carry qubit
4. `ControlledBlockGate(AddInteger(N), carry)` — add N back if overflowed
5. `AddInteger(a).Inverse()` — subtract a to restore ancilla
6. `X(carry)` — invert carry
7. `Cnot(n−1, carry)` — uncompute carry

This structure ensures the carry qubit returns to |0⟩ after use (uncomputation pattern).

## MulModulus — Modular Multiplication

```go
func NewMulModulus(x0, x1, mul, mod int) *MulModulus
// |x⟩ → |x · mul mod N⟩
```

Uses 2·size+1 qubits (input register + result register + 1 carry):

1. For each bit `i` of the input (x register): apply `ControlledBlockGate(AddIntegerModulus(mul·2ⁱ mod N), control=i)` — adds `mul·2ⁱ mod N` to the result register conditioned on x-bit `i`
2. SWAP input and result registers
3. For each bit `i`: inverse `ControlledBlockGate(AddIntegerModulus(invmul·2ⁱ mod N), i)` — uncomputes the result register (now contains x before swap)

**Modular inverse:** `getInverseModulus(a, N)` uses the Extended Euclidean Algorithm to find `a⁻¹ mod N`:
```go
func getInverseModulus(a, n int) int {
    // Extended GCD (Euclidean algorithm)
    t, newt := 0, 1
    r, newr := n, a
    for newr != 0 {
        quotient := r / newr
        t, newt = newt, t-quotient*newt
        r, newr = newr, r-quotient*newr
    }
    if t < 0 { t += n }
    return t
}
```

Returns -1 if `a` is not invertible mod `n` (gcd(a,n) ≠ 1).

## Key Points

- All arithmetic gates are `BlockGate` trees — never materialize exponentially large matrices.
- Inverse = execute steps in reverse order with each gate inverted (passed through `ApplyOptimize`).
- QFT structure: H gates + CR gates (2π/2ᵏ) + bit-reversal SWAPs.
- `AddInteger` uses matrix multiplication to combine phase rotations per qubit into single `SingleQubitMatrixGate`.
- `AddIntegerModulus` requires an ancilla carry qubit and uses uncomputation to clean it up.
- `MulModulus` uses modular inverse (Extended GCD) for the uncomputation step.
- This stack is what powers [[shors-algorithm]] — `NewMulModulus` is the "expensive" quantum kernel.

## Sources

- `core/arithmetic.go`

## References

- Draper, T.G. "Addition on a Quantum Computer." arXiv:[quant-ph/0008033](https://arxiv.org/abs/quant-ph/0008033) (2000) — the QFT-based adder (no carry qubits) implemented by `NewAdd` and `NewAddInteger`.
- Vedral, V., Barenco, A. & Ekert, A. "Quantum Networks for Elementary Arithmetic Operations." arXiv:[quant-ph/9511018](https://arxiv.org/abs/quant-ph/9511018) (1995) — earlier ripple-carry approach; provides the modular exponentiation network that `NewMulModulus` implements.
- Shor, P.W. "Polynomial-Time Algorithms for Prime Factorization and Discrete Logarithms on a Quantum Computer." arXiv:[quant-ph/9508027](https://arxiv.org/abs/quant-ph/9508027) (1995) — `NewMulModulus` is the quantum kernel Shor's algorithm requires.
