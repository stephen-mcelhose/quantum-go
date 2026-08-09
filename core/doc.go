/*
Package core defines the fundamental domain models for quantum computing.

This package provides the building blocks for constructing quantum circuits:
gates, steps, programs, and execution results. It is designed to be independent
of any specific execution environment.

# Core Concepts

Gate: A quantum gate is an operation that transforms qubit states. Gates implement
the Gate interface and provide their matrix representation along with metadata about
which qubits they affect.

Step: A step contains one or more gates that can be applied in parallel. All gates
in a step must operate on disjoint sets of qubits.

Program: A program is a sequence of steps that defines a complete quantum circuit.
It specifies the number of qubits and their initial states.

Result: An interface for execution outcomes. It provides access to the state vector
and qubit measurement probabilities.

# Result Types

CompactResult: The default implementation that stores only the final state vector.

InstrumentedResult: An implementation that stores intermediate states (state vectors
after each step). This is useful for debugging and analyzing quantum processes over time.

# Gate Types

Single-Qubit Gates:
  - Hadamard: Creates superposition (H)
  - Pauli gates: X (bit flip), Y (bit+phase flip), Z (phase flip)
  - Rotation: Phase rotation parameterized by angle
  - Identity: No-op gate
  - Measurement: Collapses quantum state

Multi-Qubit Gates:
  - CNOT: Controlled-NOT, flips target if control is |1⟩
  - CZ: Controlled-Z, phase flip if both qubits are |1⟩
  - SWAP: Exchanges states of two qubits
  - CR: Controlled rotation
  - Toffoli: CCNOT, flips target if both controls are |1⟩

Specialized Gates:
  - TimeEvolution: Unitary evolution exp(-iHt) based on a Hamiltonian matrix.
  - Oracle: Applies an arbitrary unitary matrix.

Composite Gates (BlockGate):
  - Fourier: Quantum Fourier Transform
  - Add: Quantum addition using QFT

# Example Usage

	// Create a 3-qubit GHZ state
	program := NewProgram(3)
	program.AddStep(NewStep(NewHadamard(0)))
	program.AddStep(NewStep(NewCnot(0, 1)))
	program.AddStep(NewStep(NewCnot(1, 2)))

	// The program can now be executed by an ExecutionEnvironment

# Block Gates

Complex multi-step operations can be encapsulated as blocks and used as single gates.
This allows operations like QFT to be reused and inverted efficiently:

	// Create a QFT block
	fourier := NewFourier(3, 0)  // 3-qubit QFT starting at qubit 0
	program.AddStep(NewStep(fourier))

	// Use inverse QFT
	invFourier := NewFourier(3, 0)
	invFourier.SetInverse(true)
	program.AddStep(NewStep(invFourier))

# State Vector Computation

The quantum state is represented as a complex vector of length 2^n for n qubits.
Each element represents the probability amplitude for a basis state:

	|ψ⟩ = α₀|000⟩ + α₁|001⟩ + α₂|010⟩ + ... + α₇|111⟩  (for 3 qubits)

Individual qubit probabilities are calculated by summing the squared magnitudes
of amplitudes where that qubit's bit is 1.
*/
package core
