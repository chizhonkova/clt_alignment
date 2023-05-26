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

	BothEmptyCase  = 0
	LeftEmptyCase  = 1
	RightEmptyCase = 2
)
