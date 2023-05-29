package alignment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (a *AlignmentTask) fixParents(pair TreePair, parent TreePair, childType TreePair) {
	if parent.Lhs != NoParent {
		if childType.Lhs == LeftChild {
			*a.firstResult[parent.Lhs].Left = pair.Lhs
		} else if childType.Lhs == RightChild {
			*a.firstResult[parent.Lhs].Right = pair.Lhs
		}
	}

	if parent.Rhs != NoParent {
		if childType.Rhs == LeftChild {
			*a.secondResult[parent.Rhs].Left = pair.Rhs
		} else if childType.Rhs == RightChild {
			*a.secondResult[parent.Rhs].Right = pair.Rhs
		}
	}
}

func (a *AlignmentTask) buildFirstCase(pair TreePair, parent TreePair, childType TreePair) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag

	newFirstID := a.firstCnt
	a.firstCnt++
	newSecondID := a.secondCnt
	a.secondCnt++

	firstLeft := EmptyTreeID
	firstRight := EmptyTreeID
	secondLeft := EmptyTreeID
	secondRight := EmptyTreeID

	a.firstResult[newFirstID] = &NodeDescription{
		ID:    newFirstID,
		Tag:   firstTag,
		Left:  &firstLeft,
		Right: &firstRight,
	}
	a.secondResult[newSecondID] = &NodeDescription{
		ID:    newSecondID,
		Tag:   secondTag,
		Left:  &secondLeft,
		Right: &secondRight,
	}

	a.fixParents(TreePair{Lhs: newFirstID, Rhs: newSecondID}, parent, childType)

	aID := a.first[firstID].LeftOrDefault()
	bID := a.first[firstID].RightOrDefault()
	cID := a.second[secondID].LeftOrDefault()
	dID := a.second[secondID].RightOrDefault()

	c := a.qCase[pair]

	if c == FirstACAndBDCase {
		a.buildNextPair(
			TreePair{Lhs: aID, Rhs: cID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: LeftChild, Rhs: LeftChild})
		a.buildNextPair(
			TreePair{Lhs: bID, Rhs: dID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: RightChild, Rhs: RightChild})
	} else {
		a.buildNextPair(
			TreePair{Lhs: aID, Rhs: dID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: LeftChild, Rhs: RightChild})
		a.buildNextPair(
			TreePair{Lhs: bID, Rhs: cID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: RightChild, Rhs: LeftChild})
	}
}

func (a *AlignmentTask) buildSecondCase(pair TreePair, parent TreePair, childType TreePair) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag

	newFirstID := a.firstCnt
	a.firstCnt++
	newSecondID := a.secondCnt
	a.secondCnt++

	firstLeft := EmptyTreeID
	firstRight := EmptyTreeID
	secondLeft := EmptyTreeID
	secondRight := EmptyTreeID

	a.firstResult[newFirstID] = &NodeDescription{
		ID:    newFirstID,
		Tag:   firstTag,
		Left:  &firstLeft,
		Right: &firstRight,
	}
	a.secondResult[newSecondID] = &NodeDescription{
		ID:    newSecondID,
		Tag:   secondTag,
		Left:  &secondLeft,
		Right: &secondRight,
	}

	a.fixParents(TreePair{Lhs: newFirstID, Rhs: newSecondID}, parent, childType)

	aDelID := firstID
	if aDelID < len(a.first)/2 {
		aDelID += len(a.first) / 2
	}
	bID := a.second[secondID].LeftOrDefault()
	cID := a.second[secondID].RightOrDefault()

	// Symmetric case.
	dID := a.first[firstID].LeftOrDefault()
	eID := a.first[firstID].RightOrDefault()
	fDelID := secondID
	if fDelID < len(a.second)/2 {
		fDelID += len(a.second) / 2
	}

	c := a.qCase[pair]

	if c == SecondABAndCCase {
		a.buildNextPair(
			TreePair{Lhs: aDelID, Rhs: bID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: LeftChild, Rhs: LeftChild})
		// Build deletion tree and C
		a.buildDeletionPair(
			TreePair{Lhs: EmptyTreeID, Rhs: cID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: RightChild, Rhs: RightChild})

	} else if c == SecondACAndBCase {
		a.buildNextPair(
			TreePair{Lhs: aDelID, Rhs: cID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: LeftChild, Rhs: RightChild})
		a.buildDeletionPair(
			TreePair{Lhs: EmptyTreeID, Rhs: bID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: RightChild, Rhs: LeftChild})

	} else if c == SecondDFAndECase {
		a.buildNextPair(
			TreePair{Lhs: dID, Rhs: fDelID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: LeftChild, Rhs: LeftChild})
		a.buildDeletionPair(
			TreePair{Lhs: eID, Rhs: EmptyTreeID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: RightChild, Rhs: RightChild})

	} else if c == SecondEFAndDCase {
		a.buildNextPair(
			TreePair{Lhs: eID, Rhs: fDelID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: RightChild, Rhs: LeftChild})
		a.buildDeletionPair(
			TreePair{Lhs: dID, Rhs: EmptyTreeID},
			TreePair{Lhs: newFirstID, Rhs: newSecondID},
			TreePair{Lhs: LeftChild, Rhs: RightChild})
	}
}

