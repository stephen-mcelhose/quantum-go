---
type: concept
title: Oracle Gate Implementations
description: core/oracles.go — ConstantOracle, BalancedOracle, InnerProductOracle, SimonOracle — how quantum black-box functions are implemented as permutation matrices on (input|output) register pairs.
resource: core/oracles.go
tags: [oracle, deutsch-jozsa, bernstein-vazirani, simon, permutation-matrix, ancilla]
timestamp: 2026-08-09T03:26:15Z
---

# Oracle Gate Implementations

`core/oracles.go` provides pre-built oracle circuits for standard quantum algorithm demonstrations. All use the same pattern: a permutation matrix on a (data|ancilla) register pair implementing `|x⟩|y⟩ → |x⟩|y ⊕ f(x)⟩`. This is the quantum phase-kickback pattern from [[quantum-concepts]].

## The Oracle Pattern

All oracles implement the "standard form" unitary:
```
U_f |x⟩|y⟩ = |x⟩|y ⊕ f(x)⟩
```
where ⊕ is XOR. When y = |−⟩ = (|0⟩−|1⟩)/√2 (prepared by X then H), the phase kickback effect marks the input state with a phase flip:
```
U_f |x⟩|−⟩ = (−1)^{f(x)} |x⟩|−⟩
```
This is how the oracle in [[grovers-algorithm]] and the Deutsch-Jozsa algorithm applies a phase without measuring.

## ConstantOracle

```go
func NewConstantOracle(n int, value int) *Oracle
```

f(x) = value (constant, 0 or 1 for all x).

**Matrix construction** for n input qubits + 1 output qubit (dim = 2^{n+1}):
```go
for x := 0; x < (1 << n); x++ {
    for y := 0; y < 2; y++ {
        input := (y << n) | x
        output := ((y ^ value) << n) | x
        m.Set(output, input, 1)   // permutation matrix entry
    }
}
```

- If `value=0`: identity (f never flips the output).
- If `value=1`: flips all output qubits regardless of input.
- Deutsch-Josza algorithm: ≥1 call to ConstantOracle vs BalancedOracle to determine type.

## BalancedOracle

```go
func NewBalancedOracle(n int) *Oracle
```

f(x) = x₀ (LSB of x). Exactly half of all inputs map to 0, half to 1.

```go
fx := x & 1  // just the least significant bit of x
input := (y << n) | x
output := ((y ^ fx) << n) | x
```

This is the canonical "balanced" oracle — exactly 2^{n-1} inputs produce f(x)=0 and 2^{n-1} produce f(x)=1. Deutsch-Josza can distinguish ConstantOracle from BalancedOracle with a single query.

## InnerProductOracle (Bernstein-Vazirani)

```go
func NewInnerProductOracle(s string) *Oracle
```

f(x) = s·x (mod 2) — inner product of x with hidden string s (dot product in GF(2)).

**String interpretation:** `s` is a binary string like `"101"`. The function computes `f(x) = s₀x₀ ⊕ s₁x₁ ⊕ … ⊕ s_{n-1}x_{n-1}`.

```go
sBits := make([]int, n)
for i, r := range s {
    if r == '1' { sBits[n-1-i] = 1 }  // LSB of string is at index 0
}
```

Bernstein-Vazirani algorithm finds `s` with a single oracle query (vs O(n) classical queries).

**Little-endian note:** `sBits[n-1-i]` reverses the string — the rightmost character of `s` is sBit[0], matching the qubit index convention from [[project-overview]].

## SimonOracle (Simon's Algorithm)

```go
func NewSimonOracle(s string) *Oracle
```

f(x) = f(y) iff y = x ⊕ s — "two-to-one" function with hidden XOR period s.

The implementation uses a canonical construction:
- If x < x⊕s: f(x) = x
- If x > x⊕s: f(x) = x⊕s

```go
sVal, _ := strconv.ParseInt(s, 2, 64)
fx := x
if int64(x ^ int(sVal)) < int64(x) {
    fx = x ^ int(sVal)
}
```

This ensures consistent pairing: f(x) = f(x⊕s) for all x. Simon's algorithm finds s using O(n) quantum queries instead of O(√2ⁿ) classical.

Note: This oracle maps (2n)-qubit space (input register `x` + output register `y`), so `dim = 2^{2n}`. The output register gets f(x) XOR'd in.

## Comparison

| Oracle                | Hidden structure    | Quantum advantage           | Algorithm               |
| --------------------- | ------------------- | --------------------------- | ----------------------- |
| ConstantOracle        | f is constant       | 1 query vs 2^{n-1}+1        | Deutsch-Josza           |
| BalancedOracle        | f is balanced       | 1 query vs 2^{n-1}+1        | Deutsch-Josza           |
| InnerProductOracle    | Hidden string s     | 1 query vs n queries        | Bernstein-Vazirani      |
| SimonOracle           | Hidden XOR period s | O(n) queries vs O(√2ⁿ)     | Simon's Algorithm       |

## Key Points

- All oracles use the `|y ⊕ f(x)⟩` pattern — permutation matrix on (input ⊗ ancilla).
- Phase kickback: when ancilla = |−⟩, the output is `(−1)^{f(x)}|x⟩|−⟩` — phase is imprinted on input.
- BalancedOracle uses f(x)=x₀ (LSB) — the simplest balanced function.
- InnerProductOracle's string indexing is reversed (LSB convention) — sBits[n-1-i].
- SimonOracle maps 2n qubits (output register is also n bits).
- For custom oracles, use `core.NewOracle(idx, matrix)` directly with a manually constructed unitary.

## Sources

- `core/oracles.go`
