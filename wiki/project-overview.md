---
type: concept
title: Project Overview
description: High-level framing of quantum-go — what it is, how it was born, how to install and run it, and where to find learning resources.
resource: README.md
tags: [overview, installation, quick-start, project, strange]
timestamp: 2026-08-09T03:26:15Z
---

# Project Overview

quantum-go is a **high-performance quantum circuit simulator written in Go**, ported from the [Strange](https://github.com/redfx-quantum/strange) Java implementation. It simulates quantum circuits on classical hardware via full state vector simulation — tracking all 2ⁿ complex amplitudes for n qubits. The architecture is documented in [[package-architecture]].

## What It Is

A complete, verified quantum simulator that:
- Builds circuits using a clean Go DSL (see [[quantum-dsl]])
- Runs them on a local state vector engine (see [[simulation-engine]])
- Exports to OpenQASM 2.0 for use in external tools (see [[openqasm-parser]])
- Is verified against Qiskit Aer and IBM Quantum hardware (see [[verification-tests]])

It is **not** a cloud quantum service — it runs locally, simulating quantum mechanics classically. The practical qubit ceiling is ~30 qubits (~16 GB RAM). See [[scaling-and-limits]].

## Quick Start

```bash
# Build the CLI
go build -o strange ./cmd/strange/main.go

# Run all tests
go test ./...

# Run fuzz tests
go test -fuzz=FuzzToffoli ./local
go test -fuzz=FuzzMatrixMul ./math
```

For Qiskit-based verification:
```bash
python3 -m venv venv && source venv/bin/activate
pip install qiskit qiskit-aer
```

## Gate Support

| Category      | Gates                                              |
| ------------- | -------------------------------------------------- |
| Fundamental   | Hadamard, Pauli-X/Y/Z, Identity                    |
| Multi-qubit   | CNOT, CZ, SWAP, Toffoli, Fredkin                   |
| Rotations     | Rx, Ry, Rz, PhaseShift, Universal                  |
| Phase         | S, T, V (SX), and their inverses (S†, T†, V†)     |
| Special       | Measurement                                        |
| Composite     | QFT, IQFT, Draper Adder, MulModulus (for Shor's)  |

Full gate reference: [[quantum-concepts]]. Implementations: [[gate-implementations]], [[rotation-implementations]], [[composite-gates]].

## Project Structure

```
quantum-go/
├── core/         # Gate, Step, Program, Result types   → [[quantum-dsl]]
├── local/        # Simulation engine                   → [[simulation-engine]]
├── math/         # Matrix ops, Kronecker, complex math → [[quantum-linear-algebra]]
├── qasm/         # OpenQASM 2.0 parser                 → [[openqasm-parser]]
├── examples/     # Pedagogical guides per algorithm
├── docs/         # Architecture, optimization, verification docs
└── verification/ # Qiskit verification scripts + reports
```

## Learning Guides

The `examples/` directory contains markdown guides that explain the *why* and *how* of each algorithm alongside the Go implementation:

| Guide                    | Wiki page               |
| ------------------------ | ----------------------- |
| Quantum Fundamentals     | [[quantum-fundamentals]] |
| Entanglement (Bell/GHZ)  | [[entanglement]]         |
| Shor's Algorithm         | [[shors-algorithm]]      |
| Grover's Algorithm       | [[grovers-algorithm]]    |
| Quantum Error Correction | [[error-correction]]     |
| Quantum Teleportation    | [[teleportation]]        |
| BB84 Key Distribution    | [[bb84-qkd]]             |
| Quantum Arithmetic       | [[quantum-arithmetic]]   |

## Performance Characteristics

| Metric      | Value                                                    |
| ----------- | -------------------------------------------------------- |
| Memory      | O(2ⁿ) — exponential growth; 30 qubits ≈ 16 GB           |
| Gate cost   | O(2ⁿ) per gate (must update all amplitudes)              |
| Optimizations | Bitwise tricks, identity skipping, block gate pattern |

See [[simulator-optimizations]] and [[scaling-and-limits]] for details.

## Heritage

Ported from [redfx-quantum/strange](https://github.com/redfx-quantum/strange) (BSD 3-Clause). The Go port maintains the same algorithmic core but uses Go idioms, flat-array matrices, and adds a richer verification suite. The original implementation plan is not archived in this wiki.

## Key Points

- quantum-go is a state vector simulator — it tracks all 2ⁿ amplitudes explicitly.
- The CLI tool is called `strange` (honouring the Java original).
- Qiskit Little-Endian ordering is matched, making cross-verification tractable.
- The three packages (`math`, `core`, `local`) have strict layering — see [[package-architecture]].

## Sources

- `README.md`
