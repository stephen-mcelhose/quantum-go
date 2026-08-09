package core

import (
	"fmt"
	"math"
	"math/cmplx"
	smath "github.com/stephen-mcelhose/quantum-go/math"
)

// NewBellProgram creates a 2-qubit program that generates a Bell state (|00> + |11>) / sqrt(2).
func NewBellProgram() *Program {
	p := NewProgram(2)
	p.AddStep(NewStep(NewHadamard(0)))
	p.AddStep(NewStep(NewCnot(0, 1)))
	return p
}

// ExpectedBellState returns the theoretical state vector for a Bell state.
func ExpectedBellState() []complex128 {
	val := 1.0 / math.Sqrt(2)
	return []complex128{complex(val, 0), 0, 0, complex(val, 0)}
}

// NewGHZProgram creates an n-qubit program that generates a GHZ state (|0...0> + |1...1>) / sqrt(2).
func NewGHZProgram(n int) *Program {
	p := NewProgram(n)
	p.AddStep(NewStep(NewHadamard(0)))
	for i := 0; i < n-1; i++ {
		p.AddStep(NewStep(NewCnot(i, i+1)))
	}
	return p
}

// ExpectedGHZState returns the theoretical state vector for an n-qubit GHZ state.
func ExpectedGHZState(n int) []complex128 {
	size := 1 << n
	res := make([]complex128, size)
	val := 1.0 / math.Sqrt(2)
	res[0] = complex(val, 0)
	res[size-1] = complex(val, 0)
	return res
}

// NewQFTProgram creates an n-qubit program that performs the Quantum Fourier Transform.
func NewQFTProgram(n int) *Program {
	p := NewProgram(n)
	p.AddStep(NewStep(NewFourier(n, 0)))
	return p
}

// ExpectedQFTState returns the theoretical state vector for an n-qubit QFT on |0...0>.
// For |0...0>, the QFT result is an equal superposition of all states.
func ExpectedQFTState(n int) []complex128 {
	size := 1 << n
	res := make([]complex128, size)
	val := 1.0 / math.Sqrt(float64(size))
	for i := range res {
		res[i] = complex(val, 0)
	}
	return res
}

// NewTeleportationProgram creates a 3-qubit program demonstrating quantum teleportation.
// q0: State to teleport
// q1: Alice's half of Bell pair
// q2: Bob's half of Bell pair
func NewTeleportationProgram() *Program {
	p := NewProgram(3)
	// Prepare |psi> = |+>
	p.AddStep(NewStep(NewHadamard(0)))
	// Entangle q1 and q2
	p.AddStep(NewStep(NewHadamard(1)))
	p.AddStep(NewStep(NewCnot(1, 2)))
	// Alice Bell measurement
	p.AddStep(NewStep(NewCnot(0, 1)))
	p.AddStep(NewStep(NewHadamard(0)))
	// Bob corrections
	p.AddStep(NewStep(NewCnot(1, 2)))
	p.AddStep(NewStep(NewCz(0, 2)))
	return p
}

// NewSuperdenseCodingProgram creates a 2-qubit program demonstrating superdense coding.
// It encodes two classical bits (11) into a single qubit using an entangled pair.
func NewSuperdenseCodingProgram() *Program {
	p := NewProgram(2)
	// 1. Create Bell pair (shared between Alice and Bob)
	p.AddStep(NewStep(NewHadamard(0)))
	p.AddStep(NewStep(NewCnot(0, 1)))

	// 2. Alice encodes two bits (e.g., "11") on q0
	// For "11", apply X and Z
	p.AddStep(NewStep(NewX(0)))
	p.AddStep(NewStep(NewZ(0)))

	// 3. Bob decodes the bits by measuring in Bell basis
	p.AddStep(NewStep(NewCnot(0, 1)))
	p.AddStep(NewStep(NewHadamard(0)))
	return p
}

