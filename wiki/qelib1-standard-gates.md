---
type: concept
title: qelib1 Standard Gate Library
description: The OpenQASM 2.0 standard gate library (qelib1.inc) — every gate defined, its decomposition, and how it maps to quantum-go constructors. Apache 2.0 licensed.
resource: https://raw.githubusercontent.com/Qiskit/openqasm/OpenQASM2.x/examples/qelib1.inc
tags: [openqasm, qelib1, standard-gates, gate-library, decomposition, qasm2]
timestamp: 2026-08-09T03:26:15Z
---

# qelib1 Standard Gate Library

`qelib1.inc` is the standard header included by every QASM 2.0 file via `include "qelib1.inc";`. It defines all named gates in terms of two hardware primitives: `U(θ,φ,λ)` (universal single-qubit) and `CX` (CNOT). The quantum-go parser reads this header implicitly when it sees the include statement.

> Source: Apache 2.0 — [Qiskit/openqasm @ OpenQASM2.x](https://github.com/Qiskit/openqasm/tree/OpenQASM2.x)

## Hardware Primitives (built-in, not in qelib1.inc)

| QASM name | Definition          | quantum-go          |
| --------- | ------------------- | ------------------- |
| `U(θ,φ,λ)`| 3-param Euler gate  | `NewU(θ, φ, λ, i)` |
| `CX c, t` | Controlled-NOT      | `NewCnot(c, t)`     |

## Hardware Aliases (defined in qelib1.inc)

```qasm
gate u3(theta,phi,lambda) q { U(theta,phi,lambda) q; }  // alias for U
gate u2(phi,lambda)        q { U(pi/2,phi,lambda) q; }   // U with θ=π/2
gate u1(lambda)            q { U(0,0,lambda) q; }         // phase-only, = PhaseShift
gate cx c,t                  { CX c,t; }                  // alias for CX
gate id a                    { U(0,0,0) a; }               // identity
```

| QASM name       | quantum-go                     | Notes                              |
| --------------- | ------------------------------ | ---------------------------------- |
| `u3(θ,φ,λ)`    | `NewU(θ, φ, λ, i)`             | Full 3-param universal             |
| `u2(φ,λ)`      | `NewU(π/2, φ, λ, i)`           | Parser computes θ=π/2              |
| `u1(λ)`         | `NewPhaseShift(λ, i)`          | Diagonal phase gate                |
| `cx c, t`       | `NewCnot(c, t)`                | Direct alias                       |
| `id a`          | `NewIdentity(i)`               | No-op                              |

## Standard Single-Qubit Gates

```qasm
gate x a   { u3(pi,0,pi) a; }        // Pauli-X (bit-flip)
gate y a   { u3(pi,pi/2,pi/2) a; }   // Pauli-Y (bit+phase flip)
gate z a   { u1(pi) a; }             // Pauli-Z (phase flip)
gate h a   { u2(0,pi) a; }           // Hadamard
gate s a   { u1(pi/2) a; }           // √Z  (S gate)
gate sdg a { u1(-pi/2) a; }          // S†  (S-dagger)
gate t a   { u1(pi/4) a; }           // ⁴√Z (T gate)
gate tdg a { u1(-pi/4) a; }          // T†  (T-dagger)
```

| QASM  | quantum-go                          |
| ----- | ----------------------------------- |
| `x`   | `NewX(i)`                           |
| `y`   | `NewY(i)`                           |
| `z`   | `NewZ(i)`                           |
| `h`   | `NewHadamard(i)`                    |
| `s`   | `NewS(i)`                           |
| `sdg` | `NewS(i)` + `.SetInverse(true)`     |
| `t`   | `NewT(i)`                           |
| `tdg` | `NewT(i)` + `.SetInverse(true)`     |

## Standard Rotation Gates

```qasm
gate rx(theta) a { u3(theta,-pi/2,pi/2) a; }  // Bloch X rotation
gate ry(theta) a { u3(theta,0,0) a; }          // Bloch Y rotation
gate rz(phi)   a { u1(phi) a; }                // Bloch Z rotation (= PhaseShift up to global phase)
```

| QASM       | quantum-go              | Notes                                               |
| ---------- | ----------------------- | --------------------------------------------------- |
| `rx(θ)`    | `NewRx(θ, i)`           | Direct                                              |
| `ry(θ)`    | `NewRy(θ, i)`           | Direct                                              |
| `rz(φ)`    | `NewRz(φ, i)`           | Note: `rz` ≠ `u1` up to global phase; parser maps to `NewRz` |

## Standard Two-Qubit Gates

```qasm
gate cz a,b  { h b; cx a,b; h b; }   // controlled-Z (via H sandwich)
gate cy a,b  { sdg b; cx a,b; s b; } // controlled-Y
gate ch a,b  { ... }                  // controlled-H (7-gate decomposition)
```

| QASM  | quantum-go        | Notes                                       |
| ----- | ----------------- | ------------------------------------------- |
| `cz`  | `NewCz(a, b)`     | Parser maps directly; qelib1 defines via H+CX |
| `cy`  | not implemented   | Decomposed in qelib1; no direct quantum-go gate |
| `ch`  | not implemented   | 7-gate decomposition; no direct quantum-go gate |

## Controlled Rotation Gates

```qasm
gate cu1(lambda) a,b {        // controlled phase (= CR gate)
  u1(lambda/2) a;
  cx a,b;
  u1(-lambda/2) b;
  cx a,b;
  u1(lambda/2) b;
}

gate crz(lambda) a,b {        // controlled-Rz
  u1(lambda/2) b;
  cx a,b;
  u1(-lambda/2) b;
  cx a,b;
}

gate cu3(theta,phi,lambda) c,t { ... }  // controlled-U3
```

| QASM         | quantum-go          | Notes                                         |
| ------------ | ------------------- | --------------------------------------------- |
| `cu1(λ)`     | `NewCr(c, t, λ)`    | Parser maps directly to CR gate               |
| `crz(λ)`     | not implemented     | Decomposed; no direct quantum-go gate          |
| `cu3(θ,φ,λ)` | not implemented     | General controlled-U; not in parser            |

## Three-Qubit Gates

```qasm
gate ccx a,b,c {      // Toffoli (15-gate CX decomposition)
  h c;
  cx b,c; tdg c;
  cx a,c; t c;
  cx b,c; tdg c;
  cx a,c; t b; t c; h c;
  cx a,b; t a; tdg b;
  cx a,b;
}
```

| QASM  | quantum-go            |
| ----- | --------------------- |
| `ccx` | `NewToffoli(a, b, c)` |

The qelib1 decomposition uses 15 CNOT+T gates. quantum-go's `NewToffoli` uses the optimised 8×8 matrix directly — the parser maps `ccx` to the native gate, not the decomposition.

## Canonical QASM 2.0 Examples

### QFT on 4 qubits (`qft.qasm`)
```qasm
OPENQASM 2.0;
include "qelib1.inc";
qreg q[4];
creg c[4];
x q[0]; x q[2];       // prepare |1010⟩
barrier q;
h q[0];
cu1(pi/2) q[1],q[0];  // controlled phase = CR(π/2)
h q[1];
cu1(pi/4) q[2],q[0];
cu1(pi/2) q[2],q[1];
h q[2];
cu1(pi/8) q[3],q[0];
cu1(pi/4) q[3],q[1];
cu1(pi/2) q[3],q[2];
h q[3];
measure q -> c;
```
Note: uses `cu1` (= `NewCr`) for the controlled-phase rotations — exactly how quantum-go's `NewFourier` builds the QFT.

### Teleportation (`teleport.qasm`)
```qasm
OPENQASM 2.0;
include "qelib1.inc";
qreg q[3];
creg c0[1]; creg c1[1]; creg c2[1];
u3(0.3,0.2,0.1) q[0];   // prepare arbitrary state
h q[1];
cx q[1],q[2];            // Bell pair
barrier q;
cx q[0],q[1];            // Alice Bell measurement
h q[0];
measure q[0] -> c0[0];
measure q[1] -> c1[0];
if(c0==1) z q[2];        // Bob corrections (classically conditioned)
if(c1==1) x q[2];
measure q[2] -> c2[0];
```
Note: the `if(c==1)` classically-conditioned gates are **not** supported by quantum-go's parser — the parser handles the quantum portion only.

## Gates in qelib1.inc Not Supported by quantum-go Parser

| QASM gate | Reason not supported                                  |
| --------- | ----------------------------------------------------- |
| `u2`      | Could be added as `NewU(π/2, φ, λ, i)`               |
| `cy`      | No direct `NewCy`; would need to decompose            |
| `ch`      | No controlled-H gate                                  |
| `crz`     | No controlled-Rz gate                                 |
| `cu3`     | No controlled-U3 gate                                 |
| `if(...)`  | Classical control flow — not part of unitary model    |
| `barrier`  | Timing hint — explicitly ignored by parser            |

## Key Points

- Every QASM 2.0 file starts with `include "qelib1.inc"` — the parser accepts this line and knows the gate vocabulary.
- `u1(λ)` = `PhaseShift(λ)` = `Rz(λ)` up to a global phase — the parser maps `u1` to `NewPhaseShift`.
- `cu1(λ)` = `CR(λ)` — this is the gate used in the canonical QFT circuit. Verifiable in `qft.qasm` above.
- `ccx` (Toffoli) is mapped to quantum-go's native `NewToffoli`, not re-decomposed.
- The `if(classical_reg == value)` classically-conditioned gate is defined in QASM 2.0 but unsupported in quantum-go's parser — it falls outside the unitary circuit model.

## Sources

- `https://raw.githubusercontent.com/Qiskit/openqasm/OpenQASM2.x/examples/qelib1.inc` (Apache 2.0)
- `https://raw.githubusercontent.com/Qiskit/openqasm/OpenQASM2.x/examples/qft.qasm` (Apache 2.0)
- `https://raw.githubusercontent.com/Qiskit/openqasm/OpenQASM2.x/examples/teleport.qasm` (Apache 2.0)
- `qasm/parser.go`

## References

- Cross, A.W., Bishop, L.S., Smolin, J.A. & Gambetta, J.M. "Open Quantum Assembly Language." arXiv:[1707.03429](https://arxiv.org/abs/1707.03429) (2017) — defines the full QASM 2.0 grammar including `qelib1.inc`.
- OpenQASM 2.x repository: [github.com/Qiskit/openqasm](https://github.com/Qiskit/openqasm/tree/OpenQASM2.x) (Apache 2.0)
