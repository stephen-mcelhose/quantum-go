---
type: concept
title: Scaling and Limits
description: Memory and time complexity of state vector simulation — the qubit ceiling (~30), algorithm-specific risks for Shor's and Grover's, and practical mitigation strategies.
resource: docs/algorithmic-risks.md
tags: [scaling, limits, complexity, memory, shor, grover, qubits]
timestamp: 2026-08-09T03:26:15Z
---

# Scaling and Limits

Classical simulation of quantum systems is inherently exponential. quantum-go's optimizations (see [[simulator-optimizations]]) reduce gate cost from O(4ⁿ) to O(2ⁿ), but the **state vector itself is O(2ⁿ) in memory** — irreducible for full state vector simulation. This sets the practical ceiling.

## State Vector Memory Growth

Each amplitude is a `complex128` = 16 bytes. The state vector has 2ⁿ entries:

| Qubits (n) | States (2ⁿ)      | Memory        |
| ---------- | ----------------- | ------------- |
| 10         | 1,024             | 16 KB         |
| 20         | 1,048,576         | 16 MB         |
| 25         | 33,554,432        | 512 MB        |
| 30         | 1,073,741,824     | 16 GB         |
| 34         | 17,179,869,184    | 256 GB        |

**Practical ceiling: ~30 qubits on a standard laptop.** Attempting 31+ qubits will trigger OOM before any gates run. See [[project-overview]] for the confirmed 30-qubit / 16 GB limit.

## Algorithm-Specific Risks

### Shor's Algorithm
Shor's (see [[shors-algorithm]]) is the most demanding algorithm in the codebase:

- **Qubit count**: Factoring an L-bit number needs roughly 3L qubits (L input + 2L for modular exponentiation).
  - N=15 (L=4): ~12–15 qubits — easily simulatable.
  - N=21 (L=5): ~15–18 qubits — manageable.
  - N=35 (L=6): ~18–21 qubits — borderline.
  - N>35: memory + time constraints make classical simulation impractical.
- **Circuit depth**: Modular exponentiation (`Add` + `Mul` operations) involves thousands of individual gates → long wall-clock time even within the qubit limit.
- **Risk**: Even within the memory budget, Shor's for larger N can take minutes to hours.

### Grover's Algorithm
Grover's (see [[grovers-algorithm]]) provides quadratic speedup O(√N) for search over N=2ⁿ items:

- **Qubit count**: Often small (n < 15 for practical examples).
- **Iteration count**: Requires ~(π/4)√(2ⁿ) iterations of the Grover operator.
- **Oracle complexity**: If the oracle encodes a complex predicate (e.g., 3-SAT), circuit depth explodes.
- **Risk**: Low qubit count but high iteration count → long simulation time despite the small state vector.

## Practical Mitigations Used in quantum-go

| Strategy                     | Description                                                                                   |
| ---------------------------- | --------------------------------------------------------------------------------------------- |
| Limit qubit counts           | All `examples/` use n < 15 qubits — run in seconds on typical hardware                       |
| Simplified oracles           | Predefined unitary matrices instead of full gate decompositions where oracle logic isn't the focus |
| Binary result filtering      | `PrintBinary()` shows only high-probability states — human-readable from large state vectors  |
| Block gate avoidance of O(4ⁿ) | QFT/Adder use step-by-step execution, not full matrix materialization (see [[gate-application]]) |

## Comparison: Classical vs Quantum Complexity

This is the simulator's inherent irony: the algorithms it simulates are exponentially *faster* on real quantum hardware, but the simulation itself is exponentially *slower* on classical hardware.

| Task                    | Quantum hardware | Classical simulation       |
| ----------------------- | ---------------- | -------------------------- |
| Shor's (factor N)       | O(log³ N)        | Exponential in log N       |
| Grover's (search N)     | O(√N)            | O(√N) iterations × O(2ⁿ) per |
| QFT (n qubits)          | O(n²) gates      | O(n² × 2ⁿ) state updates  |

## Key Points

- State vector simulation is O(2ⁿ) in memory — irreducible, unlike gate cost which can be optimized.
- ~30 qubits / 16 GB is the practical ceiling for laptop-class hardware.
- Shor's algorithm is both memory-intensive (3L qubits) and time-intensive (deep circuits).
- Grover's is qubit-cheap but iteration-expensive for large search spaces.
- The `examples/` directory deliberately uses n < 15 so circuits run in seconds.
- For educational purposes, the simulator is ideal up to ~20 qubits — beyond that, real QPUs or tensor-network simulators are needed.

## Sources

- `docs/algorithmic-risks.md`
