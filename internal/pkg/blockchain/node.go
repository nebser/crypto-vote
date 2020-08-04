package blockchain

type NodeType int

const (
	AlfaNodeType NodeType = iota
	RegularNodeType
)

type SaveNodeFn func(Node) error

type GetNodeFn func(string) (*Node, error)

type GetNodesFn func() (Nodes, error)

func (n NodeType) String() string {
	switch n {
	case AlfaNodeType:
		return "alfa"
	case RegularNodeType:
		return "regular"
	default:
		return ""
	}
}

type Node struct {
	ID   string   `json:"id"`
	Type NodeType `json:"type"`
}

type Nodes []Node

func (nodes Nodes) Regulars() (result Nodes) {
	for _, n := range nodes {
		if n.Type != AlfaNodeType {
			result = append(result, n)
		}
	}
	return
}
