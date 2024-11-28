package node

import (
	"errors"
	"math"
)

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// calculates the distance between two points
func (p Pos) DistanceTo(o Pos) float64 {
	dx := p.X - o.X
	dy := p.Y - o.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type ConnState int

const (
	Open ConnState = iota
	Closed
	Danger
)

// describes connection between two nodes
type ConnParams struct {
	Dist  float64   // distance to destination
	Size  int       // size of connection
	State ConnState // state of the connection
	NCars int       // amount of cars on connection
}

// Calculates weight of connection based on its parameters
func (p ConnParams) ConnectionWeight() float64 {
	if p.State == Closed || p.State == Danger {
		return math.MaxInt
	}

	return (p.Dist * (1 + (float64(p.NCars) / float64(p.Size))))
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
	root := Node{
		Id:    0,
		Pos:   root_pos,
		Conns: map[*Node]ConnParams{},
	}

	return Graph{
		Nodes:   []*Node{&root},
		next_id: 1,
		Root:    &root,
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

// links two nodes in a single direction with a connection with given connection size, fails if already linked
func (g *Graph) LinkOne(from, to *Node, size int) error {
	_, connected := from.Conns[to]
	if connected {
		return AlreadyLinkedError
	}

	from.Conns[to] = ConnParams{
		Dist:  from.Pos.DistanceTo(to.Pos),
		Size:  size,
		State: Open,
		NCars: 0,
	}
	return nil
}

// links two nodes in both directions with a connection with given size, fails if both already linked, does not fail if only one linked
func (g *Graph) LinkBoth(from, to *Node, size int) error {
	_, connected_from := from.Conns[to]
	_, connected_to := to.Conns[from]
	if connected_from && connected_to {
		return AlreadyLinkedError
	}

	params := ConnParams{
		Dist:  from.Pos.DistanceTo(to.Pos),
		Size:  size,
		State: Open,
		NCars: 0,
	}

	from.Conns[to] = params
	to.Conns[from] = params

	return nil
}
