# Learning Grover's Algorithm with Strange-Go

Grover's Algorithm is a fundamental quantum algorithm that provides a **quadratic speedup** for unstructured search problems. If you're looking for a \"needle in a haystack,\" Grover's is your quantum magnet.

In this guide, we'll explore the core concepts of the algorithm and walk through the 2-qubit implementation in `quantum-go`.

---

## 1. Introduction: The Challenge

Imagine you have a database of $N$ unsorted items and you're looking for one specific entry. 

- **Classical approach:** You have to check each item one by one. On average, it takes $N/2$ checks, and in the worst case, $N$ checks. This is $O(N)$ complexity.
- **Quantum approach (Grover's):** You can find the item in roughly $\sqrt{N}$ checks. This is $O(\sqrt{N})$ complexity.

For a database of 1 million items, a classical computer might need 1 million checks, while Grover's needs only about 1,000.

*   **Learn more about Search Complexity:** [Grover's Algorithm (Wikipedia)](https://en.wikipedia.org/wiki/Grover%27s_algorithm)
* **Classical Query Complexity:** [Query Complexity](https://en.wikipedia.org/wiki/Query_complexity)

---

## 2. Core Concepts

To understand how Grover's works, we need to define four key components:

### Superposition (The Haystack)
Before we can search, we must prepare our system so that all possible answers exist at once. By applying Hadamard gates to all qubits, we create a **Superposition** where every state has an equal probability amplitude.

### The Oracle (The Identifier)
The **Oracle** is a \"black box\" operation that can recognize the correct answer. It doesn't tell us *what* the answer is; it simply \"marks\" the target state when it sees it.

### Phase Flip (The Marker)
How does the Oracle mark the answer? It uses a **Phase Flip**. Mathematically, it multiplies the amplitude of the target state by -1, while leaving all other states unchanged. 

> **Hilbert Space Context**: A complex vector space with an inner product, providing the mathematical framework where quantum states are represented as vectors. [Learn more](https://en.wikipedia.org/wiki/Hilbert_space). In Grover's, the Oracle rotates the state vector in Hilbert space to point in the opposite direction for the target state.

### Diffusion Operator (The Amplifier)
The **Diffusion Operator** performs an **Inversion about the Mean**. By reflecting all amplitudes about their average value, it dramatically increases the amplitude of the state that was flipped by the oracle, while slightly decreasing all others.

---

## 3. Geometric Interpretation: Amplitude Amplification

Grover's Algorithm is best understood as **Amplitude Amplification**. It uses geometric rotation to move the quantum state vector closer to the target answer.

1.  **Start**: Equal superposition of all states.
2.  **Oracle Step**: Flip the phase of the target state.
3.  **Diffusion Step**: Reflect all amplitudes about their average value.

Because the target state's phase was flipped, its amplitude is now far away from the average. The reflection **causes** the target state's amplitude to grow significantly, while the other states' amplitudes shrink.

---

## 4. Walkthrough: 2-Qubit Search for $|11\rangle$

In `go/examples/algorithms/grover_test.go`, we search for the state $|11\rangle$ (index 3) out of 4 possible states.

### Step 1: Initialize Superposition
We use `Hadamard` gates to put both qubits into superposition.

```go
p := core.NewProgram(2)
p.AddStep(core.NewStep(core.NewHadamard(0), core.NewHadamard(1)))
```

### Step 2: The Oracle (Phase Flip for $|11\rangle$)
We create an Oracle gate with a matrix that only flips the phase of the last state ($|11\rangle$). This is represented as a diagonal matrix $\text{diag}(1, 1, 1, -1)$:

$$
\text{Oracle} = \begin{pmatrix}
1 & 0 & 0 & 0 \\
0 & 1 & 0 & 0 \\
0 & 0 & 1 & 0 \\
0 & 0 & 0 & -1
\end{pmatrix}
$$

```go
oracleMatrix := math.NewMatrix(4, 4)
oracleMatrix.Set(0, 0, 1)
oracleMatrix.Set(1, 1, 1)
oracleMatrix.Set(2, 2, 1)
oracleMatrix.Set(3, 3, -1) // Marks |11> with a Phase Flip
oracle := core.NewOracle(0, oracleMatrix)
p.AddStep(core.NewStep(oracle))
```

### Step 3: The Diffusion Operator (Inversion about Mean)
The **Diffusion Operator** completes the amplification. For 2 qubits, the matrix performs a reflection that boosts the \"marked\" state to 100% probability. The matrix for 2 qubits is:

$$
\text{Diffusion} = \begin{pmatrix}
-0.5 & 0.5 & 0.5 & 0.5 \\
0.5 & -0.5 & 0.5 & 0.5 \\
0.5 & 0.5 & -0.5 & 0.5 \\
0.5 & 0.5 & 0.5 & -0.5
\end{pmatrix}
$$

```go
diffMatrix := math.NewMatrix(4, 4)
for i := 0; i < 4; i++ {
    for j := 0; j < 4; j++ {
        if i == j {
            diffMatrix.Set(i, j, -0.5)
        } else {
            diffMatrix.Set(i, j, 0.5)
        }
    }
}
diffusion := core.NewOracle(0, diffMatrix)
p.AddStep(core.NewStep(diffusion))
```

---

## 5. Interpreting Results

After just one iteration, the probability of $|11\rangle$ is 1.0.

```
Quantum Result (2 qubits):
|11>: 1.0000
```

### The Risk of Overcooking (Over-rotation)
Because Grover's is a rotation in Hilbert space, if you run too many iterations, you will rotate **past** the correct answer. The probability of measuring the target will start to decrease. This is called **Overcooking**. 

- For 2 qubits, **1 iteration** is optimal.
- For $N$ states, the optimal number of iterations is approximately $\frac{\pi}{4}\sqrt{N}$.

---

## 6. Self-Assessment

Test your understanding:

1. **What is the complexity of Grover's Algorithm?**
   - [ ] $O(N)$
   - [ ] $O(\log N)$
   - [ ] $O(\sqrt{N})$

2. **What does the Oracle actually do to the target state?**
   - [ ] It measures the state.
   - [ ] It performs a **Phase Flip** (multiplies by -1).
   - [ ] It sets the amplitude of all other states to zero.

3. **What is the purpose of the Diffusion Operator?**
   - [ ] To put qubits into superposition.
   - [ ] To perform an **Inversion about the Mean** to amplify the marked state.
   - [ ] To measure the final result.

---

## 7. Hands-on Exploration

### Run the Example
Execute the automated test to see Grover's algorithm find the target state:
```bash
go test -v ./go/examples/algorithms/grover_test.go
```

### Experiment with the CLI
Use the `strange` CLI tool to export and run the built-in 2-qubit Grover circuit:

**Export to QASM:**
```bash
./strange export --circuit grover
```

**Run the simulation:**
```bash
./strange run --circuit grover
```

**Analyze the entanglement:**
```bash
./strange analyze --circuit grover
```

---

## 8. References & Further Reading

1.  **Wikipedia: Grover's Algorithm** - [Comprehensive mathematical and historical overview](https://en.wikipedia.org/wiki/Grover%27s_algorithm).
2.  **Qiskit Textbook: Grover's Algorithm** - [Detailed notebook implementation and visualization](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-algorithms/grover.ipynb).
3.  **Brilliant.org: Quantum Computing** - [Overview of important quantum algorithms including Grover's](https://brilliant.org/wiki/quantum-computing/#important-algorithms).
4.  **ArXiv: Grover's Original Paper (1996)** - [The primary source for the algorithm](https://arxiv.org/abs/quant-ph/9605043).
5.  **Wikipedia: Hilbert Space** - [The mathematical foundation of quantum states](https://en.wikipedia.org/wiki/Hilbert_space).

**Next Step:** Try modifying `grover_test.go` to search for $|10\rangle$. What changes do you need to make to the Oracle matrix?
