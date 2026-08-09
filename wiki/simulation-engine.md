---
type: concept
title: Simulation Engine
description: How SimpleExecutionEnvironment initializes the state vector from InitAlpha and drives the per-step execution loop that transforms it through a quantum circuit.
resource: local/engine.go
tags: [simulation, engine, state-vector, execution, local, initialization]
timestamp: 2026-08-09T03:26:15Z
---

# Simulation Engine

`local/engine.go` is the entry point for running any quantum program. It contains a single meaningful type — `SimpleExecutionEnvironment` — and one method — `RunProgram`. Everything interesting about *how* gates are applied lives in [[gate-application]].

## State Vector Initialization

For an n-qubit system, the state vector has 2ⁿ entries (one `complex128` amplitude per basis state). The engine allocates it, then initializes from `p.InitAlpha`:

```go
size := 1 << p.NumQubits
state := make([]complex128, size)

for i := 0; i < size; i++ {
    state[i] = 1.0
    for j := 0; j < p.NumQubits; j++ {
        if (i>>j)&1 == 0 {
            state[i] *= complex(p.InitAlpha[j], 0)          // qubit j is |0⟩ contribution
        } else {
            state[i] *= complex(sqrt(1 - InitAlpha[j]²), 0)  // qubit j is |1⟩ contribution
        }
    }
}
```

**What this does:** It constructs the tensor product state |ψ⟩ = |q_{n−1}⟩ ⊗ … ⊗ |q₀⟩ directly without matrix multiplication. Each basis index `i` has bits that encode which qubits are in |0⟩ vs |1⟩. The bit-check `(i>>j)&1` reads "is qubit j in the |1⟩ component of this basis state?".

With all `InitAlpha[j] = 1.0` (the default from `NewProgram`), the only non-zero amplitude is `state[0] = 1.0` — the |00…0⟩ state.

**Non-zero InitAlpha:** Setting `InitAlpha[j] = 0.0` starts qubit j in |1⟩. Intermediate values create superpositions — α|0⟩ + β|1⟩ where α = `InitAlpha[j]` and β = √(1−α²).

## Execution Loop

```go
for _, step := range p.Steps {
    if step.Type == core.StepNormal {
        state = CalculateNewState(step.Gates, state, p.NumQubits)
    }
}
```

Only `StepNormal` steps are executed — `StepPseudo` and `StepProbability` are skipped (they exist for visualization/analysis purposes). `CalculateNewState` is defined in [[gate-application]].

## Measurement

After all gates are applied, per-qubit probabilities are computed from the state vector by `GetQubits()` (which calls `CalculateQubitStatesFromVector` — see [[quantum-dsl]]). Each qubit's `MeasuredValue` is then randomly set:

```go
q.MeasuredValue = rand.Float64() < q.Probability
```

This is Born-rule measurement: collapse to |1⟩ with probability |amplitude|², |0⟩ otherwise. **Note:** this is a classical pseudo-random collapse — the simulator doesn't model true quantum randomness.

## Result

`RunProgram` returns a `*core.CompactResult` containing the full final state vector. If you need intermediate states (e.g., after each step), you'd need to modify the engine to use `InstrumentedResultImpl` — that's opt-in via the `Result` interface design in [[quantum-dsl]].

## Relationship to GlobalStepExecutor

The `init()` in `computations.go` sets `core.GlobalStepExecutor = CalculateNewState`. This means `Block.ApplyOptimize` (used by composite gates like QFT) calls the same `CalculateNewState` function. There's only one simulation kernel — it's shared between the top-level loop and nested block execution.

## Key Points

- State vector size is 2ⁿ complex128 values — allocate once, transform in place per step.
- InitAlpha encodes initial qubit states as tensor products; `1.0` → |0⟩, `0.0` → |1⟩.
- Only `StepNormal` steps are executed; pseudo/probability steps are display artifacts.
- Measurement is probabilistic collapse via `rand.Float64()` against Born-rule probabilities.
- The engine is stateless — it takes a Program and returns a Result; the same Program can be run repeatedly (though measurement outcomes will differ each run due to random collapse).

## Sources

- `local/engine.go`
