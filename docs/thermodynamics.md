# Quantum Thermodynamics in quantum-go

This document explains the thermodynamic capabilities of the `quantum-go` simulator, based on the research presented in [DOI 10.1126/sciadv.adw8462](https://www.science.org/doi/10.1126/sciadv.adw8462).

## Core Concepts

Quantum thermodynamics studies how energy, work, and heat are exchanged in quantum systems. Unlike classical systems, quantum systems can exist in superpositions and entanglements, which gives rise to unique thermodynamic phenomena like quantum heat engines and refrigerators.

### 1. Internal Energy and Hamiltonians

The internal energy $U$ of a quantum system is the expectation value of its Hamiltonian $H$:

$$U = \text{Tr}(\rho H)$$

In `quantum-go`, you define a Hamiltonian as a `math.Matrix` and use `math.ExpectationValue` to calculate energy.

```go
hz := math.NewMatrix(2, 2)
hz.Data = []complex128{1, 0, 0, -1} // Pauli-Z

energy := math.ExpectationValue(rho, hz)
```

### 2. Work and Unitary Evolution

In a closed quantum system, any change in energy is attributed to **Work ($W$)**. This is implemented via unitary gates:

- **Discrete Work**: Applying gates like `Hadamard` or `CNOT`.
- **Continuous Work**: Using the `TimeEvolution` gate, which evolves the state according to $U(t) = e^{-iHt}$.

```go
// Evolve a state under Hamiltonian H for time t
evolution := core.NewTimeEvolution(qubitIdx, H, t)
program.AddStep(core.NewStep(evolution))
```

### 3. Subsystems and Entanglement

A key aspect of the research is analyzing how a "system" qubit interacts with its "environment" qubits. This is done using the **Partial Trace**:

$$\rho_{sys} = \text{Tr}_{env}(\rho_{total})$$

```go
// Extract the state of qubit 0 from a 2-qubit system
rho_0 := math.PartialTrace(rho_total, 2, []int{1})
```

### 4. Entropy and Information

The von Neumann entropy $S$ measures the "disorder" or entanglement of a state:

$$S(\rho) = -\text{Tr}(\rho \ln \rho)$$

In quantum thermodynamics, **Mutual Information ($I$)** measures the total correlations (classical and quantum) between two subsystems:

$$I(A:B) = S(\rho_A) + S(\rho_B) - S(\rho_{AB})$$

And the **Relative Entropy ($D$)** (or Kullback-Leibler divergence) measures how much a state $\rho$ differs from a reference state $\sigma$ (e.g., a thermal state):

$$D(\rho || \sigma) = \text{Tr}(\rho \ln \rho) - \text{Tr}(\rho \ln \sigma)$$

```go
entropy := math.VonNeumannEntropy(rho_sub)
mutualInfo := math.MutualInformation(rho_total, 2, []int{0}, []int{1})
relativeEnt := math.RelativeEntropy(rho_sub, rho_thermal)
```

## CLI Analysis with `quantum-go analyze`

The `quantum-go analyze` command automates the extraction of thermodynamic observables from a quantum circuit.

### Mathematical Foundations

#### 1. Density Matrix Conversion
The simulator converts the final state vector $|\psi\rangle$ into a global density matrix:
$$\rho = |\psi\rangle\langle\psi|$$

#### 2. Subsystem Isolation (Partial Trace)
For each qubit $i$, the command calculates the reduced density matrix $\rho_i$ by tracing out the remaining $n-1$ qubits:
$$\rho_i = \text{Tr}_{\{j \neq i\}}(\rho)$$

#### 3. Entanglement Measurement (Entropy)
The von Neumann entropy $S$ is calculated for each qubit:
$$S(\rho_i) = -\text{Tr}(\rho_i \ln \rho_i)$$
A value of $S > 0$ indicates that the qubit is entangled with other qubits in the system.

#### 4. Energy Expectation Values
The command measures the internal energy $U$ in the computational (Z) basis using the Pauli-Z Hamiltonian $H = \sigma_z$:
$$U_{total} = \sum_{i} \text{Tr}(\rho_i \sigma_z)$$

## Example: Measuring the First Law

The generalized first law for a quantum process is $\Delta U = W + Q$. For the unitary processes simulated in `quantum-go`, $Q=0$, so $\Delta U = W$.

```go
// 1. Initial State Energy
res0 := engine.RunProgram(p)
u0 := real(math.ExpectationValue(math.ToDensityMatrix(res0.GetProbability()), H))

// 2. Perform Work
p.AddStep(core.NewStep(core.NewHadamard(0)))
res1 := engine.RunProgram(p)

// 3. Final State Energy
u1 := real(math.ExpectationValue(math.ToDensityMatrix(res1.GetProbability()), H))

// Result
workDone := u1 - u0
```

## Implementation Parity with Research

| Research Concept | `quantum-go` Implementation |
| :--- | :--- |
| Density Matrix $\rho$ | `math.ToDensityMatrix` |
| Observable $\langle A \rangle$ | `math.ExpectationValue` |
| Subsystem Analysis | `math.PartialTrace` |
| Entropy $S$ | `math.VonNeumannEntropy` |
| Mutual Information $I$ | `math.MutualInformation` |
| Relative Entropy $D$ | `math.RelativeEntropy` |
| Time Evolution $U(t)$ | `core.TimeEvolution` |
| Intermediate States | `core.InstrumentedResult` |

For a complete working example, see `go/examples/thermodynamics/engine_cycle_test.go`.