// NewGroverProgram creates a 2-qubit program searching for |11> using Grover's algorithm.
func NewGroverProgram() *Program {
	p := NewProgram(2)
	// Superposition
	p.AddStep(NewStep(NewHadamard(0), NewHadamard(1)))
	// Oracle (|11>)
	oracleMatrix := smath.NewMatrix(4, 4)
	oracleMatrix.Set(0, 0, 1)
	oracleMatrix.Set(1, 1, 1)
	oracleMatrix.Set(2, 2, 1)
	oracleMatrix.Set(3, 3, -1)
	p.AddStep(NewStep(NewOracle(0, oracleMatrix)))
	// Diffusion
	diffMatrix := smath.NewMatrix(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				diffMatrix.Set(i, j, -0.5)
			} else {
				diffMatrix.Set(i, j, 0.5)
			}
		}
	}
	p.AddStep(NewStep(NewOracle(0, diffMatrix)))
	return p
}

// NewShorProgram creates a program for Shor's algorithm (period finding).
// a^x mod N
func NewShorProgram(a, mod, precision int) *Program {
	length := 0
	tmp := mod
	for tmp > 0 {
		tmp >>= 1
		length++
	}
	p := NewProgram(2*length + 1 + precision)
	// Precision register into superposition
	for i := 0; i < precision; i++ {
		p.AddStep(NewStep(NewHadamard(i)))
	}
	// Modular exponentiation
	m := a
	for i := 0; i < precision; i++ {
		mul := NewMulModulus(precision, precision+length-1, m, mod)
		cbg := NewControlledBlockGate(mul, i)
		p.AddStep(NewStep(cbg))
		m = (m * m) % mod
	}
	// Inverse QFT
	invQFT := NewFourier(precision, 0)
	invQFT.SetInverse(true)
	p.AddStep(NewStep(invQFT))
	return p
}

// NewAdderProgram creates a program adding two numbers using Draper adder.
func NewAdderProgram(x0, x1, y0, y1, y2 int) *Program {
	p := NewProgram(5)
	// Initialize x
	if x1 == 1 {
		p.AddStep(NewStep(NewX(1)))
	}
	if x0 == 1 {
		p.AddStep(NewStep(NewX(0)))
	}
	// Initialize y
	if y0 == 1 {
		p.AddStep(NewStep(NewX(2)))
	}
	if y1 == 1 {
		p.AddStep(NewStep(NewX(3)))
	}
	if y2 == 1 {
		p.AddStep(NewStep(NewX(4)))
	}
	// Add
	p.AddStep(NewStep(NewAdd(0, 1, 2, 4)))
	return p
}

// NewDeutschJozsaProgram creates a program for the Deutsch-Jozsa algorithm.
func NewDeutschJozsaProgram(n int, balanced bool) *Program {
	p := NewProgram(n + 1)
	// Prepare auxiliary qubit in |1>
	p.AddStep(NewStep(NewX(n)))
	// Apply H to all qubits
	steps := make([]Gate, n+1)
	for i := 0; i <= n; i++ {
		steps[i] = NewHadamard(i)
	}
	p.AddStep(NewStep(steps...))
	// Oracle
	if balanced {
		p.AddStep(NewStep(NewBalancedOracle(n)))
	} else {
		p.AddStep(NewStep(NewConstantOracle(n, 0)))
	}
	// Apply H to input qubits
	inputSteps := make([]Gate, n)
	for i := 0; i < n; i++ {
		inputSteps[i] = NewHadamard(i)
	}
	p.AddStep(NewStep(inputSteps...))
	return p
}

// NewBernsteinVaziraniProgram creates a program for the Bernstein-Vazirani algorithm.
func NewBernsteinVaziraniProgram(s string) *Program {
	n := len(s)
	p := NewProgram(n + 1)
	// Prepare auxiliary qubit in |1>
	p.AddStep(NewStep(NewX(n)))
	// Apply H to all qubits
	steps := make([]Gate, n+1)
	for i := 0; i <= n; i++ {
		steps[i] = NewHadamard(i)
	}
	p.AddStep(NewStep(steps...))
	// Oracle
	p.AddStep(NewStep(NewInnerProductOracle(s)))
	// Apply H to input qubits
	inputSteps := make([]Gate, n)
	for i := 0; i < n; i++ {
		inputSteps[i] = NewHadamard(i)
	}
	p.AddStep(NewStep(inputSteps...))
	return p
}

