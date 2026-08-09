package core

import (
	gmath "math"

	"github.com/stephen-mcelhose/quantum-go/math"
)

// Fourier gate implements the Quantum Fourier Transform (QFT).
// The QFT is a linear transformation on quantum amplitudes and is a key component
// in many quantum algorithms including Shor's factoring algorithm.
type Fourier struct {
	BlockGate
}

// NewFourier creates a new QFT gate operating on 'dim' qubits starting at index 'idx'.
// The QFT is constructed as a block of Hadamard gates, controlled rotations, and swaps.
func NewFourier(dim, idx int) *Fourier {
	block := NewBlock("Fourier", dim)
	// We'll populate the block steps in a separate function to avoid circularity if possible,
	// or just do it here.
	f := &Fourier{
		BlockGate: *NewBlockGate(block, idx),
	}
	f.Caption = "QFT"
	f.Name = "Fourier"

	// Populate Fourier block
	for i := dim - 1; i >= 0; i-- {
		block.AddStep(NewStep(NewHadamard(i)))
		for j := 2; j <= i+1; j++ {
			// Cr(target, control, theta) ? No, Java: Cr(i+1-j, i, 2, j)
			// Java Cr(control, target, base, pow) -> R(2*PI / base^pow)
			theta := 2.0 * 3.141592653589793 / gmath.Pow(2.0, float64(j))
			block.AddStep(NewStep(NewCr(i+1-j, i, theta)))

		}
	}
	// Swap gates
	for i := 0; i < dim/2; i++ {
		block.AddStep(NewStep(NewSwap(i, dim-1-i)))
	}

	// Adjust affected qubits to be absolute circuit indexes
	for i := range f.AffectedQubits {
		f.AffectedQubits[i] = idx + i
	}

	return f
}

// Add gate implements quantum addition using the QFT.
// It adds two quantum registers in Fourier space, performing the operation:
// |x⟩|y⟩ → |x⟩|x+y⟩
// The gate operates on qubits from x0 to y1, where [x0,x1] is the first register
// and [y0,y1] is the second register that will hold the result.
type Add struct {
	BlockGate
}

// NewAdd creates a new quantum addition gate.
// x0, x1 define the first quantum register (input x)
// y0, y1 define the second quantum register (input y, output x+y)
func NewAdd(x0, x1, y0, y1 int) *Add {
	m := x1 - x0 + 1
	n := y1 - y0 + 1
	dim := y1 - x0 + 1
	block := NewBlock("Add", dim)

	// x_0 ----- y_0 + x_0
	// x_1 ----- y_1 + x_1
	// y_0 ----- y_0
	// y_1 ----- y_1

	block.AddStep(NewStep(NewFourier(m, 0)))
	for i := 0; i < m; i++ {
		for j := 0; j < m-i; j++ {
			cr0 := 2*m - j - i - 1
			if cr0 < m+n {
				theta := 2.0 * 3.141592653589793 / gmath.Pow(2.0, float64(1+j))
				block.AddStep(NewStep(NewCr(i, cr0, theta)))
			}
		}
	}
	invFourier := NewFourier(m, 0)
	invFourier.SetInverse(true)
	block.AddStep(NewStep(invFourier))

	a := &Add{
		BlockGate: *NewBlockGate(block, x0),
	}
	a.Caption = "ADD"
	a.Name = "Add"

	// Adjust affected qubits to be absolute circuit indexes
	for i := range a.AffectedQubits {
		a.AffectedQubits[i] = x0 + i
	}

	return a
}

// AddInteger gate adds a classical integer to a quantum register.
type AddInteger struct {
	BlockGate
}

// NewAddInteger creates a new AddInteger gate.
func NewAddInteger(x0, x1, num int) *AddInteger {
	m := x1 - x0 + 1
	block := NewBlock("AddInteger", m)
	
	// Implementation based on Java AddInteger.java
	block.AddStep(NewStep(NewFourier(m, 0)))
	
	pstep := NewStep()
	for i := 0; i < m; i++ {
		mat := math.IdentityMatrix(2)
		for j := 0; j < m-i; j++ {
			cr0 := m - j - i - 1
			if (num>>cr0)&1 == 1 {
				theta := 2.0 * 3.141592653589793 / gmath.Pow(2.0, float64(1+j))
				rGate := NewPhaseShift(theta, i)
					// Combine matrices
					gateMat := rGate.GetMatrix()
					mat = math.Mul(mat, gateMat)


			}
		}
		pstep.AddGate(NewSingleQubitMatrixGate(i, mat))
	}
	block.AddStep(pstep)
	
	invFourier := NewFourier(m, 0)
	invFourier.SetInverse(true)
	block.AddStep(NewStep(invFourier))

	res := &AddInteger{
		BlockGate: *NewBlockGate(block, x0),
	}
	res.Caption = "ADDI"
	res.Name = "AddInteger"
	
	for i := range res.AffectedQubits {
		res.AffectedQubits[i] = x0 + i
	}
	return res
}

