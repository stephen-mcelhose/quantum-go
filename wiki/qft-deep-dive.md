---
type: synthesis
title: QFT Deep Dive
description: Step-by-step construction and circuit analysis of the Quantum Fourier Transform — H gates, CR rotations at 2π/2ᵏ, bit-reversal SWAPs, inverse, and connections to QFT-dependent algorithms.
tags: [qft, fourier, cr-gates, bit-reversal, phase-accumulation, shor, arithmetic]
timestamp: 2026-08-09T03:26:15Z
---

# QFT Deep Dive

The Quantum Fourier Transform (QFT) is the quantum analogue of the Discrete Fourier Transform. It transforms between the **computational basis** (bit-strings) and the **Fourier basis** (phases). It is the common engine of Shor's algorithm, the Draper adder, phase estimation, and quantum signal processing.

## Mathematical Definition

The QFT maps computational basis states |j⟩ → Fourier basis states:

```
QFT|j⟩ = (1/√N) Σₖ e^{2πijk/N} |k⟩
```

where N = 2ⁿ. Each output state |k⟩ receives amplitude e^{2πijk/N}/√N — a phase that encodes the position j of the input in the Fourier frequency domain.

Applied to a superposition Σⱼ aⱼ|j⟩, the QFT produces Σₖ Âₖ|k⟩ where Âₖ is the DFT of the sequence (a₀, a₁, …, a_{N-1}).

## Circuit Construction

For n qubits, the QFT on qubit register [0..n-1] (see `core/arithmetic.go`):

```
For i from n-1 down to 0:
    H(i)                              ← Hadamard creates local superposition
    CR(i+1-j, i, 2π/2ʲ) for j=2..i+1  ← controlled phase from lower qubits

Then: SWAP(0, n-1), SWAP(1, n-2), ... ← bit-reversal permutation
```

### Step-by-Step for n=3 (qubits q₀, q₁, q₂)

**Starting from top qubit down:**

```
H(2) → [CR(1,2,π/2), CR(0,2,π/4)]
H(1) → [CR(0,1,π/2)]
H(0) →

Then: SWAP(0, 2)
```

Visual circuit diagram:
```
q₂ ──H──CR(π/2)──CR(π/4)──────────────×──
         │           │                  │
q₁ ──────●─────────H─────CR(π/2)──────┼──×──
                          │             │  │
q₀ ──────────────────────●─────────H───×──●──
```

### Why Controlled-R Gates at 2π/2ᵏ?

Each qubit j needs a phase e^{2πij/2^{n−i}} where j is the value of the lower qubits. The CR gate at angle 2π/2ᵏ contributes phase e^{iθ} to the |11⟩ state — when the lower qubit (control) is |1⟩, it contributes its "1 bit" to the accumulated phase. The angles are chosen so that:

- Bit position 1 (one below): CR(π/2) = e^{iπ/2} — contributes π/2 per unit
- Bit position 2 (two below): CR(π/4) — contributes π/4 per unit
- Bit position k (k below): CR(π/2^{k-1}) — contributes 2π/2^k per unit

This is exactly the structure of binary fraction representation.

### Bit-Reversal SWAP

After the H and CR steps, qubit j contains the j-th Fourier coefficient — but in **reversed** bit order compared to standard DFT. The final SWAP sequence reverses the qubit order:

```go
for i := 0; i < dim/2; i++ {
    block.AddStep(NewStep(NewSwap(i, dim-1-i)))
}
```

This corrects the bit-reversal without any phase changes.

## Implementation in quantum-go

From `NewFourier(dim, idx)` in `core/arithmetic.go`:

```go
for i := dim - 1; i >= 0; i-- {
    block.AddStep(NewStep(NewHadamard(i)))
    for j := 2; j <= i+1; j++ {
        theta := 2.0 * π / math.Pow(2.0, float64(j))
        block.AddStep(NewStep(NewCr(i+1-j, i, theta)))
    }
}
for i := 0; i < dim/2; i++ {
    block.AddStep(NewStep(NewSwap(i, dim-1-i)))
}
```

