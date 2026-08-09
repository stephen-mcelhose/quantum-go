---
type: concept
title: Quantum Error Correction — 3-Qubit Bit-Flip Code
description: The simplest quantum error correcting code — how CNOT encoding, syndrome measurement, and Toffoli correction protect a logical qubit against single bit-flip errors.
resource: examples/algorithms/error_correction.md
tags: [error-correction, bit-flip-code, syndrome, toffoli, no-cloning, qec]
timestamp: 2026-08-09T03:26:15Z
---

# Quantum Error Correction — 3-Qubit Bit-Flip Code

Quantum states are fragile — thermal noise, measurement crosstalk, and cosmic rays can all flip qubits. Classical error correction copies bits; quantum mechanics forbids copying (No-Cloning Theorem). Instead, quantum error correction **spreads information across entangled qubits** so that errors can be detected without measuring the data.

## Why No-Cloning Matters

In classical systems: backup a bit by copying it. Done.

In quantum systems: you cannot clone an unknown state |ψ⟩. Measuring to check destroys the superposition. The No-Cloning Theorem (from [[entanglement]] background) forces a different approach: use entanglement to redundantly encode the state without copying it.

## The 3-Qubit Bit-Flip Code

Protects one logical qubit against a **single X (bit-flip) error** on any of the three physical qubits. Does **not** protect against Z (phase-flip) errors.

### Encoding

Map the logical qubit |ψ⟩ = α|0⟩ + β|1⟩ to three qubits:
```
|0⟩_L → |000⟩
|1⟩_L → |111⟩
```
so |ψ⟩_L = α|000⟩ + β|111⟩.

Implemented with two CNOT gates:
```go
// Assume q0 holds the logical qubit
p.AddStep(core.NewStep(core.NewCnot(0, 1)))  // q1 mirrors q0
p.AddStep(core.NewStep(core.NewCnot(0, 2)))  // q2 mirrors q0
```

**This is not copying** — the three qubits are entangled. The logical qubit's information is distributed non-locally.

### Error Simulation

An X gate on qubit 1 (for example) flips it:
```
α|010⟩ + β|101⟩   ← after error on q1
```
The majority vote is broken — q1 disagrees with q0 and q2.

### Syndrome Extraction

Reapply the same CNOT gates:
```go
p.AddStep(core.NewStep(core.NewCnot(0, 1)))
p.AddStep(core.NewStep(core.NewCnot(0, 2)))
```
After these, the ancilla qubits (q1, q2) store the **syndrome** — parity information that indicates which qubit was flipped, without revealing whether the logical qubit is |0⟩ or |1⟩.

| Syndrome (q1, q2) | Error location |
| ------------------ | -------------- |
| 00                 | No error       |
| 10                 | Error on q1    |
| 01                 | Error on q2    |
| 11                 | Error on q0    |

### Correction

A **Toffoli gate** (CCNOT) applies the fix:
```go
p.AddStep(core.NewStep(core.NewCCNot(1, 2, 0)))
```
If both q1 and q2 are |1⟩ (syndrome 11 → error on q0), Toffoli flips q0 back. Other syndromes are handled by additional Toffoli gates targeting q1 or q2 respectively.

## Limitations

| Error type    | Protected? | Explanation                                        |
| ------------- | ---------- | -------------------------------------------------- |
| X (bit-flip)  | ✓          | Detected by parity syndrome; corrected by Toffoli  |
| Z (phase-flip)| ✗          | Phase errors are invisible in computational basis   |
| Both          | ✗          | Requires Shor Code (9 qubits) or surface codes      |

To protect against both X and Z errors, encode the logical qubit in 9 physical qubits (Shor's 9-qubit code = bit-flip code on top of phase-flip code). The Shor code is the simplest code that protects against arbitrary single-qubit errors.

## In quantum-go

```go
// Encode bit '1'
program := core.NewErrorCorrectionProgram(1)
env := local.NewSimpleExecutionEnvironment()
result := env.RunProgram(program)
result.PrintBinary()
```

The CLI: `./quantum-go run --circuit error-correction -p bit=1`

## Key Points

- No-Cloning Theorem forbids copying → must use entanglement to distribute information.
- The 3-qubit code maps |0⟩_L → |000⟩ and |1⟩_L → |111⟩ (majority vote encoding).
- Syndrome measurement reads parity without measuring the logical qubit value.
- Toffoli gate (CCNOT) performs the conditional correction based on syndrome.
- Only protects against X errors — Z errors require a separate code (dual code structure).
- The Shor 9-qubit code = 3-qubit bit-flip code applied within a 3-qubit phase-flip code.
- Real quantum computers need fault-tolerant codes with hundreds to thousands of physical qubits per logical qubit.

## Sources

- `examples/algorithms/error_correction.md`

## References

- Shor, P.W. "Scheme for reducing decoherence in quantum computer memory." *Physical Review A* 52 (1995): R2493. DOI:[10.1103/PhysRevA.52.R2493](https://doi.org/10.1103/PhysRevA.52.R2493) — introduces the 9-qubit Shor code (bit-flip + phase-flip concatenation).
- Calderbank, A.R. & Shor, P.W. "Good quantum error-correcting codes exist." *Physical Review A* 54.2 (1996): 1098. arXiv:[quant-ph/9512032](https://arxiv.org/abs/quant-ph/9512032) — CSS codes, the general framework quantum-go's 3-qubit code belongs to.
- Steane, A.M. "Multiple Particle Interference and Quantum Error Correction." *Proceedings of the Royal Society A* 452 (1996). arXiv:[quant-ph/9601029](https://arxiv.org/abs/quant-ph/9601029) — Steane codes, the other branch of CSS.
- Wikipedia: [Quantum error correction](https://en.wikipedia.org/wiki/Quantum_error_correction)
