---
type: concept
title: Quantum DSL — Program, Step, Gate
description: The core.go type hierarchy that forms the lingua franca of quantum-go — how Programs, Steps, Gates, Blocks, and Results compose into executable circuits.
resource: core/core.go
tags: [dsl, go, core, program, step, gate, block, result]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum DSL — Program, Step, Gate

`core/core.go` defines the **Program → Step → Gate** hierarchy that every algorithm in quantum-go is expressed in. Understanding this structure is a prerequisite for reading any source file.

## The Type Hierarchy

```
Program
  └── []Step
        └── []Gate  (disjoint qubits per step)
```

**Program** = the whole circuit. Holds `NumQubits`, an ordered slice of Steps, and `InitAlpha` (initial amplitude for each qubit — 1.0 means |0⟩).

**Step** = one time-slice. Gates within a single step *must operate on disjoint qubits* — this is enforced at `AddGate` time with a panic. That constraint models true quantum parallelism: independent qubits can be transformed simultaneously.

**Gate** = one quantum operation. The `Gate` interface is the central abstraction:

```go
type Gate interface {
    GetMatrix() smath.Matrix           // 2ⁿ × 2ⁿ unitary matrix
    GetAffectedQubitIndexes() []int    // which qubits this gate touches
    GetHighestAffectedQubitIndex() int
    GetCaption() string                // e.g., "H"
    GetName() string                   // e.g., "Hadamard"
    GetGroup() string
    GetSize() int                      // 1 for single-qubit, 2 for two-qubit, etc.
    SetInverse(inv bool)               // adjoint / conjugate transpose
    IsInverse() bool
    HasOptimization() bool             // true if ApplyOptimize bypasses matrix multiply
    ApplyOptimize(v []complex128) []complex128  // specialized fast path
}
```

## BaseGate

Most concrete gates embed `BaseGate` which provides default implementations: `AffectedQubits []int`, `Inverse bool`, `Caption/Name/Group string`. Specific gate types (see [[gate-implementations]]) override `GetMatrix()` and optionally `ApplyOptimize()`.

## Block and BlockGate — Composite Gates

A **Block** is a named sub-circuit with its own Steps:

```go
type Block struct {
    Steps   []*Step
    NQubits int
    Name    string
}
```

`Block.ApplyOptimize(v, inverse)` runs its Steps through `GlobalStepExecutor` — the function variable set by the `local` package to avoid an import cycle. When `inverse=true`, steps run in reverse order and each gate's inverse flag is set, implementing the adjoint circuit automatically.

This is the **block gate pattern**: QFT, IQFT, Adder, etc. are all composed as Blocks (see [[composite-gates]]). A block gate *simulates itself* — it doesn't need a 2ⁿ × 2ⁿ matrix; it just runs its sub-steps on the state vector directly.

## GlobalStepExecutor

```go
var GlobalStepExecutor StepExecutor
// type: func(gates []Gate, vector []complex128, numQubits int) []complex128
```

This is the one runtime coupling that crosses the `core` → `local` boundary. The `local` package sets this during initialization. Without it, `Block.ApplyOptimize` is a no-op (safe default). With it, composite gates can execute self-contained. See [[package-architecture]] for the dependency design rationale.

## Result Types

Two Result implementations for different use cases:

| Type                    | Stores                               | Use case                          |
| ----------------------- | ------------------------------------ | --------------------------------- |
| `CompactResult`         | Final state vector only              | Default — minimal memory          |
| `InstrumentedResultImpl`| All intermediate state vectors       | Debugging, step-by-step analysis  |

Both expose `GetProbability() []complex128` (the state vector) and `GetQubits() []*Qubit` (per-qubit marginal probabilities).

`StateVectorReference` is the JSON interchange format — used by a qiskit bridge to exchange state vectors between Go and Python.

## Qubit — Measurement Model

```go
type Qubit struct {
    Probability   float64  // P(qubit = |1⟩)
    Measured      bool
    MeasuredValue bool
}
```

`CalculateQubitStatesFromVector` derives per-qubit marginal probabilities from the full state vector by summing squared amplitudes where that qubit's bit is 1:

```
P(qubitᵢ = 1) = Σⱼ |ψⱼ|² where bit i of j is 1
```

This is a partial trace over all other qubits — see [[quantum-linear-algebra]] for the density matrix equivalent.

## Building a Circuit — Canonical Pattern

```go
p := core.NewProgram(3)                        // 3-qubit circuit, all |0⟩
p.AddStep(core.NewStep(core.NewHadamard(0)))   // H on qubit 0
p.AddStep(core.NewStep(core.NewCnot(0, 1)))    // CNOT (control=0, target=1)
p.AddStep(core.NewStep(core.NewCnot(1, 2)))    // CNOT (control=1, target=2)
// Result: GHZ state (|000⟩ + |111⟩)/√2
```

Gate constructors (NewHadamard, NewCnot, etc.) live in [[gate-implementations]]. Composite gate constructors (NewFourier, NewAdd) live in [[composite-gates]].

## Key Points

- **Step enforces disjointness** — the panic in `verifyUnique` catches circuit construction bugs early.
- **Gate interface has two apply paths**: `GetMatrix()` for the general case (matrix-vector multiply); `ApplyOptimize()` for specialized implementations that skip matrix construction entirely (see [[gate-application]]).
- **Block = reusable sub-circuit** — the QFT block is the same Block whether used standalone or inside Shor's algorithm.
- **Inverse is automatic** — setting `Inverse=true` on a Block reverses step order and conjugate-transposes each gate's matrix; no separate IQFT implementation needed.
- **InitAlpha** allows starting qubits in states other than |0⟩ — a value of 0.0 starts the qubit in |1⟩.
- **Pre-built programs** are available in [[circuits-library]] — use those rather than constructing common circuits by hand.

## Sources

- `core/core.go`
