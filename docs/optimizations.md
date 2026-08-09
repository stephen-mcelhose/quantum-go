# Simulation Optimizations in Strange-Go

This document explains the bitwise optimization strategy used in the 'Strange' Go simulator to achieve high-performance quantum circuit simulation.

## Core Strategy: Bit-Loop Optimization

Traditional quantum simulators often apply gates by constructing a full unitary matrix $U$ of size $2^n \times 2^n$ and performing matrix-vector multiplication $v' = Uv$. This operation has a time complexity of $O(2^{2n})$ (or $O(4^n)$), which is prohibitively expensive even for a moderate number of qubits.

Strange-Go avoids constructing large matrices whenever possible. Instead, it uses specialized loops that operate directly on the state vector by manipulating bit patterns in the indices.

### 1-Qubit Gate Optimization
For a single-qubit gate at index $i$, the simulator partitions the state vector into pairs of amplitudes $(v_{j}, v_{j+2^i})$ where the $i$-th bit is 0 and 1, respectively. It then applies the 2x2 gate matrix to each pair:
$$
\begin{bmatrix} v'_j \\ v'_{j+2^i} \end{bmatrix} = \begin{bmatrix} m_{00} & m_{01} \\ m_{10} & m_{11} \end{bmatrix} \begin{bmatrix} v_j \\ v_{j+2^i} \end{bmatrix}
$$
This reduces the complexity from $O(4^n)$ to $O(2^n)$.

### 2-Qubit Gate Optimization
Controlled gates (like CNOT, CZ) are implemented by checking if the control bit is 1 before performing an operation on the target qubit.
- **CNOT**: If `(index >> control) & 1 == 1`, swap `v[index]` and `v[index ^ (1 << target)]`.
- **CZ**: If `(index >> control) & 1 == 1` and `(index >> target) & 1 == 1`, multiply `v[index]` by -1.
- **SWAP**: Exchanges amplitudes where the control and target bits differ.

### 3-Qubit Gate Optimization
- **Toffoli (CCNOT)**: Flips the target bit if both control bits are 1. This is implemented using bitwise `AND` and `XOR`.

### Composite Gate (Block) Optimization
For gates composed of multiple steps (like QFT or Adder), Strange-Go uses an `ApplyOptimize` pattern. Instead of calculating the total unitary matrix for the entire block, the simulator executes each constituent step in sequence. This prevents the exponential memory and CPU cost of constructing large-dimensional matrices.

## Performance Comparison

| Operation | General Matrix Complexity | Optimized Complexity |
| :--- | :--- | :--- |
| Single Qubit Gate | $O(4^n)$ | $O(2^n)$ |
| Two Qubit Gate | $O(4^n)$ | $O(2^n)$ |
| $k$-Qubit Block | $O(4^n)$ | $O(Steps \times 2^n)$ |

By leveraging these bitwise optimizations, Strange-Go can simulate larger circuits on standard hardware than a naive matrix-based simulator.
