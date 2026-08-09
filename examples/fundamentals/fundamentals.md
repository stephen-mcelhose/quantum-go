# Learning Quantum Fundamentals with Strange-Go

Quantum computing is built on a few foundational principles: **Superposition**, **Unitary Transformations**, and **Measurement**. In this guide, we'll explore these concepts using the basic gates provided by `quantum-go`.

---

## 1. The Challenge: Beyond Binary

In classical computing, a bit is either $0$ or $1$. This binary nature limits the types of problems that can be solved efficiently. Quantum computers use qubits, which can exist in multiple states at once, allowing for massive parallel processing of information.

- **Classical approach**: Bits ($0$ or $1$).
- **Quantum approach**: Qubits (Superposition of $|0\rangle$ and $|1\rangle$).

*   **Learn more about Qubits**: [Qubit (Wikipedia)](https://en.wikipedia.org/wiki/Qubit)

---

## 2. Core Concepts

### Superposition
Unlike a classical bit, a quantum bit (qubit) can exist in a linear combination of both states simultaneously. This is known as **Superposition**. 

### The Bloch Sphere
The state of a single qubit can be visualized as a point on the surface of a unit sphere, called the **Bloch Sphere**. The north pole represents $|0\rangle$, the south pole $|1\rangle$, and points on the equator represent equal superpositions (like $|+\rangle$).

### Unitary Gates
Operations on qubits are represented by **Unitary Matrices**. These operations are reversible and preserve the total probability (the state vector always has a length of 1).

> **Hilbert Space Context**: Single-qubit operations are rotations within a 2-dimensional Hilbert Space. Any unitary gate can be thought of as a rotation of the state vector $|\psi\rangle$ around an axis on the Bloch Sphere.

---

## 3. The Foundation: Single-Qubit Gates

### The Hadamard Gate (H)
The Hadamard gate is the most important single-qubit gate. It creates an equal superposition of $|0\rangle$ and $|1\rangle$.

$$H = \frac{1}{\sqrt{2}} \begin{pmatrix} 1 & 1 \\ 1 & -1 \end{pmatrix}$$

Applying $H$ to $|0\rangle$ results in $|+\rangle = \frac{|0\rangle + |1\rangle}{\sqrt{2}}$.

### Pauli Gates (X, Y, Z)
The Pauli gates represent rotations of $180^\circ$ ($\pi$ radians) around the X, Y, and Z axes of the Bloch sphere.

- **Pauli-X (NOT)**: Flips $|0\rangle \leftrightarrow |1\rangle$.
- **Pauli-Y**: Flips $|0\rangle \leftrightarrow i|1\rangle$.
- **Pauli-Z (Phase Flip)**: Leaves $|0\rangle$ unchanged and flips the phase of $|1\rangle$.

### Phase and Square Root Gates (S, T, V)
Beyond the fundamental gates, we have several gates that perform smaller rotations:
- **S (Phase)**: Rotation of $90^\circ$ around Z.
- **T (π/8)**: Rotation of $45^\circ$ around Z.
- **V (SX)**: Square root of the X gate.

### Parameterized Rotations (Rx, Ry, Rz)
For precise control, you can rotate by any angle $\theta$:
- **Rx(θ)**: Rotation around X axis.
- **Ry(θ)**: Rotation around Y axis.
- **Rz(θ)**: Rotation around Z axis.

Learn more in [Parameterized Rotations](./rotations.md).

---

## 4. Walkthrough: Basic Transformations

### Creating Superposition
In `go/examples/fundamentals/hadamard_test.go`, we create a 50/50 superposition.

```go
p := core.NewProgram(1)
p.AddStep(core.NewStep(core.NewHadamard(0)))
```

### Advanced Rotations
In `go/examples/fundamentals/rotations_test.go`, we explore precision positioning on the Bloch Sphere.

```go
p := core.NewProgram(1)
p.AddStep(core.NewStep(core.NewRy(math.Pi/4, 0)))
```

### Universal Gates
Learn how any gate can be represented using the [Universal U Gate](./universality.md).

### Multi-Qubit Superposition
In `go/examples/fundamentals/superposition_test.go`, we apply Hadamard to 3 qubits. This creates $2^3 = 8$ possible states, each with $1/8 = 0.125$ probability.

```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(
    core.NewHadamard(0),
    core.NewHadamard(1),
    core.NewHadamard(2),
))
```

### Manipulating Bits with Pauli-X
In `go/examples/fundamentals/pauli_test.go`, we use the X gate to flip bits.

```go
p := core.NewProgram(3)
p.AddStep(core.NewStep(
    core.NewX(0),
    core.NewY(1),
    core.NewZ(2),
))
```

---

## 5. Interpreting Results

**Superposition Result:**
```
Quantum Result (1 qubits):
|0>: 0.5000
|1>: 0.5000
```
This shows the qubit has a 50% chance of being measured as $0$ or $1$.

**Pauli-X Result (on 3 qubits):**
```
|011>: 1.0000
```
Since Strange-Go uses **little-endian** notation, $q_0$ is the rightmost bit. Applying $X$ to $q_0$ and $q_1$ results in binary `011` (decimal 3).

---

## 6. Self-Assessment

1. **What happens if you apply a Hadamard gate twice to the same qubit?**
   - [ ] It stays in superposition.
   - [ ] It returns to its original state (H is its own inverse).

2. **Which Pauli gate is equivalent to a classical NOT gate?**
   - [ ] Pauli-Z
   - [ ] Pauli-X

3. **In a 2-qubit system, how many possible states exist in an equal superposition?**
   - [ ] 2
   - [ ] 4 ($2^2$)

---

## 7. Hands-on Exploration

### Run the Example
Execute the automated tests to see basic gates in action:
```bash
go test -v ./go/examples/fundamentals/...
```

	### Experiment with the CLI
	Use the `strange` CLI tool to experiment with basic gates:

	**Discover supported gates:**
	```bash
	./strange list gates
	```

	**Create a superposition qubit:**

```bash
./strange run -n 1 -s "h q[0]"
```

**Flip a bit and inspect:**
```bash
./strange inspect -n 1 -s "x q[0]"
```

**Apply a rotation gate ($R_z(\pi/4)$):**
```bash
./strange run -n 1 -s "u1(0.785398) q[0]"
```

---

## 8. References & Further Reading

1.  **Wikipedia: Quantum Logic Gate** - [Full list of common gates and their matrices](https://en.wikipedia.org/wiki/Quantum_logic_gate).
2.  **Qiskit Textbook: Single Qubit Gates** - [Visual interactive guide using the Bloch Sphere](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-states/single-qubit-gates.ipynb).
3.  **Brilliant.org: Quantum Computing** - [Foundational concepts and interactive problems](https://brilliant.org/wiki/quantum-computing/).
4.  **ArXiv: Elementary Gates for Quantum Computation** - [Foundational gate theory (1995)](https://arxiv.org/abs/quant-ph/9503016).

**Next Step:** Modify `pauli_test.go`. What happens if you apply a `Z` gate to a qubit in state $|0\rangle$? Does the measured probability change?
