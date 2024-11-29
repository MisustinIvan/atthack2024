package traffic_test

import (
	"optitraffic/node"
	"optitraffic/traffic"
	"testing"
)

func TestTraffic(t *testing.T) {
    graph := node.NewGraph(node.Pos{X: 0, Y: 0})
    root := graph.Root
    node1 := graph.NewNode(node.Pos{X: -1, Y: 1})
    node2 := graph.NewNode(node.Pos{X: -1, Y: 2})
    node3 := graph.NewNode(node.Pos{X: 1, Y: 2})
    node4 := graph.NewNode(node.Pos{X: 1, Y: 1})
    node5 := graph.NewNode(node.Pos{X: 2, Y: 2})
    node6 := graph.NewNode(node.Pos{X: 0, Y: 3})

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

    tm := traffic.NewTrafficManager(&graph)

    tm.NewRandomVehicle(node.NormalVehicle)
    tm.NewRandomVehicle(node.NormalVehicle)
    tm.NewRandomVehicle(node.NormalVehicle)
    tm.NewRandomVehicle(node.EmergencyVehicle)
    for range 100 {
        tm.Update(1)
    }
}
