package graphdb

import (
	"database/sql"
	conv "optitraffic/graphConvertor"
	"optitraffic/node"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

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
    poses := [...]node.Pos{
        {X: 0, Y: 50},
        {X: 1, Y: 50},
        {X: 1, Y: 51},
    }
    graph := node.NewGraph(poses[0])
    first, second, third := graph.Root, graph.NewNode(poses[1]), graph.NewNode(poses[2])
    graph.LinkBoth(first, second, 2); graph.LinkBoth(second, third, 1)
    // Act
    err = dao.StoreGraph(graph)
    // Assert
    if err != nil {
        t.Error(err.Error())
    }
}

func TestIndependent(t *testing.T) {
	// Connect
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        t.Error(err.Error())
    }
    db.Exec("DELETE FROM graph_paths")
    db.Exec("DELETE FROM path_ends")
    db.Exec("DELETE FROM graph_nodes")
    // dao := NewDAO(db)
    // // Setup
    // poses := [...]node.Pos{
    //     {X: 0, Y: 50},
    //     {X: 1, Y: 50},
    //     {X: 1, Y: 51},
    // }

}

func TestGetGraph(t *testing.T) {
    // Connect
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        t.Error(err.Error())
    }
    dao := NewDAO(db)
    // act & assert
    paths, points, err := dao.GetGraph()
    if err != nil {
        t.Fail()
    }
    _, err = conv.GeoJSONToGraph(paths, points)
    if err != nil {
        t.Fail()
    }

}