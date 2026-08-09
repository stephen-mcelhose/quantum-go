---
type: synthesis
title: Gate Zoo
description: Complete reference table of every gate in quantum-go — constructor, matrix, QASM mnemonic, size (qubits), caption, and which wiki page covers it in depth.
tags: [gate-zoo, reference, hadamard, pauli, rotation, cnot, toffoli, qft, oracle]
timestamp: 2026-08-09T03:26:15Z
---

# Gate Zoo

Complete reference of all gates implemented in quantum-go. Organized by category. For implementation details, see the linked wiki pages.

## Single-Qubit Fixed Gates

| Gate        | Constructor          | Caption | Matrix (2×2)             | QASM     | Details                    |
| ----------- | -------------------- | ------- | ------------------------ | -------- | -------------------------- |
| Hadamard    | `NewHadamard(i)`     | `H`     | [[1,1],[1,-1]] / √2     | `h`      | [[gate-implementations]]   |
| Pauli-X     | `NewX(i)`            | `X`     | [[0,1],[1,0]]            | `x`      | [[gate-implementations]]   |
| Pauli-Y     | `NewY(i)`            | `Y`     | [[0,-i],[i,0]]           | `y`      | [[gate-implementations]]   |
| Pauli-Z     | `NewZ(i)`            | `Z`     | [[1,0],[0,-1]]           | `z`      | [[gate-implementations]]   |
| Identity    | `NewIdentity(i)`     | `I`     | [[1,0],[0,1]]            | `id`     | [[gate-implementations]]   |
| Measurement | `NewMeasurement(i)`  | `M`     | Identity (special)       | `measure`| [[simulation-engine]]      |

## Single-Qubit Rotation Gates

| Gate         | Constructor                      | Caption | Matrix                          | QASM              | Details                         |
| ------------ | -------------------------------- | ------- | ------------------------------- | ----------------- | ------------------------------- |
| S / S†       | `NewS(i)` `.SetInverse(true)`    | `S`     | diag(1, i) / diag(1, -i)       | `s` / `sdg`       | [[rotation-implementations]]   |
| T / T†       | `NewT(i)` `.SetInverse(true)`    | `T`     | diag(1, e^{iπ/4}) / conj        | `t` / `tdg`       | [[rotation-implementations]]   |
| V (SX / √X) | `NewV(i)`                        | `V`     | (1+i)/2 [[1,-i],[-i,1]]        | `sx`              | [[rotation-implementations]]   |
| PhaseShift   | `NewPhaseShift(θ, i)`            | `PS`    | diag(1, e^{iθ})                 | `u1(θ)`           | [[rotation-implementations]]   |
| Rx(θ)        | `NewRx(θ, i)`                    | `Rx`    | Bloch X rotation                | `rx(θ)`           | [[rotation-implementations]]   |
| Ry(θ)        | `NewRy(θ, i)`                    | `Ry`    | Bloch Y rotation (real)         | `ry(θ)`           | [[rotation-implementations]]   |
| Rz(θ)        | `NewRz(θ, i)`                    | `Rz`    | diag(e^{-iθ/2}, e^{+iθ/2})    | `rz(θ)`           | [[rotation-implementations]]   |
| U (Universal)| `NewU(θ,φ,λ, i)`                 | `U`     | 3-param Euler rotation          | `u3(θ,φ,λ)`       | [[universality]]               |

## Two-Qubit Gates

| Gate    | Constructor          | Caption  | Matrix    | QASM             | Details                      |
| ------- | -------------------- | -------- | --------- | ---------------- | ---------------------------- |
| CNOT    | `NewCnot(ctrl, tgt)` | `CNOT`   | 4×4       | `cx`             | [[gate-implementations]]     |
| CZ      | `NewCz(ctrl, tgt)`   | `CZ`     | diag(1,1,1,-1) | `cz`      | [[gate-implementations]]     |
| SWAP    | `NewSwap(q1, q2)`    | `SWAP`   | 4×4       | `swap`           | [[gate-implementations]]     |
| CR(θ)   | `NewCr(ctrl,tgt, θ)` | `CR`     | diag(1,1,1,e^{iθ}) | `cu1(θ)` | [[rotation-implementations]] |

## Three-Qubit Gates

