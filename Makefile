GRAPH_OUT ?= wikigraph/wiki_graph.html
WIKI_DIR  ?= wiki

.PHONY: build test lint vet analyze graph visualize tools help

## build: compile the main quantum-go CLI
build:
	go build ./...

## test: run the full test suite
test:
	go test ./...

## vet: run go vet across all packages
vet:
	go vet ./...

## lint: run go vet (extend with golangci-lint if available)
lint: vet
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping (go vet already ran)"; \
	fi

## analyze: run wikigraph structural analysis on the wiki
analyze:
	wikigraph analyze $(WIKI_DIR)

## graph: generate the wiki link graph HTML
graph:
	wikigraph graph $(WIKI_DIR) -o $(GRAPH_OUT) \
		-s 's/width="1400" height="900"/width="100%" height="100%" style="position:fixed;top:0;left:0;"/' \
		-s 's/const W = 1400, H = 900;/const W = window.innerWidth, H = window.innerHeight;/' \
		-s 's/body { margin: 0/body { margin: 0; overflow: hidden/' \
		-s 's/.distance(d => 80 + 120 \* (1 - d.value))/.distance(d => 180 + 280 * (1 - d.value))/' \
		-s 's/.strength(d => 0.3 + 0.5 \* d.value)))/.strength(d => 0.1 + 0.2 * d.value)))/' \
		-s 's/.force("charge", d3.forceManyBody().strength(-400))/.force("charge", d3.forceManyBody().strength(-1800))/' \
		-s 's/.force("center", d3.forceCenter(W\/2, H\/2))/.force("x", d3.forceX(W\/2).strength(0.04)).force("y", d3.forceY(H\/2).strength(0.04))/'

## visualize: generate and open the wiki link graph in the browser
visualize: graph
	open $(GRAPH_OUT)

## tools: install dev tools (wikigraph)
tools:
	go install github.com/stephen-mcelhose/wikigraph@latest

## help: list available commands
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
