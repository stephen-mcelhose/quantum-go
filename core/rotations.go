package core

import (
	gmath "math"
	"math/cmplx"

	"github.com/stephen-mcelhose/quantum-go/math"
)

// PhaseShift gate applies a phase rotation to the |1⟩ state.
// It is defined by the matrix:
// PhaseShift(θ) = [[1, 0], [0, e^(iθ)]]
// This gate is parameterized by an angle theta.
type PhaseShift struct {
	BaseGate
	// Theta is the rotation angle in radians.
	Theta float64
}

// NewPhaseShift creates a new phase shift gate with the specified angle and qubit index.
func NewPhaseShift(theta float64, idx int) *PhaseShift {
	return &PhaseShift{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "PS",
			Name:           "PhaseShift",
		},
		Theta: theta,
	}
}

// GetMatrix returns the 2x2 matrix representation of the phase shift gate.
// If the gate is marked as inverse, the conjugate transpose is returned.
func (g *PhaseShift) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	m.Data = []complex128{
		1, 0,
		0, cmplx.Exp(complex(0, g.Theta)),
	}
	if g.Inverse {
		m = math.ConjugateTranspose(m)
	}
	return m
}

// Rx gate (Rotation around X-axis) is defined by the matrix:
// Rx(θ) = [[cos(θ/2), -i*sin(θ/2)], [-i*sin(θ/2), cos(θ/2)]]
type Rx struct {
	BaseGate
	Theta float64
}

// NewRx creates a new Rx gate with the specified angle and qubit index.
func NewRx(theta float64, idx int) *Rx {
	return &Rx{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "Rx",
			Name:           "Rx",
		},
		Theta: theta,
	}
}

func (g *Rx) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	c := complex(gmath.Cos(g.Theta/2), 0)
	s := complex(0, -gmath.Sin(g.Theta/2))
	if g.Inverse {
		s = -s
	}
	m.Data = []complex128{c, s, s, c}
	return m
}

// Ry gate (Rotation around Y-axis) is defined by the matrix:
// Ry(θ) = [[cos(θ/2), -sin(θ/2)], [sin(θ/2), cos(θ/2)]]
type Ry struct {
	BaseGate
	Theta float64
}

// NewRy creates a new Ry gate with the specified angle and qubit index.
func NewRy(theta float64, idx int) *Ry {
	return &Ry{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "Ry",
			Name:           "Ry",
		},
		Theta: theta,
	}
}

func (g *Ry) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	c := complex(gmath.Cos(g.Theta/2), 0)
	s := complex(gmath.Sin(g.Theta/2), 0)
	if g.Inverse {
		m.Data = []complex128{c, s, -s, c}
	} else {
		m.Data = []complex128{c, -s, s, c}
	}
	return m
}

// Rz gate (Rotation around Z-axis) is defined by the matrix:
// Rz(θ) = [[e^(-iθ/2), 0], [0, e^(iθ/2)]]
type Rz struct {
	BaseGate
	Theta float64
}

// NewRz creates a new Rz gate with the specified angle and qubit index.
func NewRz(theta float64, idx int) *Rz {
	return &Rz{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "Rz",
			Name:           "Rz",
		},
		Theta: theta,
	}
}

func (g *Rz) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	p := cmplx.Exp(complex(0, g.Theta/2))
	if g.Inverse {
		m.Data = []complex128{p, 0, 0, cmplx.Conj(p)}
	} else {
		m.Data = []complex128{cmplx.Conj(p), 0, 0, p}
	}
	return m
}

// S gate (Phase gate) is the square root of Z gate.
// It is defined by the matrix:
// S = [[1, 0], [0, i]]
type S struct {
	BaseGate
}

// NewS creates a new S gate operating on the specified qubit index.
func NewS(idx int) *S {
	return &S{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "S",
			Name:           "S",
		},
	}
}

func (g *S) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	i := math.I
	if g.Inverse {
		i = -i
	}
	m.Data = []complex128{1, 0, 0, i}
	return m
}

// T gate (π/8 gate) is the square root of S gate.
// It is defined by the matrix:
// T = [[1, 0], [0, e^(iπ/4)]]
type T struct {
	BaseGate
}

// NewT creates a new T gate operating on the specified qubit index.
func NewT(idx int) *T {
	return &T{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "T",
			Name:           "T",
		},
	}
}

func (g *T) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	p := cmplx.Exp(complex(0, gmath.Pi/4))
	if g.Inverse {
		p = cmplx.Conj(p)
	}
	m.Data = []complex128{1, 0, 0, p}
	return m
}

// V gate (Square root of X) is also known as the SX gate.
// It is defined by the matrix:
// V = (1+i)/2 * [[1, -i], [-i, 1]]
type V struct {
	BaseGate
}

// NewV creates a new V gate operating on the specified qubit index.
func NewV(idx int) *V {
	return &V{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "V",
			Name:           "V",
		},
	}
}

func (g *V) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	// (1+i)/2
	f := complex(0.5, 0.5)
	mi := -math.I
	if g.Inverse {
		f = cmplx.Conj(f)
		mi = math.I
	}
	m.Data = []complex128{f, f * mi, f * mi, f}
	return m
}

// U gate (Universal rotation) is defined by three angles: theta, phi, and lambda.
// U(θ, φ, λ) = [[cos(θ/2), -e^(iλ)sin(θ/2)], [e^(iφ)sin(θ/2), e^(i(φ+λ))cos(θ/2)]]
type U struct {
	BaseGate
	Theta  float64
	Phi    float64
	Lambda float64
}

// NewU creates a new universal rotation gate with the specified angles and qubit index.
func NewU(theta, phi, lambda float64, idx int) *U {
	return &U{
		BaseGate: BaseGate{
			AffectedQubits: []int{idx},
			Caption:        "U",
			Name:           "U",
		},
		Theta:  theta,
		Phi:    phi,
		Lambda: lambda,
	}
}

func (g *U) GetMatrix() math.Matrix {
	m := math.NewMatrix(2, 2)
	c := gmath.Cos(g.Theta / 2)
	s := gmath.Sin(g.Theta / 2)
	m.Data = []complex128{
		complex(c, 0),
		-cmplx.Exp(complex(0, g.Lambda)) * complex(s, 0),
		cmplx.Exp(complex(0, g.Phi)) * complex(s, 0),
		cmplx.Exp(complex(0, g.Phi+g.Lambda)) * complex(c, 0),
	}
	if g.Inverse {
		m = math.ConjugateTranspose(m)
	}
	return m
}

// Cr gate (Controlled Phase Shift) applies a phase rotation to the target qubit
// when the control qubit is in the |1⟩ state.
// It is a controlled version of the PhaseShift gate.
type Cr struct {
	BaseGate
	// Theta is the rotation angle in radians.
	Theta float64
}

// NewCr creates a new controlled rotation gate with the specified control qubit,
// target qubit, and rotation angle.
func NewCr(control, target int, theta float64) *Cr {
	return &Cr{
		BaseGate: BaseGate{
			AffectedQubits: []int{control, target},
			Caption:        "CR",
			Name:           "Controlled Rotation",
		},
		Theta: theta,
	}
}

// GetMatrix returns the 4x4 matrix representation of the controlled rotation gate.
// If the gate is marked as inverse, the conjugate of the rotation is applied.
func (g *Cr) GetMatrix() math.Matrix {
	m := math.NewMatrix(4, 4)
	rot := cmplx.Exp(complex(0, g.Theta))
	if g.Inverse {
		rot = cmplx.Conj(rot)
	}
	m.Data = []complex128{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, rot,
	}
	return m
}
