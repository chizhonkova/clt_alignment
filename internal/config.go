package alignment

type NodeDescription struct {
	ID          int    `json:"id"`
	Tag         string `json:"tag"`
	Left        *int   `json:"left"`
	Right       *int   `json:"right"`
	FullPenalty int    `json:"-"`
}

func (d *NodeDescription) LeftOrDefault() int {
	if d.Left == nil {
		return EmptyTreeID
	}
	return *d.Left
}

func (d *NodeDescription) RightOrDefault() int {
	if d.Right == nil {
		return EmptyTreeID
	}
	return *d.Right
}

type Config struct {
	FirstGraphDescription  []*NodeDescription
	SecondGraphDescription []*NodeDescription
	DeletionCost           int
	TagEqualityCost        int
	TagUnequalityCost      int

	ResultPath string
}

const (
	DeletionTag    = "-"
	EmptyTreeID    = -1
	DefaultPenalty = -1
	NoParent       = -1

	LeftChild  = 0
	RightChild = 1
	NotAChild  = -1

	BothEmptyCase  = 0
	LeftEmptyCase  = 1
	RightEmptyCase = 2

	// 1
	FirstACAndBDCase = 3
	FirstADAndBCCase = 4

	// 2
	SecondABAndCCase = 5
	SecondACAndBCase = 6
	SecondDFAndECase = 7
	SecondEFAndDCase = 8

	// 3
	ThirdCase = 9

	// 4
	FourthCase = 10

	// 5
	FifthACAndBCase = 11
	FifthBCAndACase = 12
	FifthDEAndFCase = 13
	FifthDFAndECase = 14

	// 6
	SixthFirstCase  = 15
	SixthSecondCase = 16

	// 7
	SeventhFirstCase  = 17
	SeventhSecondCase = 18
)
