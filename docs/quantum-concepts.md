# Quantum Concepts and Gate Reference

This document provides a deep dive into the quantum mechanics principles and gates implemented in the `quantum-go` simulator.

## 1. Quantum Bases

A **basis** is a set of reference states used to represent and measure the state of a qubit.

### Z Basis (Computational Basis)
The standard reference frame, matching classical bits.
- **States**: $|0\rangle$ and $|1\rangle$.
- **Bloch Sphere**: Vertical axis (Poles).
- **Measurement**: Standard `Measure()` operation returns bits in this basis.

### Phase Basis (X / Hadamard Basis)
States in perfect superposition, where information is stored in the **relative phase**.
- **States**: $|+\rangle$ and $|-\rangle$.
- **Mathematical Definition**:
    - $|+\rangle = \frac{|0\rangle + |1\rangle}{\sqrt{2}}$
    - $|-\rangle = \frac{|0\rangle - |1\rangle}{\sqrt{2}}$
- **Bloch Sphere**: Horizontal axis (Equator).

### Comparison Summary

| Feature | Z Basis (Computational) | Phase Basis (X / Hadamard) |
| :--- | :--- | :--- |
| **States** | $|0\rangle, |1\rangle$ | $|+\rangle, |-\rangle$ |
| **Bit Value** | Definite (0 or 1) | Indeterminate (50/50 probability) |
| **Info Source** | Probability Amplitude | Relative Phase ($+$ or $-$) |
| quantum-go Gate | $X$ (bit flip) | $Z$ (phase flip) |
| **Bloch Sphere** | Vertical Axis (Poles) | Horizontal Axis (Equator) |

---

## 2. Gate Reference

### Fundamental Single-Qubit Gates

