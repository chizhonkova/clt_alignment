# Alignment of Cell Lineage Trees
Algorithm for alignment of CLTs.

## Build
```(console)
cd clt_alignment/cmd/tree_alignment
go build
```

## Usage
```(console)
$ ./tree_alignment --help
calculates alignment for two binary trees with tags

Usage:
  tree-alignment [flags]

Flags:
      --deletion-cost int          penalty for deletion (default 2)
      --first-graph-path string    path to the first graph description
  -h, --help                       help for tree-alignment
      --result-path string         path to the resulting graph description
      --second-graph-path string   path to the second graph description
      --tag-equality-cost int      cost for equal tags (default 4)
      --tag-unequality-cost int    penalty for unequal tags (default 3)
```

## Example
All necessary data for the example can be found in cmd/tree_alignment.
```(console)
$ ./tree_alignment --first-graph-path first_graph.json --second-graph-path second_graph.json --result-path res.json
Calculating maximum quality...
Quality: 5
Building CLT alignment...
Result is written to res.json.
```
