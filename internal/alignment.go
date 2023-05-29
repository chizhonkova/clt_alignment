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

	firstResult  Tree
	firstCnt     int
	secondResult Tree
	secondCnt    int
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
			pairs = append(pairs, TreePair{Lhs: i + first, Rhs: j + second})
			pairs = append(pairs, TreePair{Lhs: i, Rhs: j + second})
			pairs = append(pairs, TreePair{Lhs: i + first, Rhs: j})
			pairs = append(pairs, TreePair{Lhs: i, Rhs: j})
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

func findMaxQuality(qualities, cases []int) (int, int) {
	maxQ := qualities[0]
	maxCase := cases[0]

	for i, quality := range qualities {
		if quality > maxQ {
			maxQ = quality
			maxCase = cases[i]
		}
	}

	return maxQ, maxCase
}

func (a *AlignmentTask) findRootTagsCost(firstTag, secondTag string) int {
	tagCost := a.config.TagEqualityCost
	if firstTag == DeletionTag && secondTag == DeletionTag {
		tagCost = 0
	}
	if firstTag != secondTag {
		if firstTag != DeletionTag && secondTag != DeletionTag {
			tagCost = -a.config.TagUnequalityCost
		} else {
			tagCost = -a.config.DeletionCost
		}
	}
	return tagCost
}

func (a *AlignmentTask) findFirstCaseQuality(pair TreePair) (int, int) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag
	tagCost := a.findRootTagsCost(firstTag, secondTag)

	aID := a.first[firstID].LeftOrDefault()
	bID := a.first[firstID].RightOrDefault()
	cID := a.second[secondID].LeftOrDefault()
	dID := a.second[secondID].RightOrDefault()
	aCAndBd := tagCost + a.quality[TreePair{Lhs: aID, Rhs: cID}] + a.quality[TreePair{Lhs: bID, Rhs: dID}]
	adAndBc := tagCost + a.quality[TreePair{Lhs: aID, Rhs: dID}] + a.quality[TreePair{Lhs: bID, Rhs: cID}]

	if aCAndBd > adAndBc {
		return aCAndBd, FirstACAndBDCase
	} else {
		return adAndBc, FirstADAndBCCase
	}
}

func (a *AlignmentTask) findSecondCaseQuality(pair TreePair) (int, int) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag
	tagCost := a.findRootTagsCost(firstTag, secondTag)

	aDelID := firstID
	if aDelID < len(a.first)/2 {
		aDelID += len(a.first) / 2
	}
	bID := a.second[secondID].LeftOrDefault()
	cID := a.second[secondID].RightOrDefault()

	aBAndC := tagCost + a.quality[TreePair{Lhs: aDelID, Rhs: bID}] - a.second[cID].FullPenalty
	aCAndB := tagCost + a.quality[TreePair{Lhs: aDelID, Rhs: cID}] - a.second[bID].FullPenalty

	// Symmetric case.
	dID := a.first[firstID].LeftOrDefault()
	eID := a.first[firstID].RightOrDefault()
	fDelID := secondID
	if fDelID < len(a.second)/2 {
		fDelID += len(a.second) / 2
	}

	dFAndE := tagCost + a.quality[TreePair{Lhs: dID, Rhs: fDelID}] - a.first[eID].FullPenalty
	eFAndD := tagCost + a.quality[TreePair{Lhs: eID, Rhs: fDelID}] - a.first[dID].FullPenalty

	qualitites := []int{aBAndC, aCAndB, dFAndE, eFAndD}
	cases := []int{SecondABAndCCase, SecondACAndBCase, SecondDFAndECase, SecondEFAndDCase}
	return findMaxQuality(qualitites, cases)
}

func (a *AlignmentTask) findThirdCaseQuality(pair TreePair) (int, int) {
	// This case is sutable only if one of the root tags is not a deletion.

	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag
	tagCost := a.findRootTagsCost(firstTag, secondTag)

	aDelID := firstID
	if aDelID < len(a.first)/2 {
		aDelID += len(a.first) / 2
	}
	bDelId := secondID
	if bDelId < len(a.second)/2 {
		bDelId += len(a.second) / 2
	}

	quality := tagCost + a.quality[TreePair{Lhs: aDelID, Rhs: bDelId}]
	return quality, ThirdCase
}

func (a *AlignmentTask) findFourthCaseQuality(pair TreePair) (int, int) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag
	tagCost := a.findRootTagsCost(firstTag, secondTag)

	aDelID := firstID
	if aDelID < len(a.first)/2 {
		aDelID += len(a.first) / 2
	}
	bDelId := secondID
	if bDelId < len(a.second)/2 {
		bDelId += len(a.second) / 2
	}

	quality := tagCost - a.first[aDelID].FullPenalty - a.second[bDelId].FullPenalty
	return quality, FourthCase
}

