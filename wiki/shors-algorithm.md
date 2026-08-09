---
type: concept
title: Shor's Algorithm
description: Period finding via QFT as an interference engine — how Shor's reduces integer factoring (exponentially hard classically) to period finding (efficient quantumly), breaking RSA.
resource: examples/algorithms/shor.md
tags: [shor, factoring, period-finding, qft, rsa, modular-exponentiation, interference]
timestamp: 2026-08-09T03:26:15Z
---

# Shor's Algorithm

Shor's (1994) is the most consequential quantum algorithm — it breaks RSA encryption in polynomial time. The classical best algorithm (General Number Field Sieve) is sub-exponential in the key size; Shor's is O(log³ N).

**Core insight:** Factoring reduces to period finding, and period finding can be solved efficiently with the Quantum Fourier Transform.

## The Reduction: Factoring → Period Finding

Given N to factor, pick random a with 1 < a < N. The function f(x) = aˣ mod N is periodic with period r:
```
a^r ≡ 1 (mod N)
```
Once r is known, compute:
```
gcd(a^{r/2} ± 1, N)
```
With high probability, this yields a non-trivial factor of N. The quantum part — finding r — is where the speedup lives.

## Why Classical Period Finding Is Hard

To find r for aˣ mod N classically, you'd evaluate f(x) for x = 0, 1, 2, … until you see a repeat. The period r can be as large as N, so this is O(N) — exponential in the number of bits of N.

## The Quantum Algorithm

### Step 1: Quantum Parallelism
Put the precision register (x) into uniform superposition — simultaneously evaluating f for all x at once:
```go
for i := 0; i < offset; i++ {
    p.AddStep(core.NewStep(core.NewHadamard(i)))
}
```
After this, the state is Σ|x⟩|0⟩ for all x.

### Step 2: Modular Exponentiation
Apply f(x) = aˣ mod N to the result register for all x simultaneously:
```go
mul := core.NewMulModulus(offset, offset+length-1, m, mod)
cbg := core.NewControlledBlockGate(mul, i)
p.AddStep(core.NewStep(cbg))
```
The state becomes Σ|x⟩|aˣ mod N⟩ — a superposition encoding f for all inputs.

### Step 3: The QFT as Interference Engine
The Inverse QFT transforms the precision register so that amplitudes at multiples of 2^precision/r interfere constructively (peaks), while all others destructively cancel:
```go
invQFT := core.NewFourier(offset, 0)
invQFT.SetInverse(true)
p.AddStep(core.NewStep(invQFT))
```

**Why this works:** The state after measuring the result register collapses to a superposition of |x⟩ values with the same f value — a periodic comb with spacing r. The QFT extracts the frequency (1/r) of this comb. Peaks appear at multiples of 2^{precision}/r.

### Step 4: Classical Post-Processing
Measure peaks at positions k ≈ j · 2^{precision}/r. Use continued fractions to extract r from k/2^{precision} ≈ j/r. Then compute gcd(a^{r/2}±1, N).

## Example: Factoring via Period of 2ˣ mod 7

The period of 2ˣ mod 7 is r=3 (because 2³=8≡1 mod 7).

With 3 bits of precision, peaks appear at: 0, 8/3≈2.6 (binary `010`), 16/3≈5.3 (binary `110`).

```
|0000001000>: 0.1153
|0000001010>: 0.0981
|0000001110>: 0.0731
```

## Complexity

| Approach          | Complexity for factoring N    |
| ----------------- | ----------------------------- |
| GNFS (classical)  | O(e^{1.9(log N)^{1/3}(log log N)^{2/3}}) |
| Shor's (quantum)  | O((log N)³)                   |

## Resource Requirements

For an L-bit number N:
- Precision register: ~L+logL qubits
- Result register: ~2L qubits
- Total: ~3L+ qubits

N=15 (L=4): ~15 qubits — simulatable. N=2048 (RSA): ~6000+ qubits — beyond any current quantum hardware. See [[scaling-and-limits]] for the simulation ceiling.

## Gate Architecture in quantum-go

```
Shor circuit = H⊗offset ⊗ MulModulus gates ⊗ IQFT
                              ↑
                    Uses [[arithmetic-gates]] (Add + MulModulus)
                    Which uses [[composite-gates]] (QFT as Block)
                    Which uses [[gate-implementations]] (H, CR gates)
```

## Key Points

- Period finding is the quantum core; factoring is classical reduction around it.
- Modular exponentiation is the expensive part — it uses [[quantum-arithmetic]] blocks repeatedly.
- The QFT is the "frequency analyzer" that extracts the period from quantum amplitudes.
- Measuring peaks at multiples of 2^precision/r allows continued-fraction reconstruction of r.
- For large N, Shor's is both memory-intensive (3L qubits) and depth-intensive — see [[scaling-and-limits]].
- See [[algorithm-comparison]] for Shor's vs Grover's and other algorithms.

## Sources

- `examples/algorithms/shor.md`
- `examples/algorithms/shor_test.go`

## References

- Shor, P.W. "Polynomial-Time Algorithms for Prime Factorization and Discrete Logarithms on a Quantum Computer." arXiv:[quant-ph/9508027](https://arxiv.org/abs/quant-ph/9508027) (1995) — expanded version of the FOCS 1994 paper "Algorithms for quantum computation: discrete logarithms and factoring."
- Wikipedia: [Shor's algorithm](https://en.wikipedia.org/wiki/Shor%27s_algorithm)
- Wikipedia: [Quantum phase estimation](https://en.wikipedia.org/wiki/Quantum_phase_estimation_algorithm) — the sub-routine at the core of Shor's period-finding.
- Wikipedia: [Continued fractions](https://en.wikipedia.org/wiki/Continued_fraction) — the classical post-processing step to recover the period r from the measured frequency.
