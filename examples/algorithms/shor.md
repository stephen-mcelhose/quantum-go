# Learning Shor's Algorithm with quantum-go

Shor's Algorithm is arguably the most famous quantum algorithm. Why? Because it has the potential to break the RSA encryption that secures almost all modern digital communication.

In this guide, we'll walk through how Shor's Algorithm works and how it's implemented using the `quantum-go` library.

---

## 1. Introduction: The Challenge

Imagine you have a large number $N = 15$. You want to find its prime factors (3 and 5). For $N=15$, it's easy. But if $N$ is 2048 bits long, a classical supercomputer would take billions of years to factor it.

This computational difficulty is the bedrock of **RSA Encryption**. If you can factor $N$, you can derive the private key.

### Quantum Advantage: Exponential Speedup
- **Classical Complexity**: The best known classical algorithm (General Number Field Sieve) runs in $O(e^{1.9 (\log N)^{1/3} (\log \log N)^{2/3}})$, which is sub-exponential but still impractical for large $N$.
- **Quantum Complexity**: Shor's Algorithm runs in **$O((\log N)^3)$**, which is polynomial. This represents an exponential speedup.

*   **Learn more about Factoring:** [Prime Factorization (GeeksforGeeks)](https://www.geeksforgeeks.org/prime-factorization-of-a-big-number/)

---

## 2. Core Concepts

### Period Finding
Shor's Algorithm turns the **Factoring Problem** into a **Period Finding Problem**. 
1. Pick a random number $a$ such that $1 < a < N$.
2. Consider the function $f(x) = a^x \pmod N$.
3. This function is **periodic**. There exists some $r$ (the period) such that $a^r \equiv 1 \pmod N$.
4. Once you find $r$, you can calculate the factors of $N$ using $\gcd(a^{r/2} \pm 1, N)$.

### Quantum Parallelism
A quantum computer creates a **superposition** of all possible inputs $x$ simultaneously. By applying the modular exponentiation operation once to this superposition, it effectively computes the function's output for all inputs at the same time.

### The Interference Engine (QFT)
The **Quantum Fourier Transform (QFT)** acts as a "Frequency Analyzer." It manipulates the quantum amplitudes such that the probabilities for states that correspond to the period $r$ reinforce each other (**Constructive Interference**), while others cancel out.

> **Hilbert Space Context**: Shor's algorithm operates by rotating the state vector in a high-dimensional complex vector space (Hilbert Space). The QFT aligns the state vector with specific basis vectors that correspond to the period of the function.

---

## 3. Walkthrough: Factoring $2^x \pmod 7$

In `go/examples/algorithms/shor_test.go`, we simulate finding the period of $2^x \pmod 7$.

### Step 1: Initialization
We create a **Precision Register** (for $x$) and a **Result Register** (for $a^x \pmod N$).

```go
length := 3 // ceil(log2(7))
offset := 3 // precision bits
p := core.NewProgram(2*length + 1 + offset)

// Put precision register into superposition
for i := 0; i < offset; i++ {
    p.AddStep(core.NewStep(core.NewHadamard(i)))
}
```

### Step 2: Modular Exponentiation
We calculate $a^x \pmod N$ for all $x$ in the superposition using the `MulModulus` gate.

```go
mul := core.NewMulModulus(offset, offset+length-1, m, mod)
cbg := core.NewControlledBlockGate(mul, i)
p.AddStep(core.NewStep(cbg))
```

### Step 3: Extracting the Period (Inverse QFT)
We use the **Inverse QFT** to shift the probability peaks to values related to $1/r$.

```go
invQFT := core.NewFourier(offset, 0)
invQFT.SetInverse(true)
p.AddStep(core.NewStep(invQFT))
```

---

## 4. Interpreting Results

When you run the simulation, you get measurement results like this:

```
Measurement Results (Binary):
|0000001000>: 0.1153
|0000001010>: 0.0981
|0000001110>: 0.0731
```

**How to read this:**
- The **first 3 bits** (on the right) are the precision register.
- In this case, the period of $2^x \pmod 7$ is $r=3$.
- Peaks appear at multiples of $2^{precision} / r$. With 3 bits of precision ($2^3 = 8$), peaks appear near $0$, $8/3 \approx 2.6$ (binary `010`), and $16/3 \approx 5.3$ (binary `110`).

---

## 5. Self-Assessment

1. **What is the complexity of Shor's Algorithm?**
   - [ ] $O(N)$
   - [ ] $O((\log N)^3)$
   - [ ] $O(e^N)$

2. **What mathematical problem is Shor's algorithm actually solving to find factors?**
   - [ ] Discrete Logarithm
   - [ ] Period Finding
   - [ ] Sorting

3. **Why do we need the Inverse QFT?**
   - [ ] To perform multiplication.
   - [ ] To act as a frequency analyzer and extract the period from the amplitudes.

---

## 6. Hands-on Exploration

### Run the Example
Execute the automated test to see the period-finding in action:
```bash
go test -v ./go/examples/algorithms/shor_test.go
```

### Experiment with the CLI
Use the `quantum-go` CLI tool to export and run the Shor circuit:

**Export to QASM:**
```bash
./quantum-go export --circuit shor
```

**Run the simulation:**
```bash
./quantum-go run --circuit shor
```

**Analyze the circuit:**
```bash
./quantum-go inspect --circuit shor
```

---

## 7. References & Further Reading

1.  **Wikipedia: Shor's Algorithm** - [Comprehensive overview of the algorithm and its math](https://en.wikipedia.org/wiki/Shor%27s_algorithm).
2.  **Qiskit Textbook: Quantum Fourier Transform** - [Developer-focused notebook on QFT](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-algorithms/quantum-fourier-transform.ipynb).
3.  **Brilliant.org: Quantum Computing** - [Overview of important quantum algorithms including Shor's](https://brilliant.org/wiki/quantum-computing/#important-algorithms).
4.  **ArXiv: Shor's Original Paper (1994)** - [The primary source for the algorithm](https://arxiv.org/abs/quant-ph/9508027).
5.  **Qiskit Textbook: Shor's Algorithm** - [Notebook-based implementation and theory](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-algorithms/shor.ipynb).
6.  **Quirk Simulator: Shor's Flow** - [Visualize the constructive interference in real-time](https://algassert.com/quirk#circuit=%7B%22cols%22%3A%5B%5B%22H%22%2C%22H%22%2C%22H%22%5D%2C%5B%22X%22%2C1%2C1%2C%22X%22%5D%5D%7D).

**Next Step:** Try modifying `shor_test.go` to find the period of $3^x \pmod 10$. How many qubits will you need?
