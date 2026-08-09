# Wiki Index

| Page                          | Summary                                                                         | Tags                                              |
| ----------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------- |
| [[quantum-concepts]]          | Core quantum principles (superposition, entanglement) and full gate reference   | quantum-mechanics, gates, qubits, basis           |
| [[package-architecture]]      | Three-layer Go package design (math→core→local) and the GlobalStepExecutor trick | architecture, go-packages, design                |
| [[quantum-dsl]]               | Program→Step→Gate hierarchy; Block, Result, and Qubit types from core.go       | dsl, go, core, program, gate, block               |
| [[project-overview]]          | High-level framing, installation, gate support, and learning guide index        | overview, installation, quick-start, project      |
| [[simulation-engine]]         | State vector initialization from InitAlpha and per-step execution loop          | simulation, engine, state-vector, local           |
| [[gate-application]]          | Bitwise bit-loop optimizations — O(2^n) gate application without large matrices | optimization, bitwise, computations, performance  |
| [[quantum-linear-algebra]]    | Flat complex128 matrix type, Kronecker products, conjugate transpose            | linear-algebra, matrix, kronecker, math           |
| [[simulator-optimizations]]   | Conceptual rationale for the bit-loop strategy and complexity table             | optimization, performance, complexity             |
| [[scaling-and-limits]]        | Memory/time ceiling (~30 qubits), Shor's and Grover's risks, mitigations        | scaling, limits, memory, complexity               |
| [[quantum-fundamentals]]      | Superposition, measurement, Bloch Sphere, single-qubit gates with worked examples | fundamentals, superposition, hadamard, pauli    |
| [[rotation-gates]]            | Rx, Ry, Rz, PhaseShift — continuous Bloch Sphere rotations at arbitrary angles  | rotation, rx, ry, rz, phaseshift, parameterized  |
| [[universality]]              | Universal gate sets, the U gate as 3-parameter Euler decomposition, Solovay-Kitaev | universality, u-gate, gate-sets, decomposition |
| [[entanglement]]              | Bell states, GHZ states, non-separability, H+CNOT construction                  | entanglement, bell-state, ghz, cnot              |
| [[quantum-arithmetic]]        | Draper adder — QFT-based addition in phase space without carry qubits           | arithmetic, draper-adder, qft, modular           |
| [[grovers-algorithm]]         | Amplitude amplification for unstructured search in O(√N) — oracle + diffusion   | grover, search, oracle, diffusion, quadratic     |
| [[shors-algorithm]]           | Period finding via QFT interference breaks RSA in O(log³N)                      | shor, factoring, period-finding, qft, rsa        |
| [[error-correction]]          | 3-qubit bit-flip code — syndrome measurement and Toffoli correction             | error-correction, bit-flip, syndrome, toffoli    |
| [[teleportation]]             | Bell measurement + 2 classical bits transfer quantum state without matter        | teleportation, bell-measurement, correction      |
| [[bb84-qkd]]                  | BB84 protocol — conjugate bases and Heisenberg uncertainty for secure key exchange | bb84, qkd, cryptography, conjugate-bases      |
| [[gate-implementations]]      | core/gates.go — H, X, Y, Z, CNOT, CZ, SWAP, Toffoli, Fredkin as BaseGate types   | gates, implementation, pauli, toffoli         |
| [[rotation-implementations]]  | core/rotations.go — Rx/Ry/Rz/U/S/T/V/PhaseShift/CR with inverse support          | rotations, u-gate, s-gate, t-gate, inverse    |
| [[composite-gates]]           | BlockGate, Oracle, ControlledGate, ControlledBlockGate, TimeEvolution              | composite, block-gate, oracle, controlled      |
| [[arithmetic-gates]]          | Fourier (QFT), Add, AddInteger, AddIntegerModulus, MulModulus — Shor's engine     | arithmetic, qft, fourier, mulmodulus, shor    |
| [[oracle-gates]]              | ConstantOracle, BalancedOracle, InnerProductOracle, SimonOracle — permutation matrices | oracle, deutsch-jozsa, bernstein-vazirani |
| [[openqasm-parser]]           | qasm/parser.go — QASM 2.0 subset parser for cross-platform circuit import         | openqasm, parser, qiskit, interoperability    |
| [[qelib1-standard-gates]]     | qelib1.inc — every standard gate, its decomposition, and quantum-go mapping        | openqasm, qelib1, standard-gates, qasm2       |
| [[testing-strategy]]          | Test suite architecture — unit, integration, verification, fuzz, example tests     | testing, go-test, fuzz, verification          |
| [[verification-tests]]        | local/verification_test.go — 15+ circuits compared against expected amplitudes     | verification, state-vector, bell, toffoli     |
| [[fuzz-testing]]              | FuzzToffoli, FuzzFourier (unitarity), FuzzAdd, FuzzMatrixMul — Go 1.18+ fuzzing   | fuzz, unitarity, toffoli, matrix              |
| [[circuits-library]]          | core/circuits.go — Program factories, ExpectedXxxState helpers, ToQASM export     | circuits, library, factory, qasm-export       |
| [[quantum-thermodynamics]]    | Density matrices, expectation values, Von Neumann entropy, TimeEvolution           | thermodynamics, density-matrix, entropy       |
| [[state-vector-model]]        | End-to-end lifecycle of the []complex128 state vector — allocation to measurement  | state-vector, amplitude, measurement, model   |
| [[gate-zoo]]                  | Complete reference table of every gate — constructor, matrix, QASM, size           | reference, gate-zoo, all-gates, constructor   |
| [[qft-deep-dive]]             | Step-by-step QFT circuit — H, CR rotations at 2π/2ᵏ, bit-reversal, IQFT           | qft, fourier, cr-gates, bit-reversal          |
| [[algorithm-comparison]]      | Side-by-side: qubit count, depth, query complexity, oracle type for all algorithms  | algorithms, comparison, complexity, reference |
| [[how-to-add-a-new-gate]]     | Contributor guide — struct, GetMatrix, QASM export, verification test checklist     | how-to, contribution, gate, testing           |