// NewSimonsProgram creates a program for Simon's algorithm.
func NewSimonsProgram(s string) *Program {
	n := len(s)
	p := NewProgram(2 * n)
	// Apply H to input qubits (0..n-1)
	inputSteps := make([]Gate, n)
	for i := 0; i < n; i++ {
		inputSteps[i] = NewHadamard(i)
	}
	p.AddStep(NewStep(inputSteps...))
	// Apply Simon oracle
	p.AddStep(NewStep(NewSimonOracle(s)))
	// Apply H to input qubits
	p.AddStep(NewStep(inputSteps...))
	return p
}

// NewErrorCorrectionProgram creates a program demonstrating the 3-qubit bit-flip code.
func NewErrorCorrectionProgram(bit int) *Program {
	p := NewProgram(3)
	// 1. Encoding
	if bit == 1 {
		p.AddStep(NewStep(NewX(0)))
	}
	p.AddStep(NewStep(NewCnot(0, 1)))
	p.AddStep(NewStep(NewCnot(0, 2)))

	// 2. Noise (bit-flip on q1)
	p.AddStep(NewStep(NewX(1)))

	// 3. Decoding / Correction
	p.AddStep(NewStep(NewCnot(0, 1)))
	p.AddStep(NewStep(NewCnot(0, 2)))
	p.AddStep(NewStep(NewToffoli(1, 2, 0)))
	return p
}

// NewSuperpositionProgram creates a program that puts n qubits into an equal superposition.
func NewSuperpositionProgram(n int) *Program {
	p := NewProgram(n)
	steps := make([]Gate, n)
	for i := 0; i < n; i++ {
		steps[i] = NewHadamard(i)
	}
	p.AddStep(NewStep(steps...))
	return p
}

// NewEngineProgram creates a program representing a quantum thermodynamic cycle.
func NewEngineProgram() *Program {
	p := NewProgram(1)
	// Step 1: Isothermal expansion (modeled by rotation)
	p.AddStep(NewStep(NewRx(math.Pi/4, 0)))
	// Step 2: Adiabatic expansion (Hadamard)
	p.AddStep(NewStep(NewHadamard(0)))
	// Step 3: Isothermal compression
	p.AddStep(NewStep(NewRy(math.Pi/4, 0)))
	return p
}

// NewToffoliProgram creates a program demonstrating the Toffoli gate.
func NewToffoliProgram() *Program {
	p := NewProgram(3)
	// Set controls to |11>
	p.AddStep(NewStep(NewX(0), NewX(1)))
	// Apply Toffoli
	p.AddStep(NewStep(NewToffoli(0, 1, 2)))
	return p
}

// NewFredkinProgram creates a program demonstrating the Fredkin gate.
func NewFredkinProgram() *Program {
	p := NewProgram(3)
	// Set control to |1> and one target to |1>
	p.AddStep(NewStep(NewX(0), NewX(1)))
	// Apply Fredkin (should swap q1 and q2, so q1 becomes 0 and q2 becomes 1)
	p.AddStep(NewStep(NewFredkin(0, 1, 2)))
	return p
}

// NewQKDProgram creates a program for a single bit QKD preparation.
func NewQKDProgram(bit, basis int) *Program {
	p := NewProgram(1)
	if bit == 1 {
		p.AddStep(NewStep(NewX(0)))
	}
	if basis == 1 {
		p.AddStep(NewStep(NewHadamard(0)))
	}
	return p
}

// ExpectedGroverState returns the theoretical state vector for 2-qubit Grover searching for |11>.
func ExpectedGroverState() []complex128 {
	return []complex128{0, 0, 0, 1}
}

// GetExpectedStateVector returns the theoretical state vector for a named circuit.
func GetExpectedStateVector(name string, n int) ([]complex128, error) {
	switch name {
	case "bell":
		return ExpectedBellState(), nil
	case "ghz":
		return ExpectedGHZState(n), nil
	case "qft":
		return ExpectedQFTState(n), nil
	case "grover":
		return ExpectedGroverState(), nil
	default:
		return nil, fmt.Errorf("no theoretical reference for circuit: %s", name)
	}
}

// CompareStateVectors compares two state vectors with a given tolerance.
func CompareStateVectors(got, want []complex128, tolerance float64) error {
	if len(got) != len(want) {
		return fmt.Errorf("dimension mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range got {
		diff := cmplx.Abs(got[i] - want[i])
		if diff > tolerance {
			return fmt.Errorf("value mismatch at index %d: got %v, want %v (diff %v)", i, got[i], want[i], diff)
		}
	}
	return nil
}
