# ADR-003 — Wiki file conventions: OKF frontmatter, process reference, and quality rubric

**Date:** 2026-08
**Status:** Accepted
**Issue:** [#17](https://github.com/stephen-mcelhose/quantum-go/issues/17)

---

## Decision

**In the context of** a growing `wiki/` directory (~30+ files) produced and maintained
by the `llm-wiki` skill, where files had partial, inconsistent frontmatter and no
standard for describing how a page was produced or how to judge its quality,

**facing the concern** that without a schema enforced at the agent-instruction level,
every new `llm-wiki ingest` continues producing non-compliant files — making any
retroactive audit immediately undone by the next write operation, and making it
impossible to reliably search, filter, or health-check wiki content with tooling,

**we decided** to establish three tiers of convention enforced primarily through
`wiki/AGENTS.md` (the governing schema file read by the llm-wiki skill on every
operation):

1. **MUST — OKF frontmatter.** Every wiki file must open with a valid
   [Open Knowledge Format](https://github.com/bayer-int/csgdaa-skills/blob/main/skills/okf.md)
   frontmatter block declaring at minimum `type`, `title`, `description`, `tags`,
   and `timestamp`.

2. **MUST — Process reference.** Every wiki file must record the template or process
   used to produce it (e.g. a `process:` frontmatter key or a note in the body),
   so the file can be reproduced or updated without guessing.

3. **SHOULD — Quality rubric.** Every wiki file should include or link to a rubric
   covering Accuracy, Completeness, Clarity, and Maintainability — making quality
   assessment objective and repeatable.

**to achieve** self-enforcing compliance: once `wiki/AGENTS.md` encodes these
conventions, every future llm-wiki operation produces conforming files by default,
with no manual correction required,

**accepting** that existing files must be audited and backfilled retroactively (tracked
in issue #17), and that the rubric requirement (SHOULD) is advisory — enforced by
convention rather than tooling.

---

## Why enforce via AGENTS.md, not CI

The primary enforcement point is `wiki/AGENTS.md`, not a CI lint step, because:

- `wiki/AGENTS.md` is read by the llm-wiki skill **before every write** — it shapes
  agent behaviour at generation time rather than catching violations after the fact.
- A CI check catches non-compliant commits but cannot prevent them; AGENTS.md
  prevents non-compliant files from being generated in the first place.
- CI validation remains a stretch goal (issue #17 acceptance criteria) for catching
  manually-authored pages or skill regressions.

---

## Why OKF over a custom schema

OKF frontmatter is already the standard emitted by the `llm-wiki` and `okf` skills.
Adopting it here means:

- **No new schema to maintain** — the field definitions live in the skill, not here.
- **Tooling compatibility** — existing health-check and lint tools that understand OKF
  work against these files without modification.
- **Consistency** — wiki files align with every other OKF-managed document in the
  project.

A custom schema would require defining and documenting field semantics that OKF
already covers.

---

## Consequences

- `wiki/AGENTS.md` `## Conventions` section is expanded to fully specify the required
  OKF schema, process-reference requirement, and rubric standard.
- A canonical template file (`wiki/_template.md`) is created so new entries are
  compliant from the first keystroke.
- All existing `wiki/*.md` files are audited and updated to meet MUST requirements
  (tracked in issue #17).
- The `llm-wiki lint` operation gains an explicit check for missing OKF fields and
  absent process references.
- Future ADRs, how-tos, or guides produced in this repo are expected to carry the
  same OKF frontmatter standard, as it is now the established pattern.
