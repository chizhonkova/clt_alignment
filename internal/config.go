package alignment

type NodeDescription struct {
	ID    int    `json:"id"`
	Tag   string `json:"tag"`
	Left  *int   `json:"left"`
	Right *int   `json:"right"`
}

type Config struct {
	FirstGraphDescription  []NodeDescription
	SecondGraphDescription []NodeDescription
	DeletionCost           int
	TagEqualityCost        int
	TagUnequalityCost      int

	ResultPath string
}

const (
	DeletionTag = "-"
	EmptyTreeID = -1
)
