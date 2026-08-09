---
type: concept
title: Entanglement
description: Bell states, GHZ states, and the mathematical definition of entanglement as non-separability — with the H+CNOT construction used throughout quantum-go.
resource: examples/entanglement/entanglement.md
tags: [entanglement, bell-state, ghz, cnot, correlation, non-separability]
timestamp: 2026-08-09T03:26:15Z
---

# Entanglement

Entanglement is the defining feature of quantum mechanics with no classical analogue. An entangled multi-qubit state **cannot** be written as a tensor product of individual qubit states — it is fundamentally non-separable.

**Mathematical definition:** |ψ⟩ is entangled if and only if |ψ⟩ ≠ |φ₁⟩ ⊗ |φ₂⟩ ⊗ … for any choice of individual states.

Entanglement is the resource behind [[teleportation]], [[bb84-qkd]], [[grovers-algorithm]], and [[shors-algorithm]].

## Bell State (|Φ⁺⟩)

The simplest maximally entangled 2-qubit state:

```
|Φ⁺⟩ = (|00⟩ + |11⟩)/√2
```

Construction — H then CNOT:

```go
p := core.NewProgram(2)
p.AddStep(core.NewStep(core.NewHadamard(0)))  // |00⟩ → (|00⟩ + |10⟩)/√2
p.AddStep(core.NewStep(core.NewCnot(0, 1)))   // → (|00⟩ + |11⟩)/√2
```

Step-by-step:
1. H on q₀ creates superposition: (|00⟩ + |10⟩)/√2 (q₀ in superposition, q₁ still |0⟩)
2. CNOT(control=q₀, target=q₁): when q₀=1, flip q₁ → (|00⟩ + |11⟩)/√2

**What entanglement means here:** Measuring q₀ as |0⟩ forces q₁ to |0⟩. Measuring q₀ as |1⟩ forces q₁ to |1⟩. The two qubits have no individual state — only the joint state exists.

**Output:**
```
|00>: 0.5000
|11>: 0.5000
```
States |01⟩ and |10⟩ have zero probability — the qubits are locked together.

## The Four Bell States

The Bell state above is |Φ⁺⟩. The other three are produced by additional single-qubit gates before/after the H+CNOT:

| State    | Ket                              | Relation to |Φ⁺⟩                |
| -------- | -------------------------------- | --------------------- |
| \|Φ⁺⟩   | (|00⟩ + |11⟩)/√2                | Base state            |
| \|Φ⁻⟩   | (|00⟩ − |11⟩)/√2                | Apply Z to q₀         |
| \|Ψ⁺⟩   | (|01⟩ + |10⟩)/√2                | Apply X to q₀         |
| \|Ψ⁻⟩   | (|01⟩ − |10⟩)/√2                | Apply XZ to q₀        |

These four states form an orthonormal basis for the 4-dimensional 2-qubit Hilbert space. Alice's Bell measurement in [[teleportation]] projects onto one of these four.

## GHZ State — 3-Qubit Entanglement

The Greenberger–Horne–Zeilinger state generalizes Bell pairs to 3+ qubits:

```
|GHZ⟩ = (|000⟩ + |111⟩)/√2
```

Construction — H then two CNOTs in a chain:

```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(core.NewHadamard(0)))
p.AddStep(core.NewStep(core.NewCnot(0, 1)))
p.AddStep(core.NewStep(core.NewCnot(1, 2)))  // now all three are correlated
```

**Output:**
```
|000>: 0.5000
|111>: 0.5000
```

All 6 remaining states (|001⟩, |010⟩, …, |110⟩) have zero probability. Measuring any one qubit instantly determines all others.

## Spooky Action at a Distance

The Bell/GHZ correlation appears instantaneous regardless of qubit separation. Einstein called this "spooky action at a distance" and believed it indicated quantum mechanics was incomplete. Bell's theorem (1964) proved that no local hidden variable theory can reproduce quantum correlations — experiments have confirmed quantum mechanics.

**Crucial caveat:** Entanglement cannot transmit information faster than light. Bob's measurement outcomes are random until he receives Alice's classical bits — no FTL communication. This is the constraint exploited in [[bb84-qkd]] and [[teleportation]].

## Entanglement as a Non-Classical Resource

| Property                        | Classical bits     | Entangled qubits                        |
| ------------------------------- | ------------------ | --------------------------------------- |
| State of pair                   | Two independent bits | One joint inseparable state           |
| Measurement correlation         | Pre-determined     | Instantaneous (but not FTL)            |
| Copying the state               | Trivial            | Impossible (No-Cloning Theorem)        |
| State after measurement         | Unchanged          | Collapsed to definite value            |

The Von Neumann entropy of a reduced density matrix (ρ_A = Tr_B|ψ⟩⟨ψ|) measures entanglement. For a Bell pair, S(ρ_A) = log(2) = 1 ebit — maximum entanglement. The `math` package can compute this via `VonNeumannEntropy`. See [[quantum-thermodynamics]].

## Key Points

- Bell state = H on one qubit + CNOT entangling the pair. This two-gate sequence appears in nearly every quantum protocol.
- GHZ state = Bell pair + CNOT chain propagating entanglement to more qubits.
- Entanglement is non-separability — the joint state cannot be factored into individual qubit states.
- Measuring one entangled qubit collapses the entire shared state.
- No-Cloning Theorem: you cannot copy an unknown quantum state — this is what makes BB84 secure.
- Entanglement entropy (Von Neumann) = 0 for separable states, 1 ebit for maximally entangled Bell pairs.

## Sources

- `examples/entanglement/entanglement.md`

## References

- Einstein, A., Podolsky, B. & Rosen, N. "Can Quantum-Mechanical Description of Physical Reality Be Considered Complete?" *Physical Review* 47 (1935): 777. DOI:[10.1103/PhysRev.47.777](https://doi.org/10.1103/PhysRev.47.777) — the EPR paradox paper that introduced entanglement as a puzzle.
- Bell, J.S. "On the Einstein Podolsky Rosen paradox." *Physics Physique Fizika* 1.3 (1964): 195. Open access (CERN): [cds.cern.ch/record/111654](https://cds.cern.ch/record/111654/files/vol1p195-200_001.pdf) — proves entanglement cannot be explained by local hidden variables.
- Wikipedia: [Quantum entanglement](https://en.wikipedia.org/wiki/Quantum_entanglement)
- Wikipedia: [Bell's theorem](https://en.wikipedia.org/wiki/Bell%27s_theorem)
- Wikipedia: [No-cloning theorem](https://en.wikipedia.org/wiki/No-cloning_theorem)
