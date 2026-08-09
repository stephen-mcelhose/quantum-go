---
type: synthesis
title: Algorithm Comparison
description: Side-by-side comparison of all quantum algorithms in quantum-go — query complexity, qubit count, key gates, circuit depth, and what they solve.
tags: [algorithms, comparison, complexity, grover, shor, deutsch-josza, bernstein-vazirani, simon, qkd, error-correction]
timestamp: 2026-08-09T03:26:15Z
---

# Algorithm Comparison

Concise side-by-side comparison of every quantum algorithm implemented in quantum-go.

## Complexity Overview

| Algorithm              | Problem              | Classical complexity | Quantum complexity | Speedup type              |
| ---------------------- | -------------------- | ------------------- | ------------------ | ------------------------- |
| Deutsch-Jozsa          | Balanced vs constant | O(2^{n-1}+1)        | O(1) queries       | Exponential (exact)       |
| Bernstein-Vazirani     | Find hidden string s | O(n) queries        | O(1) query         | Polynomial (exact)        |
| Simon's Algorithm      | Find XOR period      | O(2^{n/2})          | O(n) queries       | Exponential               |
| Grover's Search        | Unstructured search  | O(N)                | O(√N)              | Quadratic                 |
| Shor's Algorithm       | Integer factoring    | Sub-exponential     | O((log N)³)        | Exponential               |
| BB84 QKD               | Key distribution     | Insecure            | Unconditionally secure | Security proof      |
| Quantum Teleportation  | State transfer       | Impossible          | 2 classical bits + ebit | Communication     |

## Resource Requirements

| Algorithm              | Qubits           | Circuit depth      | Key gates used                    | Factory function                |
| ---------------------- | ---------------- | ------------------ | --------------------------------- | ------------------------------- |
| Deutsch-Jozsa (n-bit)  | n+1              | O(1) steps         | H, Oracle, H                      | `NewDeutschJozsaProgram(n, b)`  |
| Bernstein-Vazirani     | n+1              | O(1) steps         | X, H, InnerProductOracle, H       | `NewBernsteinVaziraniProgram(s)`|
| Simon's (n-bit)        | 2n               | O(1) steps         | H, SimonOracle, H                 | `NewSimonsProgram(s)`           |
| Grover (2-qubit)       | 2                | 3 steps            | H⊗H, Oracle(4×4), Diffusion(4×4) | `NewGroverProgram()`            |
| Shor's (mod N, prec p) | 2⌈log₂N⌉+1+p    | O(p) steps         | H, X, MulModulus, IQFT            | `NewShorProgram(a, N, p)`       |
| BB84 (1 bit)           | 1 + 1 (Bob)      | 1-2 steps          | X, H (conditionally)              | `NewQKDProgram(bit, basis)`     |
| Teleportation          | 3                | 7 steps            | H, CNOT, CNOT, H, CNOT, CZ        | `NewTeleportationProgram()`     |
| Superdense Coding      | 2                | 6 steps            | H, CNOT, X, Z, CNOT, H           | `NewSuperdenseCodingProgram()`  |
| Error Correction (3-bit)| 3               | 7 steps            | X, CNOT, Toffoli                  | `NewErrorCorrectionProgram(bit)`|
| QFT (n-qubit)          | n                | O(n²)              | H, CR, SWAP                       | `NewQFTProgram(n)`              |

## What Each Algorithm Produces

