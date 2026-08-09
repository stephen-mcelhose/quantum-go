# How-To: Wiki Health Check with wikigraph

Run a structural analysis of the wiki, interpret the results, and file actionable GitHub issues — periodically or before a release.

---

## Prerequisites

- `wikigraph` installed: `go install github.com/stephen-mcelhose/wikigraph@latest` (or check `~/go/bin/wikigraph`)
- `gh` CLI authenticated: `gh auth status`
- Wiki lives in `wiki/` at the repo root, one `.md` file per page using `[[slug]]` wikilinks

---

## Step 1 — Run the health report

```bash
wikigraph analyze wiki/ > /tmp/wikigraph-analysis.txt
cat /tmp/wikigraph-analysis.txt
```

The report has six sections:

| Section | What it tells you |
| ------- | ----------------- |
| **Overview** | Page count, edge count, entropy rate, number of communicating classes |
| **Communicating classes** | Whether any pages are structurally isolated (transient = problem) |
| **Orphan pages** | Bottom 10% by stationary distribution — hard to reach via link-following |
| **Sink pages** | Pages with zero outgoing wikilinks — reader dead ends |
| **Most central** | Top 5 hub pages (highest stationary distribution) |
| **Suggested missing links** | Unlinked pairs with low commute time — natural candidates to connect |

### Flags worth knowing

```bash
# Raise the orphan threshold (flag more pages)
wikigraph analyze wiki/ --orphan-pct 0.2

# Show more suggested links per page
wikigraph analyze wiki/ --suggest-top 5

# Suppress suggested links entirely (faster summary)
wikigraph analyze wiki/ --suggest-top 0
```

---

## Step 2 — Interpret the results

### Communicating classes
- **1 recurrent class containing all pages** = healthy baseline. Every page is reachable.
- **Transient class** = a page or cluster that can be entered but never left. Fix immediately — it means some pages have no outgoing links AND no pages link back to them.

