---
type: concept
title: BB84 Quantum Key Distribution
description: The Bennett-Brassard 1984 protocol for unconditionally secure key exchange — using conjugate bases, Heisenberg uncertainty, and sifting to create keys whose interception is physically detectable.
resource: examples/security/qkd.md
tags: [bb84, qkd, cryptography, conjugate-bases, heisenberg, sifting, security]
timestamp: 2026-08-09T03:26:15Z
---

# BB84 Quantum Key Distribution

BB84 (Bennett & Brassard, 1984) is the first quantum cryptography protocol. It creates a shared secret key between Alice and Bob whose secrecy is guaranteed by **physics**, not computational hardness. Any eavesdropping disturbs the quantum states in a detectable way.

## Why Quantum Cryptography?

Classical key exchange (RSA, Diffie-Hellman) relies on computational hardness — problems no classical computer can solve efficiently. But Shor's algorithm (see [[shors-algorithm]]) breaks RSA on a quantum computer. BB84 is immune to this threat because its security comes from Heisenberg's uncertainty principle, not math.

## Core Mechanism: Conjugate Bases

BB84 uses two **non-orthogonal** (conjugate) bases from [[quantum-concepts]]:

| Basis name     | Symbol | States              | Preparation         |
| -------------- | ------ | ------------------- | ------------------- |
| Rectilinear    | Z      | \|0⟩ and \|1⟩       | Identity or X gate  |
| Diagonal       | X      | \|+⟩ and \|−⟩       | H gate after Z prep |

Key property: measuring a Z-basis state with an X-basis detector (or vice versa) gives a **random result** — the measurement is maximally disturbing.

## The Protocol

### Alice's Preparation
For each bit:
1. Choose random bit (0 or 1)
2. Choose random basis (Z or X)
3. Prepare qubit accordingly:

```go
p := core.NewProgram(1)
if aliceBit == 1 {
    p.AddStep(core.NewStep(core.NewX(0)))      // flip to |1⟩
}
if aliceBasis == X_BASIS {
    p.AddStep(core.NewStep(core.NewHadamard(0))) // rotate to X basis
}
```

### Bob's Measurement
Bob randomly chooses his measurement basis (Z or X):
```go
if bobBasis == X_BASIS {
    p.AddStep(core.NewStep(core.NewHadamard(0))) // rotate back before measuring
}
result := engine.RunProgram(p)
bobBit := result.GetQubits()[0].Measure()
```

### Sifting (Classical Channel)
Alice and Bob publicly compare their **basis choices** (not the bits). They keep only bits where they used the **same basis** — these are perfectly correlated. Mismatched-basis bits are discarded.

Approximately 50% of bits survive sifting (expected basis match rate).

## Why Eavesdropping Is Detectable

If Eve intercepts and measures a qubit, she must choose a basis — she has a 50% chance of choosing wrong. When she re-sends, she destroys the original superposition. This introduces errors even when Alice and Bob chose matching bases.

**Effect of eavesdropping on error rate:**
- No eavesdropper: after sifting, Alice and Bob's bits are 100% correlated.
- With eavesdropper: ~25% of sifted bits are wrong (Eve's basis mismatch × 50%).

Alice and Bob detect Eve by comparing a sample of their sifted bits over the classical channel. QBER (Quantum Bit Error Rate) above ~11% indicates eavesdropping.

## Information-Theoretic Security

BB84's security is **unconditional** — it doesn't depend on Eve being computationally limited. Even an infinitely powerful classical or quantum computer cannot break it. The Heisenberg uncertainty principle makes it physically impossible to measure both Z and X basis states simultaneously.

This contrasts with RSA (broken by Shor's) and AES (weakened by Grover's, but not broken — key size doubles).

## Simulation in quantum-go

```
Alice bit=1, basis=Z → X gate → qubit is |1⟩
Bob basis=Z → no H → measures |1⟩ deterministically
Bases match → keep bit. Shared key bit = 1.
```

If Bob chose X basis instead:
```
Bob H gate → H|1⟩ = |−⟩ → measures 0 or 1 randomly
Result discarded (bases differ) during sifting.
```

## Key Points

- BB84 uses two conjugate bases (Z = computational, X = Hadamard/diagonal).
- Alice randomly mixes Z and X basis encodings; Bob randomly mixes Z and X measurements.
- Sifting: only keep bits where both chose the same basis (~50% kept).
- Eavesdropping introduces ~25% error rate in sifted bits — statistically detectable.
- Security is information-theoretic (Heisenberg), not computational — immune to quantum computers.
- The [[quantum-concepts]] Z vs X basis table is the direct foundation for understanding BB84 security.
- Contrast with RSA: [[shors-algorithm]] breaks RSA; BB84 is quantum-safe by design.

## Sources

- `examples/security/qkd.md`

## References

- Bennett, C.H. & Brassard, G. "Quantum cryptography: Public key distribution and coin tossing." *Proceedings of IEEE International Conference on Computers, Systems and Signal Processing*, Bangalore, India (1984): 175–179. — The original BB84 paper; predates arXiv. Republished in *Theoretical Computer Science* 560 (2014): 7–11. DOI:[10.1016/j.tcs.2014.05.025](https://doi.org/10.1016/j.tcs.2014.05.025)
- Wikipedia: [BB84](https://en.wikipedia.org/wiki/BB84)
- Wikipedia: [Quantum key distribution](https://en.wikipedia.org/wiki/Quantum_key_distribution)
