---
type: concept
title: Composite Gate Types
description: BlockGate, Oracle, ControlledGate, ControlledBlockGate, PermutationGate, SingleQubitMatrixGate, TimeEvolution — how quantum-go composes gates into higher-order circuit elements.
resource: core/composite.go
tags: [composite-gates, block-gate, oracle, controlled-gate, permutation, time-evolution, go]
timestamp: 2026-08-09T03:26:15Z
---

# Composite Gate Types

`core/composite.go` provides the higher-order gate abstractions that allow quantum-go to express complex, structured circuits — like QFT or modular multiplication — as single gate objects. These types all live in the `core` package and are used by [[arithmetic-gates]], [[oracle-gates]], and the [[simulation-engine]].

## BlockGate — The Recursive Circuit Primitive

```go
type BlockGate struct {
    BaseGate
    Block *Block       // contains its own Steps and Gates
}
```

A `BlockGate` wraps a `Block` (which is itself a mini-circuit with NQubits, Steps). The key insight: **`HasOptimization() → true`** — the engine never calls `GetMatrix()` on a BlockGate. Instead it calls:

```go
func (g *BlockGate) ApplyOptimize(v []complex128) []complex128 {
    return g.Block.ApplyOptimize(v, g.Inverse)
}
```

This recurses into the Block's own step execution, meaning any circuit can be composited into a single gate without materializing its full matrix (which would be 2^n × 2^n). This is the foundation for [[arithmetic-gates]] (QFT, Add, MulModulus).

**`GetMatrix()` is a stub** — returns an empty 2^n × 2^n matrix. If called by mistake (not by the normal path), it won't simulate correctly.

## Oracle — User-Defined Unitary

```go
type Oracle struct {
    BaseGate
    Matrix math.Matrix   // full unitary matrix
}

func NewOracle(idx int, m math.Matrix) *Oracle {
    nq := int(math.Log2(float64(m.Rows)))  // infer qubit count from matrix size
    // sets AffectedQubits = [idx, idx+1, ..., idx+nq-1]
}
```

Allows plugging in **any 2^n × 2^n unitary matrix** as a circuit gate. Used in:
- [[grovers-algorithm]] — Oracle phase-flip and Diffusion operator
- Custom circuit experiments

Inverse: `GetMatrix()` returns `ConjugateTranspose(g.Matrix)` when `g.Inverse = true`.

## ControlledGate — Generic Control Wrapper

```go
type ControlledGate struct {
    BaseGate
    ControlIndexes []int
    RootGate       Gate    // the gate being controlled
}
```

`NewControlledGate(g, control...)` prefixes control qubits to the gate's affected list and prepends "C" to the Caption. `HasOptimization() → true` but `ApplyOptimize` is a stub — the actual logic is handled by the [[simulation-engine]] using the `ControlIndexes` field to do conditional state-vector manipulation.

The `GetMatrix()` is a placeholder that returns the right-sized empty matrix. It's not used for simulation.

## ControlledBlockGate — Conditioned Sub-Circuit

```go
type ControlledBlockGate struct {
    BlockGate
    ControlIndex int
}
```

Created from a `BlockGate` + a single control qubit. Critical for [[arithmetic-gates]]:
- `NewAddIntegerModulus` uses `ControlledBlockGate(addN, dim)` — conditioned on the carry qubit.
- `NewMulModulus` uses `ControlledBlockGate(add, i)` — each modular addition conditioned on a bit of the input register.

Must be constructed from a `BlockGateInterface` — panics if given a plain gate. The simulation engine handles it via `ControlledBlockGate.ControlIndex`.

## PermutationGate — Bit-Level Reordering

```go
type PermutationGate struct {
    BaseGate
    Target1, Target2 int
}
```

Swaps two qubit indices across the **entire state vector** (not just a subspace). Used internally by the QFT's bit-reversal step and the `swapBits` helper:

```go
func swapBits(i, i1, i2 int) int {
    b1 := (i >> i1) & 1
    b2 := (i >> i2) & 1
    if b1 == b2 { return i }
    return i ^ (1 << i1) ^ (1 << i2)
}
```

The PermutationGate matrix is constructed by mapping each row `i` to column `swapBits(i, t1, t2)`.

## SingleQubitMatrixGate — Inline Matrix Gate

```go
type SingleQubitMatrixGate struct {
    BaseGate
    Matrix math.Matrix   // 2×2
}
```

Used in `NewAddInteger` to combine multiple phase shifts into one gate object:
```go
gateMat := rGate.GetMatrix()
mat = math.Mul(mat, gateMat)      // fold all phase shifts for this qubit
pstep.AddGate(NewSingleQubitMatrixGate(i, mat))
```

This avoids creating many separate gate objects in the step when multiple phase rotations hit the same qubit.

## TimeEvolution — Hamiltonian Simulation

```go
type TimeEvolution struct {
    Oracle   // extends Oracle
}

func NewTimeEvolution(idx int, h math.Matrix, t float64) *TimeEvolution {
    // compute m = -iHt
    // u = math.Exp(m)     ← matrix exponential
    // wrap as Oracle
}
```

Simulates quantum time evolution U = e^{−iHt} where H is a Hamiltonian. The matrix exponential (`math.Exp`) is from [[quantum-linear-algebra]]. Used in quantum chemistry and simulation contexts.

## Composition Flow in Shor's Algorithm

```
MulModulus
└── ControlledBlockGate (control=i)
    └── AddIntegerModulus
        ├── AddInteger (QFT → PhaseShifts → IQFT)
        ├── ControlledBlockGate (control=carry)
        │   └── AddInteger
        └── AddInteger (inverse)
```

Each level is a `BlockGate` whose `ApplyOptimize` recurses into sub-blocks. The full Shor circuit is a deeply nested tree of BlockGates — never materializing a single large matrix.

## Key Points

- `BlockGate.ApplyOptimize` recurses into the inner Block — never materializes the exponentially large matrix.
- `Oracle` = plug any unitary matrix in; supports inverse via ConjugateTranspose.
- `ControlledBlockGate` panics if given a non-BlockGate — the gate must implement `BlockGateInterface`.
- `swapBits` enables the QFT's bit-reversal without Swap gate objects.
- `TimeEvolution` wraps matrix exponential — enables Hamiltonian simulation in quantum-go.
- `SingleQubitMatrixGate` is a performance optimization that merges multiple phase shifts on one qubit into one gate.

## Sources

- `core/composite.go`