### Orphan pages (low stationary distribution π)
A low π means a random reader following links will rarely land here. Check:
1. How many inbound links does the page actually have? (`grep -rl '\[\[slug\]\]' wiki/`)
2. Are the inbound links only from `index.md`? (Index is excluded from the graph — those links don't count.)
3. Do the outgoing links from this page point to high-π hubs? If not, the page is doubly isolated.

**High commute time (> 200) = urgent.** The page is far from the main graph.

### Sink pages (no outgoing wikilinks)
Verify with: `grep -c '\[\[' wiki/page-name.md` — if 0, it's a sink.
These are usually reference/table pages that were written without navigation in mind.

### Suggested missing links
Low commute time = the two pages are already close in the graph (many shared neighbors) but not directly linked. These are the easiest wins — add a `[[link]]` where the prose naturally supports it.

**Do not force links.** Only add a wikilink where it is contextually meaningful to the reader. Graph metrics suggest candidates; editorial judgment makes the final call.

---

## Step 3 — Choose a response strategy

wikigraph tells you *what* is wrong; the strategy determines *how* to fix it. Different signals call for different responses — not every orphan should be fixed with "add a link."

### Strategy catalogue

---

#### Add inbound links
**Use when:** A page exists and is good, but nothing points to it.  
**wikigraph signal:** Orphan page with a reasonable π and a small number of real inbound links (not just index).  
**What it does to the graph:** Raises π of the orphan; shortens its commute time to hubs.  
**Steps:** Find existing pages whose prose naturally mentions the orphan's topic → add `[[slug]]` inline.

---

#### Add outgoing links
**Use when:** A page is a reader dead end — no forward navigation.  
**wikigraph signal:** Sink page (zero outgoing wikilinks).  
**What it does to the graph:** Turns a sink into a pass-through; readers can continue exploring.  
**Steps:** Read the page and identify every concept, gate, or algorithm mentioned → add `[[slug]]` wherever the prose already implies the relationship.

---

#### Add cross-links
**Use when:** Two pages are structurally close (low commute time) but not directly linked.  
**wikigraph signal:** Suggested missing links section — low commute time pairs.  
**What it does to the graph:** Reduces commute time between hubs; increases overall graph density.  
**Caution:** Only add where editorially justified. Don't wire pages together just because the graph says they're close.

---

#### Build a new page
**Use when:** Multiple existing pages all refer to a concept in prose — but there is no wiki page for that concept. The missing page would immediately attract inbound links and become a natural hub.  
**wikigraph signal:** *Not directly visible* — the tool only graphs pages that exist. Detect by searching for a term appearing in prose across many pages without a corresponding `[[slug]]` link. Also a signal: a cluster of orphans that all share a theme with no central concept page.  
**What it does to the graph:** Creates a new node; inbound links migrate from implicit prose references to explicit wikilinks; the new page immediately raises the π of all pages that link to it.  
**Steps:**
1. Identify the missing concept (grep for the term across `wiki/`)
2. Write the page — it should have outgoing links to related pages to avoid being a sink itself
3. Add `[[new-slug]]` to every page that was already referencing the concept in prose
4. Re-run `wikigraph analyze` — the new page should not appear as an orphan

**Example trigger:** Seeing that `error-correction`, `shors-algorithm`, and `scaling-and-limits` all discuss "syndrome measurement" without a dedicated page for it.

---

#### Build an index
**Use when:** A cluster of related pages exists but has no gateway or landing page. Readers arriving at any one page in the cluster can't easily discover the others.  
**wikigraph signal:** A group of pages with similar tags/topics that have reasonable π individually but are weakly connected to each other (high pairwise commute times within the cluster). Also useful when entropy rate is high — readers are scattered rather than directed.  
**What it does to the graph:** The index page becomes a high-π hub for its cluster; all cluster pages gain inbound links from it; readers have a clear entry point.  
**Steps:**
1. Identify the cluster (e.g., "all testing pages", "all gate-reference pages")
2. Create `wiki/[cluster-topic]-index.md` — a structured table or list of pages in the cluster with one-line descriptions
3. Link the index from each cluster page ("See also: [[cluster-topic-index]]")
4. Link the index from higher-level pages (`project-overview`, `quantum-concepts`, etc.)
5. Add the index slug to `--exclude` in wikigraph if it's purely navigational (no informational content)

**Example trigger:** `gate-zoo`, `gate-implementations`, `rotation-implementations`, `composite-gates`, and `arithmetic-gates` all exist but a reader doesn't know they form a "Gates Reference" cluster.

---

#### Build a community
**Use when:** Two or more isolated clusters exist in the graph — pages within each cluster link to each other, but the clusters barely link across. This is a structural divide, not just an orphan.  
**wikigraph signal:** High commute times between pages that *should* be related (e.g., algorithm pages and their underlying gate/math pages don't cross-link). Recognisable in the communicating classes output if a transient sub-cluster appears, or by inspecting the suggested links section and seeing that nearly all suggestions bridge the same two groups.  
**What it does to the graph:** Merges two sub-communities into one; reduces overall average commute time; increases entropy rate toward its theoretical maximum (fully connected graph).  
**Steps:**
1. Map the two clusters — list every page in each
2. Identify natural bridge points: pages in cluster A that discuss concepts owned by cluster B
3. Add cross-links at the bridge points (this is targeted "add cross-links" applied between clusters, not within one)
4. Consider whether a **bridge page** is needed — a new page that explicitly sits between the two communities (e.g., a "From algorithm to implementation" guide that links algorithm pages to their gate implementation pages)
5. Re-run `wikigraph analyze` — commute times between the two former clusters should drop significantly

**Example trigger:** The "algorithms" cluster (`grovers-algorithm`, `shors-algorithm`, `bb84-qkd`, etc.) and the "implementation" cluster (`gate-application`, `composite-gates`, `simulator-optimizations`, etc.) are navigable within themselves but don't cross-reference each other enough. A reader learning Grover's algorithm can't easily get to the gate-application page that explains how its oracle is actually applied.

---

#### Break into communities

> ⚠️ **This is the rarest intervention and the easiest to apply prematurely. Exhaust all other strategies first. Do not apply to wikis under ~100 pages unless the audience split is unambiguously real and causing measurable harm.**

**Use when:** The graph has become so densely connected that it has lost navigational structure — readers have no sense of direction because everything is one hop from everything else. Or: the wiki demonstrably serves two distinct audiences whose needs are in tension, and the current graph actively misleads one of them.  

**wikigraph signal:** This requires multiple signals in combination — no single metric justifies it:
- Entropy rate approaching log₂(N) (near-random graph — no structure)
- One hub page with disproportionately high π (e.g., > 3× the expected share of 1/N) that is central for structural rather than conceptual reasons
- Suggested links output is saturated — nearly every page pair is suggested, meaning the tool can no longer discriminate
- Audience analysis (not a graph metric): you can identify two reader types whose journeys through the wiki are genuinely opposed

**What it does to the graph:** Reduces average connectivity; raises commute times within the overall graph but lowers them *within* each community; makes learning paths meaningful again by reintroducing distance between clusters.

**What it looks like in practice** — you rarely delete links outright:
- Replace direct cross-cluster links with a **bridge page** that is the single sanctioned crossing point between communities
- Introduce **community index pages** as the only entry point into a cluster from outside
- Audit and remove links that were added for graph-metric reasons rather than editorial ones — if a link doesn't help a specific reader type, remove it

**Hard preconditions — all must be true before applying:**
1. The wiki has enough pages that communities are meaningful (rough guide: 80+ pages)
2. You can name the two audiences and describe their distinct reading goals without ambiguity
3. The density problem is confirmed by entropy rate, not just inferred from a single over-central page
4. You have tried "add cross-links" and found it made navigation *worse*, not better

**The premature-application failure mode:** Splitting a small wiki into communities before it has enough content creates the illusion of structure while actually just adding index pages that nobody links to. You end up with orphan index pages — which is the problem you were trying to solve. At small scale, the right tool is almost always "build an index" (one cluster, one entry point) rather than "break into communities" (enforced partition between clusters).

**Example trigger (legitimate):** A wiki that has grown to 200+ pages serving both end-users learning quantum algorithms and core contributors building the simulator engine. Link-following from a beginner algorithm page now reaches deep internals within 2 hops, and vice versa. Contribution docs are polluting the learner journey and the learner context is making it hard for contributors to find implementation details quickly.

**Example non-trigger (premature):** A 36-page wiki where one hub page has π=0.096. That's a monitoring signal, not an action signal. The right response is to watch it across future runs — if it climbs above ~0.15 as the wiki grows, then revisit.

---

### Strategy selection quick-reference

| wikigraph signal | First-choice strategy | Consider also |
| --- | --- | --- |
| Orphan, commute > 400, few inbound links | Add inbound links | Build a new page (if concept is missing) |
| Orphan, mutual-linking island (two pages only link each other) | Add inbound links to both | Build an index (if they're part of a larger cluster) |
| Sink page (zero outgoing links) | Add outgoing links | — |
| Suggested links (low commute, unlinked) | Add cross-links | — |
| Multiple orphans sharing a theme, no hub page | Build a new page | Build an index |
| Cluster of related pages with no entry point | Build an index | — |
| High commute between two topic areas that should be related | Build a community | Add cross-links at bridge points |
| Transient communicating class | Add outgoing links (immediate) | Build a community (structural) |
| Near-random entropy rate + distinct audiences confirmed + 80+ pages | Break into communities | Build an index per cluster first |

---

## Step 4 — Create GitHub labels (first time only)

```bash
gh label create "wiki" \
  --description "Wiki structure, content, and cross-linking" \
  --color "0052cc"

gh label create "priority: high" \
  --description "Should be addressed soon" \
  --color "b60205"

gh label create "priority: medium" \
  --description "Important but not urgent" \
  --color "e4e669"

gh label create "priority: low" \
  --description "Nice to have" \
  --color "0e8a16"

gh label create "auto-documentation" \
  --description "Improvements identified by automated documentation tooling" \
  --color "5319e7"
```

These only need creating once. Skip any that already exist.

---

## Step 5 — File issues

File one issue per distinct problem type (not one per page). Group related pages together.

### Priority guide

| Condition | Priority |
| --------- | -------- |
| Orphan with commute > 400 OR mutual-linking island | High |
| Sink page (zero outgoing links) | High |
| Orphan with commute 100–400 | Medium |
| Suggested missing links (low commute, not yet linked) | Low |

### Issue template

```
## Problem
[Quote the wikigraph output — π value, commute time, inbound link count]

## Root cause
[Why is the page isolated? Mutual island? Reference page never linked from prose?]

## Suggested fixes
[Specific pages + specific proposed wikilinks — make it mechanical for the fixer]

## Source
Identified by `wikigraph analyze wiki/` — [section name].
```

### Labels to apply

Every wiki issue should get: `wiki`, `documentation`, `auto-documentation`, and the appropriate `priority: *` label.

```bash
gh issue create \
  --title "wiki: [problem description]" \
  --label "wiki,documentation,auto-documentation,priority: high" \
  --body "..."
```

---

## Step 6 — Verify labels were applied

```bash
gh issue list --label "auto-documentation" --state open
```

---

## Worked example — quantum-go (2026-08-09)

Running `wikigraph analyze wiki/` against the 36-page wiki produced:

| Finding | Pages | Action taken |
| ------- | ----- | ------------ |
| Mutual-linking island (π ≈ 0.0016–0.0017, commute ~600+) | `gate-zoo`, `how-to-add-a-new-gate` | Filed #10 — priority: high |
| Sink pages (0 outgoing links) | `algorithm-comparison`, `fuzz-testing`, `qelib1-standard-gates`, `verification-tests` | Filed #11 — priority: high |
| Orphan — low inbound (π=0.0025, commute ~420) | `rotation-gates` | Filed #12 — priority: medium |
| Orphan — low inbound (π=0.0026, commute ~390) | `error-correction` | Filed #13 — priority: medium |
| Suggested hub cross-links (commute 18–45) | 10 page pairs | Filed #14 — priority: low |

Most central pages (healthy hubs to link toward): `composite-gates` (π=0.096), `gate-application` (π=0.090), `quantum-linear-algebra` (π=0.077).

---

## How often to run this

| Trigger | Rationale |
| ------- | --------- |
| New wiki page added | Check it isn't immediately an orphan or sink |
| Before a release | Ensure the wiki is navigable for new readers |
| Periodically (monthly) | Catch drift as pages are added over time |
| After a batch of edits | Confirm cross-links were added correctly |
