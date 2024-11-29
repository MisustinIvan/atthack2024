package graphdb

import (
	"database/sql"
	"fmt"
	"optitraffic/node"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var testPoses [3]node.Pos = [...]node.Pos{
        {X: 0, Y: 50},
        {X: 1, Y: 50},
        {X: 1, Y: 51}}

func TestStoreGraph(t *testing.T) {
    // Connect
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        t.Error(err.Error())
    }
    db.Exec("DELETE FROM graph_paths")
    db.Exec("DELETE FROM path_ends")
    db.Exec("DELETE FROM graph_nodes")
    dao := NewDAO(db)
    // Setup
    graph := node.NewGraph(testPoses[0])
    first, second, third := graph.Root, graph.NewNode(testPoses[1]), graph.NewNode(testPoses[2])
    graph.LinkBoth(first, second, 2); graph.LinkBoth(second, third, 1)
    // Act
    err = dao.StoreGraph(graph)
    // Assert
    if err != nil {
        t.Error(err.Error())
    }
}

func TestGetGraph(t *testing.T) {
    // Connect
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        t.Error(err.Error())
    }
    dao := NewDAO(db)
    // Setup
    graph := node.NewGraph(testPoses[0])
    first, second, third := graph.Root, graph.NewNode(testPoses[1]), graph.NewNode(testPoses[2])
    graph.LinkBoth(first, second, 2); graph.LinkBoth(second, third, 1)
    // act
    result, err := dao.GetGraph()
    if err != nil {
        t.Error(err.Error())
    }
    // assert
    if len(result.Nodes) != len(graph.Nodes) {
        t.Fatal("rozdilne delky")
    }
    leng := len(result.Nodes)
    for i := 0; i < leng; i++ {
        if !areEqualNodes(*result.Nodes[i], *graph.Nodes[i]) {
            t.Fatalf("wanted: %v; got: %v", *result.Nodes[i], *graph.Nodes[i])
        }
        fmt.Printf("a: %v, b: %v", (*result.Nodes[i]).Conns, (*graph.Nodes[i]).Conns)
    }
}

func areEqualNodes(a, b node.Node) bool {
    if a.Id != b.Id {
        return false
    }
    if a.Pos != b.Pos {
        return false
    }
    aNil, bNil := a.Conns == nil, b.Conns == nil
    if aNil != bNil {
        return false
    }

    return true
}