func (a *AlignmentTask) buildFifthCase(pair TreePair, parent TreePair, childType TreePair) {
	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag

	newFirstID := a.firstCnt
	a.firstCnt++
	newSecondID := a.secondCnt
	a.secondCnt++

	firstLeft := EmptyTreeID
	firstRight := EmptyTreeID
	secondLeft := EmptyTreeID
	secondRight := EmptyTreeID

	a.fixParents(TreePair{Lhs: newFirstID, Rhs: newSecondID}, parent, childType)

	c := a.qCase[pair]

	if c == FifthACAndBCase || c == FifthBCAndACase {
		a.firstResult[newFirstID] = &NodeDescription{
			ID:    newFirstID,
			Tag:   firstTag,
			Left:  &firstLeft,
			Right: &firstRight,
		}
		a.secondResult[newSecondID] = &NodeDescription{
			ID:    newSecondID,
			Tag:   DeletionTag,
			Left:  &secondLeft,
			Right: &secondRight,
		}

		aID := a.first[firstID].LeftOrDefault()
		bID := a.first[firstID].RightOrDefault()

		if c == FifthACAndBCase {
			a.buildNextPair(
				TreePair{Lhs: aID, Rhs: secondID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: LeftChild, Rhs: LeftChild})
			a.buildDeletionPair(
				TreePair{Lhs: bID, Rhs: EmptyTreeID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: RightChild, Rhs: RightChild})
		} else if c == FifthBCAndACase {
			a.buildNextPair(
				TreePair{Lhs: bID, Rhs: secondID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: RightChild, Rhs: LeftChild})
			a.buildDeletionPair(
				TreePair{Lhs: aID, Rhs: EmptyTreeID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: LeftChild, Rhs: RightChild})
		}

	} else if c == FifthDEAndFCase || c == FifthDFAndECase {
		a.firstResult[newFirstID] = &NodeDescription{
			ID:    newFirstID,
			Tag:   DeletionTag,
			Left:  &firstLeft,
			Right: &firstRight,
		}
		a.secondResult[newSecondID] = &NodeDescription{
			ID:    newSecondID,
			Tag:   secondTag,
			Left:  &secondLeft,
			Right: &secondRight,
		}

		eID := a.second[secondID].LeftOrDefault()
		fID := a.second[secondID].RightOrDefault()

		if c == FifthDEAndFCase {
			a.buildNextPair(
				TreePair{Lhs: firstID, Rhs: eID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: LeftChild, Rhs: LeftChild})
			a.buildDeletionPair(
				TreePair{Lhs: EmptyTreeID, Rhs: fID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: RightChild, Rhs: RightChild})
		} else if c == FifthDFAndECase {
			a.buildNextPair(
				TreePair{Lhs: firstID, Rhs: fID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: LeftChild, Rhs: RightChild})
			a.buildDeletionPair(
				TreePair{Lhs: EmptyTreeID, Rhs: eID},
				TreePair{Lhs: newFirstID, Rhs: newSecondID},
				TreePair{Lhs: RightChild, Rhs: LeftChild})
		}
	}
}

