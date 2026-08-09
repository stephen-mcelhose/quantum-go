WIKI_DIR  ?= ../../wiki
OUT       ?= wiki_graph.html

.PHONY: build graph

build:
	go build .

graph: build
	./wikigraph $(WIKI_DIR) -o $(OUT) \
		-s 's/width="1400" height="900"/width="100%" height="100%" style="position:fixed;top:0;left:0;"/' \
		-s 's/const W = 1400, H = 900;/const W = window.innerWidth, H = window.innerHeight;/' \
		-s 's/body { margin: 0/body { margin: 0; overflow: hidden/' \
		-s 's/.distance(d => 80 + 120 \* (1 - d.value))/.distance(d => 180 + 280 * (1 - d.value))/' \
		-s 's/.strength(d => 0.3 + 0.5 \* d.value)))/.strength(d => 0.1 + 0.2 * d.value)))/' \
		-s 's/.force("charge", d3.forceManyBody().strength(-400))/.force("charge", d3.forceManyBody().strength(-1800))/' \
		-s 's/.force("center", d3.forceCenter(W\/2, H\/2))/.force("x", d3.forceX(W\/2).strength(0.04)).force("y", d3.forceY(H\/2).strength(0.04))/'
