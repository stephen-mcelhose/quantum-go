---
type: concept
title: Gate Application — Bitwise Optimizations
description: How computations.go applies quantum gates to the state vector using bitwise index tricks instead of O(4^n) matrix multiplication, achieving O(2^n) complexity.
resource: local/computations.go
tags: [optimization, bitwise, state-vector, gate-application, simulation, performance]
timestamp: 2026-08-09T03:26:15Z
---

# Gate Application — Bitwise Optimizations

`local/computations.go` is the hot path of the quantum-go simulator. It dispatches gate application through a priority chain: identity skip → block gate → optimized path → general matrix multiply. The key insight: **you never need to build a 2ⁿ×2ⁿ matrix for most common gates**. See [[simulator-optimizations]] for the conceptual explanation; this page covers the implementation.

## Dispatch Chain in `applyGate`

```
applyGate(gate, v, controls)
  ├── Identity? → return v unchanged (O(1))
  ├── ControlledGate? → merge controls, recurse on inner gate
  ├── ControlledBlockGate? → merge controls, applyBlock with inverse flag
  ├── BlockGateInterface? → applyBlock (runs sub-steps, see [[quantum-dsl]])
  ├── no controls + HasOptimization()? → gate.ApplyOptimize(v)  ← gate-specific fast path
  └── fallback → processGate(gate, v, controls)
```

The `GlobalStepExecutor` is set to `CalculateNewState` by `computations.go`'s `init()`, making this the one kernel shared between top-level execution and nested block gates.

## Single-Qubit Gate — O(2ⁿ) via Bit Grouping

For a gate on qubit index `idx`, the state vector splits into paired amplitudes: index `j` (qubit at 0) pairs with `j + 2^idx` (qubit at 1). The 2×2 gate matrix is applied to each pair:

```go
qdelta := 1 << idx
ngroups := size / (2 * qdelta)

for group := 0; group < ngroups; group++ {
    for j := 2*group*qdelta; j < (2*group+1)*qdelta; j++ {
        if checkControls(j, controls) {
            v0 := v[j]
            v1 := v[j+qdelta]
            answer[j]       = m[0,0]*v0 + m[0,1]*v1
            answer[j+qdelta] = m[1,0]*v0 + m[1,1]*v1
        }
    }
}
```

**Why this works:** In a flat state vector indexed by an n-bit integer, bit `idx` being 0 or 1 determines which of the pair you're in. Iterating `j` over the "lower half" (bit `idx` = 0) and pairing with `j + qdelta` (bit `idx` = 1) naturally enumerates all 2^{n-1} pairs. Cost: O(2ⁿ) multiplications vs O(4ⁿ) for full matrix-vector multiply.

## Two-Qubit Gate Specializations

CNOT, CZ, SWAP, and CR each have hand-coded loops that skip matrix construction entirely:

| Gate  | Condition                                       | Action                                      |
| ----- | ----------------------------------------------- | ------------------------------------------- |
| CNOT  | `(i >> control) & 1 == 1`                       | `answer[i] = v[i ^ (1<<target)]` — XOR flips target bit |
| CZ    | both control and target bits = 1                | `answer[i] = -v[i]` — phase flip           |
| SWAP  | control bit ≠ target bit                        | `answer[i] = v[j]` where j has bits swapped |
| CR    | both bits = 1                                   | `answer[i] = v[i] * rot` — apply phase factor |

The XOR trick for CNOT is elegant: if bit `target` is currently 0, `i ^ (1<<target)` gives the index with that bit flipped to 1, and vice versa. This is a O(2ⁿ) amplitude shuffle with no multiply needed.

## Three-Qubit Gate: Toffoli (CCNOT)

```go
if (i>>a)&1 == 1 && (i>>b)&1 == 1 {
    answer[i] = v[i^(1<<c)]  // flip target c if both a and b are 1
}
```

Same XOR-flip pattern, now conditional on two control bits. Still O(2ⁿ).

## Control Qubit Propagation

`checkControls(i, controls)` verifies that all control bits in index `i` are 1:

```go
for _, c := range controls {
    if (i>>c)&1 == 0 {
        return false
    }
}
return true
```

This enables arbitrary controlled gates — a `ControlledGate` wrapper just appends to the controls list and recurses. Toffoli is effectively CNOT with an extra control; this propagation mechanism handles it generically.

## Fallback: Full Matrix Multiply

When no specialization applies (e.g., an exotic n-qubit gate with a full 2ⁿ×2ⁿ matrix):

```go
for i := 0; i < size; i++ {
    if checkControls(i, controls) {
        var sum complex128
        for j := 0; j < size; j++ {
            sum += matrix.Get(i, j) * v[j]
        }
        answer[i] = sum
    }
}
```

This is O(4ⁿ) — avoid for large n. In practice, composite gates use the Block path (which recursively applies sub-steps) rather than materializing a large matrix, so this fallback is rarely hit. See [[simulator-optimizations]] for the complexity table.

## Block Gate Execution (`applyBlock`)

For composite gates (QFT, Adder, etc. — see [[composite-gates]]):

```go
func applyBlock(block *core.Block, v []complex128, controls []int, inverse bool) []complex128 {
    if inverse {
        // iterate steps in reverse, toggle each gate's inverse flag
    } else {
        for _, step := range block.Steps {
            for _, g := range step.Gates {
                w = applyGate(g, w, controls)  // recursive!
            }
        }
    }
}
```

Inverse blocks walk steps in reverse order and toggle `IsInverse()` on each gate before applying. This automatically implements IQFT from the QFT Block definition without a separate code path.

## Key Points

- Identity gates cost O(1) — type-checked and skipped immediately.
- Single-qubit gates use paired amplitude grouping: O(2ⁿ) regardless of circuit size.
- CNOT/CZ/SWAP/Toffoli use XOR/bitwise tricks: O(2ⁿ), no matrix needed.
- Block gates are executed step-by-step (O(Steps × 2ⁿ)), avoiding O(4ⁿ) for composite gates.
- Control qubit propagation via `checkControls` makes arbitrary controlled-gate combinations work generically.
- The full matrix fallback (O(4ⁿ)) exists but should be avoided for large circuits.

## Sources

- `local/computations.go`
- `docs/optimizations.md`
