---
type: concept
title: Circuits Library
description: core/circuits.go — pre-built Program factory functions for all major quantum circuits, ExpectedXxxState helpers, GetExpectedStateVector dispatch, CompareStateVectors, and ToQASM export.
resource: core/circuits.go, core/qasm.go
tags: [circuits, library, factory, programs, bell, ghz, grover, shor, teleportation, qasm-export]
timestamp: 2026-08-09T03:26:15Z
---

# Circuits Library

`core/circuits.go` is the **program factory** — pre-built functions for every major algorithm and circuit type. This is the primary entry point for building circuits programmatically, alongside the [[quantum-dsl]] primitives. `core/qasm.go` provides QASM export.

## Program Factory Functions

### Entanglement Circuits

```go
core.NewBellProgram()      // 2 qubits: H(0) + CNOT(0,1) → |Φ⁺⟩
core.NewGHZProgram(n)      // n qubits: H(0) + CNOT chain → |GHZ⟩
```

### Quantum Algorithms

```go
core.NewGroverProgram()                    // 2-qubit Grover searching for |11⟩
core.NewShorProgram(a, mod, precision)     // Shor's period finding aˣ mod N
core.NewQFTProgram(n)                      // n-qubit Quantum Fourier Transform
core.NewDeutschJozsaProgram(n, balanced)   // Deutsch-Josza: n-bit oracle (balanced/constant)
core.NewBernsteinVaziraniProgram(s)        // BV: find hidden string s via inner product oracle
core.NewSimonsProgram(s)                   // Simon's algorithm: find XOR period s
```

### Networking / Communication

```go
core.NewTeleportationProgram()      // 3 qubits: |ψ⟩ + Bell pair + corrections
core.NewSuperdenseCodingProgram()   // 2 qubits: encode "11" into one qubit via entanglement
core.NewQKDProgram(bit, basis)      // 1 qubit: BB84 state preparation (bit=0/1, basis=Z/X)
```

### Arithmetic

```go
core.NewAdderProgram(x0, x1, y0, y1, y2)  // 5-qubit Draper adder
```

### Error Correction

```go
core.NewErrorCorrectionProgram(bit)  // 3-qubit bit-flip code for logical bit value
```

### Basic / Test Circuits

```go
core.NewSuperpositionProgram(n)    // n H gates in parallel → uniform superposition
core.NewToffoliProgram()           // |110⟩ + Toffoli → |111⟩ demo
core.NewFredkinProgram()           // |110⟩ + Fredkin → |101⟩ demo
core.NewEngineProgram()            // Quantum thermodynamic cycle (Rx + H + Ry)
```

## Reference State Vectors

Functions returning the expected theoretical state vector (as `[]complex128`) for validation:

```go
core.ExpectedBellState()      // [1/√2, 0, 0, 1/√2]
core.ExpectedGHZState(n)      // [1/√2, 0, …, 0, 1/√2] (2^n amplitudes)
core.ExpectedQFTState(n)      // [1/√N, 1/√N, …, 1/√N] (all real, uniform)
core.ExpectedGroverState()    // [0, 0, 0, 1] (|11⟩ found with probability 1)
```

Dispatch by name:
```go
vec, err := core.GetExpectedStateVector("bell", 2)
vec, err := core.GetExpectedStateVector("ghz", 3)
vec, err := core.GetExpectedStateVector("qft", 2)
vec, err := core.GetExpectedStateVector("grover", 0)
// returns error for unknown circuit names
```

## Comparison Utility

```go
err := core.CompareStateVectors(got, want, 1e-6)
```

The same logic as `compareStateVectors` in [[verification-tests]] but as an exported function, usable from any package. Returns nil on success, error describing the first mismatch.

## QASM Export

`core/qasm.go` — `Program.ToQASM()`:

```go
p := core.NewBellProgram()
qasm := p.ToQASM()
// OPENQASM 2.0;
// include "qelib1.inc";
// qreg q[2];
// creg c[2];
// h q[0];
// cx q[0], q[1];
```

**Gate export mapping:**

| Go type       | QASM output                        |
| ------------- | ---------------------------------- |
| Hadamard      | `h q[i];`                          |
| X / Y / Z     | `x q[i];` / `y q[i];` / `z q[i];` |
| S / S†        | `s q[i];` / `sdg q[i];`            |
| T / T†        | `t q[i];` / `tdg q[i];`            |
| V             | `sx q[i];`                         |
| Rx/Ry/Rz      | `rx(θ) q[i];` etc.                 |
| U             | `u3(θ,φ,λ) q[i];`                  |
| CNOT          | `cx q[c], q[t];`                   |
| CZ            | `cz q[c], q[t];`                   |
| SWAP          | `swap q[i], q[j];`                 |
| Toffoli       | `ccx q[a], q[b], q[c];`            |
| PhaseShift    | `u1(θ) q[i];`                      |
| CR            | `cu1(θ) q[c], q[t];`               |
| Identity      | `id q[i];`                         |
| Measurement   | `measure q[i] -> c[i];`            |
| BlockGates    | Recursively decomposed to primitives |
| Unknown/Oracle| `// gate name q[...]` comment      |

**Roundtrip:** `Parse(p.ToQASM())` should regenerate an equivalent program for all supported gates. BlockGates are unrolled into their primitive constituents during export.

## Key Points

- `core.NewShorProgram` uses `NewMulModulus` with a different loop order than `examples/algorithms/shor_test.go` — the example in shor_test.go manually constructs the modular exponentiation for pedagogical clarity.
- `NewEngineProgram()` is the thermodynamic cycle — Rx(π/4) + H + Ry(π/4), not a real heat engine but a Bloch-sphere rotation demo.
- `ExpectedGHZState(n)` always has 1/√2 at index 0 and 1/√2 at index 2ⁿ−1, zeros everywhere else.
- `GetExpectedStateVector` returns an error for "grover" regardless of n — the second arg is unused.
- QASM export for Oracle and ControlledGate produces comments, not valid QASM — these gates have no standard QASM representation.

## Sources

- `core/circuits.go`
- `core/qasm.go`
