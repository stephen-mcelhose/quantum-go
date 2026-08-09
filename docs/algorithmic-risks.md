# Algorithmic Scaling Risks

Quantum algorithms often promise exponential speedups, but simulating them on classical hardware involves exponential costs. This document details the scaling behavior and risks of the most complex algorithms implemented in quantum-go.

## 1. State Vector Memory Scaling

The primary constraint for the local simulator is memory. The quantum state vector grows as $2^n$, where $n$ is the number of qubits. Each amplitude is a `complex128` (16 bytes).

| Qubits ($n$) | States ($2^n$) | Memory Required |
| :--- | :--- | :--- |
| 10 | 1,024 | 16 KB |
| 20 | 1,048,576 | 16 MB |
| 25 | 33,554,432 | 512 MB |
| 30 | 1,073,741,824 | 16 GB |
| 34 | 17,179,869,184 | 256 GB |

**Risk**: Attempting to simulate more than 28-30 qubits on a standard laptop will likely result in an Out-Of-Memory (OOM) error.

## 2. Shor's Algorithm

Shor's algorithm for factoring an integer $N$ requires a significant number of qubits and circuit depth.
- **Qubit Requirement**: To factor an $L$-bit number, you typically need $L$ qubits for the input register and $2L$ qubits for the modular exponentiation register, plus auxiliary qubits.
- **Example**: Factoring $N=15$ ($L=4$) requires at least 12-15 qubits. Factoring $N=21$ ($L=5$) requires 15-18 qubits.
- **Depth**: The modular exponentiation involves many `Add` and `Mul` steps, leading to thousands of individual gates.

**Risk**: Simulating Shor's algorithm for $N > 35$ is extremely challenging on classical hardware due to both memory and time constraints.

## 3. Grover's Algorithm

Grover's search provides a quadratic speedup $O(\sqrt{N})$, where $N=2^n$ is the search space.
- **Iteration Count**: Finding a single item in a search space of $2^n$ requires approximately $\frac{\pi}{4} 2^{n/2}$ iterations of the Grover operator.
- **Oracle Complexity**: Each iteration requires an "Oracle" gate that identifies the correct item. If the oracle is complex (e.g., solving a 3-SAT problem), the circuit depth grows rapidly.

**Risk**: While the number of qubits for Grover's is often small, the high iteration count and oracle depth can lead to long simulation times.

## 4. Mitigation Strategies

To handle these risks during development and education:
1. **Limit Qubit Counts**: Standard examples in the `go/examples` directory are limited to $n < 15$ to ensure they run in seconds.
2. **Simplified Oracles**: Use predefined unitary matrices for Oracles rather than complex gate decompositions when the search logic itself is not the focus.
3. **Binary Result Filtering**: Use the `PrintBinary()` method to focus only on high-probability states, making it easier to interpret results from large state vectors.
