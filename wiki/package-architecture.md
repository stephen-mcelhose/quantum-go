---
type: concept
title: Package Architecture
description: How quantum-go is split into three packages — math, core, local — their responsibilities, and why the layering prevents circular dependencies.
resource: docs/packages.md
tags: [architecture, go-packages, math, core, local, design]
timestamp: 2026-08-09T03:26:15Z
---

# Package Architecture

quantum-go is organized into three packages with a strict dependency hierarchy: `math` ← `core` ← `local`. Nothing flows upward. The `BlockGate` circular-dependency problem is solved with a `GlobalStepExecutor` function variable (see [[quantum-dsl]]).

## The Three Layers

### `math` — Linear Algebra Foundation
The lowest layer. No quantum-specific code — just matrix operations over `complex128`.

- **Matrix type**: flat `[]complex128` slice (row-major, cache-friendly). See [[quantum-linear-algebra]].
- **Core ops**: `Mul` (matrix multiply), `Tensor` (Kronecker product), `ConjugateTranspose` (†), `IdentityMatrix`, `Exp` (matrix exponentiation for time evolution).
- **Thermodynamics**: `ToDensityMatrix`, `PartialTrace`, `ExpectationValue`, `VonNeumannEntropy` — used by [[quantum-thermodynamics]].
- **Testing**: property-based fuzz tests for matrix dimensions and multiplication correctness.

Nothing in `math` knows what a qubit is.

### `core` — Quantum Circuit Domain Model
The middle layer. Defines the *language* of quantum circuits: [[quantum-dsl]] covers this in depth.

- **Gate interface**: `GetMatrix()`, `GetAffectedQubitIndexes()`, `ApplyOptimize()`, and inverse/size metadata.
- **BaseGate**: default implementation that concrete gates embed.
- **TimeEvolution**: unitary evolution e^{−iHt} generated from a Hamiltonian matrix.
- **Step**: a container of gates that operate on disjoint qubits — can run in parallel.
- **Program**: ordered sequence of Steps plus NumQubits and InitAlpha.
- **Block / BlockGate**: composite gate pattern — a named sub-circuit that can be nested. Requires `GlobalStepExecutor` to simulate itself (set by `local`).
- **Result types**: `CompactResult` (final state only) vs `InstrumentedResult` (intermediate states per step, useful for debugging).
- **Qubit**: probability of being |1⟩ plus measured value.

`core` imports only `math`. It never imports `local`.

### `local` — Simulation Engine
The top layer. Makes circuits actually run.

- **ExecutionEnvironment** (`SimpleExecutionEnvironment`): takes a `Program`, runs each `Step`, and returns a `Result`. See [[simulation-engine]].
- **Computations**: the hot path — optimized gate application using bitwise tricks instead of full matrix exponentiation. See [[gate-application]].
- **Initialization**: sets `core.GlobalStepExecutor` on startup, enabling `BlockGate`s to execute via the local engine. This is the one coupling back to `core`, but it goes through a function variable — no import cycle.

## Dependency Graph

```
local  ──imports──▶  core  ──imports──▶  math
  │                    │
  └─sets──────────▶  GlobalStepExecutor (var in core)
```

`local` sets the executor variable; `core.BlockGate.ApplyOptimize` calls it. This is the only runtime coupling that crosses the boundary upward — and it's type-safe (function signature, not package import).

## Key Design Decisions

| Decision                            | Rationale                                                                                  |
| ----------------------------------- | ------------------------------------------------------------------------------------------ |
| Flat `[]complex128` for matrices    | Cache-friendly sequential memory access; avoids pointer indirection of 2D slice            |
| `Gate` as interface                 | Allows composite gates (BlockGate) and specialized gates (ApplyOptimize) behind one type   |
| `GlobalStepExecutor` function var   | Breaks core↔local import cycle without reflection or plugin architecture                   |
| `CompactResult` vs `InstrumentedResult` | Compact is the default (lower memory); Instrumented is opt-in for debugging          |
| Separate `math` package             | Thermodynamics and linear algebra are reusable outside quantum circuits                    |

## Key Points

- The three-layer hierarchy (`math` → `core` → `local`) is strict — no upward imports.
- The only "escape hatch" is `GlobalStepExecutor`, set by `local` at init, called by `core.Block`.
- `ApplyOptimize` on gates is the specialization hook — gates that can bypass matrix multiply use it.
- `InstrumentedResult` is expensive (stores every intermediate state vector) — use only for debugging.

## Sources

- `docs/packages.md`
