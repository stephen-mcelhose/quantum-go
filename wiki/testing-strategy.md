---
type: concept
title: Testing Strategy
description: How quantum-go's test suite is organized — unit tests, integration tests, verification tests, fuzz tests, and example tests — with commands to run them.
resource: local/engine_test.go, core/core_test.go
tags: [testing, unit-tests, integration, verification, fuzz, go-test]
timestamp: 2026-08-09T03:26:15Z
---

# Testing Strategy

quantum-go uses a layered test strategy covering unit tests, integration tests, state-vector verification tests, and fuzz tests. All tests run with standard `go test`.

## Test Layers

| Layer                   | Location                       | What it tests                                     |
| ----------------------- | ------------------------------ | ------------------------------------------------- |
| Unit — core types       | `core/core_test.go`            | BaseGate, Step, Program, gate matrices            |
| Unit — circuits         | `core/circuits_test.go`        | ExpectedBellState, GHZ, QFT helpers               |
| Unit — math             | `math/matrix_test.go`          | Mul, Tensor, ConjugateTranspose, constants        |
| Integration — engine    | `local/engine_test.go`         | Full circuit execution (Bell, GHZ, gates, add)   |
| Verification            | `local/verification_test.go`   | State vector equality to 1e-6 for 15+ circuits   |
| Fuzz — simulator        | `local/fuzz_test.go`           | Toffoli, Fourier, Add across random inputs        |
| Fuzz — math             | `math/fuzz_test.go`            | Matrix Mul, Tensor with random dimensions         |
| Parser                  | `qasm/parser_test.go`          | QASM 2.0 parsing round-trips                     |
| Examples                | `examples/**/*_test.go`        | Algorithm demos (Grover, Shor, QKD, teleportation)|

## Running Tests

```bash
# All tests
go test ./...

# Single package
go test ./local/...
go test ./core/...
go test ./math/...
go test ./qasm/...

# Verbose with output
go test -v ./local/...

# Run specific test
go test -run TestBellState ./local/...
go test -run TestVerifyStandardStates ./local/...

# Run fuzz tests (Go 1.18+)
go test -fuzz=FuzzToffoli ./local/...
go test -fuzz=FuzzFourier ./local/...
go test -fuzz=FuzzAdd ./local/...

# Run fuzz with time limit
go test -fuzz=FuzzMatrixMul -fuzztime=30s ./math/...

# Race detector (important for simulator goroutine safety)
go test -race ./...
```

## Key Testing Patterns

### State Vector Comparison

The gold standard in `local/verification_test.go` — compare the full 2ⁿ-dimensional state vector with expected amplitudes to tolerance 1e-6:

```go
const tolerance = 1e-6

func compareStateVectors(t *testing.T, name string, got, want []complex128) {
    for i := range got {
        diff := cmplx.Abs(got[i] - want[i])
        if diff > tolerance { t.Errorf(...) }
    }
}
```

**Why complex amplitudes?** `res.GetProbability()` returns the state vector (amplitudes, not probabilities). Probability = |amplitude|². The verification tests compare amplitudes directly, which is stricter than probability comparison because it catches phase errors too.

### Marginal Probability Tests

`local/engine_test.go` uses marginal probability (single-qubit probability) for entanglement tests:
```go
qubits[0].Probability  // P(qubit 0 = |1⟩)
qubits[0].Measure()    // deterministic result for known states
```

### Panic Tests

`core/core_test.go` verifies that qubit overlap in a step panics:
```go
func TestStepUnique(t *testing.T) {
    defer func() {
        if r := recover(); r == nil { t.Error("expected panic") }
    }()
    core.NewStep(core.NewHadamard(0), core.NewX(0))  // same qubit → panic
}
```

## Test Organization Principles

1. **`_test.go` in same package** (white-box): `core/core_test.go` (package `core_test` but tests internals via exported API)
2. **`_test.go` with `_test` suffix** (black-box): `local/engine_test.go` (package `local_test`)
3. **Example tests with expected output**: examples use `testing.T` assertions, not `t.Example*` output matching
4. **Each algorithm has a test**: `examples/algorithms/grover_test.go`, `shor_test.go`, etc.

## Tolerance

All floating-point comparisons use `1e-6`:
- Complex amplitude comparison: `cmplx.Abs(got - want) > 1e-6`
- Probability marginal: `qubits[i].Probability` within ±0.01 of expected
- Fuzz normalization: `math.Abs(totalProb - 1.0) > 1e-9`

The 1e-6 tolerance is chosen to be above float64 rounding (~1e-15) while below physically meaningful differences.

## Key Points

- `local/verification_test.go` is the most authoritative test file — it compares full state vectors for 15+ canonical circuits.
- `res.GetProbability()` returns amplitudes (complex), not probabilities (real) — see [[simulation-engine]].
- Fuzz tests in `local/fuzz_test.go` cover Toffoli, Fourier identity, and quantum addition.
- `TestStepUnique` confirms the panic behavior for overlapping gates — a design-by-contract check.
- Race detector (`-race`) is important if multi-circuit parallelism is added.

## Sources

- `local/engine_test.go`
- `core/core_test.go`
