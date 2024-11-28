package node

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ConnState int

const (
	Open ConnState = iota
	Closed
	Danger
)

type Conn struct {
	Dist  float64   // distance to destination
	Size  int       // size of connection
	State ConnState // state of the connection
}

type Node struct {
	Pos   Pos
	Conns map[*Node]Conn // conn characterises connection to destination
}

func (n *Node) Link(dest *Node, dist float64, size int, state ConnState) {
	n.Conns[dest] = Conn{
		Dist:  dist,
		Size:  size,
		State: state,
	}
}

func (n *Node) Unlink(dest *Node) {
	delete(n.Conns, dest)
}

type Graph struct {
	Root *Node
}

func NewGraph(root_pos Pos) Graph {
	return Graph{
		Root: &Node{
			Pos:   root_pos,
			Conns: map[*Node]Conn{},
		},
	}
}
