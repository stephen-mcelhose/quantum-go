# Learning Entanglement with Strange-Go

Quantum Entanglement is a physical phenomenon where pairs or groups of particles are generated, interact, or share spatial proximity in such a way that the quantum state of each particle cannot be described independently of the state of the others.

In this guide, we'll explore the two most famous entangled states: the **Bell State** and the **GHZ State**.

---

## 1. The Challenge: Beyond Independence

In a classical system, two bits are independent. Flipping one bit has no effect on the other unless you explicitly connect them with a wire. In quantum mechanics, qubits can be \"entangled\" such that they share a single unified state, even if they are physically separated.

- **Classical approach**: Correlation through shared history or direct connection.
- **Quantum approach**: Instantaneous correlation of measurement outcomes across the entire entangled system.

*   **Learn more about the Theory**: [Quantum Entanglement (Wikipedia)](https://en.wikipedia.org/wiki/Quantum_entanglement)

---

## 2. Core Concepts

### Bell States (The EPR Pair)
A Bell state is a maximally entangled quantum state of two qubits. The most common one used in simulations is:

$$| \Phi^+ \rangle = \frac{1}{\sqrt{2}}(|00\rangle + |11\rangle)$$

This means if you measure the first qubit and find it to be `0`, you **know** with 100% certainty that the second qubit is also `0`.

### GHZ States
The Greenberger–Horne–Zeilinger (GHZ) state is a generalization of entanglement to three or more qubits. For three qubits, it is:

$$| GHZ \rangle = \frac{1}{\sqrt{2}}(|000\rangle + |111\rangle)$$

> **Hilbert Space Context**: Entanglement represents a state vector in the multi-qubit Hilbert Space that **cannot be factored** into the tensor product of individual qubit states (i.e., $| \psi \rangle \neq | \phi_1 \rangle \otimes | \phi_2 \rangle$). This non-separability is the mathematical definition of entanglement.

### Spooky Action at a Distance
Albert Einstein famously referred to entanglement as \"spooky action at a distance.\" While the correlation is instantaneous, entanglement **cannot** be used to communicate information faster than light.

---

## 3. Walkthrough: Creating Entanglement

### Example 1: The Bell State
In `go/examples/entanglement/bell_state_test.go`, we entangle two qubits.

1.  **Initialize**: Both qubits start at $|00\rangle$.
2.  **Superposition**: Apply Hadamard to `q0`. State becomes $\frac{1}{\sqrt{2}}(|00\rangle + |01\rangle)$.
3.  **Entangle**: Apply CNOT with `q0` as control and `q1` as target. State becomes $\frac{1}{\sqrt{2}}(|00\rangle + |11\rangle)$.

```go
p := core.NewProgram(2)
p.AddStep(core.NewStep(core.NewHadamard(0))) // Superposition
p.AddStep(core.NewStep(core.NewCnot(0, 1)))   // Entanglement
```

### Example 2: The GHZ State
In `go/examples/entanglement/ghz_state_test.go`, we extend the entanglement to a third qubit.

1.  Create a Bell pair between `q0` and `q1`.
2.  Apply CNOT with `q1` as control and `q2` as target.

```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(core.NewHadamard(0)))
p.AddStep(core.NewStep(core.NewCnot(0, 1)))
p.AddStep(core.NewStep(core.NewCnot(1, 2))) // Now all three are entangled
```

---

## 4. Interpreting Results

For both examples, the output shows two distinct peaks with 50% probability each.

**Bell State Output:**
```
Quantum Result (2 qubits):
|00>: 0.5000
|11>: 0.5000
```

**GHZ State Output:**
```
Quantum Result (3 qubits):
|000>: 0.5000
|111>: 0.5000
```

Notice that states like $|01\rangle$ or $|101\rangle$ have **zero** probability. The qubits are \"locked\" together; they must either all be `0` or all be `1`.

---

## 5. Self-Assessment

1. **If you have a Bell State and measure the first qubit as $|1\rangle$, what is the state of the second qubit?**
   - [ ] Random (50/50)
   - [ ] $|1\rangle$
   - [ ] $|0\rangle$

2. **Which gate is primarily responsible for creating the \"link\" (entanglement) between qubits?**
   - [ ] Hadamard
   - [ ] CNOT (Controlled-NOT)
   - [ ] Pauli-X

3. **Can you describe a 3-qubit GHZ state by describing each qubit individually?**
   - [ ] Yes
   - [ ] No, that's the definition of entanglement.

---

## 6. Hands-on Exploration

### Run the Example
Execute the automated tests to see entanglement in action:
```bash
go test -v ./go/examples/entanglement/...
```

### Experiment with the CLI
Use the `strange` CLI tool to export and analyze the built-in entangled circuits:

**Export the Bell State:**
```bash
./strange export --circuit bell
```

**Run the GHZ state simulation:**
```bash
./strange run --circuit ghz -n 3
```

**Analyze the entropy of the GHZ state (S > 0 indicates entanglement):**
```bash
./strange analyze --circuit ghz -n 3
```

---

## 7. References & Further Reading

1.  **Wikipedia: Bell State** - [The four maximally entangled states](https://en.wikipedia.org/wiki/Bell_state).
2.  **Qiskit Textbook: Entanglement** - [Multiple-qubit gates and entanglement](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-gates/multiple-qubits-entangled-states.ipynb).
3.  **Brilliant.org: Entanglement** - [Visual guide to quantum correlations](https://brilliant.org/wiki/quantum-entanglement/).
4.  **ArXiv: GHZ Original Paper (1989)** - [Going Beyond Bell's Theorem](https://arxiv.org/abs/0712.0921).

**Next Step:** Modify `ghz_state_test.go` to add a fourth qubit. How many CNOT gates do you need?
