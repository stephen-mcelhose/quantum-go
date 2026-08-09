/*
Package local provides a local quantum circuit simulator.

This package implements quantum state vector simulation on classical hardware,
allowing quantum programs defined with the core package to be executed and tested
without requiring actual quantum hardware.

# Simulation Method

The simulator uses the state vector method, where the quantum state is represented
as a complex-valued vector of length 2^n for n qubits. Each gate operation
transforms this vector according to the gate's matrix representation.

# Usage

	import (
		"github.com/stephen-mcelhose/quantum-go/core"
		"github.com/stephen-mcelhose/quantum-go/local"
	)

	// Create a quantum program
	program := core.NewProgram(2)
	program.AddStep(core.NewStep(core.NewHadamard(0)))
	program.AddStep(core.NewStep(core.NewCnot(0, 1)))

	// Execute locally
	env := local.NewSimpleExecutionEnvironment()
	result := env.RunProgram(program)

	// Access results
	qubits := result.GetQubits()
	fmt.Printf("Q0 probability: %.3f\n", qubits[0].Probability)
	fmt.Printf("Q1 probability: %.3f\n", qubits[1].Probability)

# Initialization

All qubits start in the |0⟩ state unless specified otherwise through the
Program.InitAlpha field. The initial state vector is [1, 0, 0, ..., 0].

# Gate Application

Gates are applied to the state vector in sequence. The simulator provides
optimized implementations for common gates to avoid explicit matrix multiplication:

Single-Qubit Gates: Uses a partitioned loop structure that applies the 2x2
gate matrix to pairs of amplitudes separated by 2^(qubit_index).

CNOT: Directly swaps amplitudes when the control bit is 1, avoiding the 4x4
matrix multiplication.

CZ: Applies a phase flip to amplitudes where both control and target bits are 1.

SWAP: Exchanges amplitudes that differ in the two qubit positions.

Toffoli: Flips target bit in amplitude indices where both control bits are 1.

For gates without specialized implementations, the full matrix is applied
(not recommended for large qubit counts).

# Measurements

After execution, qubit measurement probabilities are calculated by summing
the squared magnitudes of state vector elements. Measurement outcomes are
randomly determined based on these probabilities using rand.Float64().

# Performance Considerations

Memory: The state vector size grows exponentially with qubit count (2^n complex128 values).
For n=20 qubits, this requires ~32 MB. Practical simulations are limited to ~30 qubits
on typical hardware.

Gate Optimization: Identity gates are skipped. Gates implementing HasOptimization()
use custom state transformation logic instead of matrix multiplication.

Block Gates: Composite operations (QFT, Add) apply their constituent gates
sequentially through the GlobalStepExecutor, avoiding exponential-sized matrix construction.

# Thread Safety

The SimpleExecutionEnvironment is not thread-safe. Create separate instances
for concurrent program execution.
*/
package local