| Gate | Symbol | Matrix ($2 \times 2$) | Wikipedia |
| :--- | :--- | :--- | :--- |
| **Identity** | $I$ | $\begin{bmatrix} 1 & 0 \\ 0 & 1 \end{bmatrix}$ | [Identity](https://en.wikipedia.org/wiki/Quantum_logic_gate#Identity_gate) |
| **Hadamard** | $H$ | $\frac{1}{\sqrt{2}} \begin{bmatrix} 1 & 1 \\ 1 & -1 \end{bmatrix}$ | [Hadamard](https://en.wikipedia.org/wiki/Quantum_logic_gate#Hadamard_gate) |
| **Pauli-X** | $X$ | $\begin{bmatrix} 0 & 1 \\ 1 & 0 \end{bmatrix}$ | [Pauli-X](https://en.wikipedia.org/wiki/Quantum_logic_gate#Pauli_gates) |
| **Pauli-Y** | $Y$ | $\begin{bmatrix} 0 & -i \\ i & 0 \end{bmatrix}$ | [Pauli-Y](https://en.wikipedia.org/wiki/Quantum_logic_gate#Pauli_gates) |
| **Pauli-Z** | $Z$ | $\begin{bmatrix} 1 & 0 \\ 0 & -1 \end{bmatrix}$ | [Pauli-Z](https://en.wikipedia.org/wiki/Quantum_logic_gate#Pauli_gates) |
| **S Gate** | $S$ | $\begin{bmatrix} 1 & 0 \\ 0 & i \end{bmatrix}$ | [S gate](https://en.wikipedia.org/wiki/Quantum_logic_gate#Phase_shift_gates) |
| **T Gate** | $T$ | $\begin{bmatrix} 1 & 0 \\ 0 & e^{i\pi/4} \end{bmatrix}$ | [T gate](https://en.wikipedia.org/wiki/Quantum_logic_gate#Phase_shift_gates) |
| **SX (V)** | $V$ | $\frac{1+i}{2}\begin{bmatrix} 1 & -i \\ -i & 1 \end{bmatrix}$ | [SX gate](https://en.wikipedia.org/wiki/Quantum_logic_gate#Square_root_of_not) |

### Parameterized Rotation Gates

| Gate | Symbol | Matrix ($2 \times 2$) | Description |
| :--- | :--- | :--- | :--- |
| **Rx** | $Rx(\theta)$ | $\begin{bmatrix} \cos(\theta/2) & -i\sin(\theta/2) \\ -i\sin(\theta/2) & \cos(\theta/2) \end{bmatrix}$ | Rotation around X axis. |
| **Ry** | $Ry(\theta)$ | $\begin{bmatrix} \cos(\theta/2) & -\sin(\theta/2) \\ \sin(\theta/2) & \cos(\theta/2) \end{bmatrix}$ | Rotation around Y axis. |
| **Rz** | $Rz(\theta)$ | $\begin{bmatrix} e^{-i\theta/2} & 0 \\ 0 & e^{i\theta/2} \end{bmatrix}$ | Rotation around Z axis. |
| **Universal**| $U(\theta, \phi, \lambda)$ | $\begin{bmatrix} \cos(\theta/2) & -e^{i\lambda}\sin(\theta/2) \\ e^{i\phi}\sin(\theta/2) & e^{i(\phi+\lambda)}\cos(\theta/2) \end{bmatrix}$ | Universal single-qubit gate ([Qiskit Docs](https://docs.quantum.ibm.com/api/qiskit/qiskit.circuit.library.UGate)). |
| **PhaseShift**| $PS(\theta)$ | $\begin{bmatrix} 1 & 0 \\ 0 & e^{i\theta} \end{bmatrix}$ | Arbitrary phase rotation on $|1\rangle$. |

### Two-Qubit and Controlled Gates

| Gate | Symbol | Description | Wikipedia |
| :--- | :--- | :--- | :--- |
| **CNOT** | $CX$ | Flips target if control is $|1\rangle$. | [CNOT](https://en.wikipedia.org/wiki/Quantum_logic_gate#Controlled_gates) |
| **Controlled-Z**| $CZ$ | Phase flip if both are $|1\rangle$. | [CZ](https://en.wikipedia.org/wiki/Quantum_logic_gate#Controlled_gates) |
| **SWAP** | $SWAP$ | Exchanges states of two qubits. | [SWAP](https://en.wikipedia.org/wiki/Quantum_logic_gate#Swap_gate) |

### Advanced and Multi-Qubit Gates

| Gate | Symbol | Description | Wikipedia |
| :--- | :--- | :--- | :--- |
| **Toffoli** | $CCNOT$ | Flips target if both controls are $|1\rangle$. | [Toffoli](https://en.wikipedia.org/wiki/Toffoli_gate) |
| **Rotation** | $R(\theta)$ | Phase rotation $e^{i\theta}$ on $|1\rangle$. | [Phase Shift](https://en.wikipedia.org/wiki/Quantum_logic_gate#Phase_shift_gates) |
| **Controlled-R**| $CR(\theta)$ | Controlled phase rotation. | [Phase Shift](https://en.wikipedia.org/wiki/Quantum_logic_gate#Phase_shift_gates) |

### Algorithmic Blocks

| Block | Key Components | Function | Wikipedia |
| :--- | :--- | :--- | :--- |
| **Fourier (QFT)** | $H$, $CR$ | Basis rotation (Z $\leftrightarrow$ Phase). | [QFT](https://en.wikipedia.org/wiki/Quantum_Fourier_transform) |
| **Adder (Add)** | $QFT$, $CR$ | Quantum addition: $|x, y\rangle \to |x, x+y\rangle$. | [Quantum Adder](https://en.wikipedia.org/wiki/Quantum_adder) |

---

## 3. Core Principles

### Superposition
The ability of a quantum system to be in multiple states at once. Applying a **Hadamard (H)** gate to a $|0\rangle$ state creates an equal superposition:
$$H|0\rangle = \frac{|0\rangle + |1\rangle}{\sqrt{2}}$$
When measured, this state collapses to $|0\rangle$ or $|1\rangle$ with 50% probability each.

### Entanglement
A phenomenon where the quantum states of two or more objects are linked such that one cannot be described independently of the others.
Creating a **Bell State** (the simplest entanglement) requires an $H$ gate followed by a $CNOT$:
1. $H$ on $q_0$: $\frac{|00\rangle + |10\rangle}{\sqrt{2}}$
2. $CNOT(q_0, q_1)$: $\frac{|00\rangle + |11\rangle}{\sqrt{2}}$
In this state, measuring one qubit immediately determines the state of the other, regardless of distance.
