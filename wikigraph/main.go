// wikigraph generates an interactive force-directed graph of a wiki's internal
// link structure and writes it to a self-contained HTML file.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

var wikilinkRe = regexp.MustCompile(`\[\[([a-z][a-z0-9-]*)\]\]`)

var defaultExcludes = []string{"index", "log", "AGENTS"}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "wikigraph <wiki-dir>",
	Short: "Generate an interactive wikilink graph from a Markdown wiki",
	Long: `wikigraph reads a directory of Markdown wiki pages, extracts [[wikilink]]
cross-references, and produces a self-contained interactive HTML file showing
the link graph as a force-directed network.

Wiki format expected:
  - One Markdown file per page, named <slug>.md (e.g. grovers-algorithm.md)
  - Cross-references written as [[slug]] wikilinks anywhere in the body
  - Three meta-files are excluded from the graph automatically:
      index.md   — the page catalogue
      log.md     — the append-only change log
      AGENTS.md  — the wiki schema / agent instructions
  - All other .md files become nodes in the graph
  - A [[slug]] that resolves to another page becomes a directed edge
  - Pages with no outgoing links are treated as sink nodes: they get a uniform
    probability over all pages (Google-style teleportation) so the Markov kernel
    stays well-defined and the stationary distribution can be computed

Output:
  A single self-contained HTML file. Open it in any browser — no server needed.
    Node size   ∝  stationary distribution (centrality — how often a random walk lands here)
    Node colour =  communicating class (cluster of pages that can all reach each other)
    Edge width  ∝  transition probability (how likely a random walk follows that link)
    Drag, zoom, and pan are fully supported.`,

	Args: cobra.ExactArgs(1),
	RunE: run,
}

var (
	flagOut      string
	flagTitle    string
	flagMinEdge  float64
	flagExclude  []string
	flagSed      []string
)

func init() {
	rootCmd.Flags().StringVarP(&flagOut, "out", "o", "wiki_graph.html", "output HTML file")
	rootCmd.Flags().StringVarP(&flagTitle, "title", "t", "", "graph title shown in browser tab (default: <wiki-dir> wiki)")
	rootCmd.Flags().Float64VarP(&flagMinEdge, "min-edge", "m", 0.005, "omit edges below this transition probability")
	rootCmd.Flags().StringArrayVarP(&flagExclude, "exclude", "e", defaultExcludes, "slugs to exclude from the graph (meta-pages, changelogs, etc.)")
	rootCmd.Flags().StringArrayVarP(&flagSed, "sed", "s", nil, "sed expression(s) to apply to the HTML output (repeatable, e.g. -s 's/foo/bar/')")
}

func run(cmd *cobra.Command, args []string) error {
	wikiDir := args[0]

	title := flagTitle
	if title == "" {
		title = filepath.Base(wikiDir) + " wiki"
	}

	exclude := make(map[string]bool, len(flagExclude))
	for _, s := range flagExclude {
		exclude[s] = true
	}

	pages, idx, err := loadPages(wikiDir, exclude)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Pages: %d\n", len(pages))

	adj, err := buildAdjacency(wikiDir, pages, idx)
	if err != nil {
		return err
	}

	k, err := catrace.NewRandomWalkKernel(adj, pages)
	if err != nil {
		return fmt.Errorf("building kernel: %w", err)
	}

	html, err := k.ToHTML(&catrace.VisualiseOptions{
		Title:   title,
		MinEdge: flagMinEdge,
		Width:   1400,
		Height:  900,
	})
	if err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	page := string(html)

	if len(flagSed) > 0 {
		page, err = applySed(page, flagSed)
		if err != nil {
			return err
		}
	}

	if err := os.WriteFile(flagOut, []byte(page), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", flagOut, err)
	}
	fmt.Fprintf(os.Stderr, "Written: %s\n", flagOut)
	return nil
}

// loadPages reads all non-meta .md files from dir and returns the sorted list
// of slugs and a slug→index map.
func loadPages(dir string, exclude map[string]bool) ([]string, map[string]int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading wiki dir %q: %w", dir, err)
	}
	var pages []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".md")
		if exclude[slug] {
			continue
		}
		pages = append(pages, slug)
	}
	sort.Strings(pages)
	idx := make(map[string]int, len(pages))
	for i, p := range pages {
		idx[p] = i
	}
	return pages, idx, nil
}

// buildAdjacency reads each page and extracts [[wikilinks]], returning a
// square adjacency matrix. Sink pages (no outgoing links) get uniform weight
// across all pages so the Markov kernel stays well-defined.
func buildAdjacency(dir string, pages []string, idx map[string]int) (*mat.Dense, error) {
	n := len(pages)
	adj := mat.NewDense(n, n, nil)
	for i, slug := range pages {
		raw, err := os.ReadFile(filepath.Join(dir, slug+".md"))
		if err != nil {
			return nil, fmt.Errorf("reading %s.md: %w", slug, err)
		}
		linked := map[int]bool{}
		for _, m := range wikilinkRe.FindAllSubmatch(raw, -1) {
			if j, ok := idx[string(m[1])]; ok && j != i {
				linked[j] = true
			}
		}
		if len(linked) == 0 {
			// Sink node: teleport uniformly to avoid a zero row.
			for j := 0; j < n; j++ {
				adj.Set(i, j, 1.0)
			}
		} else {
			for j := range linked {
				adj.Set(i, j, 1.0)
			}
		}
	}
	return adj, nil
}

// applySed pipes html through sed with the given expressions (each passed as a
// separate -e argument). Requires sed to be available on PATH.
func applySed(html string, exprs []string) (string, error) {
	args := make([]string, 0, len(exprs)*2)
	for _, e := range exprs {
		args = append(args, "-e", e)
	}
	cmd := exec.Command("sed", args...)
	cmd.Stdin = strings.NewReader(html)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("sed: %w", err)
	}
	return string(out), nil
}

