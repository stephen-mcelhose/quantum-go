---
type: concept
title: Simulator Optimizations
description: The conceptual rationale for quantum-go's bitwise "bit-loop" strategy — why O(4^n) general matrix multiply is avoided and how each gate family is handled instead.
resource: docs/optimizations.md
tags: [optimization, performance, bitwise, complexity, O(2^n), simulation]
timestamp: 2026-08-09T03:26:15Z
---

# Simulator Optimizations

The core insight behind quantum-go's performance: **you never need to build a 2ⁿ×2ⁿ matrix to apply a gate**. The naive approach — materialize the full unitary and do matrix-vector multiply — is O(4ⁿ), which becomes prohibitive around n=20. The bit-loop strategy keeps everything at O(2ⁿ).

## Why Naive Matrix Multiply Is Expensive

To apply a 1-qubit gate H to qubit 0 of an n-qubit system, the naive approach:
1. Builds the 2ⁿ×2ⁿ Kronecker product H ⊗ I ⊗ … ⊗ I (see [[quantum-linear-algebra]])
2. Multiplies it against the 2ⁿ state vector

Step 1 costs O(4ⁿ) time and O(4ⁿ) memory. For n=20, that's 1 trillion operations and 128 TB of memory. Clearly unworkable.

## The Bit-Loop Strategy

Instead of building the large matrix, observe that qubit `idx` corresponds to **bit `idx` of the state vector index**. States where bit `idx` = 0 pair with states where bit `idx` = 1 (all other bits identical). The 2×2 gate matrix only mixes amplitudes within each pair:

```
(v[j], v[j+2^idx]) → (m₀₀·v[j] + m₀₁·v[j+2^idx], m₁₀·v[j] + m₁₁·v[j+2^idx])
```

This scans 2^{n-1} pairs → O(2ⁿ) total. The 2ⁿ×2ⁿ matrix is never built.

For controlled gates: check the control bit of each index first (`(i>>c)&1 == 1`), then apply only if the condition holds. For CNOT, the XOR trick `v[i ^ (1<<target)]` fetches the "bit-flipped" partner index in O(1).

See [[gate-application]] for the Go implementation of each optimization.

## Complexity Comparison

| Gate family               | Naive matrix approach | Bit-loop optimized    |
| ------------------------- | --------------------- | --------------------- |
| Single-qubit gate         | O(4ⁿ)                 | O(2ⁿ)                 |
| Two-qubit gate (CNOT/CZ)  | O(4ⁿ)                 | O(2ⁿ)                 |
| k-qubit block (QFT, Add)  | O(4ⁿ)                 | O(Steps × 2ⁿ)         |

The block gate optimization (see [[composite-gates]]) is especially important for QFT and Adder, which internally contain many steps. Instead of materializing a single 2ⁿ×2ⁿ unitary for the whole block, the engine runs each sub-step individually. An n-qubit QFT has n(n+1)/2 gate operations → O(n² × 2ⁿ) total, far better than O(4ⁿ).

## Additional Micro-Optimizations

| Optimization             | Where                              | Savings                               |
| ------------------------ | ---------------------------------- | ------------------------------------- |
| Identity gate skip       | `applyGate` type check             | O(1) — no iteration at all            |
| Zero-skip in Tensor      | `math.Tensor` inner loop           | Reduces work for sparse gate matrices |
| Flat array for Matrix    | `math.Matrix` struct               | Cache locality vs 2D slice pointers   |
| No intermediate allocs   | Bit-loop reuses `answer` slice     | Reduces GC pressure                   |

## Practical Limits

These optimizations keep the simulator practical up to ~30 qubits (16 GB RAM). Beyond that, the state vector itself becomes the bottleneck — O(2ⁿ) storage is irreducible for full state vector simulation. See [[scaling-and-limits]] for the qubit/memory table.

## Key Points

- O(4ⁿ) → O(2ⁿ) by exploiting bit structure of state vector indices.
- Controlled gates need only check one bit per index — no matrix needed.
- XOR flipping (`i ^ (1<<target)`) finds the paired amplitude in O(1).
- Block gates run step-by-step → O(Steps × 2ⁿ), not O(4ⁿ).
- Beyond ~30 qubits, memory is the limit regardless of algorithmic optimizations.

## Sources

- `docs/optimizations.md`
