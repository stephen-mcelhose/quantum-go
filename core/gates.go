package core

import (
	"github.com/stephen-mcelhose/quantum-go/math"
)

// Hadamard gate creates an equal superposition of |0⟩ and |1⟩ states.
// It is one of the most fundamental quantum gates, defined by the matrix:
// H = (1/√2) * [[1, 1], [1, -1]]
type Hadamard struct {
	BaseGate
}

// NewHadamard creates a new Hadamard gate operating on the specified qubit index.
func NewHadamard(idx int) *Hadamard {
	return &Hadamard{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "H",
			Name:           "Hadamard",
		},
	}
}

// GetMatrix returns the 2x2 matrix representation of the Hadamard gate.
func (g *Hadamard) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	h := complex(math.HV, 0)
	m.Data = []complex128{h, h, h, -h}
	return m
}

// X gate (Pauli-X) is the quantum equivalent of a classical NOT gate.
// It flips |0⟩ to |1⟩ and vice versa, defined by the matrix:
// X = [[0, 1], [1, 0]]
type X struct {
	BaseGate
}

// NewX creates a new Pauli-X gate operating on the specified qubit index.
func NewX(idx int) *X {
	return &X{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "X",
			Name:           "X",
		},
	}
}

// GetMatrix returns the 2x2 matrix representation of the Pauli-X gate.
func (g *X) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	m.Data = []complex128{0, 1, 1, 0}
	return m
}

// Y gate (Pauli-Y) performs a bit flip and a phase flip.
// It is defined by the matrix:
// Y = [[0, -i], [i, 0]]
type Y struct {
	BaseGate
}

// NewY creates a new Pauli-Y gate operating on the specified qubit index.
func NewY(idx int) *Y {
	return &Y{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "Y",
			Name:           "Y",
		},
	}
}

// GetMatrix returns the 2x2 matrix representation of the Pauli-Y gate.
func (g *Y) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	m.Data = []complex128{0, -math.I, math.I, 0}
	return m
}

// Z gate (Pauli-Z) performs a phase flip, leaving |0⟩ unchanged and flipping the phase of |1⟩.
// It is defined by the matrix:
// Z = [[1, 0], [0, -1]]
type Z struct {
	BaseGate
}

// NewZ creates a new Pauli-Z gate operating on the specified qubit index.
func NewZ(idx int) *Z {
	return &Z{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "Z",
			Name:           "Z",
		},
	}
}

// GetMatrix returns the 2x2 matrix representation of the Pauli-Z gate.
func (g *Z) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	m.Data = []complex128{1, 0, 0, -1}
	return m
}

// Cnot gate (Controlled-NOT) flips the target qubit if and only if the control qubit is |1⟩.
// It is a fundamental two-qubit gate for creating entanglement.
// The matrix representation is 4x4 for the two-qubit space.
type Cnot struct {
	BaseGate
}

// NewCnot creates a new CNOT gate with the specified control and target qubit indices.
func NewCnot(control, target int) *Cnot {
	return &Cnot{
		BaseGate: BaseGate{
			AffectedQubits: []int{control, target},
			Caption:        "CNOT",
			Name:           "Cnot",
		},
	}
}

// GetMatrix returns the 4x4 matrix representation of the CNOT gate.
func (g *Cnot) GetMatrix() math.Matrix {
	m := math.NewMatrix(4, 4)
	m.Data = []complex128{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 0, 1,
		0, 0, 1, 0,
	}
	return m
}

// Identity gate leaves the qubit state unchanged.
// It is defined by the identity matrix:
// I = [[1, 0], [0, 1]]
type Identity struct {
	BaseGate
}

// NewIdentity creates a new identity gate operating on the specified qubit index.
func NewIdentity(idx int) *Identity {
	return &Identity{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "I",
			Name:           "Identity",
		},
	}
}

// GetMatrix returns the 2x2 identity matrix.
func (g *Identity) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	m.Data = []complex128{1, 0, 0, 1}
	return m
}

// Measurement gate represents a quantum measurement operation.
// When executed, it collapses the qubit's state to either |0⟩ or |1⟩
// according to its probability distribution.
type Measurement struct {
	BaseGate
}

// NewMeasurement creates a new measurement gate for the specified qubit index.
func NewMeasurement(idx int) *Measurement {
	return &Measurement{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "M",
			Name:           "Measurement",
		},
	}
}