func (a *AlignmentTask) findFifthCaseQuality(pair TreePair) (int, int) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	aID := a.first[firstID].LeftOrDefault()
	bID := a.first[firstID].RightOrDefault()
	tagCost := a.findRootTagsCost(a.first[firstID].Tag, DeletionTag)

	acAndB := a.quality[TreePair{Lhs: aID, Rhs: secondID}] + tagCost - a.first[bID].FullPenalty
	bcAndA := a.quality[TreePair{Lhs: bID, Rhs: secondID}] + tagCost - a.first[aID].FullPenalty

	eID := a.second[secondID].LeftOrDefault()
	fID := a.second[secondID].RightOrDefault()
	tagCost = a.findRootTagsCost(DeletionTag, a.second[secondID].Tag)

	deAndF := a.quality[TreePair{Lhs: firstID, Rhs: eID}] + tagCost - a.second[fID].FullPenalty
	dfAndE := a.quality[TreePair{Lhs: firstID, Rhs: fID}] + tagCost - a.second[eID].FullPenalty

	qualities := []int{acAndB, bcAndA, deAndF, dfAndE}
	cases := []int{FifthACAndBCase, FifthBCAndACase, FifthDEAndFCase, FifthDFAndECase}
	return findMaxQuality(qualities, cases)
}

func (a *AlignmentTask) findSixthCaseQuality(pair TreePair) (int, int) {
	qualities := []int{}
	cases := []int{}

	firstID := pair.Lhs
	secondID := pair.Rhs
	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag

	if firstTag != DeletionTag {
		aDelID := firstID
		if aDelID < len(a.first)/2 {
			aDelID += len(a.first) / 2
		}
		quality := a.quality[TreePair{Lhs: aDelID, Rhs: secondID}] - a.config.DeletionCost
		qualities = append(qualities, quality)
		cases = append(cases, SixthFirstCase)
	}

	if secondTag != DeletionTag {
		bDelID := secondID
		if bDelID < len(a.second)/2 {
			bDelID += len(a.second) / 2
		}
		quality := a.quality[TreePair{Lhs: firstID, Rhs: bDelID}] - a.config.DeletionCost
		qualities = append(qualities, quality)
		cases = append(cases, SixthSecondCase)
	}

	return findMaxQuality(qualities, cases)
}

func (a *AlignmentTask) findSeventhCaseQuality(pair TreePair) (int, int) {
	firstID := pair.Lhs
	secondID := pair.Rhs
	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag

	tagCost := a.findRootTagsCost(firstTag, DeletionTag)
	aDelID := firstID
	if aDelID < len(a.first)/2 {
		aDelID += len(a.first) / 2
	}
	firstQuality := tagCost - a.first[aDelID].FullPenalty - a.second[secondID].FullPenalty

	tagCost = a.findRootTagsCost(DeletionTag, secondTag)
	bDelID := secondID
	if bDelID < len(a.second)/2 {
		bDelID += len(a.second) / 2
	}
	secondQuality := tagCost - a.first[firstID].FullPenalty - a.second[bDelID].FullPenalty

	if firstQuality > secondQuality {
		return firstQuality, SeventhFirstCase
	} else {
		return secondQuality, SeventhSecondCase
	}
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

	qualities := []int{}
	cases := []int{}

	// 1
	q, c := a.findFirstCaseQuality(pair)
	qualities = append(qualities, q)
	cases = append(cases, c)

	// 2
	q, c = a.findSecondCaseQuality(pair)
	qualities = append(qualities, q)
	cases = append(cases, c)

	// 3
	if a.first[pair.Lhs].Tag != DeletionTag || a.second[pair.Rhs].Tag != DeletionTag {
		q, c = a.findThirdCaseQuality(pair)
		qualities = append(qualities, q)
		cases = append(cases, c)
	}

	// 4
	q, c = a.findFourthCaseQuality(pair)
	qualities = append(qualities, q)
	cases = append(cases, c)

	// 5
	q, c = a.findFifthCaseQuality(pair)
	qualities = append(qualities, q)
	cases = append(cases, c)

	// 6
	if a.first[pair.Lhs].Tag != DeletionTag || a.second[pair.Rhs].Tag != DeletionTag {
		q, c = a.findSixthCaseQuality(pair)
		qualities = append(qualities, q)
		cases = append(cases, c)
	}

	// 7
	q, c = a.findSeventhCaseQuality(pair)
	qualities = append(qualities, q)
	cases = append(cases, c)

	return findMaxQuality(qualities, cases)
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
	fmt.Printf("Quality: %v\n", a.quality[TreePair{Lhs: 0, Rhs: 0}])
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

	// Calculate quality of this pairs.
	fmt.Println("Calculating maximum quality...")

	a.calculateQuality()

	a.printInfo()

	// Dynamically build alignment tree.
	fmt.Println("Building alignment...")

	a.buildResult()
	a.flushResult()

	return nil
}
