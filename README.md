# quantum-go - Go Implementation

A high-performance quantum circuit simulator written in Go, ported from the [Strange](https://github.com/redfx-quantum/strange) Java implementation.

## Overview

This is a complete quantum simulator that allows you to construct and execute quantum circuits on classical hardware. It uses state vector simulation to track the full quantum state and provides optimized gate implementations for performance.

## CLI Tool

The project includes a command-line tool named `quantum-go` for running and analyzing quantum circuits.

### Installation

```bash
go build -o quantum-go ./cmd/quantum-go/main.go
```

#### Verification Environment (Optional)

To use the Qiskit-based verification features, you must set up a Python environment with `qiskit-aer`:

```bash
# Create and activate a virtual environment
python3 -m venv venv
source venv/bin/activate

# Install required packages
pip install qiskit qiskit-aer
```

	### Usage
	
	```bash
	# List all built-in quantum circuits
	./quantum-go list circuits

	# List all supported quantum gates
	./quantum-go list gates

	# Run a built-in circuit (use 'quantum-go list circuits' to see options)
		./quantum-go run --circuit deutsch-jozsa

		# Run a built-in circuit with parameters
		./quantum-go run --circuit bernstein-vazirani -p s=1011
		./quantum-go run --circuit simon -p s=11

	# Run a custom 2-qubit circuit built via flags

./quantum-go run -n 2 -s "h q[0]" -s "cx q[0], q[1]"

# Run a QASM file
./quantum-go run circuit.qasm

# Analyze thermodynamic properties (Entropy, Energy)
./quantum-go analyze --circuit grover

# Inspect circuit metadata and gate count
./quantum-go inspect --circuit teleportation

# Export a built-in circuit to QASM
./quantum-go export --circuit bell > bell.qasm

# Verify circuit calculation against reference
./quantum-go verify --circuit bell --mode theoretical
./quantum-go verify --circuit qft --mode qiskit
```

## Package Structure

The simulator is organized into three main packages:

### [`core`](./core/)
Defines quantum circuit components:
- **Gates**: 
    - **Fundamental**: Hadamard, X, Y, Z, CNOT, CZ, Toffoli, Swap
    - **Rotations**: Rx, Ry, Rz, Universal (U), PhaseShift (U1)
    - **Clifford/Phase**: S, T, V (SX)
    - **Algorithmic**: QFT, Arithmetic (Add, MulMod)
- **Steps**: Collections of gates that execute in parallel
- **Programs**: Complete quantum circuits with multiple steps
- **Results**: Execution outcomes with state vectors and measurement probabilities

### [`local`](./local/)
Local quantum simulator implementation:
- **Engine**: State vector simulation on classical hardware
- **Computations**: Optimized gate application algorithms
- Supports up to ~30 qubits on typical hardware

### [`math`](./math/)
Complex matrix operations and thermodynamics:
- Matrix multiplication, tensor products, conjugate transpose
- Density matrix formalism and partial trace
- von Neumann entropy and observable expectation values
- Common quantum computing constants

## Installation

```bash
go get github.com/stephen-mcelhose/quantum-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/stephen-mcelhose/quantum-go/core"
    "github.com/stephen-mcelhose/quantum-go/local"
)

func main() {
    // Create a Bell state: (|00⟩ + |11⟩)/√2
    program := core.NewProgram(2)
    program.AddStep(core.NewStep(core.NewHadamard(0)))
    program.AddStep(core.NewStep(core.NewCnot(0, 1)))
    
    // Execute the circuit
    env := local.NewSimpleExecutionEnvironment()
    result := env.RunProgram(program)
    
    // Display results
    qubits := result.GetQubits()
    for i, q := range qubits {
        fmt.Printf("Qubit %d: P(|1⟩) = %.3f\n", i, q.Probability)
    }
}
```

## Available Gates

For a complete and up-to-date list of all supported quantum gates, use the CLI:

```bash
./quantum-go list gates
```

The simulator supports a wide range of gates, including:
- **Fundamental**: Hadamard, Pauli-X/Y/Z, Identity
- **Multi-Qubit**: CNOT, CZ, SWAP, Toffoli, Fredkin
- **Rotations**: Rx, Ry, Rz, PhaseShift, Universal
- **Phase**: S, T, V (SX), and their inverses
- **Special**: Measurement

## Documentation

Comprehensive godoc documentation is available for all packages.

### Learning Aids (Pedagogical Guides)

These guides bridge the gap between quantum theory and the Go implementation, providing a clear "why" and "how" for each example.

- [**Quantum Fundamentals**](./examples/fundamentals/fundamentals.md)
- [**Entanglement (Bell & GHZ)**](./examples/entanglement/entanglement.md)
- [**Shor's Algorithm**](./examples/algorithms/shor.md)
- [**Simon's Algorithm**](./examples/algorithms/simon.md)
- [**Grover's Algorithm**](./examples/algorithms/grover.md)
- [**Quantum Error Correction**](./examples/algorithms/error_correction.md)
- [**Quantum Teleportation**](./examples/networking/teleportation.md)
- [**Quantum Key Distribution (BB84)**](./examples/security/qkd.md)
- [**Quantum Arithmetic (Adder)**](./examples/arithmetic/arithmetic.md)

### Interoperability and Verification

`quantum-go` is verified against industry standards to ensure mathematical correctness.
- **OpenQASM 2.0 Support**: Export programs using `p.ToQASM()` for use in external simulators.
- **Qiskit Compatible**: Uses Little-Endian qubit ordering, matching Qiskit's conventions.
- **Verification Suite**: Core gates and complex states are benchmarked against theoretical ground truth.

See [Verification and Interoperability](./docs/verification.md) for more details.

```bash
# View package documentation
go doc github.com/stephen-mcelhose/quantum-go/core
go doc github.com/stephen-mcelhose/quantum-go/local
go doc github.com/stephen-mcelhose/quantum-go/math

# View specific type documentation
go doc github.com/stephen-mcelhose/quantum-go/core.Gate
go doc github.com/stephen-mcelhose/quantum-go/core.Program
go doc github.com/stephen-mcelhose/quantum-go/local.SimpleExecutionEnvironment

# View function documentation
go doc github.com/stephen-mcelhose/quantum-go/math.Tensor
```

### Documentation Files
- [`doc.go`](./doc.go) - Main package documentation with usage examples
- [`core/doc.go`](./core/doc.go) - Core quantum circuit components
- [`local/doc.go`](./local/doc.go) - Local simulator implementation details
- [`math/doc.go`](./math/doc.go) - Matrix operations and linear algebra

## Examples

### GHZ State (3-qubit entanglement)

```go
// Create GHZ state: (|000⟩ + |111⟩)/√2
program := core.NewProgram(3)
program.AddStep(core.NewStep(core.NewHadamard(0)))
program.AddStep(core.NewStep(core.NewCnot(0, 1)))
program.AddStep(core.NewStep(core.NewCnot(1, 2)))
```

### Quantum Fourier Transform

```go
// Apply 3-qubit QFT starting at qubit 0
program := core.NewProgram(3)
qft := core.NewFourier(3, 0)
program.AddStep(core.NewStep(qft))

// Apply inverse QFT
invQft := core.NewFourier(3, 0)
invQft.SetInverse(true)
program.AddStep(core.NewStep(invQft))
```

### Quantum Addition

```go
// Add two quantum registers
// x_register: qubits 0-1
// y_register: qubits 2-3 (will contain x+y after execution)
program := core.NewProgram(4)
adder := core.NewAdd(0, 1, 2, 3)
program.AddStep(core.NewStep(adder))
```

## Testing

The project includes comprehensive tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run fuzz tests
go test -fuzz=FuzzToffoli ./local
go test -fuzz=FuzzMatrixMul ./math
```

## Performance Characteristics

- **Memory**: O(2^n) for n qubits (exponential growth)
- **Gate Application**: O(2^n) per gate (must update all amplitudes)
- **Practical Limit**: ~30 qubits on typical hardware (requires ~16GB memory)

### Optimizations
- Flat array representation for cache efficiency
- Specialized implementations for common gates (CNOT, SWAP, etc.)
- Identity gates are skipped
- Block gates avoid exponential matrix construction
- Zero elements skipped in tensor products

## Project Structure

```
quantum-go/
├── doc.go                 # Main package documentation
├── go.mod                 # Go module definition
├── core/
│   ├── doc.go            # Core package documentation
│   ├── core.go           # Program, Step, Result types
│   ├── gates.go          # Gate implementations
│   └── core_test.go      # Unit tests
├── local/
│   ├── doc.go            # Local simulator documentation
│   ├── engine.go         # Execution environment
│   ├── computations.go   # State vector operations
│   ├── engine_test.go    # Integration tests
│   └── fuzz_test.go      # Fuzz tests
├── math/
│   ├── doc.go            # Math package documentation
│   ├── matrix.go         # Matrix operations
│   ├── matrix_test.go    # Unit tests
│   └── fuzz_test.go      # Fuzz tests
└── docs/
    ├── README.md          # Project documentation index
    ├── implementation-plan.md
    ├── fuzzy-testing.md
    └── packages.md
```

## Contributing

When contributing, please:
1. Add godoc comments to all exported types, functions, and methods
2. Include examples in package documentation where appropriate
3. Write unit tests for new functionality
4. Run `go fmt` before committing
5. Ensure `go vet` and `go test` pass

## References

- Original Java implementation: [redfx-quantum/strange](https://github.com/redfx-quantum/strange)
- Quantum Computing Primer: [IBM Quantum Learning](https://learning.quantum.ibm.com/)
- Go Documentation Guidelines: [Effective Go](https://golang.org/doc/effective_go)

## License

See the [LICENSE](./LICENSE) file for details. This project is a Go port of [redfx-quantum/strange](https://github.com/redfx-quantum/strange), licensed under BSD 3-Clause.
