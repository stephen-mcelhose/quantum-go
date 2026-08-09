---
type: concept
title: Fuzz Testing
description: Go 1.18+ fuzz tests in local/fuzz_test.go and math/fuzz_test.go — FuzzToffoli, FuzzFourier (unitarity), FuzzAdd (arithmetic), FuzzMatrixMul, FuzzMatrixTensor.
resource: local/fuzz_test.go, math/fuzz_test.go
tags: [fuzz-testing, go-fuzz, toffoli, fourier, unitarity, add, matrix]
timestamp: 2026-08-09T03:26:15Z
---

# Fuzz Testing

quantum-go uses Go 1.18+ native fuzzing (`testing.F`) to discover edge cases beyond the fixed verification tests. Two packages contain fuzz tests: `local` (simulator) and `math` (linear algebra).

## Running Fuzz Tests

```bash
# Run specific fuzz targets
go test -fuzz=FuzzToffoli -fuzztime=30s ./local/...
go test -fuzz=FuzzFourier -fuzztime=30s ./local/...
go test -fuzz=FuzzAdd -fuzztime=60s ./local/...
go test -fuzz=FuzzMatrixMul -fuzztime=30s ./math/...
go test -fuzz=FuzzMatrixTensor -fuzztime=30s ./math/...

# Run seeds only (no mutation, fast verification)
go test -run=FuzzToffoli ./local/...
```

## local/fuzz_test.go

### FuzzToffoli

**Purpose:** Verify Toffoli truth table for all 8 three-qubit computational basis states.

**Seed corpus:** 8 seeds (0-7, one per initial state).

**Property tested:**
- If both control bits (q0=1 AND q1=1): target q2 is flipped
- Otherwise: q2 is unchanged

```go
expectedState := initialState
if (initialState & 3) == 3 {  // q0=1 and q1=1
    expectedState ^= 4         // flip q2 (bit 2)
}
// verify: only expectedState has probability > 0.99
```

**What it finds:** Incorrect bit ordering in Toffoli's optimized path, off-by-one errors in the `AffectedQubits` slice.

**Input constraint:** `initialState ∈ [0, 7]` — fuzzer generates integers, guarded at top.

### FuzzFourier

**Purpose:** Verify QFT × IQFT = Identity (unitarity property).

**Seeds:** 2, 3, 4 qubits.

**Property tested:**
1. **Round-trip identity:** QFT then IQFT applied to |0…0⟩ returns to |0…0⟩ with probability > 0.99.
2. **Normalization:** Total probability sums to 1.0 ± 1e-9.

```go
p.AddStep(core.NewStep(core.NewFourier(numQubits, 0)))
invFourier := core.NewFourier(numQubits, 0)
invFourier.SetInverse(true)
p.AddStep(core.NewStep(invFourier))
// Assert: GetProbability()[0] ≈ 1.0
```

**Input constraint:** `numQubits ∈ [2, 6]` — beyond 6 the state vector is too large for fast fuzzing.

**What it finds:** Phase errors in QFT, incorrect inverse flag propagation in BlockGate.

### FuzzAdd

**Purpose:** Verify quantum addition `x + y mod 2^m` for random x, y, m.

**Seeds:** (2,1,1) and (3,3,1).

**Property tested:** The x register (qubits 0..m-1) should contain `(x+y) mod 2^m` after applying `NewAdd`.

```go
expectedX := (x + y) % maxVal
for i := 0; i < m; i++ {
    expectedBit := (expectedX >> i) & 1
    // qubits[i].Probability should be ~1.0 if bit is 1, ~0.0 if bit is 0
}
```

**Input constraint:** `m ∈ [1, 4]`, x and y are clamped to `[0, 2^m)`.

**What it finds:** Phase accumulation errors in CR gates, off-by-one in register indexing.

## math/fuzz_test.go

### FuzzMatrixMul

**Purpose:** Matrix multiplication doesn't panic and produces correct output dimensions.

**Property tested:**
1. `Mul(A, B)` doesn't panic for valid dimensions (c1 = r2 enforced).
2. Output dimensions: `result.Rows == r1`, `result.Cols == c2`.

**Input constraint:** All dimensions in [1, 10], c1 forced equal to r2 to satisfy compatibility.

**What it finds:** Index out of bounds, dimension miscalculation in the flat-array multiplication.

### FuzzMatrixTensor

**Purpose:** Kronecker product doesn't panic and produces correct output dimensions.

**Property tested:** `Tensor(A, B)` produces `(r1*r2) × (c1*c2)` matrix.

**Input constraint:** All dimensions in [1, 5] — tensor product grows as product of inputs.

**What it finds:** Overflow in dimension computation, allocation errors for large tensors.

## Design Patterns

**Guard clauses first:** All fuzz functions start by validating inputs and returning early for out-of-range values. This focuses the fuzzer on the interesting domain without crashing on noise inputs.

**Panic recovery in math fuzz tests:**
```go
defer func() {
    if r := recover(); r != nil {
        t.Errorf("Mul panicked: %v", r)
    }
}()
```
The simulator fuzz tests do NOT use panic recovery — panics are expected to surface bugs in the engine.

**Seed corpus selection:** Seeds are chosen to cover boundary conditions (small qubits, unit values) that verify basic correctness before the fuzzer mutates.

## Key Points

- FuzzFourier tests **unitarity** — QFT followed by IQFT must return to |0…0⟩.
- FuzzToffoli exhaustively covers all 8 3-qubit states via 8 seeds (fuzzer only adds noise).
- FuzzAdd verifies the quantum Draper adder against classical arithmetic for m up to 4 bits.
- Math fuzz tests check that matrix operations don't panic on random dimensions.
- `numQubits > 6` is rejected in FuzzFourier — state vector would be 64+ elements (acceptable but slower).
- These fuzz tests run as unit tests with their seed corpus when `-fuzz` is not specified.

## Sources

- `local/fuzz_test.go`
- `math/fuzz_test.go`
