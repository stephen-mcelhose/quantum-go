---
type: concept
title: Grover's Algorithm
description: Amplitude amplification for unstructured search — how the Oracle phase-flips the target state and the Diffusion Operator (inversion about mean) amplifies it to near-certainty.
resource: examples/algorithms/grover.md
tags: [grover, search, amplitude-amplification, oracle, diffusion, quadratic-speedup]
timestamp: 2026-08-09T03:26:15Z
---

# Grover's Algorithm

Grover's algorithm solves **unstructured search** in O(√N) queries — a quadratic speedup over the classical O(N). For N=2ⁿ items, classical search needs N/2 checks on average; Grover's needs (π/4)√N iterations. Discovered by Lov Grover in 1996.

## The Problem

Find a marked item in an unsorted database of N = 2ⁿ items. We have access only to an "Oracle" that recognizes the answer — not a sorted structure we can binary-search.

- Classical: O(N) in worst case.
- Grover's: O(√N) quantum queries.
- For N=1,000,000: classical checks ≈ 500,000; Grover ≈ 785 iterations.

See [[scaling-and-limits]] for the simulation cost at scale.

## The Four Components

### 1. Superposition (The Haystack)
Apply H to all n qubits → equal superposition over all 2ⁿ states, each with amplitude 1/√N.

### 2. The Oracle (Phase Flip)
A black-box operation that "marks" the target state |t⟩ by flipping its phase:
```
Oracle|x⟩ = −|x⟩  if x = t (target)
Oracle|x⟩ = |x⟩   otherwise
```
Implemented as a diagonal matrix with −1 in the target position. In quantum-go, the Oracle is created as a custom `core.NewOracle(qubit, matrix)` using a manually constructed matrix. See [[oracle-gates]] for the gate type.

### 3. The Diffusion Operator (Inversion About Mean)
Also called the Grover diffusion or "amplitude amplification" operator. It reflects all amplitudes about their average value:
```
Diffusion: aᵢ → 2⟨a⟩ − aᵢ
```
Because the oracle flipped the target amplitude negative, it's far below the mean. After diffusion, it gets reflected to a large positive value while non-target amplitudes shrink slightly.

### 4. Iteration
Repeat Oracle + Diffusion approximately (π/4)√N times. **Over-iterating (overcooking)** rotates past the target — probability starts decreasing again.

## Geometric Interpretation

Grover's is a rotation in the 2D subspace spanned by:
- |t⟩ (target state)
- |s⊥⟩ (uniform superposition over non-target states)

Each iteration rotates the state vector by angle 2arcsin(1/√N) ≈ 2/√N radians toward |t⟩. After (π/4)√N rotations, the vector nearly aligns with |t⟩.

## 2-Qubit Implementation (N=4, search for |11⟩)

```go
p := core.NewProgram(2)

// Step 1: Superposition
p.AddStep(core.NewStep(core.NewHadamard(0), core.NewHadamard(1)))

// Step 2: Oracle — phase flip on |11⟩ (index 3 in 4-state system)
oracleMatrix := math.NewMatrix(4, 4)
oracleMatrix.Set(0, 0, 1); oracleMatrix.Set(1, 1, 1)
oracleMatrix.Set(2, 2, 1); oracleMatrix.Set(3, 3, -1)  // mark |11⟩
oracle := core.NewOracle(0, oracleMatrix)
p.AddStep(core.NewStep(oracle))

// Step 3: Diffusion operator (2-qubit version)
diffMatrix := math.NewMatrix(4, 4)
for i := 0; i < 4; i++ {
    for j := 0; j < 4; j++ {
        if i == j { diffMatrix.Set(i, j, -0.5) } else { diffMatrix.Set(i, j, 0.5) }
    }
}
diffusion := core.NewOracle(0, diffMatrix)
p.AddStep(core.NewStep(diffusion))
```

**Result after 1 iteration:** `|11>: 1.0000` — perfect amplification in a single step.

(For N=4, one iteration is optimal. For N=16, ~3 iterations needed.)

## Iteration Count

| N (items) | Optimal iterations | Formula          |
| --------- | ------------------ | ---------------- |
| 4         | 1                  | (π/4)√4 ≈ 1.57  |
| 16        | 3                  | (π/4)√16 ≈ 3.1  |
| 64        | 6                  | (π/4)√64 ≈ 6.3  |
| 1024      | 25                 | (π/4)√1024 ≈ 25 |

Round to nearest integer. Using ⌊(π/4)√N⌋ + 1 can overrotate.

## Key Points

- Oracle phase-flips the target; diffusion amplifies it — two rotations per iteration.
- One iteration suffices for N=4 (2 qubits). For larger N, use (π/4)√N iterations.
- Overcooking: too many iterations rotates past the target and probability decreases again.
- The Oracle is problem-specific — for a real problem (3-SAT, graph coloring), implementing the oracle efficiently is the hard part.
- Grover's is optimal — no quantum algorithm can solve unstructured search faster than O(√N).
- See [[algorithm-comparison]] for how Grover's compares to Shor's and other algorithms.

## Sources

- `examples/algorithms/grover.md`

## References

- Grover, L.K. "A fast quantum mechanical algorithm for database search." *Proceedings of the 28th Annual ACM Symposium on Theory of Computing* (1996): 212–219. arXiv:[quant-ph/9605043](https://arxiv.org/abs/quant-ph/9605043) — the original paper introducing the algorithm.
- Grover, L.K. "Quantum mechanics helps in searching for a needle in a haystack." *Physical Review Letters* 79.2 (1997): 325. arXiv:[quant-ph/9706033](https://arxiv.org/abs/quant-ph/9706033) — journal version with the amplitude amplification framing.
- Wikipedia: [Grover's algorithm](https://en.wikipedia.org/wiki/Grover%27s_algorithm)
