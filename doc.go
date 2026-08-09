// Package quantum_go is a quantum computing simulator for Go.
//
// quantum-go is a Go port of the Java quantum simulator from https://github.com/redfx-quantum/strange,
// providing high-performance quantum circuit simulation on classical hardware.
//
// # Architecture
//
// The simulator consists of three main packages:
//
//   - core: Defines quantum circuit components (gates, steps, programs)
//   - local: Implements the simulation engine using state vector methods
//   - math: Provides complex matrix operations for quantum gates
//
// # Basic Usage
//
// To create and run a simple quantum circuit:
//
//	import (
//	    "github.com/stephen-mcelhose/quantum-go/core"
//	    "github.com/stephen-mcelhose/quantum-go/local"
//	)
//
//	// Create a 2-qubit Bell state circuit
//	program := core.NewProgram(2)
//	program.AddStep(core.NewStep(core.NewHadamard(0)))
//	program.AddStep(core.NewStep(core.NewCnot(0, 1)))
//
//	// Run the simulation
//	env := local.NewSimpleExecutionEnvironment()
//	result := env.RunProgram(program)
//
//	// Access results
//	qubits := result.GetQubits()
//	for i, q := range qubits {
//	    fmt.Printf("Qubit %d: probability=%.3f, measured=%d\n",
//	        i, q.Probability, q.Measure())
//	}
//
// # Quantum Gates
//
// The simulator provides standard quantum gates:
//
//   - Fundamental: Hadamard, X, Y, Z, Identity, Measurement
//   - Rotations: Rx, Ry, Rz, Universal (U), PhaseShift
//   - Clifford/Phase: S, T, V (SX)
//   - Two-qubit: CNOT, CZ, SWAP, CR (controlled phase shift)
//   - Three-qubit: Toffoli (CCNOT)
//   - Composite: Fourier (QFT), Arithmetic (Add, AddInteger, MulModulus)
//
// # State Vector Representation
//
// For n qubits, the quantum state is represented as a complex-valued vector
// of length 2^n. Each element represents the probability amplitude for a
// basis state. The simulator optimizes gate application to avoid explicit
// matrix construction when possible.
//
// # Performance
//
// The simulator uses several optimizations:
//
//   - Flat array representation for matrices and state vectors
//   - Specialized implementations for common gates (CNOT, SWAP, etc.)
//   - Optional gate optimizations via the HasOptimization interface
//   - Block gates for composite operations (QFT, arithmetic)
//
// For more information, see the documentation in individual packages.
package quantum_go
