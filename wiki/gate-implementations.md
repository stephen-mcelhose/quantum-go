---
type: concept
title: Gate Implementations
description: How core/gates.go implements fundamental gate matrices — H, X, Y, Z, CNOT, CZ, SWAP, Toffoli, Fredkin — as concrete Go types embedding BaseGate.
resource: core/gates.go
tags: [gates, implementation, hadamard, pauli, cnot, toffoli, fredkin, go]
timestamp: 2026-08-09T03:26:15Z
---

# Gate Implementations

`core/gates.go` contains the concrete implementations of all fundamental gate types. Each gate embeds `BaseGate` (from [[quantum-dsl]]) and provides a `GetMatrix()` method returning a flat `complex128` matrix (from [[quantum-linear-algebra]]). The `[[gate-application]]` engine reads these matrices when no specialized bitwise path exists.

## Implementation Pattern

All gates follow the same structure:
```go
type Hadamard struct {
    BaseGate
}

func NewHadamard(idx int) *Hadamard {
    return &Hadamard{BaseGate: BaseGate{
        AffectedQubits: []int{idx},
        Caption: "H",
        Name: "Hadamard",
    }}
}

func (g *Hadamard) GetMatrix() math.Matrix {
    m := math.NewMatrix(2, 2)
    h := complex(math.HV, 0)  // 1/√2
    m.Data = []complex128{h, h, h, -h}
    return m
}
```

The Caption string (e.g., `"CNOT"`, `"CZ"`) is what the [[gate-application]] engine checks to route to its optimized path.

## Gate Catalog

### Single-Qubit Gates

| Type      | Caption | Matrix (flat Data slice)                      | Notes                              |
| --------- | ------- | --------------------------------------------- | ---------------------------------- |
| Hadamard  | "H"     | [HV, HV, HV, -HV] where HV=1/√2              | H²=I; maps Z↔X basis              |
| X         | "X"     | [0, 1, 1, 0]                                  | Quantum NOT; X²=I                  |
| Y         | "Y"     | [0, -i, i, 0]                                 | Bit flip + phase flip; Y²=I       |
| Z         | "Z"     | [1, 0, 0, -1]                                 | Phase flip; Z²=I                  |
| Identity  | "I"     | [1, 0, 0, 1]                                  | No-op; **skipped by engine** (O(1)) |
| Measurement| "M"    | [1, 0, 0, 1] (identity, handled specially)    | Collapse via rand.Float64()        |

### Two-Qubit Gates (4×4 matrices)

**CNOT:**
```
[1, 0, 0, 0]
[0, 1, 0, 0]
[0, 0, 0, 1]   ← |10⟩ → |11⟩
[0, 0, 1, 0]   ← |11⟩ → |10⟩
```
The engine uses XOR (`v[i ^ (1<<target)]`) instead of this matrix — see [[gate-application]].

**CZ:**
```
diag(1, 1, 1, -1)   ← only |11⟩ gets phase flip
```

**SWAP:**
```
[1, 0, 0, 0]
[0, 0, 1, 0]   ← |01⟩ → |10⟩
[0, 1, 0, 0]   ← |10⟩ → |01⟩
[0, 0, 0, 1]
```

### Three-Qubit Gates (8×8 matrices)

**Toffoli (CCNOT):**
- Identity for all states where both controls aren't |1⟩ (first 6 diagonal entries = 1).
- Swaps |110⟩ ↔ |111⟩ (off-diagonal entries at rows/cols 6,7).
- Engine uses `v[i ^ (1<<c)]` when both control bits = 1 — see [[gate-application]].

**Fredkin (CSWAP):**
- 8×8 identity with rows 3 and 5 swapped.
- When control (bit 0) = 1: swaps targets bits 1 and 2.
- `|101⟩` (index 5) ↔ `|011⟩` (index 3) — note: assumes a=0, b=1, c=2 ordering.

## Why Matrices Are Flat

`m.Data = []complex128{...}` — a single slice allocation for the entire matrix. The `Get(row, col)` method computes `Data[row*Cols + col]`. This is the flat-array design from [[quantum-linear-algebra]] for cache efficiency.

## The `HV` Constant

`math.HV = 1.0 / math.Sqrt(2.0)` is used by both Hadamard and is the most-referenced constant in quantum-go. Having it pre-computed avoids repeated sqrt calls.

## Caption Routing to Optimized Paths

When [[gate-application]] receives a gate, it checks `gate.GetCaption()` for specific values:
- `"CNOT"` → XOR flip optimization
- `"CZ"` → conditional phase flip
- `"SWAP"` → XOR exchange
- `"CR"` → scalar multiply (from [[rotation-implementations]])
- `"CCNOT"` → double-control XOR flip

All other gates fall to the general matrix-vector multiply path.

## Key Points

- Every gate is a struct embedding `BaseGate` with a `GetMatrix()` override.
- Caption is used by [[gate-application]] for dispatch — changing it breaks optimizations.
- Identity is type-checked (`if _, ok := gate.(*core.Identity)`) and skipped in O(1).
- Measurement's matrix is identity — measurement is handled by the engine post-execution.
- Fredkin (CSWAP) is the controlled-SWAP — useful for sorting networks on qubits.
- Toffoli is the syndrome correction gate in the 3-qubit bit-flip code — see [[error-correction]] for the full circuit.
- See [[rotation-implementations]] for parameterized gates (Rx/Ry/Rz/U/S/T/V/CR/PhaseShift).
- For the complete gate reference table (constructors, captions, matrices, QASM mnemonics for every gate), see [[gate-zoo]].
- To add a new gate type, see [[how-to-add-a-new-gate]].

## Sources

- `core/gates.go`
