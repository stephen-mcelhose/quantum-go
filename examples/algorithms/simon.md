# Simon's Algorithm

Simon's algorithm was one of the first examples of a quantum algorithm that is exponentially faster than any known classical algorithm for a specific (albeit somewhat artificial) problem. It served as a direct inspiration for Shor's factoring algorithm.

## The Problem

Given a function $f: \{0,1\}^n \to \{0,1\}^n$ and a secret bitstring $s \in \{0,1\}^n$, we are promised that:
$$f(x) = f(y) \iff y = x \oplus s$$

The goal is to find the secret string $s$.

Classically, finding $s$ requires checking a number of inputs proportional to $2^{n/2}$ (the birthday problem). Simon's algorithm finds $s$ using only $O(n)$ queries to the quantum oracle.

## Implementation in Strange

The "Strange" simulator provides a built-in Simon's algorithm implementation.

### Usage via CLI

```bash
# Run Simon's algorithm with a hidden string '11'
./strange run --circuit simon -p s=11
```

The output will show measurements $y$ such that $y \cdot s = 0 \pmod 2$. By repeating the execution and collecting enough independent equations, you can solve for $s$ classically using Gaussian elimination.

### Programmatic Usage

```go
package main

import (
    "fmt"
    "github.com/stephen-mcelhose/quantum-go/core"
    "github.com/stephen-mcelhose/quantum-go/local"
)

func main() {
    s := "11"
    program := core.NewSimonsProgram(s)
    
    env := local.NewSimpleExecutionEnvironment()
    result := env.RunProgram(program)
    
    // Result contains states |y> where y . s = 0
    result.PrintBinary()
}
```

## How it Works

1.  **Initialize**: Prepare $n$ input qubits and $n$ output qubits.
2.  **Superposition**: Apply Hadamard gates to the input register to create an equal superposition of all possible $2^n$ strings.
3.  **Oracle**: Apply the Simon Oracle $U_f$ which maps $|x\rangle|0\rangle \to |x\rangle|f(x)\rangle$.
4.  **Interference**: Apply Hadamard gates again to the input register.
5.  **Measure**: Measuring the input register yields a bitstring $y$ such that $y \cdot s = 0 \pmod 2$.

By running this circuit approximately $n$ times, we obtain a system of $n-1$ linear equations which uniquely determine $s$.
