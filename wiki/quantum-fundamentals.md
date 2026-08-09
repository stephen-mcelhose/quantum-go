---
type: concept
title: Quantum Fundamentals
description: Superposition, measurement, unitary transforms, and the Bloch Sphere — the conceptual foundation for all quantum circuits, with worked examples in quantum-go.
resource: examples/fundamentals/fundamentals.md
tags: [fundamentals, superposition, measurement, hadamard, pauli, bloch-sphere]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Fundamentals

The entry point to quantum computing. Three pillars: **superposition**, **unitary transforms**, and **measurement**. Every algorithm in quantum-go is built on these.

## The Three Pillars

### Superposition
A qubit is a linear combination α|0⟩ + β|1⟩ where |α|² + |β|² = 1. Unlike a classical bit, a qubit simultaneously "is" both 0 and 1 — until measured. The Hadamard gate is the canonical tool for creating equal superposition:

```
H|0⟩ = (|0⟩ + |1⟩)/√2    ← |+⟩ state (50/50)
H|1⟩ = (|0⟩ − |1⟩)/√2    ← |−⟩ state (50/50 amplitudes, but with phase difference)
```

Critically: **H² = I**. Hadamard is its own inverse — applying it twice returns you to the start. This is a key self-assessment question.

### Unitary Transforms
Every quantum gate is a unitary matrix U satisfying U†U = I. Unitarity means:
- The operation is **reversible** — no information is destroyed.
- The total probability is **preserved** — the state vector stays on the unit sphere.

All gate operations are rotations on the Bloch Sphere. See [[quantum-concepts]] for full gate matrices and [[quantum-linear-algebra]] for the math.

### Measurement
Measurement collapses the quantum state. A qubit in state α|0⟩ + β|1⟩ collapses to:
- |0⟩ with probability |α|²
- |1⟩ with probability |β|²

After measurement, the superposition is **destroyed** — you cannot recover the original α and β. The [[simulation-engine]] implements this with `rand.Float64() < q.Probability`.

## The Bloch Sphere

The state of a single qubit |ψ⟩ = cos(θ/2)|0⟩ + e^{iφ}sin(θ/2)|1⟩ maps to a point on the unit sphere:
- **North pole** → |0⟩
- **South pole** → |1⟩
- **Equator** → equal superpositions (|+⟩, |−⟩, and rotations thereof)

Gates are rotations on this sphere:
- **H**: 180° rotation about the diagonal X+Z axis → maps |0⟩↔|+⟩
- **X**: 180° rotation about X axis → |0⟩↔|1⟩ (quantum NOT)
- **Z**: 180° rotation about Z axis → leaves |0⟩ alone, flips phase of |1⟩
- **Rx/Ry/Rz(θ)**: arbitrary angle rotations → see [[rotation-gates]]

## Single-Qubit Gate Family

| Gate   | Bloch rotation       | Special property                                    |
| ------ | -------------------- | --------------------------------------------------- |
| H      | 180° about X+Z axis  | H² = I; maps Z basis ↔ X basis                    |
| X      | 180° about X         | Quantum NOT; X² = I                               |
| Y      | 180° about Y         | Y = iXZ; Y² = I                                  |
| Z      | 180° about Z         | Phase flip; Z² = I                               |
| S      | 90° about Z          | S = T²; S² = Z                                   |
| T      | 45° about Z          | T needed for universality (see [[universality]])  |
| V (SX) | 90° about X          | V² = X (square root of NOT)                      |

## Worked Examples in quantum-go

### Single-Qubit Superposition
```go
p := core.NewProgram(1)
p.AddStep(core.NewStep(core.NewHadamard(0)))
// Result: |0>: 0.5000, |1>: 0.5000
```

### Multi-Qubit Superposition (8 equal states)
```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(
    core.NewHadamard(0),
    core.NewHadamard(1),
    core.NewHadamard(2),
))
// Result: all 8 states |000>...|111> with probability 0.125 each
```

### Bit Manipulation with Pauli Gates
```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(core.NewX(0), core.NewY(1), core.NewZ(2)))
// Result: |011>: 1.0000 (little-endian: q0 is rightmost bit)
```

**Little-endian note:** quantum-go uses the same convention as Qiskit — qubit 0 is the rightmost (least significant) bit in binary notation. `|011⟩` means q₀=1, q₁=1, q₂=0. See [[project-overview]] for the Qiskit compatibility rationale.

## Key Points

- **H² = I** — Hadamard is self-inverse; two H gates cancel.
- **Z on |0⟩ is invisible** — Z only affects the |1⟩ amplitude; measuring in Z basis after Z on |0⟩ looks identical.
- Pauli gates (X, Y, Z) are 180° rotations; phase gates (S, T) are partial rotations.
- Multi-qubit superposition: n H gates on |0…0⟩ creates uniform distribution over 2ⁿ states.
- Measurement is irreversible — it collapses the state (see [[simulation-engine]] for implementation).
- For parameterized rotations at arbitrary angles, see [[rotation-gates]].

## Sources

- `examples/fundamentals/fundamentals.md`

## References

- Nielsen, M.A. & Chuang, I.L. *Quantum Computation and Quantum Information* (10th anniversary ed.). Cambridge University Press, 2010. ISBN 978-1-107-00217-3. — The canonical textbook; Chapters 1–2 cover qubits, measurement, and Dirac notation.
- Wikipedia: [Qubit](https://en.wikipedia.org/wiki/Qubit)
- Wikipedia: [Quantum superposition](https://en.wikipedia.org/wiki/Quantum_superposition)
- Wikipedia: [Bra–ket notation](https://en.wikipedia.org/wiki/Bra%E2%80%93ket_notation)
- Wikipedia: [Born rule](https://en.wikipedia.org/wiki/Born_rule)
