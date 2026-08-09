# Learning Quantum Arithmetic with quantum-go

Quantum Arithmetic performs mathematical operations on quantum registers in superposition. One of the most efficient ways to perform addition on a quantum computer is the **Draper Adder**, which uses the Quantum Fourier Transform (QFT).

---

## 1. The Challenge: Adding in Phase Space

In classical computers, adding numbers involves a series of logic gates (XOR, AND) that compute carries. On a quantum computer, we can perform addition much faster by moving into the \"phase space\" using the Fourier Transform.

- **Classical approach**: Ripple-carry adders ($O(N)$ gates).
- **Quantum approach (Draper)**: Convert the numbers into the Fourier basis, apply phase rotations, and transform back.

*   **Learn more about the Theory**: [Quantum Fourier Transform (Wikipedia)](https://en.wikipedia.org/wiki/Quantum_Fourier_transform)

---

## 2. Core Concepts

### Quantum Fourier Transform (QFT)
The QFT transforms a state from the computational basis ($|0\rangle, |1\rangle$) into the Fourier basis. In this basis, information is stored in the **phase** of the qubits rather than their bit values.

### Draper Adder
The Draper Adder works in three high-level steps:
1.  **Transform**: Apply QFT to the first register ($x$).
2.  **Add**: For every bit set in the second register ($y$), apply a conditional phase rotation to the first register.
3.  **Inverse Transform**: Apply the Inverse QFT to the first register to return to the computational basis.

> **Hilbert Space Context**: The Draper adder performs a rotation of the state vector in the Fourier-transformed Hilbert Space. Instead of manipulating bits directly, it uses the phase of the amplitudes to represent values, where addition becomes a simple phase accumulation.

#### The Phase Rotation
The core of the addition is the conditional rotation $R_k$:
$$R_k = \begin{pmatrix} 1 & 0 \\ 0 & e^{2\pi i / 2^k} \end{pmatrix}$$
This gate rotates the phase of the target qubit depending on the state of a control qubit from the second register.

### Modular Arithmetic
Quantum adders naturally perform **Modular Addition**. If the result of the addition exceeds the capacity of the register ($2^n$), it \"wraps around\" back to zero.

---

## 3. Walkthrough: Adding 2 + 3

In `go/examples/arithmetic/adder_test.go`, we add $x=2$ and $y=3$.

### Step 1: Initialize the Registers
We use 5 qubits total:
- Register $x$: qubits 0 and 1 (2-bit register, max value 3).
- Register $y$: qubits 2, 3, and 4 (3-bit register).

```go
p := core.NewProgram(5)

// Initialize x = 2 (|10> in binary)
p.AddStep(core.NewStep(core.NewX(1)))

// Initialize y = 3 (|011> in binary)
p.AddStep(core.NewStep(core.NewX(2), core.NewX(3)))
```

### Step 2: Apply the Add Gate
The `core.NewAdd` gate in quantum-go automatically handles the QFT, the phase additions, and the Inverse QFT.

```go
// Add(x_start, x_end, y_start, y_end)
// The result is stored in the first register (x)
adder := core.NewAdd(0, 1, 2, 4)
p.AddStep(core.NewStep(adder))
```

---

## 4. Interpreting Results

The final state vector shows a single peak:

```
Quantum Result (5 qubits):
|01101>: 1.0000
```

To read this result, we split the bitstring $|q_4 q_3 q_2 q_1 q_0\rangle$:
- **Qubits 4, 3, 2 (y register)**: `011` binary = 3.
- **Qubits 1, 0 (x register)**: `01` binary = 1.

**Why is x=1?**
We added $2+3=5$. However, register $x$ only has 2 bits ($2^2=4$). Thus, the result is:
$$5 \pmod 4 = 1$$

---

## 5. Self-Assessment

1. **What is the primary gate used to move into the basis where addition is easy?**
   - [ ] Hadamard
   - [ ] Quantum Fourier Transform (QFT)
   - [ ] Toffoli

2. **If you add 1 + 1 in a 1-bit register, what is the result?**
   - [ ] 2
   - [ ] 0 (due to modular wrap-around)

3. **Does the Draper Adder require an extra \"ancilla\" qubit for carries?**
   - [ ] Yes, like a classical adder.
   - [ ] No, it uses phase rotations to handle carry information.

---

## 6. Hands-on Exploration

### Run the Example
Execute the automated test to see the addition result:
```bash
go test -v ./go/examples/arithmetic/adder_test.go
```

### Experiment with the CLI
Use the `quantum-go` CLI tool to export and run the built-in quantum adder:

**Export the Adder circuit:**
```bash
./quantum-go export --circuit adder
```

**Run the addition:**
```bash
./quantum-go run --circuit adder
```

**Inspect the Adder gate count:**
```bash
./quantum-go inspect --circuit adder
```

---

## 7. References & Further Reading

1.  **Wikipedia: Quantum Fourier Transform** - [Mathematical definition and circuits](https://en.wikipedia.org/wiki/Quantum_Fourier_transform).
2.  **PennyLane: Draper Adder** - [Detailed tutorial on QFT-based arithmetic](https://pennylane.ai/qml/demos/tutorial_qft_arithmetics/).
3.  **ArXiv: Original Draper Paper (2000)** - [Addition on a Quantum Computer](https://arxiv.org/abs/quant-ph/0008033).
4.  **Qiskit Textbook: QFT** - [Interactive explanation of the transform used in the adder](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-algorithms/quantum-fourier-transform.ipynb).

**Next Step:** Try modifying `adder_test.go` to use 3 bits for register $x$. What is the result of 2 + 3 now?
