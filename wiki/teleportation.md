---
type: concept
title: Quantum Teleportation
description: The Bell-measurement protocol for transferring a quantum state from Alice to Bob via pre-shared entanglement and two classical bits — without moving any matter.
resource: examples/networking/teleportation.md
tags: [teleportation, bell-measurement, entanglement, classical-correction, no-cloning]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Teleportation

Teleportation transfers a quantum state |ψ⟩ from Alice to Bob using:
1. A pre-shared [[entanglement|Bell pair]]
2. Alice's Bell measurement (2 classical bits)
3. Bob's conditional Pauli corrections

No matter moves. The original state at Alice is **destroyed** during measurement (consistent with No-Cloning). Bob's qubit ends up in exactly the state Alice had.

## Why Not Just Send the Qubit?

- Quantum states are fragile — transmission over noise channels degrades them.
- No-Cloning means you can't backup before sending.
- Teleportation separates "sending the quantum state" from "sending the physical qubit" by using pre-established entanglement as a noiseless resource.

## The Protocol (3-Qubit, Alice=q0,q1, Bob=q2)

```
q0: Alice's state to teleport |ψ⟩
q1: Alice's half of Bell pair
q2: Bob's half of Bell pair
```

### Step 1: Prepare |ψ⟩ on q0
```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(core.NewHadamard(0)))  // |ψ⟩ = |+⟩
```

### Step 2: Create Bell Pair on (q1, q2)
```go
p.AddStep(core.NewStep(core.NewHadamard(1)))
p.AddStep(core.NewStep(core.NewCnot(1, 2)))   // |Φ⁺⟩ = (|00⟩+|11⟩)/√2
```

### Step 3: Alice's Bell Measurement
Alice entangles her state with her Bell half:
```go
p.AddStep(core.NewStep(core.NewCnot(0, 1)))   // CNOT q0→q1
p.AddStep(core.NewStep(core.NewHadamard(0)))   // H on q0
```
After these two gates, Alice measures q0 and q1 in the computational basis — projecting the joint system onto one of the four Bell states and "encoding" |ψ⟩'s information into Bob's qubit.

### Step 4: Bob's Classical Corrections
Alice sends her 2 classical bits to Bob. Bob applies:
```go
p.AddStep(core.NewStep(core.NewCnot(1, 2)))   // if q1=1: X on q2
p.AddStep(core.NewStep(core.NewCz(0, 2)))     // if q0=1: Z on q2
```
In the simulator, q0 and q1 serve as quantum controls for these gates (equivalent to classical conditional gates since they're measured). Bob's q2 now holds |ψ⟩.

## Why It Doesn't Violate No-FTL

Teleportation requires sending **2 classical bits** from Alice to Bob before Bob can perform corrections. Without these bits, Bob's qubit is in a mixed state (completely random from Bob's perspective). The classical channel caps the protocol at c. No faster-than-light communication.

## Correction Table

Depending on Alice's measurement outcome (q0, q1):

| Measurement (q0, q1) | Bob applies | Correction matrix                 |
| -------------------- | ----------- | --------------------------------- |
| 00                   | I (nothing) | Identity                          |
| 01                   | X           | Bit flip                          |
| 10                   | Z           | Phase flip                        |
| 11                   | XZ          | Both                              |

## Interpreting the Simulation Output

After teleporting |+⟩ (50/50 superposition):
```
|000>: 0.1250
|001>: 0.1250
...
|111>: 0.1250
```
All 8 states appear with equal probability. This is correct: Alice's qubits (q0, q1) can be in any of 4 measurement outcomes, and Bob's qubit (q2) is in |+⟩ (50% |0⟩, 50% |1⟩). Combined: 4 × 2 = 8 equally likely outcomes.

The teleportation success is verified by the marginal probability of q2 — 50% |0⟩ / 50% |1⟩ matches the intended |+⟩ state.

## Relationship to Other Protocols

| Protocol       | Entanglement use                          | Classical channel |
| -------------- | ----------------------------------------- | ----------------- |
| Teleportation  | Transfer quantum state                    | 2 bits            |
| [[bb84-qkd]]   | Detect eavesdropping (conjugate bases)    | Public channel    |
| Superdense coding | Send 2 classical bits via 1 qubit      | 1 qubit channel   |

## Key Points

- Requires pre-shared entanglement (Bell pair) as a resource — a "quantum link" established in advance.
- Alice's Bell measurement is: CNOT(q0→q1) then H(q0), then measure both.
- 2 classical bits from Alice → Bob's conditional X and Z corrections.
- The original state at Alice is destroyed (consistent with No-Cloning Theorem).
- Not faster-than-light: Bob needs Alice's classical bits before corrections can be applied.
- In quantum-go, CZ serves as the controlled-Z correction gate.
- For multi-qubit error protection (a complementary protocol using Toffoli correction), see [[error-correction]].

## Sources

- `examples/networking/teleportation.md`

## References

- Bennett, C.H. et al. "Teleporting an unknown quantum state via dual classical and Einstein-Podolsky-Rosen channels." *Physical Review Letters* 70.13 (1993): 1895. DOI:[10.1103/PhysRevLett.70.1895](https://doi.org/10.1103/PhysRevLett.70.1895) — the original teleportation protocol paper.
- Wikipedia: [Quantum teleportation](https://en.wikipedia.org/wiki/Quantum_teleportation)
- Wikipedia: [Superdense coding](https://en.wikipedia.org/wiki/Superdense_coding) — the dual protocol (2 classical bits via 1 qubit).
