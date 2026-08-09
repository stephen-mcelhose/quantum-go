---
type: concept
title: Quantum Concepts & Gate Reference
description: Core quantum mechanics principles — superposition, entanglement, measurement — and the complete gate reference for the quantum-go simulator.
resource: docs/quantum-concepts.md
tags: [quantum-mechanics, gates, qubits, superposition, entanglement, basis]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Concepts & Gate Reference

The conceptual bedrock for everything in quantum-go. Every gate, algorithm, and simulation choice traces back to these principles.

## Quantum Bases

A **basis** is the reference frame in which you describe and measure a qubit. quantum-go works with two:

### Z Basis (Computational Basis)
The "classical" basis — |0⟩ and |1⟩ sit at the poles of the Bloch Sphere. The standard `Measure()` operation collapses to this basis. Phase differences are invisible here; only probability amplitudes matter.

### Phase Basis (X / Hadamard Basis)
|+⟩ = (|0⟩ + |1⟩)/√2 and |−⟩ = (|0⟩ − |1⟩)/√2 — states in perfect superposition where **relative phase** carries the information. BB84 (see [[bb84-qkd]]) exploits this: measuring in the wrong basis destroys the information, which is precisely how eavesdropping is detected.

| Feature          | Z Basis (Computational)              | Phase Basis (X / Hadamard)         |
| ---------------- | ------------------------------------- | ----------------------------------- |
| States           | \|0⟩, \|1⟩                           | \|+⟩, \|−⟩                         |
| Bit value        | Definite (0 or 1)                     | Indeterminate (50/50 probability)   |
| Information in   | Probability amplitude                 | Relative phase (+ or −)             |
| "Strange" gate   | X (bit flip)                          | Z (phase flip)                      |
| Bloch sphere     | Vertical axis (poles)                 | Horizontal axis (equator)           |

## Core Principles

### Superposition
A qubit can exist in a linear combination of |0⟩ and |1⟩ simultaneously. Applying H to |0⟩ creates equal superposition:

```
H|0⟩ = (|0⟩ + |1⟩)/√2
```

Measurement collapses this to |0⟩ or |1⟩ with 50% each. Before measurement, the qubit genuinely occupies both. This is what enables quantum parallelism: an n-qubit system lives in 2ⁿ-dimensional Hilbert space, so the [[simulation-engine]] must track 2ⁿ complex amplitudes.

### Entanglement
When two qubits are entangled, measuring one instantly determines the other regardless of distance. The canonical construction — the **Bell state** — requires H followed by CNOT:

```
Step 1: H on q₀ → (|00⟩ + |10⟩)/√2
Step 2: CNOT(q₀→q₁) → (|00⟩ + |11⟩)/√2
```

Only |00⟩ and |11⟩ have amplitude. Measuring q₀ as 0 forces q₁ to 0; measuring as 1 forces q₁ to 1. See [[entanglement]] for Bell and GHZ state construction.

### Unitarity
All quantum operations (except measurement) are **unitary** — they preserve the total probability (‖ψ‖² = 1) and are reversible. Every gate matrix U satisfies U†U = I. The [[simulation-engine]] applies these as matrix-vector products on the state vector. See [[quantum-linear-algebra]] for the math.

## Gate Reference

### Fundamental Single-Qubit Gates

| Gate          | Symbol | Action                                                     |
| ------------- | ------ | ---------------------------------------------------------- |
| Identity      | I      | No-op; useful for padding in multi-qubit steps             |
| Hadamard      | H      | Rotates Z→X basis; creates/destroys superposition          |
| Pauli-X       | X      | Bit flip: \|0⟩↔\|1⟩ (quantum NOT)                         |
| Pauli-Y       | Y      | Bit flip + phase flip: \|0⟩→i\|1⟩, \|1⟩→−i\|0⟩          |
| Pauli-Z       | Z      | Phase flip: \|1⟩→−\|1⟩; leaves \|0⟩ alone                 |
| S Gate        | S      | Quarter-turn phase: \|1⟩→i\|1⟩ (= T²)                    |
| T Gate        | T      | Eighth-turn phase: \|1⟩→e^{iπ/4}\|1⟩; needed for universality |
| SX (V)        | V      | Square root of X; V² = X                                   |

### Parameterized Rotation Gates
These accept a continuous angle θ, making them essential for variational algorithms and universal decompositions (see [[universality]]):

| Gate       | Rotation axis | Key property                                           |
| ---------- | ------------- | ------------------------------------------------------ |
| Rx(θ)      | X             | Rx(π) = −iX (up to global phase)                      |
| Ry(θ)      | Y             | Real matrix — no complex entries; used in Grover's     |
| Rz(θ)      | Z             | Diagonal — efficient to apply                          |
| PhaseShift | Z (|1⟩ only)  | PS(θ) = diag(1, e^{iθ}); used in QFT                  |
| Universal  | Arbitrary     | U(θ,φ,λ) — any single-qubit gate as 3 Euler angles    |

Full matrix forms are in [[rotation-implementations]].

### Two-Qubit Gates

| Gate      | Symbol | Action                                              |
| --------- | ------ | --------------------------------------------------- |
| CNOT      | CX     | Flips target if control = \|1⟩; creates entanglement |
| CZ        | CZ     | Phase flip on \|11⟩; symmetric — no distinct target |
| SWAP      | SWAP   | Exchanges two qubit states                          |

### Advanced Multi-Qubit Gates

| Gate         | Symbol | Action                                              |
| ------------ | ------ | --------------------------------------------------- |
| Toffoli      | CCNOT  | Flips target if both controls = \|1⟩; universal for classical logic |
| Rotation     | R(θ)   | Phase e^{iθ} on \|1⟩                               |
| Controlled-R | CR(θ)  | Conditional phase; backbone of QFT                  |

### Composite Algorithmic Blocks

| Block         | Key components    | Function                                             |
| ------------- | ----------------- | ---------------------------------------------------- |
| QFT (Fourier) | H + CR gates      | Rotates Z↔Phase basis; the interference engine for Shor's |
| Adder (Add)   | QFT + CR gates    | \|x,y⟩→\|x,x+y⟩ in phase space                     |

Implementations in [[composite-gates]] and [[arithmetic-gates]]. Deep-dive in [[qft-deep-dive]].

## Key Points

- The Z basis is the measurement basis; the X/Hadamard basis is where phase information lives.
- Superposition + entanglement + interference are the three pillars enabling quantum speedup.
- Every gate is unitary; the simulation must preserve this invariant.
- Parameterized gates (Rx/Ry/Rz/U) are how you represent rotations on the Bloch sphere continuously.
- The Hadamard gate is the basis-change operator: H|0⟩=|+⟩, H|1⟩=|−⟩, and H²=I.
- CNOT is the canonical two-qubit entangling gate; combined with single-qubit gates it is universal.

## Sources

- `docs/quantum-concepts.md`
