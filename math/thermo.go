// Package math provides thermodynamic operations for quantum systems.
package math

import (
	"math"
	"math/cmplx"
)

// ToDensityMatrix converts a state vector into a density matrix ρ = |ψ⟩⟨ψ|.
func ToDensityMatrix(v []complex128) Matrix {
	dim := len(v)
	res := NewMatrix(dim, dim)
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			res.Set(i, j, v[i]*cmplx.Conj(v[j]))
		}
	}
	return res
}

// PartialTrace computes the partial trace of a density matrix over specified qubits.
// nQubits is the total number of qubits in the system.
// traceQubits are the indices of the qubits to be traced out.
func PartialTrace(m Matrix, nQubits int, traceQubits []int) Matrix {
	isTraced := make(map[int]bool)
	for _, q := range traceQubits {
		isTraced[q] = true
	}

	remQubits := []int{}
	for i := 0; i < nQubits; i++ {
		if !isTraced[i] {
			remQubits = append(remQubits, i)
		}
	}

	nRem := len(remQubits)
	dimRem := 1 << nRem
	res := NewMatrix(dimRem, dimRem)

	for rPrime := 0; rPrime < dimRem; rPrime++ {
		for cPrime := 0; cPrime < dimRem; cPrime++ {
			var sum complex128
			nTr := len(traceQubits)
			dimTr := 1 << nTr
			for t := 0; t < dimTr; t++ {
				r := 0
				for i, q := range remQubits {
					if (rPrime>>i)&1 == 1 {
						r |= (1 << q)
					}
				}
				for i, q := range traceQubits {
					if (t>>i)&1 == 1 {
						r |= (1 << q)
					}
				}

				c := 0
				for i, q := range remQubits {
					if (cPrime>>i)&1 == 1 {
						c |= (1 << q)
					}
				}
				for i, q := range traceQubits {
					if (t>>i)&1 == 1 {
						c |= (1 << q)
					}
				}
				sum += m.Get(r, c)
			}
			res.Set(rPrime, cPrime, sum)
		}
	}
	return res
}

// Trace returns the trace of a square matrix.
func Trace(m Matrix) complex128 {
	var sum complex128
	for i := 0; i < m.Rows; i++ {
		sum += m.Get(i, i)
	}
	return sum
}

// ExpectationValue computes tr(ρH), the expectation value of observable H in state ρ.
func ExpectationValue(rho Matrix, h Matrix) complex128 {
	var sum complex128
	dim := rho.Rows
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			sum += rho.Get(i, j) * h.Get(j, i)
		}
	}
	return sum
}

// VonNeumannEntropy computes the von Neumann entropy S(ρ) = -tr(ρ ln ρ).
// This uses analytical eigenvalue formulas for 1-qubit systems (dim 2).
func VonNeumannEntropy(rho Matrix) float64 {
	dim := rho.Rows
	if dim == 2 {
		// 1 qubit: λ^2 - tr(ρ)λ + det(ρ) = 0
		// tr(ρ) is assumed to be 1 for a valid density matrix.
		det := real(rho.Get(0, 0)*rho.Get(1, 1) - rho.Get(0, 1)*rho.Get(1, 0))
		
		// λ = (1 ± sqrt(1 - 4*det)) / 2
		disc := 1.0 - 4.0*det
		if disc < 0 {
			disc = 0 // Numerical stability
		}
		sq := math.Sqrt(disc)
		l1 := (1.0 + sq) / 2.0
		l2 := (1.0 - sq) / 2.0
		
		return entropyFromEigenvalues([]float64{l1, l2})
	}
	// For larger systems, a general eigensolver would be needed.
	// Note: For pure states of bipartite systems, S(A) = S(B).
	return 0
}

// MutualInformation computes I(A:B) = S(ρA) + S(ρB) - S(ρAB).
// For unitary evolution starting from a pure state, S(ρAB) = 0,
// and I(A:B) = S(ρA) + S(ρB).
func MutualInformation(rhoAB Matrix, nQubits int, qubitsA, qubitsB []int) float64 {
	rhoA := PartialTrace(rhoAB, nQubits, qubitsB)
	rhoB := PartialTrace(rhoAB, nQubits, qubitsA)
	
	sA := VonNeumannEntropy(rhoA)
	sB := VonNeumannEntropy(rhoB)
	sAB := VonNeumannEntropy(rhoAB) // Will be 0 if pure
	
	return sA + sB - sAB
}

// RelativeEntropy computes the Kullback-Leibler divergence D(ρ || σ) = tr(ρ(ln ρ - ln σ)).
// This currently only supports 1-qubit systems (dim 2).
func RelativeEntropy(rho, sigma Matrix) float64 {
	// D(ρ || σ) = -S(ρ) - tr(ρ ln σ)
	// For dim 2, ln σ can be found via diagonalization or analytical formulas.
	// σ = 1/2(I + r·σ) -> ln σ = ...
	// simpler: if σ is diagonal (thermal state), ln σ is diag(ln λi).
	
	// This is a simplified implementation for diagonal sigma:
	sRho := VonNeumannEntropy(rho)
	
	// Check if sigma is diagonal
	if cmplx.Abs(sigma.Get(0, 1)) < 1e-10 && cmplx.Abs(sigma.Get(1, 0)) < 1e-10 {
		l1 := real(sigma.Get(0, 0))
		l2 := real(sigma.Get(1, 1))
		
		trRhoLnSigma := real(rho.Get(0, 0)) * math.Log(l1) + real(rho.Get(1, 1)) * math.Log(l2)
		return -sRho - trRhoLnSigma
	}
	
	return 0
}

// Exp computes the matrix exponential e^M using a Taylor series expansion.
// This is suitable for small matrices or small norms.
// It continues iterations until the term norm is below the specified epsilon or max iterations reached.
func Exp(m Matrix) Matrix {
	res := IdentityMatrix(m.Rows)
	term := IdentityMatrix(m.Rows)
	epsilon := 1e-15

	for i := 1; i < 100; i++ {
		term = Mul(term, m)
		// scale by 1/i
		var maxNorm float64
		for j := range term.Data {
			term.Data[j] /= complex(float64(i), 0)
			norm := cmplx.Abs(term.Data[j])
			if norm > maxNorm {
				maxNorm = norm
			}
		}
		// add to res
		for j := range res.Data {
			res.Data[j] += term.Data[j]
		}
		if maxNorm < epsilon {
			break
		}
	}
	return res
}

func entropyFromEigenvalues(evs []float64) float64 {
	var s float64
	for _, l := range evs {
		if l > 1e-12 {
			s -= l * math.Log(l)
		}
	}
	return s
}
