# Learning Quantum Key Distribution (BB84) with Strange-Go

Quantum Key Distribution (QKD) is a secure communication method which implements a cryptographic protocol involving components of quantum mechanics. It enables two parties to produce a shared random secret key known only to them, which can then be used to encrypt and decrypt messages.

The most famous QKD protocol is **BB84**, named after its inventors Charles Bennett and Gilles Brassard (1984).

---

## 1. The Challenge: Secure Key Exchange

In classical cryptography, keys are often exchanged using mathematical problems that are hard to solve. However, a powerful enough computer (or a quantum computer running Shor's algorithm) could potentially break these.

- **Classical approach**: Rely on mathematical complexity (RSA, Elliptic Curves).
- **Quantum approach**: Rely on the laws of physics. Any attempt to eavesdrop on a quantum key will disturb the state, revealing the presence of the intruder.

*   **Learn more about the Theory**: [Quantum Key Distribution (Wikipedia)](https://en.wikipedia.org/wiki/Quantum_key_distribution)

---

## 2. Core Concepts

### Conjugate Bases
BB84 uses two \"non-orthogonal\" bases to encode information:
1.  **Rectilinear Basis (Z)**: $|0\rangle$ and $|1\rangle$.
2.  **Diagonal Basis (X)**: $|+\rangle = \frac{|0\rangle + |1\rangle}{\sqrt{2}}$ and $|-\rangle = \frac{|0\rangle - |1\rangle}{\sqrt{2}}$.

> **Hilbert Space Context**: The Z and X bases are rotated by $45^\circ$ relative to each other in the 2D Hilbert Space of a single qubit. Measuring a state prepared in one basis using the other basis results in a maximal uncertainty, a property we exploit for security.

### Heisenberg Uncertainty Principle
In the context of QKD, this principle implies that measuring one property of a quantum system (like the Z basis state) necessarily disturbs another property (like the X basis state). If an eavesdropper (Eve) tries to measure the key, she will change the qubits in a way that Alice and Bob can detect.

### Sifting
Alice and Bob only keep the bits where they happened to choose the **same basis**. This process is called \"sifting\" and is done over a classical (non-secure) channel.

---

## 3. Walkthrough: BB84 Simulation

In `go/examples/security/qkd_test.go`, we simulate a single bit exchange where Alice and Bob successfully match their bases.

### Step 1: Alice Prepares the Qubit
Alice chooses a secret bit (`1`) and a basis (`Z`).

```go
p := core.NewProgram(1)
// If bit is 1, apply X gate to flip |0> to |1>
if aliceBit == 1 {
    p.AddStep(core.NewStep(core.NewX(0)))
}
// If basis is X (Diagonal), apply Hadamard
if aliceBasis == 1 {
    p.AddStep(core.NewStep(core.NewHadamard(0)))
}
```

### Step 2: Bob Chooses a Basis
Bob chooses a basis to measure in. In this deterministic example, he chooses the same as Alice.

```go
// If Bob chooses X basis, he applies Hadamard before measuring
if bobBasis == 1 {
    p.AddStep(core.NewStep(core.NewHadamard(0)))
}
```

### Step 3: Measurement and Comparison
Bob measures the qubit and then compares his basis choice with Alice over a classical channel.

```go
result := engine.RunProgram(p)
bobBit := result.GetQubits()[0].Measure()

if aliceBasis == bobBasis {
    fmt.Println(\"Bases Match! Shared Key Bit Found.\")
}
```

---

## 4. Interpreting Results

```
Alice's Secret Bit: 1, Basis: Z (|0>, |1>)
Bob chooses Basis: Z (|0>, |1>)
Bob measures Bit: 1
Bases Match! Shared Key Bit Found.
SUCCESS: Shared bit matches.
```

If Bob had chosen the X basis while Alice used Z, he would have measured `0` or `1` with equal probability, and they would have discarded the bit during the sifting phase.

---

## 5. Self-Assessment

1. **What happens if an eavesdropper measures the qubit in the wrong basis?**
   - [ ] They get the secret bit perfectly.
   - [ ] They disturb the quantum state, causing Bob to have a 25% error rate even when his basis matches Alice's.

2. **Why do Alice and Bob discard bits where their bases don't match?**
   - [ ] Because the hardware failed.
   - [ ] Because the results in different bases are uncorrelated and random.

3. **Can BB84 be used to send the actual message?**
   - [ ] Yes, but it is very slow. It is primarily used to generate a secure **key**.

---

## 6. Hands-on Exploration

### Run the Example
Execute the automated test to see a successful bit exchange:
```bash
go test -v ./go/examples/security/qkd_test.go
```

### Experiment with the CLI
Use the `strange` CLI tool to prepare and analyze a QKD qubit:

**Export the QKD preparation:**
```bash
./strange export --circuit qkd
```

**Run and see the result:**
```bash
./strange run --circuit qkd
```

**Analyze the entropy:**
```bash
./strange analyze --circuit qkd
```

---

## 7. References & Further Reading

1.  **Wikipedia: BB84** - [The first quantum cryptography protocol](https://en.wikipedia.org/wiki/BB84).
2.  **Qiskit Textbook: QKD** - [Detailed walkthrough of BB84 and eavesdropping detection](https://github.com/Qiskit/textbook/blob/main/notebooks/ch-algorithms/quantum-key-distribution.ipynb).
3.  **Brilliant.org: Quantum Cryptography** - [Introduction to quantum security concepts](https://brilliant.org/wiki/quantum-cryptography/).
4.  **ArXiv: BB84 Original Paper** - [Proceedings of IEEE International Conference on Computers, Systems and Signal Processing (1984)](https://arxiv.org/abs/2003.06557).

**Next Step:** Modify `qkd_test.go` to simulate a mismatching basis. What happens to Bob's measurement?
