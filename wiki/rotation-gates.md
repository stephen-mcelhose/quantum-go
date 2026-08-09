---
type: concept
title: Rotation Gates
description: Rx, Ry, Rz, PhaseShift — parameterized single-qubit gates that rotate the state vector by an arbitrary angle on the Bloch Sphere.
resource: examples/fundamentals/rotations.md
tags: [rotation, rx, ry, rz, bloch-sphere, parameterized, phaseshift]
timestamp: 2026-08-09T03:26:15Z
---

# Rotation Gates

Parameterized rotation gates are the continuous-angle counterparts to the fixed Pauli gates. They allow precise positioning of a qubit state anywhere on the Bloch Sphere and are the building blocks of [[universality]] and variational quantum algorithms.

## The Bloch Sphere Coordinates

Any single-qubit state is:
```
|ψ⟩ = cos(θ/2)|0⟩ + e^{iφ}sin(θ/2)|1⟩
```
- θ ∈ [0, π]: polar angle from north pole (|0⟩)
- φ ∈ [0, 2π): azimuthal angle (phase)

Rotation gates move the state along great circles of this sphere.

## Rx(θ) — Rotation Around X Axis

```
Rx(θ) = [ cos(θ/2)      −i·sin(θ/2) ]
         [ −i·sin(θ/2)  cos(θ/2)    ]
```

- Mixes |0⟩ and |1⟩ amplitudes with imaginary cross-terms.
- **Rx(π) = −iX** (Pauli-X up to global phase).
- Rx(π/2) moves |0⟩ to (|0⟩ − i|1⟩)/√2.

## Ry(θ) — Rotation Around Y Axis

```
Ry(θ) = [ cos(θ/2)   −sin(θ/2) ]
         [ sin(θ/2)    cos(θ/2) ]
```

- **Real matrix** — no complex entries. This makes Ry the cleanest rotation to reason about.
- **Ry(π) = Y** (Pauli-Y up to global phase).
- **Ry(π/2)** moves |0⟩ to the equator: (|0⟩ + |1⟩)/√2 = |+⟩.
- Used in Grover's diffusion operator (see [[grovers-algorithm]]).

## Rz(θ) — Rotation Around Z Axis

```
Rz(θ) = [ e^{−iθ/2}  0          ]
         [ 0           e^{+iθ/2} ]
```

- Diagonal matrix — acts purely as a phase shift.
- Does **not** change measurement probabilities in Z basis (only changes relative phase).
- **Rz(π) = −iZ** (Pauli-Z up to global phase).
- **Rz(π/2) = −iS** (S gate up to global phase).
- Rz gates are especially efficient to implement — [[gate-application]] shows they fall into the optimized single-qubit path.

## PhaseShift(θ) — Phase Rotation on |1⟩ Only

```
PS(θ) = [ 1    0         ]
         [ 0    e^{iθ}   ]
```

- Leaves |0⟩ unchanged; rotates |1⟩ by phase e^{iθ}.
- **PS(π) = Z**, **PS(π/2) = S**, **PS(π/4) = T**.
- The controlled-R gate (CR) in QFT is a controlled PhaseShift — the backbone of the Quantum Fourier Transform. See [[qft-deep-dive]].

## Comparison Table

| Gate       | Axis | Real-valued? | Special case                                     |
| ---------- | ---- | ------------ | ------------------------------------------------ |
| Rx(θ)      | X    | No           | Rx(π) ≈ X; Rx(π/2) creates |0⟩−i|1⟩           |
| Ry(θ)      | Y    | Yes          | Ry(π/2) creates |+⟩ (same as H on |0⟩)          |
| Rz(θ)      | Z    | No (diag)    | Rz(π) ≈ Z; invisible in Z-basis measurements    |
| PhaseShift | Z    | No (diag)    | PS(π/4) = T; used in QFT controlled rotations   |

## quantum-go Usage

```go
p := core.NewProgram(1)
// Rotate 45° (π/4) around Y axis
p.AddStep(core.NewStep(core.NewRy(math.Pi/4, 0)))

// PhaseShift by π/4 (equivalent to T gate)
p.AddStep(core.NewStep(core.NewPhaseShift(math.Pi/4, 0)))
```

Implementations in [[rotation-implementations]]. The Universal U gate subsumes all three: see [[universality]].

## Key Points

- Rx(π) ≈ X, Ry(π) ≈ Y, Rz(π) ≈ Z — all Pauli gates are special cases of rotations.
- Ry is the only real-valued rotation gate — useful when working with real-amplitude circuits.
- Rz and PhaseShift both rotate phase; Rz rotates both amplitudes symmetrically while PhaseShift only rotates |1⟩.
- QFT's controlled-R gates are controlled PhaseShifts at angles 2π/2ᵏ — see [[qft-deep-dive]].
- H = U(π/2, 0, π) — Hadamard is a special case of the Universal gate.

## Sources

- `examples/fundamentals/rotations.md`

## References

- Wikipedia: [Bloch sphere](https://en.wikipedia.org/wiki/Bloch_sphere)
- Wikipedia: [Rotation operator (quantum mechanics)](https://en.wikipedia.org/wiki/Rotation_operator_(quantum_mechanics))
- Wikipedia: [Pauli matrices](https://en.wikipedia.org/wiki/Pauli_matrices)
- Nielsen, M.A. & Chuang, I.L. *Quantum Computation and Quantum Information*, Chapter 4 (quantum circuit model, single-qubit gates). Cambridge University Press, 2010.