func (a *AlignmentTask) buildDeletionPair(pair TreePair, parent TreePair, childType TreePair) {
	fmt.Printf("Deletion subtree pair. Lhs: %v, Rhs: %v.\n", pair.Lhs, pair.Rhs)

	if pair.Lhs == EmptyTreeID && pair.Rhs == EmptyTreeID {
		return
	}

	firstID := pair.Lhs
	secondID := pair.Rhs

	firstTag := a.first[firstID].Tag
	secondTag := a.second[secondID].Tag

	newFirstID := a.firstCnt
	a.firstCnt++
	newSecondID := a.secondCnt
	a.secondCnt++

	firstLeft := EmptyTreeID
	firstRight := EmptyTreeID
	secondLeft := EmptyTreeID
	secondRight := EmptyTreeID

	a.firstResult[newFirstID] = &NodeDescription{
		ID:    newFirstID,
		Tag:   firstTag,
		Left:  &firstLeft,
		Right: &firstRight,
	}
	a.secondResult[newSecondID] = &NodeDescription{
		ID:    newSecondID,
		Tag:   secondTag,
		Left:  &secondLeft,
		Right: &secondRight,
	}

	a.fixParents(TreePair{Lhs: newFirstID, Rhs: newSecondID}, parent, childType)

	aID := a.first[firstID].LeftOrDefault()
	bID := a.first[firstID].RightOrDefault()
	cID := a.second[secondID].LeftOrDefault()
	dID := a.second[secondID].RightOrDefault()

	a.buildDeletionPair(
		TreePair{Lhs: aID, Rhs: cID},
		TreePair{Lhs: newFirstID, Rhs: newSecondID},
		TreePair{Lhs: LeftChild, Rhs: LeftChild})
	a.buildDeletionPair(
		TreePair{Lhs: bID, Rhs: dID},
		TreePair{Lhs: newFirstID, Rhs: newSecondID},
		TreePair{Lhs: RightChild, Rhs: RightChild})
}

func (a *AlignmentTask) buildNextPair(pair TreePair, parent TreePair, childType TreePair) {
	fmt.Printf("Tree pair. Lhs: %v, Rhs: %v.\n", pair.Lhs, pair.Rhs)

	if pair.Lhs == EmptyTreeID && pair.Rhs == EmptyTreeID {
		return
	}

	c := a.qCase[pair]

	if c == LeftEmptyCase || c == RightEmptyCase {
		a.buildDeletionPair(pair, parent, childType)
	}

	if c == FirstACAndBDCase || c == FirstADAndBCCase {
		a.buildFirstCase(pair, parent, childType)
		return
	}

	if c >= SecondABAndCCase && c <= SecondEFAndDCase {
		a.buildSecondCase(pair, parent, childType)
		return
	}

	if c >= FifthACAndBCase && c <= FifthDFAndECase {
		a.buildFifthCase(pair, parent, childType)
	}
}

func (a *AlignmentTask) buildResult() {
	a.firstResult = NewTree()
	a.secondResult = NewTree()

	// We start with (0, 0) TreePair
	pair := TreePair{Lhs: 0, Rhs: 0}
	a.buildNextPair(
		pair,
		TreePair{Lhs: NoParent, Rhs: NoParent},
		TreePair{Lhs: NotAChild, Rhs: NotAChild})
}

func (a *AlignmentTask) flushResult() {
	fr := []*NodeDescription{}
	for _, node := range a.firstResult {
		fr = append(fr, node)
	}

	sr := []*NodeDescription{}
	for _, node := range a.secondResult {
		sr = append(sr, node)
	}

	r := [][]*NodeDescription{fr, sr}
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(a.config.ResultPath, b, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result is written to %v.\n", a.config.ResultPath)
}
