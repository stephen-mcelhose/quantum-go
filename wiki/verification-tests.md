---
type: concept
title: Verification Tests
description: local/verification_test.go — the canonical state-vector test suite for 15+ circuits with expected amplitudes, tolerance 1e-6, and the compareStateVectors helper.
resource: local/verification_test.go
tags: [verification, testing, state-vector, bell, ghz, qft, toffoli, fredkin, time-evolution]
timestamp: 2026-08-09T03:26:15Z
---

# Verification Tests

`local/verification_test.go` is the **canonical correctness reference** for quantum-go. It verifies full state-vector amplitudes (not just marginal probabilities) for 15 reference circuits using `compareStateVectors` with tolerance 1e-6. Every gate and major feature should have a corresponding case here.

## The Test Helper

```go
const tolerance = 1e-6

func compareStateVectors(t *testing.T, name string, got, want []complex128) {
    t.Helper()
    if len(got) != len(want) {
        t.Errorf("%s: dimension mismatch", name)
        return
    }
    for i := range got {
        diff := cmplx.Abs(got[i] - want[i])
        if diff > tolerance {
            t.Errorf("%s: mismatch at index %d: got %v, want %v (diff %v)",
                name, i, got[i], want[i], diff)
        }
    }
}
```

Called as: `compareStateVectors(t, tt.name, res.GetProbability(), tt.expected)`

**Note:** `res.GetProbability()` is misnamed — it returns the complex **amplitude** vector, not probabilities. Probability of state i = |GetProbability()[i]|². The naming follows a historical convention; do not be misled.

## Verified Circuits

| Test case                  | Expected state vector                                         | What it verifies                    |
| -------------------------- | ------------------------------------------------------------- | ----------------------------------- |
| Bell State                 | [1/√2, 0, 0, 1/√2]                                           | H + CNOT entanglement               |
| GHZ State (3 qubits)       | [1/√2, 0, 0, 0, 0, 0, 0, 1/√2]                              | 3-qubit entanglement chain          |
| QFT on \|00⟩               | [0.5, 0.5, 0.5, 0.5]                                         | Fourier transform of |00⟩           |
| Time Evolution exp(-iXπ/4) | [1/√2, -i/√2]                                                | Matrix exponential (Hamiltonian)    |
| Rx(π) ≈ X                  | [0, -i]                                                       | Rx(π) = -iX (global phase)         |
| Ry(π) ≈ Y                  | [0, 1]                                                        | Ry(π) = Y (real case)              |
| V·V = X                    | [0, 1]                                                        | √X gate composition                 |
| H then S                   | [1/√2, i/√2]                                                  | S gate phase rotation               |
| 2-Qubit Rx(π)              | [0, 0, 0, -1]                                                 | Parallel rotation (tensor product)  |
| 2-Qubit Ry(π)              | [0, 0, 0, 1]                                                  | Parallel Ry                         |
| 2-Qubit H then S           | [0.5, 0.5i, 0.5i, -0.5]                                      | Multi-qubit phase composition        |
| Toffoli (110 → 111)        | [0, 0, 0, 0, 0, 0, 0, 1]                                     | CCNOT target flip                  |
| Fredkin (110 → 101)        | [0, 0, 0, 0, 0, 1, 0, 0]                                     | CSWAP target exchange              |
| Superdense Coding (11)     | [0, 0, 0, 1]                                                  | `core.NewSuperdenseCodingProgram()`  |

## Key State Vector Facts

**Little-endian ordering (same as Qiskit):**
- Index 0 → |000⟩, Index 1 → |001⟩, Index 7 → |111⟩
- Index i: bit j of i = state of qubit j

**Bell State:**
```
[1/√2, 0, 0, 1/√2]
 |00⟩       |11⟩   ← only these two have amplitude
```

**GHZ State:**
```
[1/√2, 0, 0, 0, 0, 0, 0, 1/√2]
 |000⟩                   |111⟩  ← only poles
```

**QFT on |00⟩:**
```
[0.5, 0.5, 0.5, 0.5]   ← uniform superposition (real parts only, starting from |00⟩)
```

**Toffoli test:**
- Initialize |110⟩: X(0) and X(1) → state index `3` = 0b011 (little-endian: q0=1, q1=1, q2=0)
- Toffoli(0, 1, 2): both controls=1 → flip q2 → index `7` = 0b111 → expected=[0,0,0,0,0,0,0,1]

**Rx(π) result:**
- Rx(π) matrix: `[[0, -i], [-i, 0]]`
- Applied to |0⟩ = [1, 0]: result = [0, -i] = -i|1⟩
- Global phase -i is physically irrelevant but measured here for mathematical precision

## How to Add a New Verification Test

```go
{
    name: "My New Gate",
    program: func() *core.Program {
        p := core.NewProgram(n)
        // ... build circuit
        return p
    }(),
    expected: []complex128{
        // 2^n amplitudes in little-endian order
    },
},
```

Then add it to the `tests` slice in `TestVerifyStandardStates`. The framework handles all comparison automatically.

## Key Points

- `res.GetProbability()` returns amplitudes, not probabilities — divide by |·|² for actual probability. See [[state-vector-model]] for the full naming history.
- All expected amplitudes must be normalized (sum of squares = 1).
- Tolerance 1e-6 is tight enough to catch gate implementation bugs but forgiving of float64 rounding.
- Toffoli and Fredkin tests use little-endian index encoding — verify index arithmetic when debugging.
- The `Superdense Coding` test exercises `core.NewSuperdenseCodingProgram()` — a built-in program helper from [[circuits-library]].
- Adding new gates → add a verification test case here first, then implement the gate. The full gate addition workflow is in [[how-to-add-a-new-gate]].
- For the broader test suite context see [[testing-strategy]]; for randomised property tests see [[fuzz-testing]].

## Sources

- `local/verification_test.go`
