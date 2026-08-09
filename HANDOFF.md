# Agent Handoff — quantum-go

## What This Repo Is

`quantum-go` is a Go port of the [Strange](https://github.com/redfx-quantum/strange)
quantum circuit simulator by Johan Vos. It lives at:

- **Local:** `~/repos/quantum-go`
- **Remote:** https://github.com/stephen-mcelhose/quantum-go (public, `stephen-mcelhose` account)

The Go module path is `github.com/stephen-mcelhose/quantum-go`.

Everything builds and all 78 tests pass (`go test ./...`).

---

## Immediate Next Steps (do these first)

### 1. Commit and push the skill file

The llm-wiki skill was added locally but not yet pushed:

```bash
cd ~/repos/quantum-go
git add .agents/
git commit -m "chore: add repo-scoped llm-wiki skill for wiki ingestion"
git push
```

### 2. Initialize the wiki

```bash
# The wiki lives at wiki/ in the repo root
# Run the llm-wiki skill init operation:
# /llm-wiki init
```

The skill at `.agents/skills/llm-wiki/SKILL.md` has repo-specific configuration.
WIKI_ROOT is **always** `wiki/` — the skill hardcodes this, no need to ask.

---

## What Needs to Be Built: The Wiki

The primary goal is to ingest the repo's docs, source files, and examples into a
compounding LLM wiki at `wiki/` so Stephen can use it as a personal learning reference
for quantum computing.

### Ingestion Plan (execute in phase order)

**Phase 0 — Init**
Run `/llm-wiki init` to create `wiki/AGENTS.md`, `wiki/index.md`, `wiki/log.md`.

**Phase 1 — Foundation** (do first; everything else references these)
| Source | Wiki Slug |
|--------|-----------|
| `docs/quantum-concepts.md` | `quantum-concepts` |
| `docs/packages.md` | `package-architecture` |
| `core/core.go` | `quantum-dsl` |
| `README.md` | `project-overview` |

**Phase 2 — Simulation Internals**
| Source | Wiki Slug |
|--------|-----------|
| `local/engine.go` | `simulation-engine` |
| `local/computations.go` | `gate-application` |
| `math/matrix.go` | `quantum-linear-algebra` |
| `docs/optimizations.md` | `simulator-optimizations` |
| `docs/algorithmic-risks.md` | `scaling-and-limits` |

**Phase 3 — Algorithms & Pedagogical Guides** (richest content)
| Source | Wiki Slug |
|--------|-----------|
| `examples/fundamentals/fundamentals.md` | `quantum-fundamentals` |
| `examples/fundamentals/rotations.md` | `rotation-gates` |
| `examples/fundamentals/universality.md` | `universality` |
| `examples/entanglement/entanglement.md` | `entanglement` |
| `examples/arithmetic/arithmetic.md` | `quantum-arithmetic` |
| `examples/algorithms/grover.md` | `grovers-algorithm` |
| `examples/algorithms/shor.md` | `shors-algorithm` |
| `examples/algorithms/error_correction.md` | `error-correction` |
| `examples/networking/teleportation.md` | `teleportation` |
| `examples/security/qkd.md` | `bb84-qkd` |

**Phase 4 — Gate & Composite Implementations**
| Source | Wiki Slug |
|--------|-----------|
| `core/gates.go` | `gate-implementations` |
| `core/rotations.go` | `rotation-implementations` |
| `core/composite.go` | `composite-gates` |
| `core/arithmetic.go` | `arithmetic-gates` |
| `core/oracles.go` | `oracle-gates` |
| `qasm/parser.go` | `openqasm-parser` |

**Phase 5 — Verification & Interoperability**
| Source | Wiki Slug |
|--------|-----------|
| `docs/verification.md` | `verification-strategy` |
| `docs/verify-with-qiskit.md` | `qiskit-verification` |
| `docs/verify-with-ibm.md` | `ibm-quantum-verification` |
| `docs/verifying-in-quirk.md` | `quirk-verification` |
| `verification/VERIFICATION_SUMMARY.md` | `verification-results` |
| `verification/QISKIT_VERIFICATION_REPORT.md` | `qiskit-report` |

**Phase 6 — Supplemental**
| Source | Wiki Slug |
|--------|-----------|
| `docs/thermodynamics.md` | `quantum-thermodynamics` |
| `docs/future-roadmap.md` | `future-roadmap` |
| `REFERENCES.md` | `references` |

### Cross-Cutting Synthesis Pages to Create
After Phase 3+, synthesize these from multiple sources:
- `state-vector-model` — engine.go + computations.go + quantum-concepts.md
- `gate-zoo` — gates.go + rotations.go + quantum-concepts.md
- `qft-deep-dive` — composite.go + arithmetic.md + shor.md
- `how-to-add-a-new-gate` — core.go + gates.go + computations.go

---

## Repo Structure

```
quantum-go/
├── .agents/skills/llm-wiki/SKILL.md   ← repo-scoped wiki skill (needs push)
├── wiki/                               ← wiki lives here (create via init)
├── core/       Gate definitions, Program/Step/Gate DSL
├── local/      Simulation engine (state vector)
├── math/       Complex matrix library
├── qasm/       OpenQASM 2.0 parser
├── examples/   Algorithm tests + pedagogical .md guides
├── docs/       Concept guides, verification docs
└── verification/ Qiskit reports and bridge script
```

---

## Key Facts

- BSD 3-Clause license (Johan Vos, original Strange project)
- Go 1.24.4, module: `github.com/stephen-mcelhose/quantum-go`
- CLI binary: `go build -o quantum-go ./cmd/quantum-go/main.go` — then `./quantum-go list circuits`
- 78 tests passing, including fuzz tests in `local/` and `math/`
- Qiskit bridge at `verification/verify_against_qiskit.py` (needs Python venv + qiskit-aer)
- The `quantum-go` binary and `venv/` are gitignored

---

*Delete this file after the new agent session is underway.*
