# Verification and Interoperability

This document describes how the `quantum-go` simulator is verified for mathematical correctness and how it interacts with industry-standard quantum computing tools.

## Verification Strategy

Correctness in `quantum-go` is established by comparing the **State Vector** (the collection of complex amplitudes representing the system's state) against theoretical "ground truth" and external reference simulators.

### 1. State Vector Comparison
The automated test harness in `go/local/verification_test.go` executes quantum circuits and validates that the resulting amplitudes match expected values within a numerical tolerance of $10^{-6}$.

| Benchmark | Verification Method | Reference |
| :--- | :--- | :--- |
| **Bell State** | Theoretical Amplitudes | $\frac{1}{\sqrt{2}}(|00\rangle + |11\rangle)$ |
| **GHZ State** | Theoretical Amplitudes | $\frac{1}{\sqrt{2}}(|000\rangle + |111\rangle)$ |
| **Simon's Algorithm** | Qiskit Parity Check | $y \cdot s = 0 \pmod 2$ |
| **Error Correction** | Syndrome Identity | Logical state recovery after $X$ noise |
| **QFT** | Superposition Identity | Discrete Fourier Transform matrix |
| **Time Evolution** | Hamiltonian Identity | $e^{-i\theta\sigma} = \cos(\theta)I - i\sin(\theta)\sigma$ |
| **Fredkin (CSWAP)** | State Permutation | Controlled target swap |

### 2. Qubit Ordering (Endianness)
quantum-go uses **Little-Endian** qubit ordering.
- **Qubit 0** is the Least Significant Bit (LSB).
- **Qubit N-1** is the Most Significant Bit (MSB).

This matches the convention used by [Qiskit](https://qiskit.org/), making direct comparisons straightforward. Note that this is the inverse of some visualization tools like the IBM Quantum Composer, which may display the MSB on top.

## Standards and Interoperability

To facilitate comparison with external platforms, `quantum-go` supports the **OpenQASM 2.0** standard.

### OpenQASM 2.0 Exporter
Any `core.Program` can be exported to OpenQASM using the `ToQASM()` method. This allows you to verify quantum-go circuits on high-performance backends or hardware.

```go
p := core.NewProgram(2)
// ... construct circuit ...
fmt.Println(p.ToQASM())
```

### Verified Reference Technologies

The following standards and tools are used to verify the outputs of this simulator:

- **[OpenQASM 2.0 Specification](https://github.com/Qiskit/openqasm/tree/master/spec-v2.0)**: The standard for quantum circuit interchange.
- **[Qiskit Aer](https://qiskit.org/ecosystem/aer/)**: Used as the primary reference for state vector ordering and gate decomposition.
- **[IBM Quantum Platform](https://quantum.ibm.com/)**: Used for manual verification of circuit behavior and visualization.
- **[Quirk Simulator](https://algassert.com/quirk)**: An open-source, visual quantum simulator used for rapid verification of small circuits.

## Integrated Verification (CLI)

The `quantum-go` CLI provides a `verify` subcommand to validate results against theoretical truth or external references.

```bash
# Verify against theoretical math constants (built-in)
./quantum-go verify --circuit bell --mode theoretical

# Verify against Qiskit Aer (requires python3 and qiskit-aer)
./quantum-go verify --circuit ghz --qubits 3 --mode qiskit
```

## Modular Verification (Piping)

For more flexible workflows, you can pipe QASM output from `quantum-go` directly into the standalone `qiskit_verify.py` tool.

```bash
# Workflow: Export QASM -> Simulate in Qiskit -> Print State Vector
./quantum-go export --circuit bell | python3 ../qiskit_verify.py

# Workflow: Run quantum-go -> Save JSON -> Compare in Qiskit
./quantum-go run circuit.qasm --json > result.json
./quantum-go export circuit.qasm | python3 ../qiskit_verify.py --compare result.json
```

## Running Verification Tests

To execute the verification suite and ensure your local environment produces correct results:

```bash
cd go
go test -v ./local/verification_test.go
```
