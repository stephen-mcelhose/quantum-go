---
type: concept
title: Quantum Arithmetic — The Draper Adder
description: How the Draper adder performs addition in phase space using QFT + controlled phase rotations + IQFT, and why this is more efficient than classical ripple-carry addition.
resource: examples/arithmetic/arithmetic.md
tags: [arithmetic, draper-adder, qft, phase-space, addition, modular]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Arithmetic — The Draper Adder

The **Draper adder** performs quantum addition without classical-style carry chains. Instead of manipulating bits directly, it moves into *phase space* using the Quantum Fourier Transform, adds via phase rotations, then transforms back. This is the arithmetic engine behind [[shors-algorithm]].

## Classical vs Quantum Addition

| Approach   | Method                           | Carry handling              |
| ---------- | -------------------------------- | --------------------------- |
| Classical  | XOR/AND gate ripple-carry        | Explicit carry qubits       |
| Draper     | QFT → phase rotations → IQFT    | Phase accumulation (no ancilla) |

The Draper adder requires **no ancilla qubits for carries** — the phase of each qubit encodes carry information implicitly in the Fourier basis.

## How It Works

Three stages for computing |x, y⟩ → |x+y mod 2ⁿ, y⟩:

### Stage 1: QFT on x register
Transform the x register from computational basis to Fourier basis:
```
QFT|x⟩ = Σ e^{2πi·xk/2ⁿ}|k⟩  (phase encodes the value x)
```
The phase of each qubit in the x register now encodes the integer value x.

### Stage 2: Conditional Phase Rotations
For each bit j in y register that is 1, apply phase rotation Rₖ to each qubit in x register:
```
Rₖ = diag(1, e^{2πi/2ᵏ})
```
This adds e^{2πi·y/2ⁿ} to the phase — equivalent to adding y to x in the Fourier basis.

### Stage 3: IQFT on x register
The Inverse QFT transforms the accumulated phase back to a binary-encoded integer: x + y (mod 2ⁿ).

## Worked Example: 2 + 3 = 5 (mod 4 = 1)

```go
p := core.NewProgram(5)  // 5 qubits: x[0..1], y[2..4]

// Initialize x = 2 (|10⟩ in binary → flip qubit 1)
p.AddStep(core.NewStep(core.NewX(1)))

// Initialize y = 3 (|011⟩ → flip qubits 2 and 3)
p.AddStep(core.NewStep(core.NewX(2), core.NewX(3)))

// Add: x register stores x+y
adder := core.NewAdd(0, 1, 2, 4)  // NewAdd(x_start, x_end, y_start, y_end)
p.AddStep(core.NewStep(adder))
```

**Result:** `|01101⟩: 1.0000`

Reading little-endian: q₄q₃q₂ = 011 (y=3, unchanged), q₁q₀ = 01 (x register = 1).

Why 1? 2 + 3 = 5. The x register has 2 bits → max value 3 (2²=4). So 5 mod 4 = 1. **Modular arithmetic is natural** — the Fourier basis wraps around automatically.

## Modular Arithmetic

The Draper adder inherently computes addition modulo 2ⁿ (where n = size of x register). This is exactly what [[shors-algorithm]] needs: modular exponentiation aˣ mod N. The `core.NewMulModulus` gate (see [[arithmetic-gates]]) extends this to full modular multiplication, built from the Add primitive.

## Quantum Advantage

For an n-qubit register, the Draper adder uses O(n²) gate operations (from the O(n²/2) phase rotations in QFT). Classical ripple-carry also uses O(n) gates but requires O(n) ancilla qubits. The Draper approach is ancilla-free and maps naturally to the QFT's phase representation used throughout quantum-go.

See [[qft-deep-dive]] for a detailed breakdown of the QFT component.

## Key Points

- Draper adder = QFT(x) + controlled phase rotations (y bits as controls) + IQFT(x).
- No carry qubits — carries are implicit in phase accumulation.
- Result is x+y mod 2ⁿ, where n = size of x register.
- The core primitive is the controlled-R gate (CR) from [[quantum-concepts]].
- `core.NewAdd(x_start, x_end, y_start, y_end)` encapsulates the full QFT/add/IQFT as a Block gate.
- The add gate is the building block for `MulModulus`, which is the building block for Shor's.

## Sources

- `examples/arithmetic/arithmetic.md`

## References

- Draper, T.G. "Addition on a Quantum Computer." arXiv:[quant-ph/0008033](https://arxiv.org/abs/quant-ph/0008033) (2000) — introduces the QFT-based adder (Draper adder) implemented in `NewAdd`.
- Vedral, V., Barenco, A. & Ekert, A. "Quantum Networks for Elementary Arithmetic Operations." *Physical Review A* 54.1 (1996): 147. arXiv:[quant-ph/9511018](https://arxiv.org/abs/quant-ph/9511018) — earlier ripple-carry approach; Draper's QFT method is the alternative.
- Wikipedia: [Quantum arithmetic](https://en.wikipedia.org/wiki/Quantum_arithmetic)
