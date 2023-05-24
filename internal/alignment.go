package alignment

type Tree map[int]*NodeDescription

type TreePair struct {
	Lhs int
	Rhs int
}

type AlignmentTask struct {
	config      *Config
	first       Tree
	second      Tree
	sortedPairs []TreePair
	quality     map[TreePair]int
}

func NewTree() Tree {
	return make(Tree)
}

func buildTree(description []NodeDescription) Tree {
	tree := NewTree()
	for _, node := range description {
		tree[node.ID] = &node
		tree[node.ID+len(description)] = &NodeDescription{
			ID:    node.ID,
			Tag:   DeletionTag,
			Left:  node.Left,
			Right: node.Right,
		}
	}
	return tree
}

func getSortedTreePairs(first, second int) []TreePair {
	pairs := []TreePair{TreePair{Lhs: -1, Rhs: -1}}

	// Empty first subtree.
	for j := second - 1; j >= 0; j-- {
		pairs = append(pairs, TreePair{Lhs: EmptyTreeID, Rhs: j})
	}
	// Empty second subtree.
	for i := first - 1; i >= 0; i-- {
		pairs = append(pairs, TreePair{Lhs: i, Rhs: EmptyTreeID})
	}
	// Non-empty subtrees.
	for i := first - 1; i >= 0; i-- {
		for j := second - 1; j >= 0; j-- {
			pairs = append(pairs, TreePair{Lhs: i, Rhs: j})
		}
	}
	return pairs
}

func (a *AlignmentTask) calculatePairQuality(pair TreePair) int {
	// Empty trees.
	if pair.Lhs == EmptyTreeID && pair.Rhs == EmptyTreeID {
		return 0
	}

	// Empty fisrt tree.
	if pair.Lhs == EmptyTreeID {
		id := pair.Rhs
		leftID := a.second[id].Left
		rightID := a.second[id].Right
		quality := 0
		if a.second[id].Tag != DeletionTag {
			quality -= a.config.DeletionCost
		}
		if leftID != nil {

		}
	}

	// Empty second tree.

	//
}

func (a *AlignmentTask) calculateQuality() {
	a.quality[-1]
	for _, pair := range a.sortedPairs {

	}
}

func CalculateAlignment(c *Config) error {
	// Build trees.
	firstTree := buildTree(c.FirstGraphDescription)
	secondTree := buildTree(c.SecondGraphDescription)

	// Get sorted pairs of subtrees.
	sortedPairs := getSortedTreePairs(len(firstTree), len(secondTree))

	a := AlignmentTask{
		config:      c,
		first:       firstTree,
		second:      secondTree,
		sortedPairs: sortedPairs,
		quality:     make(map[TreePair]int),
	}

	// Calculate quality of this pairs.
	a.calculateQuality()

	// Dynamic.
	return nil
}
