---
type: concept
title: Quantum Thermodynamics
description: Density matrices, expectation values, Von Neumann entropy — how quantum-go models energy, work, and information using math.ToDensityMatrix, ExpectationValue, and VonNeumannEntropy.
resource: examples/thermodynamics/engine_cycle_test.go
tags: [thermodynamics, density-matrix, expectation-value, von-neumann-entropy, hamiltonian, time-evolution]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Thermodynamics

quantum-go includes density matrix and entropy tools in the `math` package for connecting quantum circuits to thermodynamic quantities. This is distinct from running circuits — these tools analyze the *statistical* properties of quantum states.

## Key Concepts

### Hamiltonian = Energy Landscape

In quantum mechanics, the Hamiltonian H is a Hermitian matrix encoding the energy of each eigenstate. For a qubit in a magnetic field (H = −Bσz):

```go
hamiltonian := smath.NewMatrix(2, 2)
hamiltonian.Data = []complex128{-1, 0, 0, 1}  // Energy |0⟩ = -1, Energy |1⟩ = +1
```

### Density Matrix

A density matrix ρ = |ψ⟩⟨ψ| encodes all observable properties of a pure state. For a state vector `v`:

```go
rho := smath.ToDensityMatrix(res.GetProbability())
```

`ToDensityMatrix` computes the outer product of the amplitude vector with its conjugate. For a pure state |ψ⟩ with amplitudes [α, β]:
```
ρ = |ψ⟩⟨ψ| = [[|α|², αβ*], [α*β, |β|²]]
```

- Pure state: ρ² = ρ (trace = 1)
- Mixed state: ρ² ≠ ρ (trace < 1)

### Expectation Value

The expected energy (internal energy U) = Tr(ρH):

```go
u := real(smath.ExpectationValue(rho, hamiltonian))
```

For state |0⟩ in H = −σz: U = Tr(|0⟩⟨0| × diag(−1, +1)) = −1.  
For state |+⟩: U = 0.5(−1) + 0.5(+1) = 0.

### Von Neumann Entropy

Quantum analog of Shannon entropy: S = −Tr(ρ log ρ)

```go
entropy := smath.VonNeumannEntropy(rho)
```

| State         | Entropy      | Physical meaning                            |
| ------------- | ------------ | ------------------------------------------- |
| Pure |0⟩     | 0            | Zero entropy — perfect knowledge             |
| Pure |+⟩     | 0            | Superposition is still a pure state          |
| Mixed I/2    | ln(2) ≈ 0.693| Maximum entropy — no information about state |

**For a Bell pair (2-qubit):** The reduced density matrix of one qubit has entropy 1 ebit = ln(2) — maximum possible for a qubit. This is the standard measure of [[entanglement]].

## TestHadamardWork — Work as Energy Change

A gate performing work on a quantum system:

```go
// Initial: |0⟩, U₀ = ⟨0|H|0⟩ = -1
res0 := e.RunProgram(p)
rho0 := smath.ToDensityMatrix(res0.GetProbability())
u0 := real(smath.ExpectationValue(rho0, hamiltonian))  // -1.0

// Apply Hadamard: |0⟩ → |+⟩
p.AddStep(core.NewStep(core.NewHadamard(0)))
res1 := e.RunProgram(p)
rho1 := smath.ToDensityMatrix(res1.GetProbability())
u1 := real(smath.ExpectationValue(rho1, hamiltonian))  // 0.0

workDone := u1 - u0  // = 1.0 Joule (normalized)
```

**Physical interpretation:** The Hadamard gate performed 1 unit of work, rotating |0⟩ (low energy in magnetic field Hamiltonian) to |+⟩ (zero energy, equal superposition).

## TestContinuousWork — Time Evolution

Time evolution under a Hamiltonian H for time t: U(t) = e^{−iHt}

```go
// Hamiltonian Hx = σx (drives |0⟩ → |1⟩)
hx := smath.NewMatrix(2, 2)
hx.Data = []complex128{0, 1, 1, 0}

// Observable Hz = σz (what we measure as energy)
hz := smath.NewMatrix(2, 2)
hz.Data = []complex128{1, 0, 0, -1}

// Evolve for t = π/4
evolution := core.NewTimeEvolution(0, hx, 3.14159/4.0)
```

At t = π/4, |0⟩ evolves to cos(π/4)|0⟩ − i·sin(π/4)|1⟩.
Expected ⟨Hz⟩ = cos²(π/4) − sin²(π/4) = 0.5 − 0.5 = 0.

`NewTimeEvolution` computes U = math.Exp(−iHt) (matrix exponential) and wraps it as an Oracle gate. See [[composite-gates]].

## TestEntropy — Pure vs Mixed States

```go
// Pure state |0⟩: ρ = [[1,0],[0,0]]
rhoPure := smath.NewMatrix(2, 2)
rhoPure.Data = []complex128{1, 0, 0, 0}
s0 := smath.VonNeumannEntropy(rhoPure)  // = 0

// Maximally mixed state I/2: ρ = [[0.5,0],[0,0.5]]
rhoMixed := smath.NewMatrix(2, 2)
rhoMixed.Data = []complex128{0.5, 0, 0, 0.5}
s1 := smath.VonNeumannEntropy(rhoMixed)  // = ln(2) ≈ 0.6931
```

## The Engine Cycle

`core.NewEngineProgram()` models a quantum thermodynamic cycle:
1. Rx(π/4) — isothermal expansion (rotation under Hamiltonian)
2. H — adiabatic expansion (Bloch sphere rotation)
3. Ry(π/4) — isothermal compression

This maps loosely to a quantum analogue of a Carnot cycle, though the mapping to thermodynamic quantities (temperature, heat) requires more infrastructure than quantum-go currently provides.

## math Package API

```go
smath.ToDensityMatrix(ampVec []complex128) math.Matrix     // |ψ⟩⟨ψ|
smath.ExpectationValue(rho, H math.Matrix) complex128      // Tr(ρH)
smath.VonNeumannEntropy(rho math.Matrix) float64           // -Tr(ρ log ρ)
```

These are in `github.com/stephen-mcelhose/quantum-go/math` alongside the core matrix operations from [[quantum-linear-algebra]].

## Key Points

- `GetProbability()` returns complex amplitudes — use `ToDensityMatrix` to compute ρ from them.
- Von Neumann entropy = 0 for pure states regardless of superposition; only mixed states have positive entropy.
- Work = ΔU = ⟨H⟩_final − ⟨H⟩_initial — a gate performs work when it changes the expectation value of the Hamiltonian.
- `NewTimeEvolution(idx, H, t)` computes e^{−iHt} via matrix exponential — used when the circuit is described by a Hamiltonian, not gates.
- The density matrix formalism subsumes the state vector formalism: for a pure state, ρ = |ψ⟩⟨ψ|.

## Sources

- `examples/thermodynamics/engine_cycle_test.go`
