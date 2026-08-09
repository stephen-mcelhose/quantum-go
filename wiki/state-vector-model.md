---
type: synthesis
title: State Vector Model
description: End-to-end tour of how quantum-go represents and evolves quantum state — from the initial []complex128 allocation through gate application to measurement and result readout.
tags: [state-vector, amplitude, probability, measurement, simulation-model]
timestamp: 2026-08-09T03:26:15Z
---

# State Vector Model

The state vector is the central data structure of quantum-go's simulation. This page traces its lifecycle from creation to measurement.

## What the State Vector Is

For an n-qubit system, the state vector is a slice of 2ⁿ complex numbers:
```go
state := make([]complex128, 1<<numQubits)
state[0] = complex(1, 0)  // initial state |0...0⟩
```

Each element `state[i]` is the **amplitude** for computational basis state `|i⟩`. The probability of measuring `|i⟩` is `|state[i]|²`.

**Invariant:** Sum of squared magnitudes = 1.0 (normalization). Every gate is unitary — it preserves this invariant.

## Little-Endian Index Convention

Index i encodes qubit values as bit fields:
```
bit j of i  =  value of qubit j
```

| Index | Binary   | Qubit values (q₂q₁q₀) |
| ----- | -------- | ---------------------- |
| 0     | 000      | q₀=0, q₁=0, q₂=0      |
| 1     | 001      | q₀=1, q₁=0, q₂=0      |
| 2     | 010      | q₀=0, q₁=1, q₂=0      |
| 3     | 011      | q₀=1, q₁=1, q₂=0      |
| 7     | 111      | q₀=1, q₁=1, q₂=1      |

This matches Qiskit's convention and the [[project-overview]] compatibility rationale. The rightmost bit in binary notation is qubit 0 (LSB = first qubit).

## Lifecycle

### 1. Allocation (NewSimpleExecutionEnvironment)

```go
state := make([]complex128, 1 << p.NumQubits)
state[0] = 1  // |0...0⟩
```

All amplitudes zero except index 0 = `|0…0⟩`. See [[simulation-engine]].

### 2. Gate Application (ApplyGate per Step)

For each Step, for each Gate, the engine calls `ApplyGate(state, gate)`.

**Fast path (single-qubit):** For a 2×2 gate acting on qubit j:
```
For each i where bit j is 0:
    state0 = state[i]           (amplitude of |...0...⟩)
    state1 = state[i | (1<<j)]  (amplitude of |...1...⟩)
    state[i]          = m00*state0 + m01*state1
    state[i|(1<<j)]   = m10*state0 + m11*state1
```
Pairs related by bit j are updated in place. O(2ⁿ) per gate. See [[gate-application]].

**Bitwise shortcuts (optimized paths):**
- CNOT: `state[i] ↔ state[i ^ (1<<target)]` only when bit control = 1
- CZ: `state[i] *= -1` only when both control and target bits = 1
- Identity: skip entirely

### 3. Measurement

After all steps, the simulator inspects each qubit's marginal probability:
```go
prob1 := 0.0
for i, amp := range state {
    if (i >> qubitIndex) & 1 == 1 {
        prob1 += real(amp)*real(amp) + imag(amp)*imag(amp)
    }
}
qubit.Probability = prob1
```

**Explicit measurement gate:** When a `Measurement` gate is encountered, the state is collapsed:
```go
if rand.Float64() < prob1 {
    // Collapse: zero out all states where qubit is 0, renormalize
} else {
    // Collapse: zero out all states where qubit is 1, renormalize
}
```

The state vector is **modified in place** by measurement — subsequent gates see the collapsed state.

### 4. Result Readout

```go
res.GetProbability()   // returns the amplitude vector ([]complex128) — NOT probabilities
res.GetQubits()        // returns marginal Qubit structs with .Probability and .Measure()
```

**Naming confusion:** `GetProbability()` returns amplitudes, not |amplitude|². Named for historical reasons. Always square-magnitude when you need actual probability. See [[verification-tests]] for the correct usage pattern.

`qubit.Measure()` — returns 0 or 1 deterministically if probability is near 0 or 1; probabilistic otherwise.

## Summary: State at Each Stage

```
Initial state:         [1, 0, 0, 0, ...]   |0...0⟩
After H on q0:         [1/√2, 1/√2, 0, ...]  superposition
After CNOT(0,1):       [1/√2, 0, 0, 1/√2]   entanglement
After measurement q0:  [0, 0, 0, 1] or [1, 0, 0, 0]  collapse
Result: GetProbability() = the current complex amplitude slice
```

## Key Points

- State vector = 2ⁿ complex128 values; index i = basis state |i⟩ in little-endian qubit ordering.
- Qubit j controls bit j of the index: `bit j of i = (i >> j) & 1`.
- All gates preserve normalization (unitarity) — no renormalization step needed until explicit measurement.
- `GetProbability()` is misnamed — it returns **amplitudes**, not probabilities.
- Measurement is the only non-unitary operation — it collapses the state and modifies the slice in-place.
- For n=30 qubits: 2³⁰ × 16 bytes ≈ 17 GB RAM. See [[scaling-and-limits]] for practical ceilings.
