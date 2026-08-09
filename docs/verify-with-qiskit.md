# Verifying quantum-go with Qiskit Aer

This guide explains how to use [Qiskit Aer](https://qiskit.org/ecosystem/aer/) to verify the state vectors produced by `quantum-go`.

## Prerequisites

- Python 3.x
- Qiskit and Qiskit Aer: `pip install qiskit qiskit-aer`

## Verification Workflow

### 1. Export Circuit from quantum-go
Generate the OpenQASM 2.0 representation of your circuit in Go:

```go
qasm := program.ToQASM()
os.WriteFile("circuit.qasm", []byte(qasm), 0644)
```

### 2. Simulate in Qiskit Aer
Use the following Python script to load the QASM and retrieve the state vector:

```python
from qiskit import QuantumCircuit
from qiskit_aer import Aer
import numpy as np

# Load the QASM file
circuit = QuantumCircuit.from_qasm_file("circuit.qasm")

# Use StatevectorSimulator
backend = Aer.get_backend('statevector_simulator')
job = backend.run(circuit)
result = job.result()
statevector = result.get_statevector()

print("Qiskit Statevector:")
print(statevector.data)
```

### 3. Compare State Vectors
quantum-go and Qiskit both use **Little-Endian** qubit ordering, meaning the indices in the state vector match exactly.

- **quantum-go**: `result.GetProbability()` returns a `[]complex128`.
- **Qiskit**: `statevector.data` is a numpy array of complex numbers.

Compare the amplitudes at each index. They should match within a tolerance of $10^{-6}$.

## Handling Global Phase
If the amplitudes differ by a constant factor of $e^{i\theta}$ across the entire vector, it is likely a global phase difference. In quantum mechanics, states $|\psi\rangle$ and $e^{i\theta}|\psi\rangle$ are physically equivalent.
