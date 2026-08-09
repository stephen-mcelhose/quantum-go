# Verifying Strange-Go with IBM Quantum Platform

This guide explains how to use the [IBM Quantum Platform](https://quantum.ibm.com/) for cloud-based verification and visualization of `quantum-go` circuits.

## Verification Workflow

### 1. Export Circuit to OpenQASM
In Strange-Go, generate the QASM 2.0 string for your circuit:

```go
fmt.Println(program.ToQASM())
```

### 2. Use IBM Quantum Composer (Visual)
The Composer is useful for visual inspection of small circuits.
1.  Log in to [IBM Quantum](https://quantum.ibm.com/).
2.  Open the **Quantum Composer**.
3.  On the right-hand side, click the **"Code Editor"** tab.
4.  Paste your OpenQASM 2.0 string.
5.  Check the **"State Vector"** or **"Probabilities"** panels at the bottom.

> [!IMPORTANT]
> **Endianness Warning**: The IBM Quantum Composer visualization displays qubits in **Big-Endian** order (Qubit 0 is the Most Significant Bit) in the state labels (e.g., $|q_n...q_0\rangle$). Strange-Go and Qiskit use **Little-Endian**. Be careful to flip the labels mentally when comparing bitstrings.

### 3. Use IBM Quantum Lab (Jupyter)
For precise numerical comparison of larger circuits:
1.  Open **Quantum Lab**.
2.  Create a new Jupyter Notebook.
3.  Use the `qiskit` library to load your QASM and run it on an IBM simulator (e.g., `ibmq_qasm_simulator`).

```python
from qiskit import QuantumCircuit, execute, IBMQ

# Paste your QASM string
qasm_str = """..."""
circuit = QuantumCircuit.from_qasm_str(qasm_str)

# Run on IBM backend
# (Requires account setup: IBMQ.load_account())
backend = IBMQ.get_provider().get_backend('ibmq_qasm_simulator')
job = execute(circuit, backend)
# ... retrieve and compare results ...
```

## Why Verify with IBM?
Verifying with IBM provides the highest level of assurance, as it confirms your simulator's logic matches the hardware-backed standards of the industry leader. It is particularly useful for verifying **Gate Decomposition** of complex `BlockGates` (like QFT).