**CR gate args:** `NewCr(control, target, theta)`. The control is the lower-order qubit; the target is the current (higher) qubit being rotated.

**Gate count:** O(n²/2) CR gates + n H gates + n/2 SWAP gates = O(n²) total. Classical FFT needs O(n·2ⁿ) operations for the same transform. Quantum QFT is exponentially more efficient — but you cannot read out all 2ⁿ amplitudes without sampling.

## Inverse QFT

```go
invQFT := core.NewFourier(n, 0)
invQFT.SetInverse(true)
```

The `SetInverse(true)` flag causes `BlockGate.ApplyOptimize` to execute steps in **reverse order** with each gate also inverted. This correctly implements U† = QFT†.

The IQFT is used in [[shors-algorithm]] to extract the period from the accumulated phase, and in [[arithmetic-gates]] (`NewAdd`, `NewAddInteger`) to convert back from the Fourier basis after arithmetic.

## Key QFT Properties

| Property                    | Value                                            |
| --------------------------- | ------------------------------------------------ |
| Input basis                 | Computational (|j⟩ = binary states)              |
| Output basis                | Fourier (phases e^{2πijk/N})                    |
| QFT on |0…0⟩               | Uniform superposition [1/√N, …, 1/√N] (all real)|
| QFT unitarity               | QFT† × QFT = I (tested in [[fuzz-testing]])     |
| Circuit depth               | O(n²) gates, O(n²) CNOT gates                   |
| CR angle for k below        | 2π/2ᵏ                                           |
| Number of CRs for n qubits  | n(n-1)/2                                         |

## Algorithmic Uses

| Algorithm           | QFT role                                                   | Page                          |
| ------------------- | ---------------------------------------------------------- | ----------------------------- |
| Shor's              | IQFT at end — extracts period from phase-encoded register  | [[shors-algorithm]]           |
| Phase estimation    | IQFT on ancilla — extracts eigenphase of unitary           | (quantum chemistry extension) |
| Draper Adder        | QFT on x → phase addition → IQFT                          | [[quantum-arithmetic]]        |
| AddInteger          | QFT on register → classical phase shifts → IQFT           | [[arithmetic-gates]]          |

## Key Points

- QFT circuit: H on each qubit + controlled-phase rotations from lower qubits + bit-reversal SWAPs.
- CR angle for k-levels below: 2π/2ᵏ — captures bit k's contribution to the binary fraction.
- `NewFourier(dim, idx)` returns a BlockGate — it applies recursively via `ApplyOptimize`, never materializes the full matrix.
- `SetInverse(true)` on a Fourier gate = IQFT. The block executes its steps in reverse.
- QFT on |00⟩ = uniform superposition [0.5, 0.5, 0.5, 0.5] for n=2 — verified in [[verification-tests]].
- CR gates in QFT are the same `NewCr` as in [[rotation-implementations]] — just used at specific angles. CR is a controlled PhaseShift — see [[rotation-gates]] for the phase rotation geometry.

## References

- Coppersmith, D. "An approximate Fourier transform useful in quantum factoring." IBM Research Report RC19642 (1994); arXiv:[quant-ph/0201067](https://arxiv.org/abs/quant-ph/0201067) (uploaded 2002) — introduces the quantum circuit for the QFT.
- Draper, T.G. "Addition on a Quantum Computer." arXiv:[quant-ph/0008033](https://arxiv.org/abs/quant-ph/0008033) (2000) — demonstrates QFT as an arithmetic primitive; the Draper adder implemented in quantum-go.
- Nielsen, M.A. & Chuang, I.L. *Quantum Computation and Quantum Information*, Chapter 5 (quantum Fourier transform and phase estimation). Cambridge University Press, 2010.
- Wikipedia: [Quantum Fourier transform](https://en.wikipedia.org/wiki/Quantum_Fourier_transform)
