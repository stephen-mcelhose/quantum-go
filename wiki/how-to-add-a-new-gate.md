---
type: synthesis
title: How to Add a New Gate
description: Step-by-step guide for implementing a new quantum gate in quantum-go — from struct definition through GetMatrix, registration in the engine, QASM export, and verification test.
tags: [how-to, gate, implementation, contribution, testing]
timestamp: 2026-08-09T03:26:15Z
---

# How to Add a New Gate

This guide walks through every file you need to touch to add a new gate, using the pattern established by existing gates.

## 1. Decide the gate type

| Gate kind                                | Where to add it         | Struct embeds    |
| ---------------------------------------- | ----------------------- | ---------------- |
| Fixed single-qubit (no params)           | `core/gates.go`         | `BaseGate`       |
| Parameterised single-qubit (angle etc.)  | `core/rotations.go`     | `BaseGate`       |
| Fixed multi-qubit (CNOT-like)            | `core/gates.go`         | `BaseGate`       |
| Composite / built from primitives        | `core/composite.go`     | `BlockGate`      |
| Arithmetic (QFT-based)                   | `core/arithmetic.go`    | `BlockGate`      |
| Custom matrix (oracle-style)             | `core/oracles.go`       | `BaseGate`       |

## 2. Define the struct and constructor (core/gates.go or core/rotations.go)

```go
// MyGate is a single-qubit gate that does X.
type MyGate struct {
    BaseGate
    // add parameters here if needed
    Theta float64
}

// NewMyGate creates a MyGate acting on qubit at index idx.
func NewMyGate(theta float64, idx int) *MyGate {
    g := &MyGate{Theta: theta}
    g.name = "MyGate"
    g.caption = "MY"
    g.size = 1                   // number of qubits
    g.affectedQubits = []int{idx}
    return g
}
```

**Rules:**
- `name` — human-readable, used in error messages (e.g. `"Hadamard"`)
- `caption` — short display label, used as the engine dispatch key (e.g. `"H"`)
- `size` — number of qubits in `AffectedQubits`; must equal `len(affectedQubits)`
- Multi-qubit: list all qubits in `affectedQubits`, highest index last

## 3. Implement GetMatrix()

```go
// GetMatrix returns the 2×2 unitary matrix for MyGate.
func (g *MyGate) GetMatrix() math.Matrix {
    if g.IsInverse() {
        // return conjugate transpose
        return math.ConjugateTranspose(g.matrix())
    }
    return g.matrix()
}

func (g *MyGate) matrix() math.Matrix {
    m := math.NewMatrix(2, 2)
    // fill in the 2×2 unitary
    m.Set(0, 0, complex(math.Cos(g.Theta/2), 0))
    m.Set(0, 1, complex(-math.Sin(g.Theta/2), 0))
    m.Set(1, 0, complex(math.Sin(g.Theta/2), 0))
    m.Set(1, 1, complex(math.Cos(g.Theta/2), 0))
    return m
}
```

**Requirements:**
- Matrix must be **unitary**: M†M = I (checked by [[fuzz-testing]] round-trips)
- For parameterised gates, check `g.IsInverse()` and return the conjugate transpose

## 4. Add an engine optimisation (optional, local/engine.go)

For simple single-qubit gates, the engine's generic path (`applySingleQubitGate`) handles them automatically via `GetMatrix()`. No changes needed.

For performance-critical or structurally special gates (like CNOT, CZ, SWAP), add a fast path in `applyGate`:

```go
case "MY":
    // fast bitwise implementation
    applyMyGateFast(state, gate)
```

Only do this if the gate will appear in tight loops (e.g. QFT's CR gates). See [[simulator-optimizations]] for the pattern.

## 5. Add QASM export (core/qasm.go)

In `gateToQASM`, add a case:

```go
case *MyGate:
    return fmt.Sprintf("mygate(%.6f) q[%d];\n", gate.Theta, indices[0])
```

If there's no QASM 2.0 equivalent, emit a comment:
```go
default:
    return fmt.Sprintf("// gate %s q%v\n", g.GetName(), indices)
```

## 6. Add a verification test (local/verification_test.go)

Add a test case to `TestVerifyStandardStates`:

```go
{
    name: "MyGate(π/2) on |0⟩",
    program: func() *core.Program {
        p := core.NewProgram(1)
        p.AddStep(core.NewStep(core.NewMyGate(math.Pi/2, 0)))
        return p
    }(),
    expected: []complex128{
        complex(math.Cos(math.Pi/4), 0),  // amplitude of |0⟩
        complex(math.Sin(math.Pi/4), 0),  // amplitude of |1⟩
    },
},
```

Calculate `expected` analytically: apply your gate matrix to |0⟩ = [1, 0]ᵀ by hand.

## 7. Add to TestGatesMatrix (core/core_test.go)

```go
core.NewMyGate(0.5, 0),
```

This checks that `GetMatrix()` returns a non-empty matrix.

## 8. Add to the gate-zoo

Update `wiki/gate-zoo.md` table with the new gate's constructor, caption, matrix summary, and QASM mnemonic.

## Checklist

- [ ] Struct + constructor in the right `.go` file
- [ ] `GetMatrix()` returns a unitary matrix; `IsInverse()` respected
- [ ] `caption` is unique (used as dispatch key in engine)
- [ ] QASM export case added in `gateToQASM`
- [ ] `TestVerifyStandardStates` case with analytically computed expected state
- [ ] `TestGatesMatrix` entry
- [ ] [[gate-zoo]] table updated
- [ ] Run `go test ./...` — all tests pass

## Key Points

- The engine dispatches via `g.GetCaption()` — every caption must be unique across all gates.
- `BaseGate.affectedQubits` controls which qubit indices the engine operates on; wrong indices = silent wrong results.
- BlockGates (composite) never implement `GetMatrix()` — they implement `ApplyOptimize()` instead. See [[composite-gates]].
- If the gate has no standard QASM equivalent, the export emits a comment. That is acceptable.
- New fuzz targets can be added to `local/fuzz_test.go` for parameterised gates — see [[fuzz-testing]].
- For the full picture of how tests are organized across the codebase, see [[testing-strategy]].
