# Learning Quantum Teleportation with quantum-go

Quantum Teleportation is a protocol that allows the transfer of a quantum state over a distance using a classical communication channel and a pre-shared entangled pair.

Contrary to its name, it does **not** move matter. Instead, it moves the *information* (the exact superposition) from one qubit to another.

---

## 1. The Challenge: Moving "Fragile" States

Quantum states are incredibly fragile. Due to the **No-Cloning Theorem**, you cannot simply copy a quantum state. Furthermore, measuring a state collapses it, losing the original superposition.

- **Classical approach**: Copy the data and send it over a wire.
- **Quantum approach**: Destroy the state at the source (Alice) and reconstruct it at the destination (Bob) using entanglement as a "bridge."

*   **Learn more about the Theory**: [Quantum Teleportation (Wikipedia)](https://en.wikipedia.org/wiki/Quantum_teleportation)

---

## 2. Core Concepts

### Entanglement (The Bridge)
Teleportation requires a **Bell Pair** (or EPR pair) shared between Alice and Bob. This pair is in a maximally entangled state:

$$| \Phi^+ \rangle = \frac{1}{\sqrt{2}}(|00\rangle + |11\rangle)$$

In this state, neither qubit has a definite value, but they are perfectly correlated.

### Bell Measurement
Alice performs a joint measurement on the qubit she wants to teleport and her half of the entangled pair. This measurement \"projects\" her qubits into one of the four Bell states, effectively destroying the original state but \"encoding\" its information into Bob's qubit.

### Unitary Reconstruction (Corrections)
After Alice's measurement, Bob's qubit is in a state that is a specific rotation away from the original. Alice sends two classical bits (the results of her measurement) to Bob. Bob then applies the corresponding **Pauli gates** (X or Z) to \"fix\" his qubit and restore the original state.

> **Hilbert Space Context**: The correction step involves conditional rotations in Bob's subspace of the global Hilbert Space. By applying $X$ or $Z$, Bob aligns his state vector exactly with the original $|\psi\rangle$ state prepared by Alice.

#### The Correction Matrix
Depending on Alice's outcome, Bob applies one of four operations:
- **00**: $I = \begin{pmatrix} 1 & 0 \\ 0 & 1 \end{pmatrix}$ (No change)
- **01**: $X = \begin{pmatrix} 0 & 1 \\ 1 & 0 \end{pmatrix}$ (Bit flip)
- **10**: $Z = \begin{pmatrix} 1 & 0 \\ 0 & -1 \end{pmatrix}$ (Phase flip)
- **11**: $XZ = \begin{pmatrix} 0 & -1 \\ 1 & 0 \end{pmatrix}$ (Both)

---

## 3. Walkthrough: Teleporting $|+\rangle$

In `go/examples/networking/teleportation_test.go`, we teleport the state $|+\rangle = H|0\rangle$ from qubit 0 to qubit 2.

### Step 1: Initialize the System
We use 3 qubits:
- `q0`: The state to be teleported ($|\psi\rangle$).
- `q1`: Alice's half of the entangled pair.
- `q2`: Bob's half of the entangled pair (The Destination).

```go
p := core.NewProgram(3)
// Prepare |psi> = |+> on q0
p.AddStep(core.NewStep(core.NewHadamard(0)))
```

### Step 2: Create the Entangled \"Bridge\"
We put `q1` and `q2` into a Bell state.

```go
p.AddStep(core.NewStep(core.NewHadamard(1)))
p.AddStep(core.NewStep(core.NewCnot(1, 2)))
```

### Step 3: Alice's Bell Measurement
Alice interacts her state (`q0`) with her bridge qubit (`q1`).

```go
p.AddStep(core.NewStep(core.NewCnot(0, 1)))
p.AddStep(core.NewStep(core.NewHadamard(0)))
```

### Step 4: Bob's Corrections
In the simulator, we use Alice's qubits as controls for Bob's gates.

```go
// If q1 is 1, Bob applies X to q2
p.AddStep(core.NewStep(core.NewCnot(1, 2)))

// If q0 is 1, Bob applies Z to q2
p.AddStep(core.NewStep(core.NewCz(0, 2)))
```

---

## 4. Interpreting Results

The simulation output shows all 8 possible measurement combinations.

```
Quantum Result (3 qubits):
|000>: 0.1250
|001>: 0.1250
...
|111>: 0.1250
```

Since we teleported the $|+\rangle$ state (50% $|0\rangle$, 50% $|1\rangle$), Bob's qubit (`q2`) should also be $|+\rangle$. In kets $|q_2 q_1 q_0\rangle$, $q_2=0$ occurs in half the results, and $q_2=1$ in the other half. The equal distribution confirms that Bob reconstructed the $|+\rangle$ state successfully.

---

## 5. Self-Assessment

1. **Does Quantum Teleportation allow faster-than-light communication?**
   - [ ] Yes, because entanglement is instantaneous.
   - [ ] No, because Bob still needs Alice's classical bits to perform the correction.

2. **What happens to Alice's original qubit after teleportation?**
   - [ ] It remains in the original state.
   - [ ] It is destroyed/collapsed by the measurement.

3. **Why are the X and Z gates called \"corrections\"?**
   - [ ] They fix hardware errors.
   - [ ] They rotate Bob's qubit back to Alice's original state.

---

## 6. Hands-on Exploration

### Run the Example
Execute the automated test to see the state reconstruction:
```bash
go test -v ./go/examples/networking/teleportation_test.go
```

### Experiment with the CLI
Use the `quantum-go` CLI tool to export and run the teleportation circuit:

**Export to QASM:**
```bash
./quantum-go export --circuit teleportation
```

**Run the simulation:**
```bash
./quantum-go run --circuit teleportation
```

**Inspect the gates:**
```bash
./quantum-go inspect --circuit teleportation
```

---

## 7. References & Further Reading

1.  **Wikipedia: Quantum Teleportation** - [Comprehensive protocol overview](https://en.wikipedia.org/wiki/Quantum_teleportation).
2.  **Qiskit Textbook: Teleportation** - [Hands-on guide with circuit visualizations](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-algorithms/teleportation.ipynb).
3.  **Brilliant.org: Quantum Teleportation** - [Interactive explanation of the Bell measurement](https://brilliant.org/wiki/quantum-teleportation/).
4.  **ArXiv: Original Bennett et al. Paper (1993)** - [The foundational paper defining the protocol](https://arxiv.org/abs/quant-ph/0006004).

**Next Step:** Try modifying `teleportation_test.go` to teleport a different state, like $|1\rangle$. What do you expect the final probabilities to look like?
