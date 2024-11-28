package node

import (
	"errors"
)

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

type ConnParams struct {
	Dist  float64   // distance to destination
	Size  int       // size of connection
	State ConnState // state of the connection
}

// node should never be created directly, because the id is going to be invalid
type Node struct {
	Id    int                  // assigned automaticaly
	Pos   Pos                  // pos assigned at creation
	Conns map[*Node]ConnParams // ConnParams characterises connection to destination
}

type Graph struct {
	next_id int     // internal counter for id of next node
	Nodes   []*Node // array of all nodes
	Root    *Node   // pointer to root node
}

// creates a new graph with a root node position
func NewGraph(root_pos Pos) Graph {
	return Graph{
		next_id: 1,
		Root: &Node{
			Id:    0,
			Pos:   root_pos,
			Conns: map[*Node]ConnParams{},
		},
	}
}

// creates a new node at a position in graph and returns it with the correct id
func (g *Graph) NewNode(pos Pos) *Node {
	n := Node{
		Id:    g.next_id,
		Pos:   pos,
		Conns: map[*Node]ConnParams{},
	}

	g.Nodes = append(g.Nodes, &n)
	g.next_id += 1

	return &n
}

// already linked error
var AlreadyLinkedError = errors.New("Nodes already linked")

// links two nodes in a single direction with a connection with given parameters, fails if already linked
func (g *Graph) LinkOne(from, to *Node, params ConnParams) error {
	_, connected := from.Conns[to]
	if connected {
		return AlreadyLinkedError
	}

	from.Conns[to] = params
	return nil
}

// links two nodes in both directions with a connection with given parameters, fails if both already linked, does not fail if only one linked
func (g *Graph) LinkBoth(from, to *Node, params ConnParams) error {
	_, connected_from := from.Conns[to]
	_, connected_to := to.Conns[from]
	if connected_from && connected_to {
		return AlreadyLinkedError
	}

	from.Conns[to] = params
	to.Conns[from] = params

	return nil
}