| Algorithm              | Output / What you read               | Key readout               |
| ---------------------- | ------------------------------------- | ------------------------- |
| Deutsch-Jozsa          | Qubits 0..n-1: all |0⟩ → constant, any |1⟩ → balanced | `Measure()` each qubit |
| Bernstein-Vazirani     | Qubits 0..n-1: decode bit-string = s | `Measure()` each qubit |
| Simon's                | n equations s·y = 0 mod 2 per run; solve classically | `Measure()` + linear algebra |
| Grover                 | Winner state with P ≈ 1 (O(√N) iterations) | `GetProbability()` highest amplitude |
| Shor's                 | Precision register: peaks at multiples of N/r | `GetProbability()`, extract peak → continued fractions |
| BB84                   | 1-qubit state encoding classical bit/basis | `Measure()` on receiver qubit |
| Teleportation          | q₂ now holds Alice's state           | `GetQubits()[2]`           |
| Superdense Coding      | State encodes 2 classical bits        | `Measure()` q₀ and q₁     |
| Error Correction       | q₀ = original logical bit (corrected)| `GetQubits()[0].Measure()` |
| QFT                    | Fourier-transformed amplitude vector  | `GetProbability()` (full)  |

## Oracle Implementation Patterns

| Algorithm              | Oracle gate                      | Matrix encoding                                      |
| ---------------------- | -------------------------------- | ---------------------------------------------------- |
| Deutsch-Jozsa          | `NewBalancedOracle` / `NewConstantOracle` | Hardcoded CNOT-style; bottom qubit = ancilla |
| Grover (2-qubit)       | Custom 4×4 matrix via `NewOracle`| diag(1,1,1,-1) marks |11⟩                           |
| Bernstein-Vazirani     | `NewInnerProductOracle(s)`       | XOR s·x into ancilla qubit                          |
| Simon's                | `NewSimonOracle(s)`              | f(x) = f(x⊕s) in 2n-qubit block                   |

For oracle gate construction and the full pre-built oracle catalogue, see [[oracle-gates]].

## Phase Kickback (Used By All Query Algorithms)

All query algorithms from DJ through Grover use **phase kickback**:

1. Prepare ancilla qubit in |−⟩ = (|0⟩ − |1⟩)/√2 via X+H
2. Oracle acts on input ⊗ |−⟩: flips ancilla iff f(x) = 1
3. Phase kickback: state |x⟩|−⟩ → (−1)^{f(x)} |x⟩|−⟩

The oracle's effect is written as a phase on the input qubit — ancilla state unchanged. This is the mechanism that allows f(x) to be extracted from a superposition.

## Decision Guide

| You want to...                           | Use                                    |
| ---------------------------------------- | -------------------------------------- |
| Distinguish constant vs balanced function| Deutsch-Jozsa                          |
| Find a hidden linear boolean function    | Bernstein-Vazirani                     |
| Find XOR period of a 2-to-1 function     | Simon's Algorithm                      |
| Search N items for 1 winner              | Grover (√N speedup)                    |
| Factor a large integer                   | Shor's (exponential speedup)           |
| Distribute secure cryptographic keys     | BB84 QKD                               |
| Transfer quantum state across space      | Teleportation                          |
| Transmit 2 bits using 1 qubit            | Superdense Coding                      |
| Protect against bit-flip errors          | 3-qubit Error Correction               |
| Add two quantum integers                 | Draper Adder (NewAdd)                  |

## Key Points

- Grover requires O(√N) *iterations* — `NewGroverProgram()` does 1 iteration for n=2 (optimal for N=4). See [[grovers-algorithm]].
- Shor's requires repeated quantum runs + classical continued fractions to extract the period; the quantum part outputs frequencies, not factors directly. See [[shors-algorithm]].
- Simon's algorithm is BQP-complete relative to the Simon oracle — it is the quantum algorithm that "proves" exponential speedup.
- BB84 security is information-theoretic, not computational — not vulnerable to quantum computers. See [[bb84-qkd]].
- Teleportation requires a pre-shared Bell pair AND 2 bits of classical communication per qubit transferred. See [[teleportation]].
- Error correction, arithmetic, and QFT deep-dives: [[error-correction]], [[arithmetic-gates]], [[qft-deep-dive]].
- All factory functions (`NewGroverProgram`, `NewShorProgram`, etc.) are defined in [[circuits-library]].
