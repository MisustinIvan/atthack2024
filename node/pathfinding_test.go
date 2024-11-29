package node_test

import (
	"fmt"
	"optitraffic/node"
	"testing"
)

func TestPathfinding(t *testing.T) {
    graph := node.NewGraph(node.Pos{0, 0})
    root := graph.Root
    node1 := graph.NewNode(node.Pos{-1, 1})
    node2 := graph.NewNode(node.Pos{-1, 2})
    node3 := graph.NewNode(node.Pos{1, 2})
    node4 := graph.NewNode(node.Pos{1, 1})
    node5 := graph.NewNode(node.Pos{2, 2})
    node6 := graph.NewNode(node.Pos{0, 3})

    graph.LinkBoth(root, node1, 1)
    graph.LinkBoth(root, node4, 1)
    graph.LinkBoth(node1, node4, 1)
    graph.LinkBoth(node1, node6, 1)
    graph.LinkBoth(node1, node2, 1)
    graph.LinkBoth(node4, node5, 1)
    graph.LinkBoth(node4, node3, 1)
    graph.LinkBoth(node5, node6, 1)
    graph.LinkBoth(node3, node6, 1)
    graph.LinkBoth(node3, node2, 1)

    pf := node.NewPathfinder(&graph)
    _, err := pf.Path(root, node6, node.NormalVehicle)
    if err != nil {
        for _, node := range graph.Nodes {
            fmt.Printf("%v\n", *node)
        }
        t.Fatal(err)
    }

    //fmt.Printf("node.PathLength(path): %v\n", node.PathLength(path))
}
