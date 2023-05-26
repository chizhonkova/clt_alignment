package alignment

import "fmt"

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
	qCase       map[TreePair]int
}

func NewTree() Tree {
	return make(Tree)
}

func buildTree(description []*NodeDescription) Tree {
	tree := NewTree()
	for _, node := range description {
		node.FullPenalty = DefaultPenalty

		tree[node.ID] = node

		extraNode := &NodeDescription{
			ID:          node.ID + len(description),
			Tag:         DeletionTag,
			Left:        node.Left,
			Right:       node.Right,
			FullPenalty: DefaultPenalty,
		}
		tree[extraNode.ID] = extraNode
	}
	tree[EmptyTreeID] = &NodeDescription{
		ID:          EmptyTreeID,
		Tag:         DeletionTag,
		FullPenalty: 0,
	}

	return tree
}

func getSortedTreePairs(first, second int) []TreePair {
	pairs := []TreePair{{Lhs: -1, Rhs: -1}}

	// Empty first subtree.
	for j := second - 1; j >= 0; j-- {
		pairs = append(pairs, TreePair{Lhs: EmptyTreeID, Rhs: j})
		pairs = append(pairs, TreePair{Lhs: EmptyTreeID, Rhs: j + second})
	}
	// Empty second subtree.
	for i := first - 1; i >= 0; i-- {
		pairs = append(pairs, TreePair{Lhs: i, Rhs: EmptyTreeID})
		pairs = append(pairs, TreePair{Lhs: i + first, Rhs: EmptyTreeID})
	}
	// Non-empty subtrees.
	for i := first - 1; i >= 0; i-- {
		for j := second - 1; j >= 0; j-- {
			pairs = append(pairs, TreePair{Lhs: i, Rhs: j})
			pairs = append(pairs, TreePair{Lhs: i, Rhs: j + second})
			pairs = append(pairs, TreePair{Lhs: i + first, Rhs: j})
			pairs = append(pairs, TreePair{Lhs: i + first, Rhs: j + second})
		}
	}
	return pairs
}

func (a *AlignmentTask) calculateForEmptyLeft(pair TreePair) (int, int) {
	if pair.Lhs != EmptyTreeID {
		return 0, 0
	}

	id := pair.Rhs
	leftID := a.second[id].LeftOrDefault()
	rightID := a.second[id].RightOrDefault()

	if a.second[id].FullPenalty == DefaultPenalty {
		penalty := 0

		if a.second[id].Tag != DeletionTag {
			penalty += a.config.DeletionCost
		}

		penalty += a.second[leftID].FullPenalty
		penalty += a.second[rightID].FullPenalty

		a.second[id].FullPenalty = penalty
	}

	return -a.second[id].FullPenalty, LeftEmptyCase
}

func (a *AlignmentTask) calculateForEmptyRight(pair TreePair) (int, int) {
	if pair.Rhs != EmptyTreeID {
		return 0, 0
	}

	id := pair.Lhs
	leftID := a.first[id].LeftOrDefault()
	rightID := a.first[id].RightOrDefault()

	if a.first[id].FullPenalty == DefaultPenalty {
		penalty := 0

		if a.first[id].Tag != DeletionTag {
			penalty += a.config.DeletionCost
		}

		penalty += a.first[leftID].FullPenalty
		penalty += a.first[rightID].FullPenalty

		a.first[id].FullPenalty = penalty
	}

	return -a.first[id].FullPenalty, RightEmptyCase
}

func (a *AlignmentTask) calculatePairQuality(pair TreePair) (int, int) {
	// Empty trees.
	if pair.Lhs == EmptyTreeID && pair.Rhs == EmptyTreeID {
		return 0, BothEmptyCase
	}

	// Empty fisrt tree.
	if pair.Lhs == EmptyTreeID {
		return a.calculateForEmptyLeft(pair)
	}

	// Empty second tree.
	if pair.Rhs == EmptyTreeID {
		return a.calculateForEmptyRight(pair)
	}

	return 0, 0
}

func (a *AlignmentTask) calculateQuality() {
	for _, pair := range a.sortedPairs {
		quality, qCase := a.calculatePairQuality(pair)
		a.quality[pair] = quality
		a.qCase[pair] = qCase
		fmt.Printf("Pair: %v, Quality: %v, Case: %v\n", pair, a.quality[pair], a.qCase[pair])
	}
}

func (a *AlignmentTask) printInfo() {
	for _, node := range a.first {
		fmt.Printf("First tree. ID: %v, Full penalty: %v\n", node.ID, node.FullPenalty)
	}

	for _, node := range a.second {
		fmt.Printf("Second tree. ID: %v, Full penalty: %v\n", node.ID, node.FullPenalty)
	}
}

func CalculateAlignment(c *Config) error {
	// Build trees.
	firstTree := buildTree(c.FirstGraphDescription)
	secondTree := buildTree(c.SecondGraphDescription)

	// Get sorted pairs of subtrees.
	sortedPairs := getSortedTreePairs(
		len(c.FirstGraphDescription),
		len(c.SecondGraphDescription))

	a := AlignmentTask{
		config:      c,
		first:       firstTree,
		second:      secondTree,
		sortedPairs: sortedPairs,
		quality:     make(map[TreePair]int),
		qCase:       make(map[TreePair]int),
	}

	a.printInfo()

	// Calculate quality of this pairs.
	a.calculateQuality()

	a.printInfo()

	// Dynamically build alignment tree.
	return nil
}
