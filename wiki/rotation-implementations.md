---
type: concept
title: Rotation Gate Implementations
description: How core/rotations.go implements Rx/Ry/Rz/U/S/T/V/PhaseShift/CR — parameterized gate types with angle fields, cmplx.Exp, and inverse support via ConjugateTranspose.
resource: core/rotations.go
tags: [rotations, rx, ry, rz, u-gate, s-gate, t-gate, v-gate, cr, phaseshift, inverse]
timestamp: 2026-08-09T03:26:15Z
---

# Rotation Gate Implementations

`core/rotations.go` implements all parameterized single-qubit gates. Each stores angle(s) as fields and computes matrix entries via `math/cmplx`. Key pattern: **inverse support** — calling `g.SetInverse(true)` applies `ConjugateTranspose` (hermitian conjugate = complex conjugate + transpose = U†) to the output matrix, reversing the gate.

## Implementation Pattern

```go
type PhaseShift struct {
    BaseGate
    Theta float64   // stored as field, used in GetMatrix()
}

func (g *PhaseShift) GetMatrix() math.Matrix {
    m := math.NewMatrix(2, 2)
    p := cmplx.Exp(complex(0, g.Theta))  // e^{iθ}
    if g.Inverse {
        p = cmplx.Conj(p)               // e^{-iθ} = (e^{iθ})†
    }
    m.Data = []complex128{1, 0, 0, p}
    return m
}
```

## Gate Catalog

### PhaseShift(θ)

```
[ 1    0        ]
[ 0    e^{iθ}  ]
```
- `NewPhaseShift(theta, idx)` — theta in radians
- PS(π) = Z, PS(π/2) = S, PS(π/4) = T
- QASM: `u1(θ)` → `core.NewPhaseShift(theta, q0)` (see [[openqasm-parser]])

### S Gate (Phase / √Z)

```
[ 1  0 ]
[ 0  i ]
```
- `NewS(idx)` — stored as `math.I` constant
- Inverse: `S† = [ 1, 0; 0, -i ]` (via `if g.Inverse { i = -i }`)
- QASM: `s` / `sdg` (inverse)

### T Gate (π/8 / √S)

```
[ 1  0           ]
[ 0  e^{iπ/4}  ]
```
- `NewT(idx)` — computed via `cmplx.Exp(complex(0, π/4))`
- Inverse: `T†` — conjugate (QASM: `tdg`)
- T gate is the key non-Clifford gate needed for [[universality]]

### V Gate (SX / √X)

```
(1+i)/2 * [ 1  -i ]
           [ -i  1 ]
```
- Factor: `f = complex(0.5, 0.5)`, cross-term: `-i`
- Inverse: `f = cmplx.Conj(complex(0.5, 0.5)) = (0.5-0.5i)`, cross-term: `+i`
- `NewV(idx)` — QASM maps `sx` → V
- V² = X (two V gates = Pauli-X)

### Rx(θ)

```
[ cos(θ/2)      -i·sin(θ/2) ]
[ -i·sin(θ/2)  cos(θ/2)    ]
```
- Computed: `c = gmath.Cos(theta/2)`, `s = gmath.Sin(theta/2)`, off-diagonals `∓complex(0,s)`
- Inverse: `ConjugateTranspose` on the 2×2 matrix
- QASM: `rx(θ)`

### Ry(θ)

```
[ cos(θ/2)   -sin(θ/2) ]
[ sin(θ/2)    cos(θ/2) ]
```
- Real matrix — no imaginary entries
- QASM: `ry(θ)`

### Rz(θ)

```
[ e^{-iθ/2}   0          ]
[ 0            e^{+iθ/2} ]
```
- Computed as `cmplx.Exp(complex(0, ±theta/2))`
- Rz only affects phase; [[gate-application]] uses the CR optimized path for controlled-Rz
- QASM: `rz(θ)`

### U Gate (Universal)

```
[ cos(θ/2)              -e^{iλ}·sin(θ/2)         ]
[ e^{iφ}·sin(θ/2)       e^{i(φ+λ)}·cos(θ/2)     ]
```
- `NewU(theta, phi, lambda, idx)`
- Inverse: `ConjugateTranspose` on the 2×2 result
- QASM: `u3(θ, φ, λ)` and `u(θ, φ, λ)` → both map to `NewU`
- Every named gate is a special case (see [[universality]])

### CR Gate (Controlled Phase Shift)

```
[ 1  0  0   0         ]
[ 0  1  0   0         ]
[ 0  0  1   0         ]
[ 0  0  0   e^{iθ}   ]
```
- `NewCr(control, target, theta)`
- Applies phase rotation only to |11⟩ — the foundation of the [[qft-deep-dive]] (QFT uses CR at 2π/2ᵏ angles)
- Inverse: conjugate of rotation factor
- QASM: `cu1(θ)` → `NewCr(q0, q1, theta)`

## Inverse Support

All rotation gates support `SetInverse(true)`, which causes `GetMatrix()` to return U† instead of U:
- Phase gates (S, T, PS, CR): conjugate the scalar `e^{iθ}` factor → `e^{-iθ}`
- Rotation gates (Rx, Ry, Rz): `ConjugateTranspose` on the result matrix
- U gate: `math.ConjugateTranspose(m)`

`ConjugateTranspose` is implemented in `math/matrix.go` — see [[quantum-linear-algebra]].

## Key Points

- All angles stored as `float64` fields, computed fresh on each `GetMatrix()` call (no caching).
- Inverse = complex conjugate on phase factor (scalars) or full ConjugateTranspose (matrices).
- V gate (SX) maps to QASM `sx`; V² = X.
- T gate is the critical non-Clifford gate for [[universality]].
- CR is the 4×4 controlled phase shift — the engine has an optimized path for it (scalar multiply on the |11⟩ amplitude only).
- U gate subsumes all single-qubit rotations: H = U(π/2, 0, π).

## Sources

- `core/rotations.go`