// AddIntegerModulus gate adds an integer mod N to a quantum register.
type AddIntegerModulus struct {
	BlockGate
}

// NewAddIntegerModulus creates a new AddIntegerModulus gate.
func NewAddIntegerModulus(x0, x1, a, N int) *AddIntegerModulus {
	n := x1 - x0 + 1
	block := NewBlock("AddIntegerModulus", n+1)
	dim := n // carry qubit index in block
	
	// Java logic
	add := NewAddInteger(0, n-1, a)
	block.AddStep(NewStep(add))
	
	min := NewAddInteger(0, n-1, N)
	min.SetInverse(true)
	block.AddStep(NewStep(min))
	
	block.AddStep(NewStep(NewCnot(n-1, dim)))
	
	addN := NewAddInteger(0, n-1, N)
	cbg := NewControlledBlockGate(addN, dim)
	block.AddStep(NewStep(cbg))
	
	add2 := NewAddInteger(0, n-1, a)
	add2.SetInverse(true)
	block.AddStep(NewStep(add2))
	
	block.AddStep(NewStep(NewX(dim)))
	block.AddStep(NewStep(NewCnot(n-1, dim)))
	block.AddStep(NewStep(NewX(dim)))
	
	add3 := NewAddInteger(0, n-1, a)
	block.AddStep(NewStep(add3))
	
	res := &AddIntegerModulus{
		BlockGate: *NewBlockGate(block, x0),
	}
	res.Caption = "ADDIM"
	res.Name = "AddIntegerModulus"
	
	affected := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		affected[i] = x0 + i
	}
	res.AffectedQubits = affected
	return res
}

// MulModulus gate performs modular multiplication.
type MulModulus struct {
	BlockGate
}

// NewMulModulus creates a new MulModulus gate.
func NewMulModulus(x0, x1, mul, mod int) *MulModulus {
	size := x1 - x0 + 1
	block := NewBlock("MulModulus", 2*size+1) // Java has 2*size+2? Let's check.
	// Java: Block answer = new Block("MulModulus", 2 * size+2);
	
	n := size
	for i := 0; i < n; i++ {
		m := (mul * (1 << i)) % mod
		add := NewAddIntegerModulus(0, n-1, m, mod)
		// Wait, Java: ControlledBlockGate cbg = new ControlledBlockGate(add, n, i);
		// In Javabg, idx is n, control is i.
		// Our NewControlledBlockGate(gate, control)
		cbg := NewControlledBlockGate(add, i)
		block.AddStep(NewStep(cbg))
	}
	
	for i := 0; i < n; i++ {
		block.AddStep(NewStep(NewSwap(i, i+n)))
	}
	
	invmul := getInverseModulus(mul, mod)
	for i := 0; i < n; i++ {
		m := (invmul * (1 << i)) % mod
		add := NewAddIntegerModulus(0, n-1, m, mod)
		cbg := NewControlledBlockGate(add, i)
		cbg.SetInverse(true)
		block.AddStep(NewStep(cbg))
	}
	
	res := &MulModulus{
		BlockGate: *NewBlockGate(block, x0),
	}
	res.Caption = "MULM"
	res.Name = "MulModulus"
	
	affected := make([]int, 2*size+1)
	for i := 0; i < 2*size+1; i++ {
		affected[i] = x0 + i
	}
	res.AffectedQubits = affected
	return res
}

func getInverseModulus(a, n int) int {
	// Extended Euclidean Algorithm
	t := 0
	newt := 1
	r := n
	newr := a
	for newr != 0 {
		quotient := r / newr
		t, newt = newt, t-quotient*newt
		r, newr = newr, r-quotient*newr
	}
	if r > 1 {
		return -1 // a is not invertible
	}
	if t < 0 {
		t = t + n
	}
	return t
}