// GetMatrix returns the identity matrix since measurement is handled specially.
func (g *Measurement) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	m.Data = []complex128{1, 0, 0, 1}
	return m
}

// Cz gate (Controlled-Z) applies a phase flip to the target qubit if the control qubit is |1⟩.
// This is a symmetric two-qubit gate where control and target can be interchanged.
type Cz struct {
	BaseGate
}

// NewCz creates a new Controlled-Z gate with the specified control and target qubit indices.
func NewCz(control, target int) *Cz {
	return &Cz{
		BaseGate: BaseGate{
			AffectedQubits: []int{control, target},
			Caption:        "CZ",
			Name:           "Cz",
		},
	}
}

// GetMatrix returns the 4x4 matrix representation of the Controlled-Z gate.
func (g *Cz) GetMatrix() math.Matrix {
	m := math.NewMatrix(4, 4)
	m.Data = []complex128{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, -1,
	}
	return m
}

// Swap gate exchanges the quantum states of two qubits.
// It swaps |01⟩ ↔ |10⟩ while leaving |00⟩ and |11⟩ unchanged.
type Swap struct {
	BaseGate
}

// NewSwap creates a new SWAP gate operating on the two specified qubit indices.
func NewSwap(q1, q2 int) *Swap {
	return &Swap{
		BaseGate: BaseGate{
			AffectedQubits: []int{q1, q2},
			Caption:        "SWAP",
			Name:           "Swap",
		},
	}
}

// GetMatrix returns the 4x4 matrix representation of the SWAP gate.
func (g *Swap) GetMatrix() math.Matrix {
	m := math.NewMatrix(4, 4)
	m.Data = []complex128{
		1, 0, 0, 0,
		0, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 0, 1,
	}
	return m
}

// Toffoli gate (Controlled-Controlled-NOT or CCNOT) is a three-qubit gate.
// It flips the target qubit if and only if both control qubits are |1⟩.
// This is a universal gate for classical reversible computing.
type Toffoli struct {
	BaseGate
}

// NewToffoli creates a new Toffoli gate with two control qubits (a, b) and one target qubit (c).
func NewToffoli(a, b, c int) *Toffoli {
	return &Toffoli{
		BaseGate: BaseGate{
			AffectedQubits: []int{a, b, c},
			Caption:        "CCNOT",
			Name:           "Toffoli",
		},
	}
}

// GetMatrix returns the 8x8 matrix representation of the Toffoli gate.
func (g *Toffoli) GetMatrix() math.Matrix {
	m := math.NewMatrix(8, 8)
	for i := 0; i < 6; i++ {
		m.Set(i, i, 1)
	}
	m.Set(6, 7, 1)
	m.Set(7, 6, 1)
	return m
}

// Fredkin gate (Controlled-SWAP or CSWAP) is a three-qubit gate.
// It swaps the target qubits if and only if the control qubit is |1⟩.
type Fredkin struct {
	BaseGate
}

// NewFredkin creates a new Fredkin gate with one control qubit (a) and two target qubits (b, c).
func NewFredkin(a, b, c int) *Fredkin {
	return &Fredkin{
		BaseGate: BaseGate{
			AffectedQubits: []int{a, b, c},
			Caption:        "CSWAP",
			Name:           "Fredkin",
		},
	}
}

// GetMatrix returns the 8x8 matrix representation of the Fredkin gate.
func (g *Fredkin) GetMatrix() math.Matrix {
	m := math.NewMatrix(8, 8)
	// Identity for all states where control bit (affected[0]) is 0.
	// Identity for states where target bits (affected[1], affected[2]) are both 0 or both 1.
	for i := 0; i < 8; i++ {
		m.Set(i, i, 1)
	}

	// Swap |101> and |110> where bit 0 is control, bit 1 and 2 are targets.
	// But we need to use the actual affected qubit indices.
	// For simplicity, if we assume a=0, b=1, c=2:
	// Control q0=1 (odd indices).
	// Targets q1, q2.
	// Swap |q2=0, q1=1, q0=1> (index 3) and |q2=1, q1=0, q0=1> (index 5).
	
	// Reset the indices we want to swap
	m.Set(3, 3, 0)
	m.Set(5, 5, 0)
	
	m.Set(3, 5, 1)
	m.Set(5, 3, 1)
	return m
}