| Gate           | Constructor               | Caption  | Matrix  | QASM    | Details                      |
| -------------- | ------------------------- | -------- | ------- | ------- | ---------------------------- |
| Toffoli (CCNOT)| `NewToffoli(a, b, c)`     | `CCNOT`  | 8×8     | `ccx`   | [[gate-implementations]]     |
| Fredkin (CSWAP)| `NewFredkin(ctrl, b, c)`  | `CSWAP`  | 8×8     | —       | [[gate-implementations]]     |

## Composite / Structural Gates

| Gate                 | Constructor                          | Caption | Qubits   | Details                   |
| -------------------- | ------------------------------------ | ------- | -------- | ------------------------- |
| QFT                  | `NewFourier(dim, idx)`               | `QFT`   | dim      | [[arithmetic-gates]]      |
| IQFT                 | `NewFourier(dim,idx).SetInverse(true)` | `QFT` | dim      | [[arithmetic-gates]]      |
| Add                  | `NewAdd(x0,x1,y0,y1)`               | `ADD`   | y1-x0+1  | [[arithmetic-gates]]      |
| AddInteger           | `NewAddInteger(x0,x1, num)`          | `ADDI`  | m        | [[arithmetic-gates]]      |
| AddIntegerModulus    | `NewAddIntegerModulus(x0,x1, a, N)`  | `ADDIM` | m+1      | [[arithmetic-gates]]      |
| MulModulus           | `NewMulModulus(x0,x1, mul, mod)`     | `MULM`  | 2m+1     | [[arithmetic-gates]]      |
| Oracle (custom)      | `NewOracle(idx, matrix)`             | `Oracle`| log₂(dim)| [[composite-gates]]       |
| BlockGate            | `NewBlockGate(block, idx)`           | `B`     | block.n  | [[composite-gates]]       |
| ControlledGate       | `NewControlledGate(g, ctrl...)`      | `C+`    | varies   | [[composite-gates]]       |
| ControlledBlockGate  | `NewControlledBlockGate(bg, ctrl)`   | `CB`    | varies   | [[composite-gates]]       |
| Permutation          | `NewPermutationGate(t1, t2, nq)`     | `P`     | nq       | [[composite-gates]]       |
| SingleQubitMatrix    | `NewSingleQubitMatrixGate(i, m)`     | `M`     | 1        | [[composite-gates]]       |
| TimeEvolution        | `NewTimeEvolution(idx, H, t)`        | `U(t)`  | n        | [[quantum-thermodynamics]]|

## Pre-Built Oracles

| Oracle               | Constructor                      | Input dim | Oracle type       | Details             |
| -------------------- | -------------------------------- | --------- | ----------------- | ------------------- |
| ConstantOracle       | `NewConstantOracle(n, value)`    | n+1 bits  | f(x) = const      | [[oracle-gates]]    |
| BalancedOracle       | `NewBalancedOracle(n)`           | n+1 bits  | f(x) = x₀        | [[oracle-gates]]    |
| InnerProductOracle   | `NewInnerProductOracle(s)`       | n+1 bits  | f(x) = s·x mod 2 | [[oracle-gates]]    |
| SimonOracle          | `NewSimonOracle(s)`              | 2n bits   | f(x)=f(x⊕s)      | [[oracle-gates]]    |

## Gate Size Reference

Gate "size" = number of qubits in `AffectedQubits`. The state space for a gate of size k is 2^k × 2^k.

| Size  | State space | Example gates                    |
| ----- | ----------- | -------------------------------- |
| 1     | 2×2         | H, X, Y, Z, S, T, V, Rx, Ry, Rz |
| 2     | 4×4         | CNOT, CZ, SWAP, CR               |
| 3     | 8×8         | Toffoli, Fredkin                 |
| n     | 2^n × 2^n   | QFT, Add, MulModulus, Oracle     |

## Key Facts for Implementation

- Caption is the routing key in [[gate-application]] — do not rename without updating dispatch logic.
- All gates support `SetInverse(true)` / `IsInverse()` — checked in each `GetMatrix()`.
- Single-qubit gates: engine uses `GetMatrix()` → applies 2×2 to paired amplitudes.
- Multi-qubit fixed gates: engine uses optimized bitwise paths (CNOT XOR, CZ conditional phase).
- BlockGates: engine calls `ApplyOptimize()` — never calls `GetMatrix()`.
- Adding a new gate: see [[how-to-add-a-new-gate]].
