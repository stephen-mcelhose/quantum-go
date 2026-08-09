# Quantum Error Correction (3-Qubit Bit-Flip Code)

Quantum computers are highly sensitive to noise. Unlike classical bits, quantum states cannot be copied (No-Cloning Theorem), so we cannot simply keep multiple copies of a qubit. Instead, we use entanglement to spread the information of one logical qubit across multiple physical qubits.

The **3-qubit bit-flip code** is the simplest error-correcting code. It can detect and correct a single bit-flip error ($X$ gate) on any one of the three qubits.

## Implementation in Strange

The "Strange" simulator includes a built-in demonstration of this code.

### Usage via CLI

```bash
# Encode bit 1, simulate an error, and correct it
./strange run --circuit error-correction -p bit=1
```

Even though the circuit simulates an $X$ error on one of the qubits, the final measurement will reliably show the original bit.

### The Circuit Logic

1.  **Encoding**: We map a logical $|0\rangle_L \to |000\rangle$ and $|1\rangle_L \to |111\rangle$.
    *   This is done using two CNOT gates: `CNOT(0, 1)` and `CNOT(0, 2)`.
2.  **Error Simulation**: An $X$ gate is applied to one of the qubits (e.g., `X(1)`).
3.  **Syndrome Extraction & Decoding**:
    *   We apply `CNOT(0, 1)` and `CNOT(0, 2)` again.
    *   If no error occurred, the state returns to $|bit, 0, 0\rangle$.
    *   If an error occurred, the auxiliary qubits (1 and 2) will store the "syndrome".
4.  **Correction**:
    *   A **Toffoli gate** (CCNOT) is used to flip qubit 0 if and only if both syndromes indicate an error that affects qubit 0 logic.
    *   In the simplified 3-qubit version, the Toffoli gate `CCNOT(1, 2, 0)` flips the data qubit back to its correct state based on the parity of the other two.

## Programmatic Example

```go
package main

import (
    "fmt"
    "github.com/stephen-mcelhose/quantum-go/core"
    "github.com/stephen-mcelhose/quantum-go/local"
)

func main() {
    // Encode bit '1'
    program := core.NewErrorCorrectionProgram(1)
    
    env := local.NewSimpleExecutionEnvironment()
    result := env.RunProgram(program)
    
    // The result will show |011>? 
    // Wait, let's check the binary:
    // Q0 is the data qubit. If it was 1 and corrected, Q0=1.
    // Q1 and Q2 are syndromes.
    result.PrintBinary()
}
```

## Limitations

The bit-flip code only protects against $X$ errors (bit flips). It does **not** protect against $Z$ errors (phase flips). To protect against both, a more complex code like the **Shor Code** (which uses 9 qubits) is required.
