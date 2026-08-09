# Strange Simulator Documentation

This directory contains research, planning, and implementation documentation for the Go port of the [Strange](https://github.com/redfx-quantum/strange) quantum simulator.

## Project Overview

Strange is a quantum simulator originally written in Java. This port aims to provide a high-performance Go implementation with native complex number support and optimized gate application.

## CLI Usage

The `strange` CLI tool provides built-in discovery features:
- `strange list circuits`: View all pre-defined quantum algorithms.
- `strange list gates`: View all supported quantum gates and their usage.

## Contents

- [Research Summary](research-summary.md): Requirements and initial research for the Go port.
- [Implementation Plan](implementation-plan.md): The three-phase plan for porting core functionality.
- [Fuzzy Testing Strategy](fuzzy-testing.md): Strategy for verifying advanced gates and matrix operations.
- [Package Documentation](packages.md): Overview of the Go package structure.
- [Quantum Concepts](quantum-concepts.md): Deep dive into quantum bases and gate theory.
- [Simulation Optimizations](optimizations.md): Details on the bitwise optimizations used for high-performance simulation.
- [Algorithmic Scaling Risks](algorithmic-risks.md): Analysis of memory and time constraints for large quantum circuits.
- [Quantum Thermodynamics](thermodynamics.md): Analysis of energy, work, and entropy in quantum systems.
- [Verification and Interoperability](verification.md): Mathematical correctness and OpenQASM support.
- [Verifying with Qiskit](verify-with-qiskit.md): Automated verification with Qiskit Aer.
- [Verifying with IBM Quantum](verify-with-ibm.md): Cloud-based verification guide.
- [Verifying in Quirk](verifying-in-quirk.md): Manual visual verification guide.

## Pedagogical Guides

These learning aids bridge the gap between quantum theory and the Go implementation.

- [Documentation Standards](documentation-standards.md): Rubric used for creating these guides.
- **Algorithms**:
  - [Shor's Algorithm](../examples/algorithms/shor.md)
  - [Grover's Algorithm](../examples/algorithms/grover.md)
- **Networking & Security**:
  - [Quantum Teleportation](../examples/networking/teleportation.md)
  - [Quantum Key Distribution (BB84)](../examples/security/qkd.md)
- **Mathematics & Arithmetic**:
  - [Quantum Arithmetic (Adder)](../examples/arithmetic/arithmetic.md)
  - [Entanglement (Bell & GHZ)](../examples/entanglement/entanglement.md)
  - [Quantum Fundamentals](../examples/fundamentals/fundamentals.md)

## Physics Research

- [sciadv.adw8462 Analysis](sciadv.adw8462_sm.md): Detailed notes on the thermodynamic benchmarks for the simulator.
- [sciadv.adw8462 Coverage](sciadv.adw8462-coverage.md): Evaluation of implementation completeness vs. the original research paper.

## Original Reference

The original Java implementation can be found at: [redfx-quantum/strange](https://github.com/redfx-quantum/strange)
