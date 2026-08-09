# Wiki Log

<!-- Append-only. Never edit existing entries. -->

## [2026-08-09] init | wiki initialized at quantum-go/wiki/
## [2026-08-09] ingest | quantum-concepts (docs/quantum-concepts.md)
## [2026-08-09] ingest | package-architecture (docs/packages.md)
## [2026-08-09] ingest | quantum-dsl (core/core.go)
## [2026-08-09] ingest | project-overview (README.md)
## [2026-08-09] ingest | simulation-engine (local/engine.go)
## [2026-08-09] ingest | gate-application (local/computations.go)
## [2026-08-09] ingest | quantum-linear-algebra (math/matrix.go)
## [2026-08-09] ingest | simulator-optimizations (docs/optimizations.md)
## [2026-08-09] ingest | scaling-and-limits (docs/algorithmic-risks.md)
## [2026-08-09] ingest | quantum-fundamentals (examples/fundamentals/fundamentals.md)
## [2026-08-09] ingest | rotation-gates (examples/fundamentals/rotations.md)
## [2026-08-09] ingest | universality (examples/fundamentals/universality.md)
## [2026-08-09] ingest | entanglement (examples/entanglement/entanglement.md)
## [2026-08-09] ingest | quantum-arithmetic (examples/arithmetic/arithmetic.md)
## [2026-08-09] ingest | grovers-algorithm (examples/algorithms/grover.md)
## [2026-08-09] ingest | shors-algorithm (examples/algorithms/shor.md)
## [2026-08-09] ingest | error-correction (examples/algorithms/error_correction.md)
## [2026-08-09] ingest | teleportation (examples/networking/teleportation.md)
## [2026-08-09] ingest | bb84-qkd (examples/security/qkd.md)
## [2026-08-09] ingest | gate-implementations (core/gates.go)
## [2026-08-09] ingest | rotation-implementations (core/rotations.go)
## [2026-08-09] ingest | composite-gates (core/composite.go)
## [2026-08-09] ingest | arithmetic-gates (core/arithmetic.go)
## [2026-08-09] ingest | oracle-gates (core/oracles.go)
## [2026-08-09] ingest | openqasm-parser (qasm/parser.go)
## [2026-08-09] ingest | testing-strategy (local/engine_test.go, core/core_test.go)
## [2026-08-09] ingest | verification-tests (local/verification_test.go)
## [2026-08-09] ingest | fuzz-testing (local/fuzz_test.go, math/fuzz_test.go)
## [2026-08-09] ingest | circuits-library (core/circuits.go, core/qasm.go)
## [2026-08-09] ingest | quantum-thermodynamics (examples/thermodynamics/engine_cycle_test.go)
## [2026-08-09] synthesize | state-vector-model (cross-cutting: simulation-engine, gate-application, verification-tests)
## [2026-08-09] synthesize | gate-zoo (cross-cutting: gates.go, rotations.go, composite.go, arithmetic.go, oracles.go)
## [2026-08-09] synthesize | qft-deep-dive (cross-cutting: arithmetic.go, grovers-algorithm, shors-algorithm)
## [2026-08-09] synthesize | algorithm-comparison (cross-cutting: all algorithm pages)
## [2026-08-09] synthesize | how-to-add-a-new-gate (cross-cutting: gates.go, engine.go, qasm.go, verification_test.go)
## [2026-08-09] lint | fixed 5 broken wikilinks; added references to 12 pages; all external links verified OK
## [2026-08-09] policy | AGENTS.md updated: no wiki/raw/ dir; external URLs in ## Sources are the raw pointer
## [2026-08-09] ingest | qelib1-standard-gates (curl: github.com/Qiskit/openqasm OpenQASM2.x — qelib1.inc, qft.qasm, teleport.qasm; Apache 2.0)
## [2026-08-09] citations | upgraded 10 arXiv references via export.arxiv.org/api/query — fixed Steane wrong title/venue; fixed Shor FOCS vs arXiv title; expanded Cross et al. to full authors; added References to arithmetic-gates; clarified Coppersmith IBM Report year
## [2026-08-09] lint | 36 pages checked, 3 orphans found, 3 fixed — added inbound links: circuits-library (quantum-dsl, project-overview), state-vector-model (simulation-engine), testing-strategy (project-overview, how-to-add-a-new-gate); no broken links; OKF frontmatter complete on all pages
